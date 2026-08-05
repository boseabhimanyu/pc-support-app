import {api} from "../../../app/api";

import type {
    JobDetails,
    JobListResponse,
} from "../types/job";

export const jobApi = {


    async getMyJobs(): Promise<JobListResponse> {

        const response =
            await api.get(
                "/me/jobs"
            );

        return response.data;

    },


    async getJob(
        id: string
    ): Promise<JobDetails> {


        const response =
            await api.get(
                `/jobs/${id}`
            );


        return response.data;

    },


};