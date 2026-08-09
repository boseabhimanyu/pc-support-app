import { api } from "../../app/api";
import { useEffect, useState } from "react";
import {
    Alert,
    Badge,
    Button,
    Card,
    Col,
    Modal,
    Row,
    Spinner,
} from "react-bootstrap";
import {
    useNavigate,
    useParams,
} from "react-router-dom";

import { useAuth } from "../auth/hooks/useAuth";

import { fetchJobByNumber, addJobNote, updateJobStatus } from "./services/jobApi";

import type { Job } from "./jobTypes";

export type AssignableStaff = {
    id: string;
    firstName: string;
    lastName: string;
    role: string;
    state: string;
};

export async function searchAssignableStaff(): Promise<AssignableStaff[]> {
    const response = await api.get<AssignableStaff[]>(
        "/staff/search?q=tech",
    );

    return response.data;
}

export async function assignJob(
    jobId: string,
    staffId: string,
) {
    const response = await api.patch(
        `/jobs/${jobId}/assign`,
        {
            staffId,
        },
    );

    return response.data;
}

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

function formatRole(role: string) {
    const roleLabels: Record<string, string> = {
        receptionist: "Receptionist",
        technician: "Technician",
        head_technician: "Head Technician",
        admin: "Admin",
        super_admin: "Super Admin",
    };

    return roleLabels[role] ?? role;
}

