import { Card, Col, Container, Row } from "react-bootstrap";

import { useAuth } from "../auth/hooks/useAuth";

export default function StaffProfile() {

    const { user } = useAuth();

    return (

        <Container fluid>

            <h2 className="mb-4">

                My Profile

            </h2>

            <Card>

                <Card.Body>

                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>First Name</strong>
                        </Col>

                        <Col md={9}>
                            {user?.firstName}
                        </Col>

                    </Row>

                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>Last Name</strong>
                        </Col>

                        <Col md={9}>
                            {user?.lastName}
                        </Col>

                    </Row>

                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>Email</strong>
                        </Col>

                        <Col md={9}>
                            {user?.email}
                        </Col>

                    </Row>

                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>Phone</strong>
                        </Col>

                        <Col md={9}>
                            {user?.phone}
                        </Col>

                    </Row>

                    <Row>

                        <Col md={3}>
                            <strong>Role</strong>
                        </Col>

                        <Col md={9}>
                            {user?.role}
                        </Col>

                    </Row>

                </Card.Body>

            </Card>

        </Container>

    );

}