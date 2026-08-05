import { Card } from "react-bootstrap";

import type { Device } from "../../auth/types";

interface Props {
    device: Device;
}

export default function DeviceCard({
    device,
}: Props) {
    return (
    <Card className="mb-3 shadow-sm">
        <Card.Body>

            <Card.Title>
                {device.brand} {device.model}
            </Card.Title>

            <Card.Text>

                <strong>Device Type:</strong> {device.type.charAt(0).toUpperCase() + device.type.slice(1)}

                <br />

                <strong>Brand:</strong> {device.brand.toUpperCase()}

                <br />

                <strong>Model:</strong> {device.model}

                <br />

                <strong>Serial No.:</strong> {device.serialNumber}

                <br />

                <strong>Notes:</strong>{" "}

                {device.notes || "-"}

            </Card.Text>

        </Card.Body>
    </Card>
    );
}