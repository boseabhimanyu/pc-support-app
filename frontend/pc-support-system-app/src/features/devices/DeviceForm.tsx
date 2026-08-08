
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
import { useParams, useLocation } from "react-router-dom";

import { api } from "../../app/api";

import {
    DEVICE_TYPES,
    DEVICE_CONDITIONS,
} from "./types";

import type {
    Device,
    DeviceFormData,
} from "./types";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth/hooks/useAuth";

export default function DeviceForm() {
    const { customerId, deviceId } = useParams();
    const navigate = useNavigate();
        const location = useLocation();
           const customerName =
        location.state?.customerName ?? "";

    const [form, setForm] = useState<DeviceFormData>({
        customerId: customerId ?? "",
        type: "",
        brand: "",
        model: "",
        serialNumber: "",
        notes: "",
        condition: "",
    });


    const { user } = useAuth();

const canEditSerialNumber =
    user?.role === "head_technician" ||
    user?.role === "admin";

    const [loadingDevice, setLoadingDevice] =
        useState(false);

    const [saving, setSaving] =
        useState(false);

    const [error, setError] =
        useState("");

    const [success, setSuccess] =
        useState("");

    useEffect(() => {
        async function loadDevice() {
            if (!deviceId) {
                return;
            }

            try {
                setLoadingDevice(true);
                setError("");

                console.log("LOADING DEVICE:", deviceId);
                const response =
                    await api.get<Device>(
                    `/devices/${deviceId}`,
                       );

                console.log("DEVICE API RESPONSE:", response.data);

                const device = response.data;
                

                setForm({
                    customerId:
                        device.customer?.id ?? "",
                    type: device.type ?? "",
                    brand: device.brand ?? "",
                    model: device.model ?? "",
                    serialNumber:
                        device.serialNumber ?? "",
                    notes: device.notes ?? "",
                    condition:
                        device.condition ?? "",
                });


            } catch (err: any) {
                console.error(
                    "Device details error:",
                    err.response?.data,
                );

                setError(
                    err.response?.data?.error ??
                    err.response?.data?.message ??
                    "Unable to load device.",
                );
            } finally {
                setLoadingDevice(false);
            }
        }

        loadDevice();
    }, [deviceId]);



    function handleChange(
        field: keyof DeviceFormData,
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

        if (!deviceId && !customerId) {
            setError("Customer ID is missing.");
            return;
        }

        if (!form.type) {
            setError("Device type is required.");
            return;
        }

if (form.type !== "desktop") {
    if (!form.brand.trim()) {
        setError("Brand is required for this device type.");
        return;
    }

    if (form.brand.trim().length < 2) {
        setError("Brand must contain at least 2 characters.");
        return;
    }

    if (!form.model.trim()) {
        setError("Model is required for this device type.");
        return;
    }

    if (form.model.trim().length < 2) {
        setError("Model must contain at least 2 characters.");
        return;
    }

    if (!form.serialNumber.trim()) {
        setError(
            "Serial number is required for this device type.",
        );
        return;
    }
}
if (!form.condition) {
    setError("Device condition is required.");
    return;
}

       const data = {
    customerId: deviceId
        ? form.customerId
        : customerId,
    type: form.type,
    brand: form.brand.trim(),
    model: form.model.trim(),
    serialNumber: form.serialNumber.trim(),
    notes: form.notes.trim(),
    condition: form.condition,
};

        try {
            setSaving(true);

            if (deviceId) {
                await api.patch(
                    `/devices/${deviceId}`,
                    data,
                );

                setSuccess(
                    "Device updated successfully.",
                );
            } else {
                await api.post(
                    "/devices",
                    data,
                );

                setSuccess(
                    "Device created successfully.",
                );

                setForm({
    customerId: customerId ?? "",
    type: "",
    brand: "",
    model: "",
    serialNumber: "",
    notes: "",
    condition: "",
});

            }
        } catch (err: any) {
            console.error(
                "Device save error:",
                err.response?.data,
            );

            setError(
                err.response?.data?.error ??
                err.response?.data?.message ??
                (
                    deviceId
                        ? "Unable to update device."
                        : "Unable to create device."
                ),
            );
        } finally {
            setSaving(false);
        }
    }



    if (loadingDevice) {
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
                    {deviceId
                        ? "Update Device"
                        : "Create Device"}
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
                    <Form.Group className="mb-3">
    <Form.Label>
        Customer
    </Form.Label>

 <Form.Control
    value={customerName}
    disabled
/>
</Form.Group>

                    <Row className="mb-3">
                        <Col md={6}>
                            <Form.Group>
                                <Form.Label>
                                    Device Type
                                </Form.Label>

                                <Form.Select
                                    value={form.type}
                                    disabled={saving}
                                    onChange={(event) =>
                                        handleChange(
                                            "type",
                                            event.target
                                                .value,
                                        )
                                    }
                                >
                                    <option value="">
                                        Select device type
                                    </option>

                                    {DEVICE_TYPES.map(
                                        (type) => (
                                            <option
                                                key={type}
                                                value={type}
                                            >
                                                {type
                                                    .replace(
                                                        "_",
                                                        " ",
                                                    )
                                                    .replace(
                                                        /\b\w/g,
                                                        (
                                                            letter,
                                                        ) =>
                                                            letter.toUpperCase(),
                                                    )}
                                            </option>
                                        ),
                                    )}
                                </Form.Select>
                            </Form.Group>
                        </Col>

                        <Col md={6}>
                            <Form.Group>
                                <Form.Label>
                                    Condition
                                </Form.Label>

                                <Form.Select
                                    value={
                                        form.condition
                                    }
                                    disabled={saving}
                                    onChange={(event) =>
                                        handleChange(
                                            "condition",
                                            event.target
                                                .value,
                                        )
                                    }
                                >
                                    <option value="">
                                        Select condition
                                    </option>

                                    {DEVICE_CONDITIONS.map(
                                        (condition) => (
                                            <option
                                                key={
                                                    condition
                                                }
                                                value={
                                                    condition
                                                }
                                            >
                                                {condition
                                                    .replace(
                                                        /_/g,
                                                        " ",
                                                    )
                                                    .replace(
                                                        /\b\w/g,
                                                        (
                                                            letter,
                                                        ) =>
                                                            letter.toUpperCase(),
                                                    )}
                                            </option>
                                        ),
                                    )}
                                </Form.Select>
                            </Form.Group>
                        </Col>
                    </Row>

                    <Row className="mb-3">
                        <Col md={6}>
                            <Form.Group>
                                <Form.Label>
                                    Brand
                                </Form.Label>

                                <Form.Control
                                    value={form.brand}
                                    disabled={saving}
                                    onChange={(event) =>
                                        handleChange(
                                            "brand",
                                            event.target
                                                .value,
                                        )
                                    }
                                />
                            </Form.Group>
                        </Col>

                        <Col md={6}>
                            <Form.Group>
                                <Form.Label>
                                    Model
                                </Form.Label>

                                <Form.Control
                                    value={form.model}
                                    disabled={saving}
                                    onChange={(event) =>
                                        handleChange(
                                            "model",
                                            event.target
                                                .value,
                                        )
                                    }
                                />
                            </Form.Group>
                        </Col>
                    </Row>

                    <Form.Group className="mb-3">
                        <Form.Label>
                            Serial Number
                        </Form.Label>

                       <Form.Control
                        value={form.serialNumber}
                        disabled={
                                    saving ||
                                    (!!deviceId && !canEditSerialNumber)
                                  }
                        onChange={(event) =>
                            handleChange(
                                "serialNumber",
                                event.target.value,
                            )
                        }
                    />
                    </Form.Group>

                    <Form.Group className="mb-4">
                        <Form.Label>
                            Notes
                        </Form.Label>

                        <Form.Control
                            as="textarea"
                            rows={4}
                            value={form.notes}
                            disabled={saving}
                            onChange={(event) =>
                                handleChange(
                                    "notes",
                                    event.target.value,
                                )
                            }
                        />
                    </Form.Group>

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
                                deviceId
                                    ? "Update"
                                    : "Save"
                            )}
                        </Button>

                       <Button
                     variant="secondary"
                    onClick={() => navigate(-1)}
                                >
                    Back
                </Button>
                    </div>
                </Form>
            </Card.Body>
        </Card>
    );
}

