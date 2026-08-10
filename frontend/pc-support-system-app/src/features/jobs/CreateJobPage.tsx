import { useLocation, useParams } from "react-router";
import JobForm from "./JobForm";

type CreateJobLocationState = {
    customerName: string;
    deviceType: string;
    deviceBrand: string;
    deviceModel: string;
    deviceSerialNumber: string;
};

export default function CreateJobPage() {
    const { customerId, deviceId } = useParams();

    const location = useLocation();

    const state =
        location.state as CreateJobLocationState | null;

    if (!customerId || !deviceId) {
        return <div>Invalid customer or device.</div>;
    }

    if (!state) {
        return (
            <div>
                Customer or device information is missing.
            </div>
        );
    }

    return (
        <JobForm
            customerId={customerId}
            deviceId={deviceId}
            customerName={state.customerName}
            deviceType={state.deviceType}
            deviceBrand={state.deviceBrand}
            deviceModel={state.deviceModel}
            deviceSerialNumber={state.deviceSerialNumber}
            onCreated={(jobNumber) => {
                console.log("Created job:", jobNumber);
            }}
        />
    );
}