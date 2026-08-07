import {api} from "../../app/api";

export async function updateStaff(
    staffId: string,
    data: {
        firstName: string;
        lastName: string;
        email: string;
        phone: string;
    },
) {
    const response = await api.patch(
        `/staff/${staffId}`,
        data,
    );

    return response.data;
}

export async function resetStaffPassword(
    staffId: string,
    password: string,
) {
    const response = await api.patch(
        `/staff/${staffId}/password`,
        {
            password,
        },
    );

    return response.data;
}