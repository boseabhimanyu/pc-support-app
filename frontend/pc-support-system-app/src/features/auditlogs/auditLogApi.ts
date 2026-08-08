import { api } from "../../app/api";

import type {
    //AuditLog,
    AuditLogResponse,
} from "./auditLogTypes";

export async function fetchAuditLogs(
    entity?: "customer" | "device" | "job" | "user",
): Promise<AuditLogResponse> {
    const response = await api.get<AuditLogResponse>(
        "/audit-logs",
        {
            params: entity
                ? { entity }
                : undefined,
        },
    );

    return response.data;
}



export type AuditLogEntity =
    | "customer"
    | "device"
    | "job"
    | "user";



export async function fetchStaff(
    staffId: string,
) {
    const response = await api.get(
        `/staff/${staffId}`,
    );

    return response.data;
}

export async function fetchCustomer(
    customerId: string,
) {
    const response = await api.get(
        `/customers/${customerId}`,
    );

    return response.data;
}

export async function fetchDevice(
    deviceId: string,
) {
    const response = await api.get(
        `/devices/${deviceId}`,
    );

    return response.data;
}

export async function fetchJob(
    jobId: string,
) {
    const response = await api.get(
        `/jobs/${jobId}`,
    );

    return response.data;
}

