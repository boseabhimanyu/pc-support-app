
import { useEffect, useState } from "react";
import {
    Alert,
    Button,
    Card,
    Col,
    Row,
    Spinner,
    Table,
} from "react-bootstrap";
import { useNavigate } from "react-router-dom";

import { useAuth } from "../auth/hooks/useAuth";

import {
    fetchAssignedJobs,
    fetchInProgressJobs,
    fetchOpenJobs,
    fetchResumedJobs,
    fetchWaitingCustomerJobs,
} from "../jobs/services/jobApi";

import type {
    Job,
    JobQueueResponse,
} from "./jobTypes";

type QueueType =
    | "open"
    | "assigned"
    | "in-progress"
    | "waiting-customer"
    | "resumed";

type Queue = {
    key: QueueType;
    title: string;
};


const QUEUES: Queue[] = [
    {
        key: "open",
        title: "Open",
    },
    {
        key: "assigned",
        title: "Assigned",
    },
    {
        key: "in-progress",
        title: "In Progress",
    },
    {
        key: "waiting-customer",
        title: "Waiting Customer",
    },
    {
        key: "resumed",
        title: "Resumed",
    },
];

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

function isAssignedQueue(queue: QueueType) {
    return (
        queue === "assigned" ||
        queue === "in-progress" ||
        queue === "waiting-customer" ||
        queue === "resumed"
    );
}

