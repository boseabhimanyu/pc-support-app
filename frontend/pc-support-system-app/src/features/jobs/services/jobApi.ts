import {api} from "../../../app/api";

import type {
    JobDetails,
    JobListResponse,
} from "../types/job";
import type { Job, CreateJobRequest, AddJobNoteRequest, AddJobNoteResponse, JobQueueResponse } from "../jobTypes";

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
                `/jobs/number/${id}`
            );


        return response.data;

    },


};


export async function createJob(
    data: CreateJobRequest,
) {
    const response = await api.post<Job>(
        "/jobs",
        data,
    );

    return response.data;
}

export async function fetchJob(
    jobId: string,
) {
    const response = await api.get<Job>(
        `/jobs/${jobId}`,
    );

    return response.data;
}

export async function addJobNote(
    jobId: string,
    data: AddJobNoteRequest,
) {
    const response =
        await api.post<AddJobNoteResponse>(
            `/jobs/${jobId}/notes`,
            data,
        );

    return response.data;
}

export async function fetchOpenJobs() {
    const response = await api.get<JobQueueResponse>(
        "/jobs/open",
    );

    return response.data;
}

export async function fetchInProgressJobs() {
    const response = await api.get<JobQueueResponse>(
        "/jobs/in-progress",
    );

    return response.data;
}

export async function fetchWaitingCustomerJobs() {
    const response = await api.get<JobQueueResponse>(
        "/jobs/waiting-customer",
    );

    return response.data;
}

export async function fetchResumedJobs() {
    const response = await api.get<JobQueueResponse>(
        "/jobs/resumed",
    );

    return response.data;
}

export async function fetchAssignedJobs() {
    const response = await api.get<JobQueueResponse>(
        "/jobs/assigned",
    );

    return response.data;
}

export async function fetchJobByNumber(
    jobNumber: string,
) {
    const response = await api.get(
        `/jobs/number/${jobNumber}`,
    );

    return response.data;
}