function formatDate(date: string) {
    return new Date(date).toLocaleString();
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

export default function StaffJobDetails() {
    const { jobNumber } = useParams();

    console.log(
        "StaffJobDetails jobNumber:",
        jobNumber,
    );

    const navigate = useNavigate();

    const { user } = useAuth();

    const [job, setJob] =
        useState<Job | null>(null);

    const [loading, setLoading] =
        useState(true);

    const [savingNote, setSavingNote] =
        useState(false);

    const [note, setNote] =
        useState("");

    const [error, setError] =
        useState("");

    const [noteError, setNoteError] =
        useState("");

    const [selectedStatus, setSelectedStatus] =
    useState("");

    const [changingStatus, setChangingStatus] =
        useState(false);

    const [statusError, setStatusError] =
        useState("");

    const [statusSuccess, setStatusSuccess] =
        useState("");
    /*
     * Note permissions:
     *
     * Receptionists can add notes to any job.
     *
     * Technicians can add notes only when
     * they are assigned to this job.
     *
     * Head technicians can add notes only
     * when they are assigned to this job.
     *
     * Admins and super admins can view notes
     * but cannot add notes.
     */
    const canAddNote =
    job?.status !== "closed" &&
    (
        user?.role === "receptionist" ||
        (
            user?.role === "technician" &&
            job?.assignedTo?.id === user.id
        ) ||
        (
            user?.role === "head_technician" &&
            job?.assignedTo?.id === user.id
        )
    );

    const canChangeJobStatus =
    job?.status !== "closed" &&
    (
        user?.role === "technician" ||
        user?.role === "head_technician"
    ) &&
    job?.assignedTo?.id === user.id;

    const canAssignJob =
    job?.status !== "closed" &&
    (
        user?.role === "admin" ||
        user?.role === "head_technician"
    );

    async function loadJob() {
        if (!jobNumber) {
            setError("Job number is missing.");
            setLoading(false);
            return;
        }

        try {
            setLoading(true);
            setError("");

            const response =
                await fetchJobByNumber(
                    jobNumber,
                );

            setJob(response);
        } catch (err: any) {
            console.error(
                "Job details error:",
                err.response?.data,
            );

            setError(
                err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to load job.",
            );
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        loadJob();
    }, [jobNumber]);

    useEffect(() => {
    if (!canAssignJob) {
        return;
    }

    async function loadAssignableStaff() {
        try {
            setLoadingStaff(true);
            setAssignmentError("");

            const results =
                await searchAssignableStaff();

            setStaff(results);
        } catch (err: any) {
            console.error(
                "Assignable staff error:",
                err.response?.data,
            );

            setAssignmentError(
                err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to load technicians.",
            );
        } finally {
            setLoadingStaff(false);
        }
    }

    loadAssignableStaff();
}, [canAssignJob]);

    async function handleAddNote() {
        if (!job) {
            return;
        }

        const trimmedNote =
            note.trim();

        if (!trimmedNote) {
            setNoteError(
                "Please enter a note.",
            );
            return;
        }

        try {
            setSavingNote(true);
            setNoteError("");

            await addJobNote(job.id, {
                note: trimmedNote,
            });

            setNote("");

            await loadJob();
        } catch (err: any) {
            console.error(
                "Add job note error:",
                err.response?.data,
            );

            setNoteError(
                err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to add note.",
            );
        } finally {
            setSavingNote(false);
        }
    }

    async function handleStatusChange() {
    if (!job || !selectedStatus) {
        return;
    }

    try {
        setChangingStatus(true);
        setStatusError("");
        setStatusSuccess("");

        const updatedJob =
            await updateJobStatus(
                job.id,
                selectedStatus,
            );

        setJob(updatedJob);

        setStatusSuccess(
            "Job status updated successfully.",
        );

        setSelectedStatus("");

    } catch (err: any) {
        console.error(
            "Update job status error:",
            err.response?.data,
        );

        setStatusError(
            err.response?.data?.error ??
                err.response?.data?.message ??
                "Unable to update job status.",
        );
    } finally {
        setChangingStatus(false);
    }
}
    const [staff, setStaff] =
    useState<AssignableStaff[]>([]);

    const [selectedStaffId, setSelectedStaffId] =
        useState("");

    const [loadingStaff, setLoadingStaff] =
        useState(false);

    const [assigning, setAssigning] =
        useState(false);

    const [assignmentError, setAssignmentError] =
        useState("");

    const [assignmentSuccess, setAssignmentSuccess] =
        useState("");
    const [showAssignConfirmation, setShowAssignConfirmation] =
        useState(false);

    async function handleAssignJob() {
    if (!job || !selectedStaffId) {
        return;
    }

    try {
        setAssigning(true);
        setAssignmentError("");
        setAssignmentSuccess("");

        await assignJob(
            job.id,
            selectedStaffId,
        );

        setAssignmentSuccess(
            "Job assigned successfully.",
        );

        setSelectedStaffId("");

        // Reload job so assignedTo and status
        // immediately reflect the backend.
        await loadJob();
    } catch (err: any) {
        console.error(
            "Assign job error:",
            err.response?.data,
        );

        setAssignmentError(
            err.response?.data?.error ??
                err.response?.data?.message ??
                "Unable to assign job.",
        );
    } finally {
        setAssigning(false);
    }
}

    if (loading) {
        return (
            <div className="text-center p-5">
                <Spinner />
            </div>
        );
    }

    if (error) {
        return (
            <Alert variant="danger">
                {error}
            </Alert>
        );
    }

    if (!job) {
        return (
            <Alert variant="warning">
                Job not found.
            </Alert>
        );
    }

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4">
                <div>
                    <h2 className="mb-1">
                        Job Details
                    </h2>

                    <div className="text-muted">
                        {job.jobNumber}
                    </div>
                </div>

                <Button
                    variant="outline-secondary"
                    onClick={() =>
                        navigate(-1)
                    }
                >
                    Back
                </Button>
            </div>

            <Row className="g-4">
                <Col xs={12} lg={8}>
                    {/* Job information */}
                    <Card className="mb-4">
                        <Card.Body>
                            <Card.Title className="mb-4">
                                Job Information
                            </Card.Title>

                            <Row className="g-3">
                                <Col xs={12} md={6}>
                                    <div className="text-muted small">
                                        Job Number
                                    </div>

                                    <div className="fw-semibold">
                                        {job.jobNumber}
                                    </div>
                                </Col>

                                <Col xs={12} md={6}>
                                    <div className="text-muted small">
                                        Status
                                    </div>

                                    <Badge bg="primary">
                                        {formatStatus(
                                            job.status,
                                        )}
                                    </Badge>
                                </Col>

                                <Col xs={12}>
                                    <div className="text-muted small">
                                        Problem Description
                                    </div>

                                    <div>
                                        {
                                            job.problemDescription
                                        }
                                    </div>
                                </Col>

                                <Col xs={12} md={6}>
                                    <div className="text-muted small">
                                        Created
                                    </div>

                                    <div>
                                        {formatDate(
                                            job.createdAt,
                                        )}
                                    </div>
                                </Col>
                            </Row>
                        </Card.Body>
                    </Card>

                    {/* Customer */}
                    <Card className="mb-4">
                        <Card.Body>
                            <Card.Title className="mb-4">
                                Customer
                            </Card.Title>

                            <Row className="g-3">
                                <Col xs={12} md={6}>
                                    <div className="text-muted small">
                                        Name
                                    </div>

                                    <div className="fw-semibold">
                                        {fullName(
                                            job.customer,
                                        )}
                                    </div>
                                </Col>

                                <Col xs={12} md={6}>
                                    <div className="text-muted small">
                                        Phone
                                    </div>

                                    <div>
                                        {
                                            job.customer
                                                .phone
                                        }
                                    </div>
                                </Col>
                            </Row>
                        </Card.Body>
                    </Card>

                    {/* Device */}
                    <Card className="mb-4">
                        <Card.Body>
                            <Card.Title className="mb-4">
                                Device
                            </Card.Title>

                            <Row className="g-3">
                                <Col xs={12} md={4}>
                                    <div className="text-muted small">
                                        Type
                                    </div>

                                    <div className="fw-semibold">
                                        {
                                            job.device
                                                .type
                                        }
                                    </div>
                                </Col>

                                <Col xs={12} md={4}>
                                    <div className="text-muted small">
                                        Brand
                                    </div>

                                    <div>
                                        {
                                            job.device
                                                .brand ||
                                            "--"
                                        }
                                    </div>
                                </Col>

                                <Col xs={12} md={4}>
                                    <div className="text-muted small">
                                        Model
                                    </div>

                                    <div>
                                        {
                                            job.device
                                                .model ||
                                            "--"
                                        }
                                    </div>
                                </Col>

                                <Col xs={12}>
                                    <div className="text-muted small">
                                        Serial Number
                                    </div>

                                    <div>
                                        {
                                            job.device
                                                .serialNumber ||
                                            "--"
                                        }
                                    </div>
                                </Col>
                            </Row>
                        </Card.Body>
                    </Card>

                    {/* Notes */}
                    <Card>
                        <Card.Body>
                            <Card.Title className="mb-4">
                                Notes
                            </Card.Title>

                            {/* Existing notes are visible to everyone */}
                            {job.notes &&
                            job.notes.length > 0 ? (
                                <div className="d-flex flex-column gap-3 mb-4">
                                    {job.notes.map(
                                        (jobNote) => (
                                            <Card
                                                key={
                                                    jobNote.id
                                                }
                                                className="border"
                                            >
                                                <Card.Body>
                                                    <div className="d-flex justify-content-between align-items-start gap-3">
                                                        <div>
                                                            <div className="fw-semibold">
                                                                {fullName(
                                                                    jobNote.author,
                                                                )}
                                                            </div>

                                                            <div className="text-muted small">
                                                                {formatRole(
                                                                    jobNote
                                                                        .author
                                                                        .role,
                                                                )}
                                                            </div>
                                                        </div>

                                                        <div className="text-muted small text-end">
                                                            {formatDate(
                                                                jobNote.createdAt,
                                                            )}
                                                        </div>
                                                    </div>

                                                    <div className="mt-3">
                                                        {
                                                            jobNote.note
                                                        }
                                                    </div>
                                                </Card.Body>
                                            </Card>
                                        ),
                                    )}
                                </div>
                            ) : (
                                <Alert variant="light">
                                    No notes have
                                    been added to
                                    this job.
                                </Alert>
                            )}

                            {/* Only permitted users can add notes */}
                            {canAddNote && (
                                <>
                                    <hr />

                                    <h6 className="mb-3">
                                        Add Note
                                    </h6>

                                    {noteError && (
                                        <Alert variant="danger">
                                            {
                                                noteError
                                            }
                                        </Alert>
                                    )}

                                    <textarea
                                        className="form-control mb-3"
                                        rows={4}
                                        value={note}
                                        disabled={
                                            savingNote
                                        }
                                        placeholder="Enter a note about this job..."
                                        onChange={(
                                            event,
                                        ) =>
                                            setNote(
                                                event
                                                    .target
                                                    .value,
                                            )
                                        }
                                    />

                                    <Button
                                        variant="primary"
                                        disabled={
                                            savingNote ||
                                            !note.trim()
                                        }
                                        onClick={
                                            handleAddNote
                                        }
                                    >
                                        {savingNote ? (
                                            <>
                                                <Spinner
                                                    animation="border"
                                                    size="sm"
                                                    className="me-2"
                                                />
                                                Adding...
                                            </>
                                        ) : (
                                            "Add Note"
                                        )}
                                    </Button>
                                </>
                            )}
                        </Card.Body>
                    </Card>
                </Col>

                <Col xs={12} lg={4}>
                    {/* Assignment */}
                    <Card className="mb-4">
                        <Card.Body>
                            <Card.Title className="mb-3">
                                Assigned To
                            </Card.Title>

                            {job.assignedTo ? (
                                <>
                                    <div className="fw-semibold">
                                        {fullName(
                                            job.assignedTo,
                                        )}
                                    </div>

                                    <div className="text-muted">
                                        {formatRole(
                                            job.assignedTo
                                                .role,
                                        )}
                                    </div>
                                </>
                            ) : (
                                <div className="text-muted">
                                    Not assigned
                                </div>
                            )}
                        </Card.Body>
                    </Card>
                    {/* Change Job Status */}
                    {canChangeJobStatus && (
                    <Card className="mb-4">
                        <Card.Body>
                            <Card.Title className="mb-3">
                                Change Job Status
                            </Card.Title>

                            {statusSuccess && (
                                <Alert variant="success">
                                    {statusSuccess}
                                </Alert>
                            )}

                            {statusError && (
                                <Alert variant="danger">
                                    {statusError}
                                </Alert>
                            )}

                            {job.status === "closed" ? (
                                <div className="text-muted">
                                    This job is closed.
                                </div>
                            ) : (
                                <>
                                    <div className="text-muted small mb-2">
                                        Current Status
                                    </div>

                                    <div className="fw-semibold mb-3">
                                        {formatStatus(job.status)}
                                    </div>

                                    <label
                                        htmlFor="jobStatus"
                                        className="form-label"
                                    >
                                        New Status
                                    </label>

                                    <select
                                        id="jobStatus"
                                        className="form-select mb-3"
                                        value={selectedStatus}
                                        disabled={changingStatus}
                                        onChange={(event) => {
                                            setSelectedStatus(
                                                event.target.value,
                                            );

                                            setStatusError("");
                                            setStatusSuccess("");
                                        }}
                                    >
                                        <option value="">
                                            Select status
                                        </option>

                                        <option value="in_progress">
                                            In progress
                                        </option>

                                        <option value="waiting_customer">
                                            Waiting for customer
                                        </option>

                                        <option value="resumed">
                                            Resumed
                                        </option>
                                    </select>

                                    <Button
                                        variant="primary"
                                        disabled={
                                            !selectedStatus ||
                                            changingStatus
                                        }
                                        onClick={handleStatusChange}
                                    >
                                        {changingStatus
                                            ? "Updating..."
                                            : "Update Status"}
                                    </Button>
                                </>
                            )}
                        </Card.Body>
                    </Card> )}
                    {/* Created By */}
                    <Card className="mb-4">
                        <Card.Body>
                            <Card.Title className="mb-3">
                                Created By
                            </Card.Title>

                            <div className="fw-semibold">
                                {fullName(
                                    job.createdBy,
                                )}
                            </div>

                            <div className="text-muted">
                                {formatRole(
                                    job.createdBy.role,
                                )}
                            </div>
                        </Card.Body>
                    </Card>
                                
                    <Card className="mb-4">
    <Card.Body>
        <Card.Title className="mb-3">
            Assign Job
        </Card.Title>

        {/* Current assignment */}
        {job.assignedTo ? (
            <div className="mb-3">
                <div className="text-muted small">
                    Currently Assigned To
                </div>

                <div className="fw-semibold">
                    {fullName(
                        job.assignedTo,
                    )}
                </div>

                <div className="text-muted">
                    {formatRole(
                        job.assignedTo.role,
                    )}
                </div>
            </div>
        ) : (
            <div className="text-muted mb-3">
                Not assigned
            </div>
        )}

        {/* Assignment controls */}
        {canAssignJob && (
            <>
                <hr />

                {assignmentSuccess && (
                    <Alert variant="success">
                        {assignmentSuccess}
                    </Alert>
                )}

                {assignmentError && (
                    <Alert variant="danger">
                        {assignmentError}
                    </Alert>
                )}

                <label
                    htmlFor="assignStaff"
                    className="form-label"
                >
                    Select Technician
                </label>

                <select
                    id="assignStaff"
                    className="form-select mb-3"
                    value={selectedStaffId}
                    disabled={
                        loadingStaff ||
                        assigning
                    }
                    onChange={(event) =>
                        setSelectedStaffId(
                            event.target.value,
                        )
                    }
                >
                    <option value="">
                        {loadingStaff
                            ? "Loading staff..."
                            : "Select staff"}
                    </option>

                    {staff.map((person) => (
                        <option
                            key={person.id}
                            value={person.id}
                        >
                            {person.firstName}{" "}
                            {person.lastName}
                            {" — "}
                            {formatRole(
                                person.role,
                            )}
                        </option>
                    ))}
                </select>

               <Button
                variant="primary"
                disabled={
                    !selectedStaffId ||
                    assigning ||
                    loadingStaff
                }
                onClick={() =>
                    setShowAssignConfirmation(true)
                }
            >
                Assign
                </Button>
            </>
        )}
    </Card.Body>
</Card>
                </Col>
            </Row>

           <Modal
    show={showAssignConfirmation}
    onHide={() =>
        setShowAssignConfirmation(false)
    }
    centered
>
    <Modal.Header closeButton>
        <Modal.Title>
            {job.assignedTo
                ? "Reassign Job"
                : "Assign Job"}
        </Modal.Title>
    </Modal.Header>

    <Modal.Body>

        {[
            "assigned",
            "in_progress",
            "waiting_customer",
            "resumed",
        ].includes(job.status) && (
            <Alert variant="warning">
                This job is currently{" "}
                <strong>
                    {formatStatus(job.status)}
                </strong>
                .
                Reassigning it will change
                the staff member responsible
                for the job.
            </Alert>
        )}

        {job.assignedTo ? (
            <>
                <p>
                    This job is currently assigned
                    to{" "}
                    <strong>
                        {fullName(
                            job.assignedTo,
                        )}
                    </strong>
                    .
                </p>

                <p className="mb-0">
                    Are you sure you want to
                    reassign this job to{" "}
                    <strong>
                        {
                            staff.find(
                                (person) =>
                                    person.id ===
                                    selectedStaffId,
                            )?.firstName
                        }{" "}
                        {
                            staff.find(
                                (person) =>
                                    person.id ===
                                    selectedStaffId,
                            )?.lastName
                        }
                    </strong>
                    ?
                </p>
            </>
        ) : (
            <p className="mb-0">
                Are you sure you want to
                assign this job to{" "}
                <strong>
                    {
                        staff.find(
                            (person) =>
                                person.id ===
                                selectedStaffId,
                        )?.firstName
                    }{" "}
                    {
                        staff.find(
                            (person) =>
                                person.id ===
                                selectedStaffId,
                        )?.lastName
                    }
                </strong>
                ?
            </p>
        )}

    </Modal.Body>

    <Modal.Footer>

        <Button
            variant="secondary"
            onClick={() =>
                setShowAssignConfirmation(false)
            }
            disabled={assigning}
        >
            Cancel
        </Button>

        <Button
            variant="primary"
            onClick={async () => {
                setShowAssignConfirmation(false);
                await handleAssignJob();
            }}
            disabled={assigning}
        >
            {assigning
                ? "Assigning..."
                : job.assignedTo
                    ? "Yes, Reassign"
                    : "Yes, Assign"}
        </Button>

    </Modal.Footer>
</Modal>
        </div>
    );
}

