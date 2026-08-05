import { useState } from "react";
import { Button, Card, Form } from "react-bootstrap";

interface LoginFormProps {
    onSubmit: (data: {
        email: string;
        password: string;
    }) => Promise<void>;

    loading?: boolean;
}

export default function LoginForm({
    onSubmit,
    loading = false,
}: LoginFormProps) {

    const [email, setEmail] = useState("");

    const [password, setPassword] = useState("");

    async function handleSubmit(
        e: React.FormEvent
    ) {
        e.preventDefault();

        await onSubmit({
            email,
            password,
        });
    }

    return (
        <Card className="shadow-sm">

            <Card.Body>

                <h3 className="text-center mb-4">
                    Login
                </h3>

                <Form onSubmit={handleSubmit}>

                    <Form.Group className="mb-3">

                        <Form.Label>
                            Email
                        </Form.Label>

                        <Form.Control
                            required
                            type="email"
                            value={email}
                            onChange={(e) =>
                                setEmail(e.target.value)
                            }
                        />

                    </Form.Group>

                    <Form.Group className="mb-4">

                        <Form.Label>
                            Password
                        </Form.Label>

                        <Form.Control
                            required
                            type="password"
                            value={password}
                            onChange={(e) =>
                                setPassword(e.target.value)
                            }
                        />

                    </Form.Group>

                    <Button
                        type="submit"
                        className="w-100"
                        disabled={loading}
                    >
                        {loading
                            ? "Signing In..."
                            : "Login"}
                    </Button>

                </Form>

            </Card.Body>

        </Card>
    );
}