import type { UserResponse } from "../../features/auth/types";

export function canEditCustomer(
    user?: UserResponse | null,
): boolean {

    if (!user) {
        return false;
    }

    return [
        "receptionist",
        "admin",
    ].includes(user.role);

}


export function canCreateCustomer(
    user?: UserResponse | null,
): boolean {

    if (!user) {
        return false;
    }

    return [
        "receptionist",
        "admin",
        "head_technician"
    ].includes(user.role);

}


export function canViewStaff(
    user?: UserResponse | null,
): boolean {

    if (!user) {
        return false;
    }

    return [
        "head_technician",
        "admin",
        "super_admin",
    ].includes(user.role);

}


export function canManageStaff(
    user?: UserResponse | null,
): boolean {

    if (!user) {
        return false;
    }

    return [
        "admin",
        "super_admin",
    ].includes(user.role);

}

export function canResetCustomerPassword(
    user?: UserResponse | null,
): boolean {

    if (!user) {
        return false;
    }

    return [
        "receptionist",
        "admin",
    ].includes(user.role);

}

export type Role =
    | "customer"
    | "receptionist"
    | "technician"
    | "head_technician"
    | "admin"
    | "super_admin";

export const staffPermissions: Record<
    Role,
    {
        canView: boolean;
        canCreate: boolean;
        canEdit: boolean;
        assignableRoles: Role[];
    }
> = {
    customer: {
        canView: false,
        canCreate: false,
        canEdit: false,
        assignableRoles: [],
    },

    receptionist: {
        canView: false,
        canCreate: false,
        canEdit: false,
        assignableRoles: [],
    },

    technician: {
        canView: false,
        canCreate: false,
        canEdit: false,
        assignableRoles: [],
    },

    head_technician: {
        canView: true,
        canCreate: false,
        canEdit: false,
        assignableRoles: [],
    },

    admin: {
        canView: true,
        canCreate: true,
        canEdit: true,
        assignableRoles: [
            "receptionist",
            "technician",
            "head_technician",
        ],
    },

    super_admin: {
        canView: true,
        canCreate: true,
        canEdit: true,
        assignableRoles: [
            "admin",
        ],
    },
};