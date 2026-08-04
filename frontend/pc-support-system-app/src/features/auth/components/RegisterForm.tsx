import { useState } from "react";
import {
  Button,
  Card,
  Col,
  Form,
  Row,
  Alert,
} from "react-bootstrap";

import type { RegisterRequest } from "../types";

interface RegisterFormProps {
  onSubmit: (data: RegisterRequest) => Promise<void>;
  loading?: boolean;
}

export default function RegisterForm({
  onSubmit,
  loading = false,
}: RegisterFormProps) {
  const [form, setForm] = useState<RegisterRequest>({
    firstName: "",
    lastName: "",
    phone: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
const [errors, setErrors] = useState<string[]>([]);
 

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    setForm({
      ...form,
      [e.target.name]: e.target.value,
    });
  };

const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const validationErrors: string[] = [];

    if (form.password.length < 8)
        validationErrors.push("Password must be at least 8 characters.");

    if (!/[A-Z]/.test(form.password))
        validationErrors.push("Password must contain at least one uppercase letter.");

    if (!/[a-z]/.test(form.password))
        validationErrors.push("Password must contain at least one lowercase letter.");

    if (!/\d/.test(form.password))
        validationErrors.push("Password must contain at least one number.");

    if (!/[!@#$%^&*(),.?":{}|<>]/.test(form.password))
        validationErrors.push("Password must contain at least one special character.");

    if (form.password !== form.confirmPassword)
        validationErrors.push("Passwords do not match.");

    if (validationErrors.length > 0) {
        setErrors(validationErrors);
        return;
    }

    setErrors([]);

    await onSubmit(form);
};

  return (
    <Card className="shadow-sm">
      <Card.Body>

        <h3 className="mb-4 text-center">
          Customer Registration
        </h3>

        <Form onSubmit={handleSubmit}>

          <Row>

            <Col md={6}>
              <Form.Group className="mb-3">

                <Form.Label>
                  First Name
                </Form.Label>

                <Form.Control
                  required
                  name="firstName"
                  value={form.firstName}
                  onChange={handleChange}
                />

              </Form.Group>
            </Col>

            <Col md={6}>
              <Form.Group className="mb-3">

                <Form.Label>
                  Last Name
                </Form.Label>

                <Form.Control
                  required
                  name="lastName"
                  value={form.lastName}
                  onChange={handleChange}
                />

              </Form.Group>
            </Col>

          </Row>

          <Form.Group className="mb-3">

            <Form.Label>
              Phone
            </Form.Label>

            <Form.Control
              required
              name="phone"
              value={form.phone}
              onChange={handleChange}
            />

          </Form.Group>

          <Form.Group className="mb-3">

            <Form.Label>
              Email
            </Form.Label>

            <Form.Control
              required
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
            />

          </Form.Group>



<Form.Group className="mb-3">

    <Form.Label>Password</Form.Label>

    <Form.Control
        required
        type="password"
        name="password"
        value={form.password}
        onChange={handleChange}
    />

</Form.Group>

<Form.Group className="mb-4">

    <Form.Label>Confirm Password</Form.Label>

    <Form.Control
        required
        type="password"
        name="confirmPassword"
        value={form.confirmPassword}
        onChange={handleChange}
    />

</Form.Group>

{errors.length > 0 && (
    <Alert variant="danger" className="mb-3">
        <ul className="mb-0">
            {errors.map((err, index) => (
                <li key={index}>{err}</li>
            ))}
        </ul>
    </Alert>
)}

<Button
    type="submit"
    className="w-100"
    disabled={loading}
>
    {loading ? "Creating Account..." : "Register"}
</Button>

        </Form>

      </Card.Body>
    </Card>
  );
}