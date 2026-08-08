
export type AuditLogEntity =
    | "customer"
    | "device"
    | "job"
    | "user";

export type AuditLog = {
    ID: string;
    EntityType: AuditLogEntity;
    EntityID: string;
    EventType: string;
    PerformedByID: string;
    Metadata?: Record<string, unknown>;
    CreatedAt: string;
};

export type AuditLogResponse = {
    count: number;
    logs: AuditLog[];
};

export type AuditLogEntityDetails = {
    id: string;
    name?: string;
    jobNumber?: string;
    [key: string]: unknown;
};

export type AuditLogDisplayItem = AuditLog & {
    performedByName?: string;
    entityDisplay?: string;
};

