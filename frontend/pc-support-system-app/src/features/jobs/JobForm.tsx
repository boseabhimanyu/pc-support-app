
import { useState } from "react";
import { Alert, Button, Form } from "react-bootstrap";

import { api } from "../../app/api";

type JobFormProps = {
    customerId: string;
    deviceId: string;

    customerName: string;

    deviceType: string;
    deviceBrand: string;
    deviceModel: string;
    deviceSerialNumber: string;

    onCreated: (jobNumber: string) => void;
};

type CreateJobResponse = {
    id: string;
    jobNumber: string;
    status: string;

    customer: {
        id: string;
        firstName: string;
        lastName: string;
        phone: string;
    };

    device: {
        id: string;
        type: string;
        brand: string;
        model: string;
        serialNumber: string;
    };

    problemDescription: string;
    notes: string | null;
    createdAt: string;

    createdBy: {
        id: string;
        firstName: string;
        lastName: string;
        role: string;
    };
};

export default function JobForm({
    customerId,
    deviceId,
    customerName,
    deviceType,
    deviceBrand,
    deviceModel,
    deviceSerialNumber,
    onCreated,
}: JobFormProps) {
    const [problemDescription, setProblemDescription] =
        useState("");

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    async function handleSubmit(
        event: React.FormEvent<HTMLFormElement>,
    ) {
        event.preventDefault();

        setError("");
        setSuccess("");

        const description = problemDescription.trim();

        if (!description) {
            setError("Please enter the problem description.");
            return;
        }

        try {
            setLoading(true);

           const response = await api.post<CreateJobResponse>("/jobs", {
            customerId,
            deviceId,
            problemDescription: description,
            });

            const createdJob = response.data;

            if (!createdJob.jobNumber) {
                throw new Error(
                    "Job was created, but the server did not return a job number.",
                );
            }

            setSuccess(
                `Job No. ${createdJob.jobNumber} was created successfully.`,
            );



            /*
             * Pass the ACTUAL job number returned by POST /jobs
             * to the parent page.
             *
             * The parent decides where to navigate based on
             * the user's role.
             */
            onCreated(createdJob.jobNumber);
        } catch (err) {
            setError(
                err instanceof Error
                    ? err.message
                    : "Failed to create job.",
            );
        } finally {
            setLoading(false);
        }
    }

    return (
        <Form onSubmit={handleSubmit}>
            <h3 className="mb-4">Create Job</h3>

            {success && (
                <Alert variant="success">
                    {success}
                </Alert>
            )}

            {error && (
                <Alert variant="danger">
                    {error}
                </Alert>
            )}

            <Form.Group className="mb-3">
                <Form.Label>Customer</Form.Label>

                <Form.Control
                    type="text"
                    value={customerName}
                    readOnly
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>Device</Form.Label>

                <Form.Control
                    type="text"
                    value={`${deviceType} ${deviceBrand} ${deviceModel}`}
                    readOnly
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>Serial Number</Form.Label>

                <Form.Control
                    type="text"
                    value={deviceSerialNumber}
                    readOnly
                />
            </Form.Group>

            <Form.Group className="mb-4">
                <Form.Label>
                    Problem Description
                </Form.Label>

                <Form.Control
                    as="textarea"
                    rows={5}
                    value={problemDescription}
                    onChange={(event) =>
                        setProblemDescription(
                            event.target.value,
                        )
                    }
                    placeholder="Describe the customer's problem"
                    disabled={loading}
                />
            </Form.Group>

            <Button
                type="submit"
                variant="primary"
                disabled={loading}
            >
                {loading ? "Creating Job..." : "Create Job"}
            </Button>
        </Form>
    );
}

