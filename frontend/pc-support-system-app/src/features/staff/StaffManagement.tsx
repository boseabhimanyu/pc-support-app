import { useState } from "react";
import {
    Alert,
    Button,
    Card,
    Col,
    Form,
    Row,
    Table,
} from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth/hooks/useAuth";


import {api} from "../../app/api";

type Staff = {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    role: string;
    state: string;
};

type StaffFormData = {
    firstName: string;
    lastName: string;
    phone: string;
    email: string;
    role: string;
};

const staffRoles = [
    "receptionist",
    "technician",
    "head_technician",
    "admin",
    "super_admin",
];

export default function StaffManagement() {
    
    const { user } = useAuth();

    const canCreateStaff =
    user?.role === "admin" ||
    user?.role === "super_admin";

    const navigate = useNavigate();
    const [search, setSearch] = useState("");
    const [staff, setStaff] = useState<Staff[]>([]);
    const [searching, setSearching] = useState(false);

    const [form, setForm] = useState<StaffFormData>({
        firstName: "",
        lastName: "",
        phone: "",
        email: "",
        role: "technician",
    });

    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [saving, setSaving] = useState(false);
    function handleChange(
        field: keyof StaffFormData,
        value: string,
    ) {

        setForm((current) => ({
            ...current,
            [field]: value,
        }));

    }

    async function handleSearch(
        event: React.FormEvent<HTMLFormElement>,
    ) {

        event.preventDefault();

        setError("");
        setMessage("");

        const query = search.trim();

        if (query.length < 3) {
            setError(
                "Enter at least 3 characters to search.",
            );
            setStaff([]);
            return;
        }

        try {

            setSearching(true);

            const response = await api.get<Staff[]>(
                `/staff/search?q=${encodeURIComponent(query)}`,
            );

            setStaff(response.data);

        } catch (error: any) {

            console.error(
                "Staff search error:",
                error.response?.data,
            );

            setStaff([]);

            setError(
                error.response?.data?.error ??
                error.response?.data?.message ??
                "Unable to search staff.",
            );

        } finally {

            setSearching(false);

        }

    }

    async function handleSubmit(
        event: React.FormEvent<HTMLFormElement>,
    ) {

        event.preventDefault();

        setMessage("");
        setError("");

        if (!form.firstName.trim()) {
            setError("First name is required.");
            return;
        }

        if (!form.lastName.trim()) {
            setError("Last name is required.");
            return;
        }

        if (!form.phone.trim()) {
            setError("Phone is required.");
            return;
        }

        if (!form.email.trim()) {
            setError("Email is required.");
            return;
        }

        if (!form.role) {
            setError("Role is required.");
            return;
        }

        try {

            setSaving(true);

            await api.post(
                "/staff",
                {
                    firstName: form.firstName.trim(),
                    lastName: form.lastName.trim(),
                    phone: form.phone.trim(),
                    email: form.email.trim(),
                    role: form.role,
                },
            );

            setMessage(
                "Staff member created successfully.",
            );

            setForm({
                firstName: "",
                lastName: "",
                phone: "",
                email: "",
                role: "technician",
            });

        } catch (error: any) {

            console.error(
                "Staff creation error:",
                error.response?.data,
            );

            setError(
                error.response?.data?.error ??
                error.response?.data?.message ??
                "Unable to create staff member.",
            );

        } finally {

            setSaving(false);

        }

    }

    function formatRole(role: string) {

        return role
            .replace(/_/g, " ")
            .replace(/\b\w/g, (letter) =>
                letter.toUpperCase(),
            );

    }

    return (

        <>

            <h2 className="mb-4">
                Staff Management
            </h2>


            {error && (
                <Alert variant="danger">
                    {error}
                </Alert>
            )}

            {message && (
                <Alert variant="success">
                    {message}
                </Alert>
            )}


            <Card className="mb-4">

                <Card.Body>

                    <Card.Title className="mb-3">
                        Search Staff
                    </Card.Title>

                    <Form
                        onSubmit={handleSearch}
                    >

                        <Row>

                            <Col md={8}>

                                <Form.Control
                                    value={search}
                                    placeholder="Search by name, email or phone"
                                    onChange={(event) =>
                                        setSearch(
                                            event.target.value,
                                        )
                                    }
                                />

                            </Col>

                            <Col md={4}>

                                <Button
                                    type="submit"
                                    disabled={searching}
                                >
                                    {searching
                                        ? "Searching..."
                                        : "Search"}
                                </Button>

                            </Col>

                        </Row>

                    </Form>

                </Card.Body>

            </Card>


            {staff.length > 0 && (

                <Card className="mb-4">

                    <Card.Body>

                        <Card.Title className="mb-3">
                            Search Results
                        </Card.Title>

                        <Table
                            responsive
                            hover
                        >

                            <thead>

                                <tr>
                                    <th>Name</th>
                                    <th>Email</th>
                                    <th>Phone</th>
                                    <th>Role</th>
                                    <th>State</th>
                                </tr>

                            </thead>

                            <tbody>

                                {staff.map(
                                    (member) => (

                                        <tr
                                            key={member.id}
                                        >

                                            <td>
                                                <Button
                                                variant="link"
                                                className="p-0"
                                                onClick={() =>
                                                    navigate(`${member.id}`)
                                                }
                                            >
                                                {member.firstName}{" "}
                                                {member.lastName}
                                             </Button>
                                            </td>

                                            <td>
                                                {member.email}
                                            </td>

                                            <td>
                                                {member.phone}
                                            </td>

                                            <td>
                                                {formatRole(
                                                    member.role,
                                                )}
                                            </td>

                                            <td>
                                                {formatRole(
                                                    member.state,
                                                )}
                                            </td>

                                        </tr>

                                    ),
                                )}

                            </tbody>

                        </Table>

                    </Card.Body>

                </Card>

            )}

{canCreateStaff && (
            <Card>

                <Card.Body>

                    <Card.Title className="mb-4">
                        Create Staff Member
                    </Card.Title>

                    <Form onSubmit={handleSubmit}>

                        <Row className="mb-3">

                            <Col md={6}>

                                <Form.Group>

                                    <Form.Label>
                                        First Name
                                    </Form.Label>

                                    <Form.Control
                                        value={form.firstName}
                                        disabled={saving}
                                        onChange={(event) =>
                                            handleChange(
                                                "firstName",
                                                event.target.value,
                                            )
                                        }
                                    />

                                </Form.Group>

                            </Col>


                            <Col md={6}>

                                <Form.Group>

                                    <Form.Label>
                                        Last Name
                                    </Form.Label>

                                    <Form.Control
                                        value={form.lastName}
                                        disabled={saving}
                                        onChange={(event) =>
                                            handleChange(
                                                "lastName",
                                                event.target.value,
                                            )
                                        }
                                    />

                                </Form.Group>

                            </Col>

                        </Row>


                        <Row className="mb-3">

                            <Col md={6}>

                                <Form.Group>

                                    <Form.Label>
                                        Phone
                                    </Form.Label>

                                    <Form.Control
                                        value={form.phone}
                                        disabled={saving}
                                        onChange={(event) =>
                                            handleChange(
                                                "phone",
                                                event.target.value,
                                            )
                                        }
                                    />

                                </Form.Group>

                            </Col>


                            <Col md={6}>

                                <Form.Group>

                                    <Form.Label>
                                        Email
                                    </Form.Label>

                                    <Form.Control
                                        type="email"
                                        value={form.email}
                                        disabled={saving}
                                        onChange={(event) =>
                                            handleChange(
                                                "email",
                                                event.target.value,
                                            )
                                        }
                                    />

                                </Form.Group>

                            </Col>

                        </Row>


                        <Row className="mb-4">

                            <Col md={6}>

                                <Form.Group>

                                    <Form.Label>
                                        Role
                                    </Form.Label>

                                    <Form.Select
                                        value={form.role}
                                        disabled={saving}
                                        onChange={(event) =>
                                            handleChange(
                                                "role",
                                                event.target.value,
                                            )
                                        }
                                    >

                                        {staffRoles.map(
                                            (role) => (

                                                <option
                                                    key={role}
                                                    value={role}
                                                >
                                                    {formatRole(
                                                        role,
                                                    )}
                                                </option>

                                            ),
                                        )}

                                    </Form.Select>

                                </Form.Group>

                            </Col>

                        </Row>


                        <Button
                            type="submit"
                            disabled={saving}
                        >
                            {saving
                                ? "Creating..."
                                : "Create Staff"}
                        </Button>

                    </Form>

                </Card.Body>

            </Card>
            )}

        </>

    );
}