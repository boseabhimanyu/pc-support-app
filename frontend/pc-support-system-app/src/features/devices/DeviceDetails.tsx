
import { useEffect, useState } from "react";
import {
    Alert,
    Button,
    Card,
    Col,
    Row,
    Spinner,
} from "react-bootstrap";
import {
    useNavigate,
    useParams,
} from "react-router-dom";

import { api } from "../../app/api";
import { useAuth } from "../auth/hooks/useAuth";
import type { Device } from "./types";


export default function DeviceDetails() {




  
    const { deviceId } = useParams();

    const navigate = useNavigate();

    const { user } = useAuth();

    const [device, setDevice] =
        useState<Device | null>(null);

    const [loading, setLoading] =
        useState(true);

    const [error, setError] =
        useState("");

    const canEdit =
    user?.role === "head_technician" ||
    user?.role === "admin" ||
    user?.role === "receptionist";

    useEffect(() => {

        async function loadDevice() {

            if (!deviceId) {
                setError("Device not found.");
                setLoading(false);
                return;
            }

            try {

                setLoading(true);
                setError("");

                const response =
                    await api.get<Device>(
                        `/devices/${deviceId}`,
                    );

                setDevice(response.data);

            } catch (error: any) {

                console.error(
                    "Device details error:",
                    error.response?.data,
                );

                setError(
                    error.response?.data?.error ??
                    error.response?.data?.message ??
                    "Unable to load device details.",
                );

            } finally {

                setLoading(false);

            }

        }

        loadDevice();

    }, [deviceId]);


    function formatValue(
        value: string,
    ) {

        if (!value) {
            return "--";
        }

        return value
            .replace(/_/g, " ")
            .replace(/\b\w/g, (letter) =>
                letter.toUpperCase(),
            );

    }


    if (loading) {

        return (
            <div className="text-center p-4">

                <Spinner />

            </div>
        );

    }


    if (error) {

        return (
            <Alert variant="danger">
                {error}
            </Alert>
        );

    }


    if (!device) {

        return (
            <Alert variant="warning">
                Device not found.
            </Alert>
        );

    }


    return (

        <>

            <div className="d-flex justify-content-between align-items-center mb-4">

    <h2 className="mb-0">
        Device Details
    </h2>

    <div className="d-flex gap-2">

        <Button
            variant="secondary"
            onClick={() => navigate(-1)}
        >
            Back
        </Button>

        {canEdit && (
            <Button
                onClick={() => navigate("edit")}
            >
                Edit
            </Button>
        )}

    </div>

</div>


           <Card>
    <Card.Body>

        <Row className="mb-3">
            <Col md={3}>
                <strong>Type</strong>
            </Col>
            <Col md={9}>
                {device.type}
            </Col>
        </Row>

        <Row className="mb-3">
            <Col md={3}>
                <strong>Brand</strong>
            </Col>
            <Col md={9}>
                {device.brand || "--"}
            </Col>
        </Row>

        <Row className="mb-3">
            <Col md={3}>
                <strong>Model</strong>
            </Col>
            <Col md={9}>
                {device.model || "--"}
            </Col>
        </Row>

        <Row className="mb-3">
            <Col md={3}>
                <strong>Serial Number</strong>
            </Col>
            <Col md={9}>
                {device.serialNumber || "--"}
            </Col>
        </Row>

        <Row className="mb-3">
            <Col md={3}>
                <strong>Condition</strong>
            </Col>
            <Col md={9}>
                {device.condition}
            </Col>
        </Row>

        <Row className="mb-3">
            <Col md={3}>
                <strong>Notes</strong>
            </Col>
            <Col md={9}>
                {device.notes || "--"}
            </Col>
        </Row>

        <Row>
            <Col md={3}>
                <strong>Active</strong>
            </Col>
            <Col md={9}>
                {device.isActive ? "Yes" : "No"}
            </Col>
        </Row>

    </Card.Body>
</Card>

        </>

    );

}

