import { Card, Container } from "react-bootstrap";
import { Link } from "react-router";
export default function Home() {
  return <Container className="home-page d-flex vh-100 justify-content-center align-items-center"><Card className="home-card text-center" style={{ maxWidth: 680 }}><Card.Body className="p-5"><div className="home-logo"><span>P</span></div><div className="auth-kicker">PC SUPPORT</div><h2 className="mb-3">Service support, beautifully organized.</h2><p className="text-muted mb-4">Manage repair jobs, customer devices and technician workflow through a centralized support platform.</p><div className="d-flex justify-content-center gap-4 mt-4"><Link to="/login" className="btn btn-primary btn-lg px-4">Login</Link><Link to="/register" className="btn btn-outline-primary btn-lg px-4">Register</Link></div></Card.Body></Card></Container>;
}
