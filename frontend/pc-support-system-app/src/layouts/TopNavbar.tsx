import { Navbar, Container, Nav, Button } from "react-bootstrap";
import { useAuth } from "../features/auth/hooks/useAuth";
export default function TopNavbar() {
  const { user, logout } = useAuth();
  return <Navbar className="app-header"><Container fluid><Navbar.Brand className="brand-mark"><span className="brand-icon">P</span><span>PC Support<small>Customer portal</small></span></Navbar.Brand><Nav className="ms-auto align-items-center"><span className="me-3 user-summary">Hello, {user?.firstName}</span><Button size="sm" variant="outline-danger" onClick={() => void logout()}>Logout</Button></Nav></Container></Navbar>;
}
