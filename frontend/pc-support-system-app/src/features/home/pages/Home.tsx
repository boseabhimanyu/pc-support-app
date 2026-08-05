import {  Card, Container } from "react-bootstrap";
import { Link } from "react-router-dom";

export default function Home() {
  return (
    <Container
    className="d-flex vh-100 justify-content-center align-items-center"
    style={{ transform: "translateY(-80px)" }}
>
      <Card className="shadow-sm text-center" style={{ maxWidth: 600 }}>
        <Card.Body className="p-5">

          <h2 className="mb-3">
            PC Support Management System
          </h2>

          <p className="text-muted mb-4">
            Manage repair jobs, customer devices and technician workflow
            through a centralized support platform.
          </p>

          <div className="d-flex justify-content-center gap-4 mt-4">

            <Link to="/login" className="btn btn-primary btn-lg px-4">
                Login
            </Link>

            <Link to="/register" className="btn btn-outline-primary btn-lg px-4">
                Register
            </Link>

          </div>

        </Card.Body>
      </Card>
    </Container>
  );
}