export default function JobManagement() {
    const { user } = useAuth();
    const navigate = useNavigate();

    const canManageQueues =
    user?.role === "head_technician" ||
    user?.role === "admin" ||
    user?.role === "super_admin";   


    const [counts, setCounts] = useState<
        Record<QueueType, number>
    >({
        open: 0,
        assigned: 0,
        "in-progress": 0,
        "waiting-customer": 0,
        resumed: 0,
    });

    const [selectedQueue, setSelectedQueue] =
        useState<QueueType | null>(null);

    const [jobs, setJobs] = useState<Job[]>([]);

    const [loadingCounts, setLoadingCounts] =
        useState(false);

    const [loadingJobs, setLoadingJobs] =
        useState(false);

    const [error, setError] = useState("");

    /*
     * Get the five queue responses.
     *
     * The backend uses:
     * - openJobsCount for Open
     * - assignedJobsCount for Assigned
     * - jobsCount for the other queues
     */
    async function loadQueueCounts() {
        try {
            setLoadingCounts(true);
            setError("");

            const [
                openResponse,
                assignedResponse,
                inProgressResponse,
                waitingCustomerResponse,
                resumedResponse,
            ] = await Promise.all([
                fetchOpenJobs(),
                fetchAssignedJobs(),
                fetchInProgressJobs(),
                fetchWaitingCustomerJobs(),
                fetchResumedJobs(),
            ]);

            setCounts({
                open:
                    openResponse.openJobsCount ??
                    openResponse.jobsCount ??
                    openResponse.jobs.length,

                assigned:
                    assignedResponse.assignedJobsCount ??
                    assignedResponse.jobsCount ??
                    assignedResponse.jobs.length,

                "in-progress":
                    inProgressResponse.jobsCount ??
                    inProgressResponse.jobs.length,

                "waiting-customer":
                    waitingCustomerResponse.jobsCount ??
                    waitingCustomerResponse.jobs.length,

                resumed:
                    resumedResponse.jobsCount ??
                    resumedResponse.jobs.length,
            });
        } catch (err: any) {
            console.error(
                "Job queue count error:",
                err.response?.data,
            );

            setError(
                err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to load job queues.",
            );
        } finally {
            setLoadingCounts(false);
        }
    }

useEffect(() => {
    if (!canManageQueues) {
        return;
    }

    loadQueueCounts();
}, [canManageQueues]);

    async function loadQueue(
        queue: QueueType,
    ) {
        try {
            setSelectedQueue(queue);
            setLoadingJobs(true);
            setError("");
            setJobs([]);

            let response: JobQueueResponse;

            switch (queue) {
                case "open":
                    response =
                        await fetchOpenJobs();
                    break;

                case "assigned":
                    response =
                        await fetchAssignedJobs();
                    break;

                case "in-progress":
                    response =
                        await fetchInProgressJobs();
                    break;

                case "waiting-customer":
                    response =
                        await fetchWaitingCustomerJobs();
                    break;

                case "resumed":
                    response =
                        await fetchResumedJobs();
                    break;
            }

            setJobs(response.jobs);
        } catch (err: any) {
            console.error(
                "Job queue error:",
                err.response?.data,
            );

            setError(
                err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to load jobs.",
            );
        } finally {
            setLoadingJobs(false);
        }
    }

    function customerName(job: Job) {
        return [
            job.customer.firstName,
            job.customer.lastName,
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


    return (
        <div>
            <h2 className="mb-4">
                Job Management
            </h2>

            {error && (
                <Alert variant="danger">
                    {error}
                </Alert>
            )}
            {canManageQueues && (
            <Row className="g-3 mb-4">
                {QUEUES.map((queue) => (
                    <Col
                        key={queue.key}
                        xs={12}
                        sm={6}
                        lg={4}
                        xl={2}
                    >
                        <Card
                            className={`h-100 ${
                                selectedQueue ===
                                queue.key
                                    ? "border-primary"
                                    : ""
                            }`}
                        >
                            <Card.Body className="d-flex flex-column">
                                <Card.Title>
                                    {queue.title}
                                </Card.Title>

                                <div
                                    className="display-5 fw-bold mb-3"
                                >
                                    {loadingCounts ? (
                                        <Spinner
                                            animation="border"
                                            size="sm"
                                        />
                                    ) : (
                                        counts[
                                            queue.key
                                        ]
                                    )}
                                </div>

                                <Button
                                    variant={
                                        selectedQueue ===
                                        queue.key
                                            ? "primary"
                                            : "outline-primary"
                                    }
                                    className="mt-auto"
                                    onClick={() =>
                                        loadQueue(
                                            queue.key,
                                        )
                                    }
                                    disabled={
                                        loadingJobs
                                    }
                                >
                                    View Jobs
                                </Button>
                            </Card.Body>
                        </Card>
                    </Col>
                ))}
            </Row>
                )}
            {canManageQueues && selectedQueue && (
                <Card>
                    <Card.Body>
                        <div className="d-flex justify-content-between align-items-center mb-3">
    <Card.Title className="mb-0">
        {
            QUEUES.find(
                (queue) =>
                    queue.key === selectedQueue,
            )?.title
        }{" "}
        Jobs
    </Card.Title>

    <div className="d-flex gap-2">
        <Button
            variant="outline-secondary"
            size="sm"
            onClick={() =>
                loadQueue(selectedQueue)
            }
            disabled={loadingJobs}
        >
            Refresh
        </Button>

        <Button
            variant="outline-secondary"
            size="sm"
            onClick={() => {
                setSelectedQueue(null);
                setJobs([]);
            }}
        >
            Close
        </Button>
    </div>
</div>

                        {loadingJobs ? (
                            <div className="text-center p-4">
                                <Spinner />
                            </div>
                        ) : jobs.length ===
                          0 ? (
                            <Alert variant="info">
                                No jobs in this
                                queue.
                            </Alert>
                        ) : (
                            <div className="table-responsive">
                                <Table
                                    bordered
                                    hover
                                    responsive
                                    className="align-middle"
                                    style={{ tableLayout: "fixed", width: "100%" }}
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
                                            {isAssignedQueue(selectedQueue) && (
                                            <th>
                                                Assigned To
                                            </th>
                                                )}
                                            <th>
                                                Status
                                            </th>
                                            <th>
                                                Created
                                            </th>
                                            <th>
                                            </th>
                                        </tr>
                                    </thead>

                                    <tbody>
                                        {jobs.map(
                                            (
                                                job,
                                            ) => (
                                                <tr
                                                    key={
                                                        job.id
                                                    }
                                                >
                                                    <td>
                                                        {
                                                            job.jobNumber
                                                        }
                                                    </td>

                                                        <td>
                                            {customerName(job)}
                                                         </td>



                                                    <td>
                                                        {deviceName(job) || job.device.type}
                                                    </td>

                                                    <td>
                                                        {
                                                            job.problemDescription
                                                        }
                                                    </td>
                                                            {isAssignedQueue(selectedQueue) && (
                                                    <td>
                                                        {job.assignedTo
                                                            ? `${job.assignedTo.firstName} ${job.assignedTo.lastName}`
                                                            : "--"}
                                                    </td>
                                                            )}
                                                    <td>
                                                        {
                                                            formatStatus(job.status)
                                                        }
                                                    </td>

                                                    <td>
                                                        {new Date(
                                                            job.createdAt,
                                                        ).toLocaleString()}
                                                    </td>

                                                    <td>
                                                        <Button
                                                            size="sm"
                                                            variant="outline-primary"
                                                            onClick={() =>
                                                                navigate(
                                                                    `${job.jobNumber}`,
                                                                )
                                                            }
                                                        >
                                                            View
                                                        </Button>
                                                    </td>
                                                </tr>
                                            ),
                                        )}
                                    </tbody>
                                </Table>
                            </div>
                        )}
                    </Card.Body>
                </Card>
            )}
        </div>
    );
}

