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