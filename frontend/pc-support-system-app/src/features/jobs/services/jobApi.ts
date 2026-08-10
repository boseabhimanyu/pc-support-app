import {api} from "../../../app/api";

import type {
    JobDetails,
    JobListResponse,
} from "../types/job";
import type { Job, CreateJobRequest, AddJobNoteRequest, AddJobNoteResponse, JobQueueResponse } from "../jobTypes";
import type { CustomerJobsResponse } from "../../jobs/jobTypes";
import type { JobCustomerProfile } from "../jobTypes";
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
export async function searchJobs(query: string) {
    const response = await api.get(
        `/jobs/search`,
        {
            params: {
                q: query,
            },
        },
    );

    return response.data;
}

export async function fetchMyJobs() {
    const response = await api.get("/jobs/my");

    return response.data;
}

export type MyJobsResponse = {
    myJobsCount: number;
    jobs: Job[];
};

export async function getJobsByCustomer(
    customerId: string,
): Promise<CustomerJobsResponse> {
    const response = await api.get<CustomerJobsResponse>(
        `/jobs/customer/${customerId}/`,
    );

    return response.data;
}

export async function assignJob(
    jobId: string,
    staffId: string,
): Promise<JobCustomerProfile> {
    const response = await api.patch<JobCustomerProfile>(
        `/jobs/${jobId}/assign`,
        {
            staffId,
        },
    );

    return response.data;
}

export type AssignableStaff = {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    role: string;
    state: string;
};

export async function searchStaff(
    query: string,
): Promise<AssignableStaff[]> {
    const response = await api.get<AssignableStaff[]>(
        `/staff/search?q=${encodeURIComponent(query)}`,
    );

    return response.data;
}
export async function updateJobStatus(
    jobId: string,
    status: string,
): Promise<Job> {
    const response = await api.patch<Job>(
        `/jobs/${jobId}/status`,
        {
            status,
        },
    );

    return response.data;
}

export interface CloseJobRequest {
    reason: JobCloseReason;
    closureNotes: string;
    internalClosureNotes?: string;
}


export type JobCloseReason =
    | "completed"
    | "not_repairable"
    | "customer_cancelled"
    | "customer_no_response"
    | "customer_declined_repair"
    | "duplicate_job";


export async function closeJob(
    jobId: string,
    data: CloseJobRequest,
): Promise<Job> {
    const response = await api.post<Job>(
        `/jobs/${jobId}/close`,
        data,
    );

    return response.data;
}
