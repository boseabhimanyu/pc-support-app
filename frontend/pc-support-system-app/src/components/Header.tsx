import { Navbar, Container, Nav, Button } from "react-bootstrap";
import { useAuth } from "../features/auth/hooks/useAuth";

export default function Header() {
    const { user, logout } = useAuth();
    async function handleLogout() { await logout(); }

    return (
        <Navbar className="app-header">
            <Container fluid>
                <Navbar.Brand className="brand-mark">
                    <span className="brand-icon">P</span>
                    <span>PC Support<small>Service desk</small></span>
                </Navbar.Brand>
                <Nav className="ms-auto align-items-center">
                    <div className="me-3 text-end user-summary">
                        <div>Welcome back, {user?.firstName}</div>
                        <small className="text-muted">{user?.role}</small>
                    </div>
                    <Button size="sm" variant="outline-danger" onClick={handleLogout}>Logout</Button>
                </Nav>
            </Container>
        </Navbar>
    );
}
