import { useEffect, useState } from "react";
import {
    Alert,
    Button,
    Card,
    Col,
    Form,
    Row,
    Spinner,
} from "react-bootstrap";
import { useParams } from "react-router-dom";

import {
    createCustomer,
    updateCustomer,
} from "../customer/services/customerApi";

import { fetchCustomer } from "../customer/services/customerApi";

type CustomerFormData = {
    firstName: string;
    lastName: string;
    phone: string;
    email: string;
};

export default function CustomerForm() {

    const { customerId } = useParams();

    const [form, setForm] = useState<CustomerFormData>({
        firstName: "",
        lastName: "",
        phone: "",
        email: "",
    });

    const [loadingCustomer, setLoadingCustomer] =
        useState(false);

    const [saving, setSaving] =
        useState(false);

    const [error, setError] =
        useState("");

    const [success, setSuccess] =
        useState("");

    useEffect(() => {

        async function loadCustomer() {

            if (!customerId) {
                return;
            }

            try {

                setLoadingCustomer(true);
                setError("");

                const response =
                    await fetchCustomer(customerId);

                setForm({
                    firstName: response.firstName ?? "",
                    lastName: response.lastName ?? "",
                    phone: response.phone ?? "",
                    email: response.email ?? "",
                });

            } catch (err) {

                console.error(err);

                setError(
                    "Unable to load customer.",
                );

            } finally {

                setLoadingCustomer(false);

            }
        }

        loadCustomer();

    }, [customerId]);

    function handleChange(
        field: keyof CustomerFormData,
        value: string,
    ) {

        setForm((current) => ({
            ...current,
            [field]: value,
        }));

    }

async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
) {

    event.preventDefault();

    setError("");
    setSuccess("");

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

    const email = form.email.trim();

    const emailPattern =
        /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

    if (!emailPattern.test(email)) {
        setError("Enter a valid email address.");
        return;
    }

    try {

        setSaving(true);

        const data = {
            firstName: form.firstName.trim(),
            lastName: form.lastName.trim(),
            phone: form.phone.trim(),
            email,
        };

        if (customerId) {

            await updateCustomer(
                customerId,
                data,
            );

            setSuccess(
                "Customer updated successfully.",
            );

        } else {

            await createCustomer(data);

            setSuccess(
                "Customer created successfully.",
            );

            setForm({
                firstName: "",
                lastName: "",
                phone: "",
                email: "",
            });

        }

    } catch (err: any) {

        console.error(
            "Customer save error:",
            err.response?.data,
        );

        setError(
            err.response?.data?.error ??
            err.response?.data?.message ??
            (
                customerId
                    ? "Unable to update customer."
                    : "Unable to create customer."
            ),
        );

    } finally {

        setSaving(false);

    }

}

    function handleReset() {

        setError("");
        setSuccess("");

        if (customerId) {

            return;

        }

        setForm({
            firstName: "",
            lastName: "",
            phone: "",
            email: "",
        });

    }

    if (loadingCustomer) {

        return (
            <div className="text-center p-4">

                <Spinner />

            </div>
        );

    }

    return (

        <Card>

            <Card.Body>

                <Card.Title className="mb-4">

                    {customerId
                        ? "Update Customer"
                        : "Create Customer"}

                </Card.Title>

                {error && (
                    <Alert variant="danger">
                        {error}
                    </Alert>
                )}

                {success && (
                    <Alert variant="success">
                        {success}
                    </Alert>
                )}

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

                    <Row className="mb-4">

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
    required
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

                    <div className="d-flex gap-2">

                        <Button
                            type="submit"
                            disabled={saving}
                        >

                            {saving ? (
                                <>
                                    <Spinner
                                        size="sm"
                                        className="me-2"
                                    />

                                    Saving...
                                </>
                            ) : (
                                customerId
                                    ? "Update"
                                    : "Save"
                            )}

                        </Button>

                        <Button
                            type="button"
                            variant="secondary"
                            disabled={saving}
                            onClick={handleReset}
                        >
                            Reset
                        </Button>

                    </div>

                </Form>

            </Card.Body>

        </Card>
    );
}