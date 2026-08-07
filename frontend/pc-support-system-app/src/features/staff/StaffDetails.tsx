import { useEffect, useState } from "react";
import {
    Alert,
    Button,
    Card,
    Col,
    Row,
    Spinner,
} from "react-bootstrap";
import {
    useNavigate,
    useParams,
} from "react-router-dom";

import {api} from "../../app/api";

type Staff = {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    role: string;
    state: string;
};

export default function StaffDetails() {

    const { staffId } = useParams();
    const navigate = useNavigate();

    const [staff, setStaff] =
        useState<Staff | null>(null);

    const [loading, setLoading] =
        useState(true);

    const [error, setError] =
        useState("");

    useEffect(() => {

        async function loadStaff() {

            if (!staffId) {
                setError("Staff ID is missing.");
                setLoading(false);
                return;
            }

            try {

                setLoading(true);
                setError("");

                const response =
                    await api.get<Staff>(
                        `/staff/${staffId}`,
                    );

                setStaff(response.data);

            } catch (error: any) {

                console.error(
                    "Staff details error:",
                    error.response?.data,
                );

                setError(
                    error.response?.data?.error ??
                    error.response?.data?.message ??
                    "Unable to load staff details.",
                );

            } finally {

                setLoading(false);

            }

        }

        loadStaff();

    }, [staffId]);


    function formatRole(
        value: string,
    ) {

        return value
            .replace(/_/g, " ")
            .replace(/\b\w/g, (letter) =>
                letter.toUpperCase(),
            );

    }


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


    if (!staff) {

        return (
            <Alert variant="warning">
                Staff member not found.
            </Alert>
        );

    }


    return (

        <>

            <div className="d-flex justify-content-between align-items-center mb-4">

                <h2 className="mb-0">
                    Staff Details
                </h2>

                <Button
                    variant="secondary"
                    onClick={() =>
                        navigate(-1)
                    }
                >
                    Back
                </Button>

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
                            {staff.firstName}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Last Name
                            </strong>
                        </Col>

                        <Col md={9}>
                            {staff.lastName}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Email
                            </strong>
                        </Col>

                        <Col md={9}>
                            {staff.email}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Phone
                            </strong>
                        </Col>

                        <Col md={9}>
                            {staff.phone}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Role
                            </strong>
                        </Col>

                        <Col md={9}>
                            {formatRole(
                                staff.role,
                            )}
                        </Col>

                    </Row>


                    <Row>

                        <Col md={3}>
                            <strong>
                                State
                            </strong>
                        </Col>

                        <Col md={9}>
                            {formatRole(
                                staff.state,
                            )}
                        </Col>

                    </Row>

                </Card.Body>

            </Card>

        </>

    );

}