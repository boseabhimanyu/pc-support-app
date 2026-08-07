import { Card, Col, Row } from "react-bootstrap";

import { useAuth } from "../auth/hooks/useAuth";

export default function Dashboard() {

    const { user } = useAuth();

    return (

        <>

            <h2 className="mb-4">

                Dashboard

            </h2>

            <p className="text-muted mb-4">

                Welcome back, {user?.firstName}.

            </p>

            <Row className="g-3">

                <Col md={3}>

                    <Card>

                        <Card.Body>

                            <Card.Title>

                                Customers

                            </Card.Title>

                            <Card.Text>

                                Customer management

                            </Card.Text>

                        </Card.Body>

                    </Card>

                </Col>

                <Col md={3}>

                    <Card>

                        <Card.Body>

                            <Card.Title>

                                Jobs

                            </Card.Title>

                            <Card.Text>

                                Job management

                            </Card.Text>

                        </Card.Body>

                    </Card>

                </Col>

                {(user?.role === "admin" ||
                  user?.role === "super_admin" ||
                  user?.role === "head_technician") && (

                    <Col md={3}>

                        <Card>

                            <Card.Body>

                                <Card.Title>

                                    Staff

                                </Card.Title>

                                <Card.Text>

                                    Staff management

                                </Card.Text>

                            </Card.Body>

                        </Card>

                    </Col>

                )}

                <Col md={3}>

                    <Card>

                        <Card.Body>

                            <Card.Title>

                                My Profile

                            </Card.Title>

                            <Card.Text>

                                View your profile

                            </Card.Text>

                        </Card.Body>

                    </Card>

                </Col>

            </Row>

        </>

    );

}