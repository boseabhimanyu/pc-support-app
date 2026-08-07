import { useEffect, useState } from "react";
import {
    Alert,
    Button,
    Card,
    Col,
    Row,
    Form,
    Spinner,
} from "react-bootstrap";
import {
    useNavigate,
    useParams,
} from "react-router-dom";

import {api} from "../../app/api";




import { useAuth } from "../auth/hooks/useAuth";
import {  updateStaff, resetStaffPassword } from "./staffApi";

type Staff = {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    role: string;
    state: string;
};

export default function StaffDetails() {


const { user } = useAuth();

const canEdit =
    user?.role === "admin" ||
    user?.role === "super_admin";

const [editing, setEditing] = useState(false);

const [form, setForm] = useState({
    firstName: "",
    lastName: "",
    email: "",
    phone: "",
});

const [saving, setSaving] = useState(false);
const [success, setSuccess] = useState("");
const [password, setPassword] = useState("");
const [showPassword, setShowPassword] = useState(false);
const [passwordMessage, setPasswordMessage] = useState("");
const [resettingPassword, setResettingPassword] =
    useState(false);

    const { staffId } = useParams();
    const navigate = useNavigate();

    const [staff, setStaff] =
        useState<Staff | null>(null);

    const [loading, setLoading] =
        useState(true);

    const [error, setError] =
        useState("");

    useEffect(() => {

        async function loadStaff() {

            if (!staffId) {
                setError("Staff ID is missing.");
                setLoading(false);
                return;
            }

            try {

                setLoading(true);
                setError("");

                const response =
                    await api.get<Staff>(
                        `/staff/${staffId}`,
                    );

                const staffData = response.data;

setStaff(staffData);

setForm({
    firstName: staffData.firstName ?? "",
    lastName: staffData.lastName ?? "",
    email: staffData.email ?? "",
    phone: staffData.phone ?? "",
});

            } catch (error: any) {

                console.error(
                    "Staff details error:",
                    error.response?.data,
                );

                setError(
                    error.response?.data?.error ??
                    error.response?.data?.message ??
                    "Unable to load staff details.",
                );

            } finally {

                setLoading(false);

            }

        }

        loadStaff();

    }, [staffId]);


    function formatRole(
        value: string,
    ) {

        return value
            .replace(/_/g, " ")
            .replace(/\b\w/g, (letter) =>
                letter.toUpperCase(),
            );

    }

    async function handlePasswordReset() {

    if (!staffId) {
        return;
    }

    setPasswordMessage("");

    const passwordPattern =
        /^(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/;

    if (!passwordPattern.test(password)) {
        setPasswordMessage(
            "Password must contain at least 8 characters, one uppercase letter, one number and one special character.",
        );
        return;
    }

    try {

        setResettingPassword(true);

        await resetStaffPassword(
            staffId,
            password,
        );

        setPasswordMessage(
            "Staff password reset successfully.",
        );

        setPassword("");

    } catch (error) {

        console.error(
            "Staff password reset error:",
            error,
        );

        setPasswordMessage(
            "Unable to reset staff password.",
        );

    } finally {

        setResettingPassword(false);

    }
}

async function handleUpdate(
    event: React.FormEvent<HTMLFormElement>,
) {
    event.preventDefault();

    if (!staffId) {
        return;
    }

    setSaving(true);
    setError("");
    setSuccess("");

    try {

        const updatedStaff =
            await updateStaff(
                staffId,
                {
                    firstName:
                        form.firstName.trim(),
                    lastName:
                        form.lastName.trim(),
                    email:
                        form.email.trim(),
                    phone:
                        form.phone.trim(),
                },
            );

        setStaff(updatedStaff);

        setSuccess(
            "Staff details updated successfully.",
        );

        setEditing(false);

    } catch (error: any) {

        console.error(
            "Staff update error:",
            error.response?.data,
        );

        setError(
            error.response?.data?.error ??
            error.response?.data?.message ??
            "Unable to update staff details.",
        );

    } finally {

        setSaving(false);

    }
}
    if (loading) {

        return (
            <div className="text-center p-4">
                <Spinner />
            </div>
        );

    }


   {error && (
    <Alert variant="danger" className="mb-3">
        {error}
    </Alert>
)}

{success && (
    <Alert variant="success" className="mb-3">
        {success}
    </Alert>
)}


    if (!staff) {

        return (
            <Alert variant="warning">
                Staff member not found.
            </Alert>
        );

    }


    return (

        <>

            <div className="d-flex justify-content-between align-items-center mb-4">

    <h2 className="mb-0">
        Staff Details
    </h2>

   <div className="d-flex gap-2">

    {canEdit && !editing && (
        <>

            <Button
                variant="primary"
                onClick={() => setEditing(true)}
            >
                Edit
            </Button>
        </>
    )}

    <Button
        variant="secondary"
        onClick={() => navigate(-1)}
    >
        Back
    </Button>

</div>

</div>

{error && (
    <Alert variant="danger" className="mb-3">
        {error}
    </Alert>
)}

{success && (
    <Alert variant="success" className="mb-3">
        {success}
    </Alert>
)}
            <Card>

                <Card.Body>

            {editing ? (

    <Form onSubmit={handleUpdate}>

        <Form.Group className="mb-3">
            <Form.Label>
                First Name
            </Form.Label>

            <Form.Control
                value={form.firstName}
                disabled={saving}
                onChange={(event) =>
                    setForm({
                        ...form,
                        firstName: event.target.value,
                    })
                }
            />
        </Form.Group>

        <Form.Group className="mb-3">
            <Form.Label>
                Last Name
            </Form.Label>

            <Form.Control
                value={form.lastName}
                disabled={saving}
                onChange={(event) =>
                    setForm({
                        ...form,
                        lastName: event.target.value,
                    })
                }
            />
        </Form.Group>

        <Form.Group className="mb-3">
            <Form.Label>
                Email
            </Form.Label>

            <Form.Control
                type="email"
                value={form.email}
                disabled={saving}
                onChange={(event) =>
                    setForm({
                        ...form,
                        email: event.target.value,
                    })
                }
            />
        </Form.Group>

        <Form.Group className="mb-4">
            <Form.Label>
                Phone
            </Form.Label>

            <Form.Control
                value={form.phone}
                disabled={saving}
                onChange={(event) =>
                    setForm({
                        ...form,
                        phone: event.target.value,
                    })
                }
            />
        </Form.Group>

        <Button
            type="submit"
            disabled={saving}
        >
            {saving ? "Saving..." : "Save Changes"}
        </Button>

        <Button
            type="button"
            variant="secondary"
            className="ms-2"
            disabled={saving}
            onClick={() => setEditing(false)}
        >
            Cancel
        </Button>

    </Form>

) : (

    <>
                 <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                First Name
                            </strong>
                        </Col>

                        <Col md={9}>
                            {staff.firstName}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Last Name
                            </strong>
                        </Col>

                        <Col md={9}>
                            {staff.lastName}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Email
                            </strong>
                        </Col>

                        <Col md={9}>
                            {staff.email}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Phone
                            </strong>
                        </Col>

                        <Col md={9}>
                            {staff.phone}
                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={3}>
                            <strong>
                                Role
                            </strong>
                        </Col>

                        <Col md={9}>
                            {formatRole(
                                staff.role,
                            )}
                        </Col>

                    </Row>


                    <Row>

                        <Col md={3}>
                            <strong>
                                State
                            </strong>
                        </Col>

                        <Col md={9}>
                            {formatRole(
                                staff.state,
                            )}
                        </Col>

                    </Row>
    </>

)}
                   

                </Card.Body>

            </Card>

            {canEdit && (
    <Card className="mt-3">

        <Card.Body>

            <Card.Title>
                Reset Staff Password
            </Card.Title>

            <Form onSubmit={(event) => {
                event.preventDefault();
                handlePasswordReset();
            }}>

                <Form.Group className="mb-3">

                    <Form.Label>
                        New Password
                    </Form.Label>

                    <div className="d-flex gap-2">

                        <Form.Control
                            type={
                                showPassword
                                    ? "text"
                                    : "password"
                            }
                            value={password}
                            disabled={resettingPassword}
                            onChange={(event) =>
                                setPassword(
                                    event.target.value,
                                )
                            }
                        />

                        <Button
                            type="button"
                            variant="outline-secondary"
                            onClick={() =>
                                setShowPassword(
                                    !showPassword,
                                )
                            }
                            disabled={resettingPassword}
                        >
                            {showPassword
                                ? "Hide"
                                : "Show"}
                        </Button>

                    </div>

                </Form.Group>

                <div className="text-muted small mb-3">
                    Password must contain at least 8
                    characters, one uppercase letter,
                    one number and one special character.
                </div>

                {passwordMessage && (
                    <Alert
                        variant={
                            passwordMessage.includes(
                                "successfully",
                            )
                                ? "success"
                                : "danger"
                        }
                    >
                        {passwordMessage}
                    </Alert>
                )}

                <Button
                    type="submit"
                    variant="warning"
                    disabled={
                        resettingPassword ||
                        !password
                    }
                >
                    {resettingPassword
                        ? "Resetting..."
                        : "Reset Password"}
                </Button>

            </Form>

        </Card.Body>

    </Card>
)}

        </>

    );

}