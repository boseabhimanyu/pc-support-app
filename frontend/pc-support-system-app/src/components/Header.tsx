import { Navbar, Container, Nav, Button } from "react-bootstrap";

import { useAuth } from "../features/auth/hooks/useAuth";

export default function Header() {

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
            className="border-bottom shadow-sm"
        >

            <Container fluid>

                <Navbar.Brand>

                    PC Support

                </Navbar.Brand>

                <Nav className="ms-auto align-items-center">

                    <div className="me-3 text-end">

                        <div>

                            Hello, {user?.firstName} 👋

                        </div>

                        <small className="text-muted">

                            {user?.role}

                        </small>

                    </div>

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