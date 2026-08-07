import { useEffect, useState } from "react";
import {
    Alert,
    Card,
    Col,
    Row,
    Spinner,
} from "react-bootstrap";
import { useParams } from "react-router-dom";

import {api} from "../../app/api";

type Customer = {
    id: string;
    firstName: string;
    lastName: string;
    phone: string;
    email: string;
};

export default function CustomerDetails() {

    const { customerId } = useParams();

    const [customer, setCustomer] =
        useState<Customer | null>(null);

    const [loading, setLoading] =
        useState(true);

    const [error, setError] =
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

    return (

        <>

            <h2 className="mb-4">
                Customer Profile
            </h2>

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

        </>

    );
}