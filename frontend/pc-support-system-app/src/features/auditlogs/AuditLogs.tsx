import { useEffect, useMemo, useState } from "react";
import {
    Alert,
    Card,
    Form,
    Pagination,
    Spinner,
    Table,
} from "react-bootstrap";

import { useAuth } from "../auth/hooks/useAuth";

import {
    fetchAuditLogs,
    fetchCustomer,
    fetchDevice,
    fetchJob,
    fetchStaff,
} from "./auditLogApi";

import type {
    AuditLog,
    AuditLogEntity,
} from "./auditLogTypes";

const RECORDS_PER_PAGE = 30;

const ENTITY_OPTIONS: {
    value: AuditLogEntity | "";
    label: string;
}[] = [
    {
        value: "",
        label: "All",
    },
    {
        value: "customer",
        label: "Customer",
    },
    {
        value: "device",
        label: "Device",
    },
    {
        value: "job",
        label: "Job",
    },
    {
        value: "user",
        label: "User",
    },
];

type DisplayLog = AuditLog & {
    performedByDisplay: string;
    entityDisplay: string;
};

export default function AuditLogs() {
    const { user } = useAuth();

    const [logs, setLogs] = useState<AuditLog[]>([]);
    const [displayLogs, setDisplayLogs] = useState<DisplayLog[]>(
        [],
    );

    const [selectedEntity, setSelectedEntity] =
        useState<AuditLogEntity | "">("");

    const [loading, setLoading] = useState(true);
    const [resolving, setResolving] = useState(false);

    const [error, setError] = useState("");

    const [currentPage, setCurrentPage] = useState(1);

    /*
     * Admin and Super Admin only.
     *
     * Change these role strings only if your backend uses
     * different role values.
     */
    const allowed =
        user?.role === "admin" ||
        user?.role === "super_admin";

    async function loadLogs() {
        try {
            setLoading(true);
            setError("");

            const response = await fetchAuditLogs(
                selectedEntity || undefined,
            );

            setLogs(response.logs);
            setCurrentPage(1);
        } catch (err: any) {
            console.error(
                "Audit logs error:",
                err.response?.data,
            );

            setError(
                err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to load audit logs.",
            );
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        if (!allowed) {
            setLoading(false);
            return;
        }

          if (!selectedEntity) {
        setLogs([]);
        setDisplayLogs([]);
        setLoading(false);
        return;
    }
    
        loadLogs();
    }, [selectedEntity, allowed]);

    /*
     * Cache resolved objects so the same BSON ID is not requested
     * repeatedly.
     */
    const staffCache = useMemo(
        () => new Map<string, string>(),
        [],
    );

    const customerCache = useMemo(
        () => new Map<string, string>(),
        [],
    );

    const deviceCache = useMemo(
        () => new Map<string, string>(),
        [],
    );

    const jobCache = useMemo(
        () => new Map<string, string>(),
        [],
    );

    async function resolvePerformedBy(
        staffId: string,
    ): Promise<string> {
        if (!staffId) {
            return "--";
        }

        const cached = staffCache.get(staffId);

        if (cached) {
            return cached;
        }

        try {
            const staff = await fetchStaff(staffId);

            const name =
                [
                    staff.firstName,
                    staff.lastName,
                ]
                    .filter(Boolean)
                    .join(" ") ||
                staff.name ||
                staff.username ||
                staff.email ||
                staffId;

            staffCache.set(staffId, name);

            return name;
        } catch {
            return staffId;
        }
    }

    async function resolveEntity(
        entityType: AuditLogEntity,
        entityId: string,
    ): Promise<string> {
        if (!entityId) {
            return "--";
        }

        try {
            if (entityType === "customer") {
                const cached =
                    customerCache.get(entityId);

                if (cached) {
                    return cached;
                }

                const customer =
                    await fetchCustomer(entityId);

                const name =
                    [
                        customer.firstName,
                        customer.lastName,
                    ]
                        .filter(Boolean)
                        .join(" ") ||
                    customer.name ||
                    entityId;

                customerCache.set(
                    entityId,
                    name,
                );

                return name;
            }

            if (entityType === "device") {
                const cached =
                    deviceCache.get(entityId);

                if (cached) {
                    return cached;
                }

                const device =
                    await fetchDevice(entityId);

                /*
                 * Display what the device API gives us,
                 * rather than inventing a new identifier.
                 */
                const parts = [
                    device.type,
                    device.brand,
                    device.model,
                    device.serialNumber,
                ].filter(Boolean);

                const display =
                    parts.length > 0
                        ? parts.join(" - ")
                        : entityId;

                deviceCache.set(
                    entityId,
                    display,
                );

                return display;
            }

            if (entityType === "job") {
                const cached =
                    jobCache.get(entityId);

                if (cached) {
                    return cached;
                }

                const job =
                    await fetchJob(entityId);

                const display =
                    job.jobNumber ??
                    job.JobNumber ??
                    job.number ??
                    job.Number ??
                    entityId;

                jobCache.set(
                    entityId,
                    String(display),
                );

                return String(display);
            }

            if (entityType === "user") {
                /*
                 * User audit entities use the staff endpoint
                 * for now.
                 */
                return resolvePerformedBy(
                    entityId,
                );
            }

            return entityId;
        } catch {
            /*
             * If an entity cannot be resolved, keep the
             * original BSON ID visible.
             */
            return entityId;
        }
    }

    useEffect(() => {
        let cancelled = false;

        async function resolveLogs() {
            if (!logs.length) {
                setDisplayLogs([]);
                return;
            }

            try {
                setResolving(true);

                const resolved =
                    await Promise.all(
                        logs.map(async (log) => {
                            const [
                                performedByDisplay,
                                entityDisplay,
                            ] =
                                await Promise.all([
                                    resolvePerformedBy(
                                        log.PerformedByID,
                                    ),
                                    resolveEntity(
                                        log.EntityType,
                                        log.EntityID,
                                    ),
                                ]);

                            return {
                                ...log,
                                performedByDisplay,
                                entityDisplay,
                            };
                        }),
                    );

                if (!cancelled) {
                    setDisplayLogs(resolved);
                }
            } finally {
                if (!cancelled) {
                    setResolving(false);
                }
            }
        }

        resolveLogs();

        return () => {
            cancelled = true;
        };
    }, [logs]);

    const totalPages = Math.ceil(
        displayLogs.length / RECORDS_PER_PAGE,
    );

    const startIndex =
        (currentPage - 1) *
        RECORDS_PER_PAGE;

    const paginatedLogs =
        displayLogs.slice(
            startIndex,
            startIndex + RECORDS_PER_PAGE,
        );

    function formatDate(
        value: string,
    ) {
        return new Date(value).toLocaleString();
    }

    function formatEventType(
        value: string,
    ) {
        return value
            .replace(/_/g, " ")
            .replace(/\b\w/g, (letter) =>
                letter.toUpperCase(),
            );
    }

    function handleEntityChange(
        value: string,
    ) {
        setSelectedEntity(
            value as AuditLogEntity | "",
        );
        setCurrentPage(1);
    }

    if (!allowed) {
        return (
            <Alert variant="danger">
                You do not have permission to view
                audit logs.
            </Alert>
        );
    }

    if (loading) {
        return (
            <div className="text-center p-4">
                <Spinner />
            </div>
        );
    }

    return (
        <Card>
            <Card.Body>
                <div className="d-flex justify-content-between align-items-center mb-4">
                    <Card.Title className="mb-0">
                        Audit Logs
                    </Card.Title>

                    <Form.Select
                        style={{
                            maxWidth: "220px",
                        }}
                        value={selectedEntity}
                        onChange={(event) =>
                            handleEntityChange(
                                event.target.value,
                            )
                        }
                    >
                        {ENTITY_OPTIONS.map(
                            (option) => (
                                <option
                                    key={
                                        option.value
                                    }
                                    value={
                                        option.value
                                    }
                                >
                                    {option.label}
                                </option>
                            ),
                        )}
                    </Form.Select>
                </div>

                {error && (
                    <Alert variant="danger">
                        {error}
                    </Alert>
                )}

                {resolving && (
                    <Alert variant="info">
                        Resolving audit log details...
                    </Alert>
                )}

                {!paginatedLogs.length ? (
                    <Alert variant="info">
                        Select the log type from the drop down
                    </Alert>
                ) : (
                    <>
                        <div className="table-responsive">
                            <Table
                                bordered
                                hover
                                responsive
                            >
                                <thead>
                                    <tr>
                                        <th>
                                            Date
                                        </th>
                                        <th>
                                            Performed By
                                        </th>
                                        <th>
                                            Entity Type
                                        </th>
                                        <th>
                                            Entity
                                        </th>
                                        <th>
                                            Event
                                        </th>
                                        <th>
                                            Details
                                        </th>
                                    </tr>
                                </thead>

                                <tbody>
                                    {paginatedLogs.map(
                                        (log) => (
                                            <tr
                                                key={
                                                    log.ID
                                                }
                                            >
                                                <td>
                                                    {formatDate(
                                                        log.CreatedAt,
                                                    )}
                                                </td>

                                                <td>
                                                    {
                                                        log.performedByDisplay
                                                    }
                                                </td>

                                                <td>
                                                    {
                                                        log.EntityType
                                                    }
                                                </td>

                                                <td>
                                                    {
                                                        log.entityDisplay
                                                    }
                                                </td>

                                                <td>
                                                    {formatEventType(
                                                        log.EventType,
                                                    )}
                                                </td>

                                                <td>
                                                    <pre
                                                        className="mb-0"
                                                        style={{
                                                            whiteSpace:
                                                                "pre-wrap",
                                                        }}
                                                    >
                                                        {JSON.stringify(
                                                            log.Metadata ??
                                                                {},
                                                            null,
                                                            2,
                                                        )}
                                                    </pre>
                                                </td>
                                            </tr>
                                        ),
                                    )}
                                </tbody>
                            </Table>
                        </div>

                        <div className="d-flex justify-content-between align-items-center mt-3">
                            <div className="text-muted">
                                Showing{" "}
                                {startIndex +
                                    1}
                                –
                                {Math.min(
                                    startIndex +
                                        RECORDS_PER_PAGE,
                                    displayLogs.length,
                                )}{" "}
                                of{" "}
                                {
                                    displayLogs.length
                                }
                            </div>

                            {totalPages > 1 && (
                                <Pagination className="mb-0">
                                    <Pagination.Prev
                                        disabled={
                                            currentPage ===
                                            1
                                        }
                                        onClick={() =>
                                            setCurrentPage(
                                                (
                                                    page,
                                                ) =>
                                                    page -
                                                    1,
                                            )
                                        }
                                    />

                                    {Array.from(
                                        {
                                            length:
                                                totalPages,
                                        },
                                        (
                                            _,
                                            index,
                                        ) => {
                                            const page =
                                                index +
                                                1;

                                            return (
                                                <Pagination.Item
                                                    key={
                                                        page
                                                    }
                                                    active={
                                                        page ===
                                                        currentPage
                                                    }
                                                    onClick={() =>
                                                        setCurrentPage(
                                                            page,
                                                        )
                                                    }
                                                >
                                                    {
                                                        page
                                                    }
                                                </Pagination.Item>
                                            );
                                        },
                                    )}

                                    <Pagination.Next
                                        disabled={
                                            currentPage ===
                                            totalPages
                                        }
                                        onClick={() =>
                                            setCurrentPage(
                                                (
                                                    page,
                                                ) =>
                                                    page +
                                                    1,
                                            )
                                        }
                                    />
                                </Pagination>
                            )}
                        </div>
                    </>
                )}
            </Card.Body>
        </Card>
    );
}