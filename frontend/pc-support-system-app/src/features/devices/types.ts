export const DEVICE_TYPES = [
    "laptop",
    "desktop",
    "printer",
    "monitor",
    "ups",
    "router",
] as const;

export const DEVICE_CONDITIONS = [
    "working",
    "not_working",
    "partially_working",
    "unknown",
] as const;

export type DeviceType =
    (typeof DEVICE_TYPES)[number];

export type DeviceCondition =
    (typeof DEVICE_CONDITIONS)[number];

export type DeviceFormData = {
    customerId: string;
    type: DeviceType | "";
    brand: string;
    model: string;
    serialNumber: string;
    notes: string;
    condition: DeviceCondition | "";
};

export type CustomerSearchResult = {
    id: string;
    firstName: string;
    lastName: string;
    phone: string;
    email?: string;
};

export type Device = {
    id: string;
    type: DeviceType;
    condition: DeviceCondition;
    brand: string;
    model: string;
    serialNumber: string;
    notes: string;
    isActive: boolean;
    customer?: {
        id: string;
        firstName: string;
        lastName: string;
        phone: string;
    };
};
