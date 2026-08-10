
export type JobStatus =
    | "created"
    | "assigned"
    | "in_progress"
    | "waiting_customer"
    | "resumed"
    | "closed";

export type JobCustomer = {
    id: string;
    firstName: string;
    lastName: string;
    phone: string;
};

export type JobDevice = {
    id: string;
    type: string;
    brand: string;
    model?: string;
    serialNumber?: string;
};

export type JobAuthor = {
    id: string;
    firstName: string;
    lastName: string;
    role: string;
};

export type JobNote = {
    id: string;
    author: JobAuthor;
    note: string;
    createdAt: string;
};

export type Job = {
    id: string;
    jobNumber: string;
    status: JobStatus | string;

    customer: JobCustomer;

    device: JobDevice;

    problemDescription: string;

    notes: JobNote[] | null;

    createdAt: string;
    closedAt?: string;
    createdBy: JobAuthor;
    assignedTo?: JobAssignedTo | null;
    closeReason?: string;
    closureNotes?: string;
    internalClosureNotes?: string;

};

export type CreateJobRequest = {
    customerId: string;
    deviceId: string;
    problemDescription: string;
};

export type AddJobNoteRequest = {
    note: string;
};

export type AddJobNoteResponse = {
    message: string;
};

export type JobQueueResponse = {
    openJobsCount?: number;
    assignedJobsCount?: number;
    jobsCount?: number;
    jobs: Job[];
};

export type JobAssignedTo = {
    id: string;
    firstName: string;
    lastName: string;
    role: string;
};

export type JobCustomerProfile = {
    id: string;
    jobNumber: string;
    status: string;

    customer: {
        id: string;
        firstName: string;
        lastName: string;
        phone: string;
    };

    device: {
        id: string;
        type: string;
        brand?: string;
        model?: string;
        serialNumber?: string;
    };

    problemDescription: string;
    notes: string | null;
    createdAt: string;

    createdBy: {
        id: string;
        firstName: string;
        lastName: string;
        role: string;
    };

    assignedTo?: {
        id: string;
        firstName: string;
        lastName: string;
        role: string;
    };

    closeReason?: string;
    closureNotes?: string;
    closedAt?: string;
    internalClosureNotes?: string;
    
};

export type CustomerJobsResponse = {
    jobsCount: number;
    jobs: JobCustomerProfile[];
};
