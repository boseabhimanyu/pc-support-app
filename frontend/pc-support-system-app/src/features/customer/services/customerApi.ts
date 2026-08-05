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