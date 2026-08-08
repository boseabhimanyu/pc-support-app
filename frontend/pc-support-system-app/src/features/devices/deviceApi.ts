import { api } from "../../app/api";
import type {Device} from "./types"


export async function fetchCustomerDevices(
    customerId: string,
) {
    const response = await api.get<Device[]>(
        `/customers/${customerId}/devices`,
    );

    return response.data;
}