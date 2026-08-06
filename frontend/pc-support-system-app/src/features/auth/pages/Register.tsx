import { useState } from "react";
import { Alert, Container } from "react-bootstrap";
import { useNavigate } from "react-router-dom";

import RegisterForm from "../components/RegisterForm";
import { authApi } from "../services/authApi";
import type { RegisterRequest } from "../types";
import { mapBackendError } from "../../../shared/utils/errorMapper";

export default function Register() {
  const navigate = useNavigate();

  const [loading, setLoading] = useState(false);

  const [error, setError] = useState("");

  const [success, setSuccess] = useState("");

  async function handleRegister(data: RegisterRequest) {
    try {
      setLoading(true);
      setError("");
      setSuccess("");

      await authApi.register(data);

      setSuccess("Registration successful.");

      setTimeout(() => {
        navigate("/login");
      }, 1500);

    } catch (err: any) {

     setError(
        mapBackendError(
        err.response?.data?.error ?? "Registration failed."
    )
);
        setTimeout(() => {
            setError("");
        }, 3500);

    } finally {
      setLoading(false);
    }
  }

  return (
    <Container
      style={{
        maxWidth: "600px",
        marginTop: "50px",
      }}
    >
      {error &&

        <Alert variant="danger">

          {error}

        </Alert>

      }

      {success &&

        <Alert variant="success">

          {success}

        </Alert>

      }

      <RegisterForm

        loading={loading}

        onSubmit={handleRegister}

      />

    </Container>
  );
}
