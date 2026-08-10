import { useEffect, useState } from "react";
import {
    Alert,
    Badge,
    Button,
    Card,
    Spinner,
    Table,
} from "react-bootstrap";
import { useNavigate } from "react-router";

import { useAuth } from "../auth/hooks/useAuth";
import { fetchMyJobs } from "./services/jobApi";

import type { Job } from "./jobTypes";

function formatStatus(status: string) {
    const statusLabels: Record<string, string> = {
        created: "Created",
        assigned: "Assigned",
        in_progress: "In progress",
        waiting_customer: "Waiting for customer",
        resumed: "Resumed",
        closed: "Closed",
    };

    return statusLabels[status] ?? status;
}

function getStatusVariant(status: string) {
    const variants: Record<string, string> = {
        created: "secondary",
        assigned: "primary",
        in_progress: "warning",
        waiting_customer: "info",
        resumed: "success",
        closed: "dark",
    };

    return variants[status] ?? "secondary";
}

function fullName(
    person?: {
        firstName: string;
        lastName: string;
    } | null,
) {
    if (!person) {
        return "--";
    }

    return [
        person.firstName,
        person.lastName,
    ]
        .filter(Boolean)
        .join(" ");
}

function deviceName(job: Job) {
    return [
        job.device.brand,
        job.device.model,
    ]
        .filter(Boolean)
        .join(" ");
}

function formatDate(date: string) {
    return new Date(date).toLocaleString();
}

export default function MyJobs() {
    const { user } = useAuth();
    const navigate = useNavigate();

    const [jobs, setJobs] = useState<Job[]>([]);
    const [jobCount, setJobCount] = useState(0);

    const [loading, setLoading] =
        useState(true);

    const [error, setError] =
        useState("");

    async function loadMyJobs() {
        try {
            setLoading(true);
            setError("");

            const response =
                await fetchMyJobs();

            setJobs(response.jobs ?? []);
            setJobCount(
                response.myJobsCount ??
                    response.jobs?.length ??
                    0,
            );
        } catch (err: any) {
            console.error(
                "My jobs error:",
                err.response?.data,
            );

            setError(
                err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to load your jobs.",
            );

            setJobs([]);
            setJobCount(0);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        if (
            user?.role !== "technician" &&
            user?.role !== "head_technician"
        ) {
            setLoading(false);
            return;
        }

        loadMyJobs();
    }, [user?.role]);

    if (
        user?.role !== "technician" &&
        user?.role !== "head_technician"
    ) {
        return (
            <Alert variant="danger">
                You are not authorized to
                access My Jobs.
            </Alert>
        );
    }

    if (loading) {
        return (
            <div className="text-center p-5">
                <Spinner />
            </div>
        );
    }

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4">
                <div>
                    <h2 className="mb-1">
                        My Jobs
                    </h2>

                    <div className="text-muted">
                        Jobs assigned to you
                    </div>
                </div>

                <Card>
                    <Card.Body className="py-2 px-4 text-center">
                        <div className="text-muted small">
                            My Jobs
                        </div>

                        <div className="fs-4 fw-bold">
                            {jobCount}
                        </div>
                    </Card.Body>
                </Card>
            </div>

            {error && (
                <Alert variant="danger">
                    {error}
                </Alert>
            )}

            {!error && jobs.length === 0 && (
                <Card>
                    <Card.Body className="text-center p-5">
                        <h5>
                            No jobs assigned
                        </h5>

                        <div className="text-muted">
                            You currently have no
                            jobs assigned to you.
                        </div>
                    </Card.Body>
                </Card>
            )}

            {jobs.length > 0 && (
                <Card>
                    <Card.Body>
                        <Card.Title className="mb-4">
                            Assigned Jobs
                        </Card.Title>

                        <Table
                            bordered
                            hover
                            responsive
                            className="align-middle"
                        >
                            <thead>
                                <tr>
                                    <th>
                                        Job Number
                                    </th>

                                    <th>
                                        Customer
                                    </th>

                                    <th>
                                        Device
                                    </th>

                                    <th>
                                        Problem
                                    </th>

                                    <th>
                                        Status
                                    </th>

                                    <th>
                                        Created
                                    </th>

                                    <th>
                                        Action
                                    </th>
                                </tr>
                            </thead>

                            <tbody>
                                {jobs.map((job) => (
                                    <tr
                                        key={job.id}
                                    >
                                        <td>
                                            <span className="fw-semibold">
                                                {
                                                    job.jobNumber
                                                }
                                            </span>
                                        </td>

                                        <td>
                                            <div className="fw-semibold">
                                                {fullName(
                                                    job.customer,
                                                )}
                                            </div>

                                            <div className="small text-muted">
                                                {
                                                    job
                                                        .customer
                                                        .phone
                                                }
                                            </div>
                                        </td>

                                        <td>
                                            <div>
                                                {
                                                    job
                                                        .device
                                                        .type
                                                }
                                            </div>

                                            <div className="small text-muted">
                                                {deviceName(
                                                    job,
                                                ) || "--"}
                                            </div>
                                        </td>

                                        <td>
                                            <div
                                                style={{
                                                    maxWidth:
                                                        "280px",
                                                }}
                                            >
                                                {
                                                    job.problemDescription
                                                }
                                            </div>
                                        </td>

                                        <td>
                                            <Badge
                                                bg={getStatusVariant(
                                                    job.status,
                                                )}
                                            >
                                                {formatStatus(
                                                    job.status,
                                                )}
                                            </Badge>
                                        </td>

                                        <td>
                                            <span className="small">
                                                {formatDate(
                                                    job.createdAt,
                                                )}
                                            </span>
                                        </td>

                                        <td>
                                            <Button
                                                size="sm"
                                                variant="outline-primary"
                                                onClick={() =>
                                                    navigate(`../jobs/${job.jobNumber}`)
                                                }
                                            >
                                                View Job
                                            </Button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </Table>    
                    </Card.Body>
                </Card>
            )}
        </div>
    );
}