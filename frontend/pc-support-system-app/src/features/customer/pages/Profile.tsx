import { useEffect, useState } from "react";
import { Alert } from "react-bootstrap";

import {
    Button,
    Card,
    Col,
    Container,
    Form,
    Row,
} from "react-bootstrap";

import { useAuth } from "../../auth/hooks/useAuth";
import { authApi } from "../../auth/services/authApi";
import CustomerLayout from "../../../layouts/CustomerLayout";
import { mapBackendError } from "../../../shared/utils/errorMapper";

export default function Profile() {

const { user, refreshUser } = useAuth();

const [editing, setEditing] = useState(false);

const [saving, setSaving] = useState(false);

const [success, setSuccess] = useState("");

const [error, setError] = useState("");

const [form, setForm] = useState({

    firstName: "",

    lastName: "",

    email: "",

    phone: "",

});

useEffect(() => {

    if (!user) return;

    setForm({

        firstName: user.firstName,

        lastName: user.lastName,

        email: user.email,

        phone: user.phone,


    });

}, [user]);

async function handleSave() {
    try {
        setSaving(true);

        setError("");
        setSuccess("");

        await authApi.updateProfile(form);

        await refreshUser();

        setEditing(false);

        setSuccess("Profile updated successfully.");

    } catch (err: any) {

        console.log(err.response?.data);

        setError(
            mapBackendError(
                err.response?.data?.error ??
                "Unable to update profile."
            )
        );

    } finally {
        setSaving(false);
    }
}




function handleCancel() {

    if (!user) return;

    setEditing(false);



    setForm({

        firstName: user.firstName,

        lastName: user.lastName,

        email: user.email,

        phone: user.phone,


    });

}

    return (

        

         <CustomerLayout>
        <Container className="py-4">
                {success && (
                    <Alert variant="success">
                        {success}
                    </Alert>
                )}

                {error && (
                    <Alert variant="danger">
                        {error}
                    </Alert>
                )}
            <h3 className="mb-4">
                Profile
            </h3>

            <Card className="shadow-sm">

                <Card.Body>

                    <Row className="mb-3">

                        <Col md={4}>
                            <strong>
                                First Name
                            </strong>
                        </Col>

                        <Col md={8}>

                            <Form.Control
    disabled={!editing}
    value={form.firstName}
    onChange={(e) =>
        setForm({
            ...form,
            firstName: e.target.value,
        })}
        />


                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={4}>
                            <strong>
                                Last Name
                            </strong>
                        </Col>

                        <Col md={8}>

                            <Form.Control
    disabled={!editing}
    value={form.lastName}
    onChange={(e) =>
        setForm({
            ...form,
            lastName: e.target.value,
        })}
/>



                        </Col>

                    </Row>


                    <Row className="mb-3">

                        <Col md={4}>
                            <strong>
                                Email
                            </strong>
                        </Col>

                        <Col md={8}>

                           <Form.Control
    disabled={!editing}
    value={form.email}
    onChange={(e) =>
        setForm({
            ...form,
            email: e.target.value,
        })}
/>



                        </Col>

                    </Row>


                    <Row className="mb-4">

                        <Col md={4}>
                            <strong>
                                Phone
                            </strong>
                        </Col>

                        <Col md={8}>

                            <Form.Control
    disabled={!editing}
    value={form.phone}
    onChange={(e) =>
        setForm({
            ...form,
            phone: e.target.value,
        })}
/>


                        </Col>

                    </Row>


                       <div className="text-end">

    {!editing ? (
        <Button
            onClick={() => setEditing(true)}
        >
            Edit Profile
        </Button>
    ) : (
        <>
            <Button
                className="me-2"
                onClick={handleSave}
                disabled={saving}
            >
                {saving ? "Saving..." : "Save"}
            </Button>

            <Button
                variant="secondary"
                onClick={handleCancel}
            >
                Cancel
            </Button>
        </>
    )}

</div>

                </Card.Body>

            </Card>

        </Container>

        </CustomerLayout>
    );
}