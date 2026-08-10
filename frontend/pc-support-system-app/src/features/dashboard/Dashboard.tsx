import { Card, Col, Row } from "react-bootstrap";
import { Link } from "react-router";
import { useAuth } from "../auth/hooks/useAuth";


const getRoleBasePath = (role?: string) => {
    switch (role) {
        case "admin":
        case "super_admin":
            return "/admin";

        case "head_technician":
            return "/head-technician";

        case "technician":
            return "/technician";

        case "receptionist":
            return "/receptionist";

        default:
            return "/";
    }
};



export default function Dashboard() {

    const { user } = useAuth();
    const basePath = getRoleBasePath(user?.role);
    return (

        <>

            <h2 className="mb-4">

                Dashboard

            </h2>

            <p className="text-muted mb-4">

                Welcome back, {user?.firstName}.

            </p>

            <Row className="g-3">
{(user?.role === "admin" ||
                  user?.role === "super_admin" ||
                  user?.role === "head_technician") && (

<Col md={3}>
    <Card
        as={Link}
        to={`${basePath}/customers`}
        className="text-decoration-none text-dark h-100"
    >
        <Card.Body>
            <Card.Title>Customers</Card.Title>
            <Card.Text>
                Customer management
            </Card.Text>
        </Card.Body>
    </Card>
</Col>
                  )}
<Col md={3}>
    <Card
        as={Link}
        to={`${basePath}/jobs`}
        className="text-decoration-none text-dark h-100"
    >
        <Card.Body>
            <Card.Title>Jobs</Card.Title>
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
    <Card
        as={Link}
        to={`${basePath}/staff`}
        className="text-decoration-none text-dark h-100"
    >
        <Card.Body>
            <Card.Title>Staff</Card.Title>
            <Card.Text>
                Staff management
            </Card.Text>
        </Card.Body>
    </Card>
</Col>
  )}
<Col md={3}>
    <Card
        as={Link}
        to={`${basePath}/profile`}
        className="text-decoration-none text-dark h-100"
    >
        <Card.Body>
            <Card.Title>My Profile</Card.Title>
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