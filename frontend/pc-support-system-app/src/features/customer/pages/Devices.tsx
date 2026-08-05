import { useEffect, useState } from "react";

import { Alert, Container, Spinner } from "react-bootstrap";

import CustomerLayout from "../../../layouts/CustomerLayout";

import DeviceCard from "../../devices/components/DeviceCard";

import { customerApi } from "../services/customerApi";

import type { Device } from "../../auth/types";

export default function Devices() {

    const [devices, setDevices] =
        useState<Device[]>([]);

    const [loading, setLoading] =
        useState(true);

    const [error, setError] =
        useState("");

    useEffect(() => {

        loadDevices();

    }, []);

    async function loadDevices() {

        try {

            setLoading(true);

            const data =
                await customerApi.getMyDevices();

            setDevices(data);

        } catch {

            setError(
                "Unable to load devices."
            );

        } finally {

            setLoading(false);

        }

    }

    return (

        <CustomerLayout>

            <Container className="py-4">

                <h2 className="mb-4">
                    My Devices
                </h2>

                {loading && (

                    <Spinner animation="border" />

                )}

                {error && (

                    <Alert variant="danger">

                        {error}

                    </Alert>

                )}

                {!loading &&
                    devices.length === 0 && (

                    <Alert variant="info">

                        You have no registered devices.

                    </Alert>

                )}

                {devices.map((device) => (

                    <DeviceCard
                        key={device.id}
                        device={device}
                    />

                ))}

            </Container>

        </CustomerLayout>

    );

}