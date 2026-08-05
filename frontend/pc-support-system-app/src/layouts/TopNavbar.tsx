import { Navbar, Container, Nav, Button } from "react-bootstrap";
//import { useNavigate } from "react-router-dom";

import { useAuth } from "../features/auth/hooks/useAuth";

export default function TopNavbar() {

   // const navigate = useNavigate();

    const {
        user,
        logout,
    } = useAuth();

    async function handleLogout() {

        await logout();

    }

    return (

        <Navbar
            bg="light"
            className="border-bottom"
        >

            <Container fluid>

                <Navbar.Brand>

                    PC Support

                </Navbar.Brand>

                <Nav className="ms-auto align-items-center">

                    <span className="me-3">

                        Hello, {user?.firstName} 👋

                    </span>

                    <Button
                        size="sm"
                        variant="outline-danger"
                        onClick={handleLogout}
                    >
                        Logout
                    </Button>

                </Nav>

            </Container>

        </Navbar>

    );
}