import { useState } from "react";
import { Alert, Container } from "react-bootstrap";
import { authApi } from "../services/authApi";
import { mapBackendError } from "../../../shared/utils/errorMapper";
import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";

import LoginForm from "../components/LoginForm";

export default function Login() {
    const navigate = useNavigate();
   const [loading, setLoading] = useState(false);
    const { setUser } = useAuth();
    const [error, setError] = useState("");
    
async function handleLogin(data: {
    email: string;
    password: string;
}) {
    try {
        setLoading(true);
        setError("");

        await authApi.login(data);

        console.log("Login successful");

        const user = await authApi.getCurrentUser();
       

        setUser(user);

        switch (user.role) {
    case "customer":
        navigate("/customer");
        console.log("Navigating to:", user.role);
        break;

    case "receptionist":
        navigate("/receptionist");
        break;

    case "technician":
        navigate("/technician");
        break;

    case "head_technician":
        navigate("/head-technician");
        break;

    case "admin":
        navigate("/admin");
        break;

    case "super_admin":
        navigate("/super-admin");
        break;

    default:
        navigate("/");
}

    } catch (err: any) {

        setError(
            mapBackendError(
                err.response?.data?.error ??
                "Login failed."
            )
        );

    } finally {
        setLoading(false);
    }
}

    return (
        <Container
            style={{
                maxWidth: "500px",
                marginTop: "50px",
            }}
        >
            {error && (
                <Alert variant="danger">
                    {error}
                </Alert>
            )}

            <LoginForm
                loading={loading}
                onSubmit={handleLogin}
            />
        </Container>
    );
}