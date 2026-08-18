import { useState } from "react";
import { Button, Card, Form } from "react-bootstrap";
import { Link } from "react-router";
interface LoginFormProps { onSubmit: (data: { email: string; password: string }) => Promise<void>; loading?: boolean }
export default function LoginForm({ onSubmit, loading = false }: LoginFormProps) {
  const [email, setEmail] = useState(""); const [password, setPassword] = useState("");
  async function handleSubmit(e: React.FormEvent) { e.preventDefault(); await onSubmit({ email, password }); }
  return <div className="auth-panel"><Card className="auth-card"><Card.Body><div className="auth-kicker">PC SUPPORT</div><h3 className="mb-2">Welcome back</h3><p className="text-muted mb-4">Sign in to manage your service desk.</p><Form onSubmit={handleSubmit}><Form.Group className="mb-3"><Form.Label>Email</Form.Label><Form.Control required type="email" value={email} onChange={(e) => setEmail(e.target.value)} /></Form.Group><Form.Group className="mb-4"><Form.Label>Password</Form.Label><Form.Control required type="password" value={password} onChange={(e) => setPassword(e.target.value)} /></Form.Group><Button type="submit" className="w-100" disabled={loading}>{loading ? "Signing In..." : "Login"}</Button></Form></Card.Body></Card><div className="mb-3 pt-3"><Link to="/" className="text-decoration-none">← Go back to Home</Link></div></div>;
}
