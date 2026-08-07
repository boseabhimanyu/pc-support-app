import {api} from "../../../app/api";


export const customerApi = {


    getMyJobs: async () => {
        const response =
            await api.get("/me/jobs");

        return response.data;
    },


   getMyDevices: async () => {
    const response = await api.get("/me/devices");
    return response.data.devices;
    },


    updateJob: async (
        jobId: string,
        data: any
    ) => {

        const response =
            await api.patch(
                `/jobs/${jobId}/`,
                data
            );

        return response.data;
    },

};

export interface CustomerSummary {
    id: string;
    firstName: string;
    lastName: string;
    phone: string;
    email: string;
}

export async function searchCustomers(
    query: string,
): Promise<CustomerSummary[]> {

    const response = await api.get<CustomerSummary[]>(
        "/customers-search",
        {
            params: {
                q: query,
            },
        },
    );

    return response.data;
}

export async function createCustomer(
    data: Omit<CustomerSummary, "id">,
): Promise<CustomerSummary> {

    const response = await api.post<CustomerSummary>(
        "/customers",
        data,
    );

    return response.data;
}

export async function updateCustomer(
    customerId: string,
    data: Omit<CustomerSummary, "id">,
): Promise<CustomerSummary> {

    const response = await api.patch<CustomerSummary>(
        `/customers/${customerId}`,
        data,
    );

    return response.data;
}

export async function fetchCustomer(
    customerId: string,
): Promise<CustomerSummary> {

    const response = await api.get<CustomerSummary>(
        `/customers/${customerId}`,
    );

    return response.data;
}

export async function resetCustomerPassword(
    customerId: string,
    password: string,
): Promise<void> {

    await api.patch(
        `/customers/${customerId}/password`,
        {
            password,
        },
    );

}