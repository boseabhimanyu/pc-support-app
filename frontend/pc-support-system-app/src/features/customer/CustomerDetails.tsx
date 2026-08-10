import { useEffect, useState } from "react";
import {
    Alert,
    Card,
    Col,
    Row,
    Spinner,
    Form
} from "react-bootstrap";
import { useParams, useNavigate  } from "react-router";
import { Button, Table } from "react-bootstrap";
import { useAuth } from "../auth/hooks/useAuth";
import { canEditCustomer, canResetCustomerPassword } from "../../shared/utils/permissions";
import {api} from "../../app/api";
import {
    fetchCustomerDevices,
} from "../devices/deviceApi";

import type {
    Device,
} from "../devices/types";



import {
    resetCustomerPassword,
} from "./services/customerApi";


import { getJobsByCustomer } from "../jobs/services/jobApi";
import type { JobCustomerProfile } from "../jobs/jobTypes";

type Customer = {
    id: string;
    firstName: string;
    lastName: string;
    phone: string;
    email: string;
};

export default function CustomerDetails() {

    const [jobs, setJobs] = useState<JobCustomerProfile[]>([]);
    const [jobsLoading, setJobsLoading] = useState(true);
    const [jobsError, setJobsError] = useState("");

    const [showResetBox, setShowResetBox] =
    useState(false);

    const [newPassword, setNewPassword] =
        useState("");

    const [showPassword, setShowPassword] =
        useState(false);

    const [passwordMessage, setPasswordMessage] =
        useState("");

    const { user } = useAuth();

    const allowEdit =
    canEditCustomer(user);

    const navigate = useNavigate();

    const { customerId } = useParams();

    const [customer, setCustomer] =
        useState<Customer | null>(null);

    const [loading, setLoading] =
        useState(true);

    const [error, setError] =
        useState("");

    const [devices, setDevices] =
    useState<Device[]>([]);

    const [loadingDevices, setLoadingDevices] =
        useState(false);

    const [deviceError, setDeviceError] =
        useState("");

    useEffect(() => {

        async function loadCustomer() {

            if (!customerId) {
                setError("Customer not found.");
                setLoading(false);
                return;
            }

            try {

                setLoading(true);
                setError("");

                const response =
                    await api.get<Customer>(
                        `/customers/${customerId}`,
                    );

                setCustomer(response.data);

            } catch (err) {

                console.error(err);

                setError(
                    "Unable to load customer.",
                );

            } finally {

                setLoading(false);

            }
        }

        loadCustomer();

    }, [customerId]);

    useEffect(() => {

    async function loadDevices() {

        if (!customerId) {
            return;
        }

        try {

            setLoadingDevices(true);
            setDeviceError("");

            const response =
                await fetchCustomerDevices(
                    customerId,
                );

            setDevices(response);

        } catch (err) {

            console.error(
                "Customer devices error:",
                err,
            );

            setDeviceError(
                "Unable to load customer devices.",
            );

        } finally {

            setLoadingDevices(false);

        }

    }

    loadDevices();

}, [customerId]);


   useEffect(() => {
    async function loadCustomerJobs() {
        if (!customerId) {
            setJobsError("Customer not found.");
            setJobsLoading(false);
            return;
        }

        try {
            setJobsLoading(true);
            setJobsError("");

            const response = await getJobsByCustomer(customerId);

            setJobs(response.jobs);
        } catch (error) {
            setJobsError(
                error instanceof Error
                    ? error.message
                    : "Failed to load customer jobs.",
            );
        } finally {
            setJobsLoading(false);
        }
    }

    loadCustomerJobs();
}, [customerId]);

    if (loading) {

        return (
            <div className="text-center p-4">

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

    if (!customer) {
        return null;
    }

    async function handlePasswordReset() {

    if (!customerId || !newPassword.trim()) {
        return;
    }

    try {

        await resetCustomerPassword(
            customerId,
            newPassword,
        );

        setPasswordMessage(
            "Password reset successfully.",
        );

        setNewPassword("");

    } catch (error) {

        console.error(error);

        setPasswordMessage(
            "Password must contain at least 8 characters, one uppercase letter, one number and one special character.",
        );

    }

}

    return (

        <>

<div className="d-flex justify-content-between align-items-center mb-4">

    <h2 className="mb-0">
        Customer Profile
    </h2>

    {allowEdit && (

        <div className="d-flex gap-2">

            <Button
                onClick={() =>
                    navigate("edit")
                }
            >
                Edit
            </Button>


{canResetCustomerPassword(user) && (

    <Button
        variant="warning"
        onClick={() =>
            setShowResetBox(
                !showResetBox,
            )
        }
    >
        Reset Password
    </Button>

)}

        </div>

    )}

    {showResetBox && (

    <Card className="mt-3">

        <Card.Body>

            <Card.Title>
                Reset Customer Password
            </Card.Title>


            <div className="d-flex gap-2">

                <Form.Control
                    type={
                        showPassword
                            ? "text"
                            : "password"
                    }
                    placeholder="Enter new password"
                    value={newPassword}
                    onChange={(e) =>
                        setNewPassword(
                            e.target.value,
                        )
                    }
                />


                <Button
                    variant="secondary"
                    onClick={() =>
                        setShowPassword(
                            !showPassword,
                        )
                    }
                >
                    {showPassword
                        ? "Hide"
                        : "Show"}
                </Button>


                <Button
                    onClick={
                        handlePasswordReset
                    }
                >
                    Save
                </Button>

            </div>


            {passwordMessage && (

                <Alert
                    className="mt-3"
                    variant="info"
                >
                    {passwordMessage}
                </Alert>

            )}

        </Card.Body>

    </Card>

)}

</div>

            <Card>

                <Card.Body>

                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                First Name
                            </strong>
                        </Col>

                        <Col md={9}>
                            {customer.firstName}
                        </Col>

                    </Row>

                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Last Name
                            </strong>
                        </Col>

                        <Col md={9}>
                            {customer.lastName}
                        </Col>

                    </Row>

                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Phone
                            </strong>
                        </Col>

                        <Col md={9}>
                            {customer.phone}
                        </Col>

                    </Row>

                    <Row>

                        <Col md={3}>
                            <strong>
                                Email
                            </strong>
                        </Col>

                        <Col md={9}>
                            {customer.email || "--"}
                        </Col>

                    </Row>

                </Card.Body>

            </Card>
            <Card className="mt-4">

    <Card.Body>

       
<div className="d-flex justify-content-between align-items-center mb-3">

    <Card.Title className="mb-0">
        Devices
    </Card.Title>

    {allowEdit && (
    <Button
        onClick={() =>
            navigate("devices/create", {
                state: {
                    customerName: `${customer.firstName} ${customer.lastName}`,
                },
            })
        }
    >
        Add Device
    </Button>
)}

</div>



        {loadingDevices && (
            <div className="text-center p-3">
                <Spinner size="sm" />
            </div>
        )}

        {deviceError && (
            <Alert variant="danger">
                {deviceError}
            </Alert>
        )}

        {!loadingDevices &&
            !deviceError &&
            devices.length === 0 && (
                <Alert variant="info">
                    No devices registered for this customer.
                </Alert>
            )}

        {!loadingDevices &&
            devices.length > 0 && (

                <div className="table-responsive">

                    <table className="table table-hover">

                       <thead>
                                <tr>
                                    <th>Type</th>
                                    <th>Brand</th>
                                    <th>Model</th>
                                    <th>Serial Number</th>
                                    <th>Condition</th>
                                    <th>Actions</th>
                                </tr>
                        </thead>

                    
                        <tbody>
                            {devices.map((device) => (
                                <tr key={device.id}>
                                    <td>{device.type}</td>
                                    <td>{device.brand || "--"}</td>
                                    <td>{device.model || "--"}</td>
                                    <td>{device.serialNumber || "--"}</td>
                                    <td>{device.condition}</td>

                                    <td>
                                        <div className="d-flex gap-2">
                                            <Button
                                                size="sm"
                                                variant="outline-primary"
                                                onClick={() =>
                                                    navigate(`devices/${device.id}`)
                                                }
                                            >
                                                View
                                            </Button>

                                            <Button
                                                size="sm"
                                                variant="primary"
                                                onClick={() =>
                                                    navigate(`devices/${device.id}/create-job`, {
                                                        state: {
                                                            customerName: `${customer.firstName} ${customer.lastName}`,
                                                            deviceType: device.type,
                                                            deviceBrand: device.brand,
                                                            deviceModel: device.model,
                                                            deviceSerialNumber: device.serialNumber,
                                                        },
                                                    })
                                                }
                                            >
                                                Create Job
                                            </Button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>



                    </table>
                        
                </div>

            )}

    </Card.Body>

</Card>

<div className="mt-5">
    <h4>Jobs</h4>

    {jobsLoading && <p>Loading jobs...</p>}

    {jobsError && (
        <Alert variant="danger">
            {jobsError}
        </Alert>
    )}

    {!jobsLoading && !jobsError && (
        <>
            {/* Ongoing Jobs */}
            <h5 className="mt-4">Ongoing Jobs</h5>

            {jobs.filter((job) => job.status !== "closed").length === 0 ? (
                <p className="text-muted">
                    No ongoing jobs.
                </p>
            ) : (
                <Table bordered hover responsive>
                      
                    <thead>
                        <tr>
                            <th>Job Number</th>
                            <th>Device</th>
                            <th>Problem</th>
                            <th>Status</th>
                            <th>Created</th>
                            <th>Action</th>
                        </tr>
                    </thead>

                    <tbody>
                        {jobs
                            .filter((job) => job.status !== "closed")
                            .map((job) => (
                                <tr key={job.id}>
                                    <td>{job.jobNumber}</td>

                                    <td>
                                        {job.device.type}
                                        {" - "}
                                        {job.device.brand || "--"}
                                        {job.device.model
                                            ? ` ${job.device.model}`
                                            : ""}
                                    </td>

                                    <td>
                                        {job.problemDescription}
                                    </td>

                                    <td>{job.status}</td>

                                    <td>
                                        {new Date(
                                            job.createdAt,
                                        ).toLocaleDateString()}
                                    </td>

                                    <td>
                                        <Button
                                            size="sm"
                                            variant="outline-primary"
                                            onClick={() =>
                                                navigate(
                                                    `/technician/jobs/${job.jobNumber}`,
                                                )
                                            }
                                        >
                                            View
                                        </Button>
                                    </td>
                                </tr>
                            ))}
                    </tbody>
                </Table>
            )}

            {/* Past Jobs */}
            <h5 className="mt-5">Past Jobs</h5>

            {jobs.filter((job) => job.status === "closed").length === 0 ? (
                <p className="text-muted">
                    No past jobs.
                </p>
            ) : (
                <Table bordered hover responsive>

                    <thead>
                        <tr>
                            <th>Job Number</th>
                            <th>Device</th>
                            <th>Problem</th>
                            <th>Assigned To</th>
                            <th>Closed</th>
                            <th>Close Reason</th>
                            <th>Action</th>
                        </tr>
                    </thead>

                    <tbody>
                        {jobs
                            .filter((job) => job.status === "closed")
                            .map((job) => (
                                <tr key={job.id}>
                                    <td>{job.jobNumber}</td>

                                    <td>
                                        {job.device.type}
                                        {" - "}
                                        {job.device.brand || "--"}
                                        {job.device.model
                                            ? ` ${job.device.model}`
                                            : ""}
                                    </td>

                                    <td>
                                        {job.problemDescription}
                                    </td>

                                    <td>
                                        {job.assignedTo
                                            ? `${job.assignedTo.firstName} ${job.assignedTo.lastName}`
                                            : "--"}
                                    </td>

                                    <td>
                                        {job.closedAt
                                            ? new Date(
                                                  job.closedAt,
                                              ).toLocaleDateString()
                                            : "--"}
                                    </td>

                                    <td>
                                        {job.closeReason || "--"}
                                    </td>

                                    <td>
                                        <Button
                                            size="sm"
                                            variant="outline-primary"
                                            onClick={() =>
                                                navigate(
                                                    `/technician/jobs/${job.jobNumber}`,
                                                )
                                            }
                                        >
                                            View
                                        </Button>
                                    </td>
                                </tr>
                            ))}
                    </tbody>
                </Table>
            )}
        </>
    )}
</div>    
        </>

    );
}