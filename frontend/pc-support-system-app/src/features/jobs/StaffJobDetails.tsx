
import { useEffect, useState } from "react";
import {
    Alert,
    Badge,
    Button,
    Card,
    Col,
    Row,
    Spinner,
} from "react-bootstrap";
import { useNavigate, useParams } from "react-router-dom";

import {
    fetchJobByNumber,
    addJobNote,
} from "./services/jobApi";

import type {
    Job,
} from "./jobTypes";

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

function fullName(person?: {
    firstName: string;
    lastName: string;
} | null) {
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

    const navigate = useNavigate();

    const [job, setJob] = useState<Job | null>(
        null,
    );

    const [loading, setLoading] =
        useState(true);

    const [savingNote, setSavingNote] =
        useState(false);

    const [note, setNote] = useState("");

    const [error, setError] = useState("");

    const [noteError, setNoteError] =
        useState("");

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
            await fetchJobByNumber(jobNumber);

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

    async function handleAddNote() {
    if (!job) {
        return;
    }

    const trimmedNote = note.trim();

    if (!trimmedNote) {
        setNoteError("Please enter a note.");
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

                            {job.notes && job.notes.length > 0 ? (
    <div className="d-flex flex-column gap-3 mb-4">
        {job.notes.map((jobNote) => (
            <Card
                key={jobNote.id}
                className="border"
            >
                <Card.Body>
                    <div className="d-flex justify-content-between align-items-start gap-3">
                        <div>
                            <div className="fw-semibold">
                                {fullName(jobNote.author)}
                            </div>

                            <div className="text-muted small">
                                {formatRole(
                                    jobNote.author.role,
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
                        {jobNote.note}
                    </div>
                </Card.Body>
            </Card>
        ))}
    </div>
) : (
    <Alert variant="light">
        No notes have been added to this job.
    </Alert>
)}

                            <hr />

                            <h6 className="mb-3">
                                Add Note
                            </h6>

                            {noteError && (
                                <Alert variant="danger">
                                    {noteError}
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
                                onChange={(event) =>
                                    setNote(
                                        event.target
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

                    {/* Created By */}
                    <Card>
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
                </Col>
            </Row>
        </div>
    );
}

