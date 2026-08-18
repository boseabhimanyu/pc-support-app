import { Nav } from "react-bootstrap";
import { NavLink } from "react-router";
export default function CustomerSidebar() {
  return <aside className="customer-sidebar"><Nav className="flex-column"><span className="sidebar-label">My account</span><Nav.Link as={NavLink} to="/customer/devices">My Devices</Nav.Link><Nav.Link as={NavLink} to="/customer/jobs">My Jobs</Nav.Link><Nav.Link as={NavLink} to="/customer/profile">Profile</Nav.Link></Nav></aside>;
}
