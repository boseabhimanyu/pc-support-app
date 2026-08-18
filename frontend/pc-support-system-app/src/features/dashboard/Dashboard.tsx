import { Card, Col, Row } from "react-bootstrap";
import { Link } from "react-router";
import { useAuth } from "../auth/hooks/useAuth";

const getRoleBasePath = (role?: string) => {
  switch (role) { case "admin": case "super_admin": return "/admin"; case "head_technician": return "/head-technician"; case "technician": return "/technician"; case "receptionist": return "/receptionist"; default: return "/"; }
};
export default function Dashboard() {
  const { user } = useAuth(); const basePath = getRoleBasePath(user?.role);
  const canManage = ["admin", "super_admin", "head_technician"].includes(user?.role ?? "");
  const cards = [
    ...(canManage ? [{ title: "Customers", text: "Customer management", path: "customers" }] : []),
    { title: "Jobs", text: "Job management", path: "jobs" },
    ...(canManage ? [{ title: "Staff", text: "Staff management", path: "staff" }] : []),
    { title: "My Profile", text: "View your profile", path: "profile" },
  ];
  return <><div className="page-heading"><div><span className="eyebrow">OVERVIEW</span><h2>Dashboard</h2></div><div className="heading-date">Service workspace</div></div><p className="text-muted mb-4">Welcome back, {user?.firstName}.</p><Row className="g-3 dashboard-grid">{cards.map((card) => <Col key={card.path} md={3} sm={6}><Card as={Link} to={`${basePath}/${card.path}`} className="text-decoration-none text-dark h-100"><Card.Body><Card.Title>{card.title}</Card.Title><Card.Text>{card.text}</Card.Text></Card.Body></Card></Col>)}</Row></>;
}
