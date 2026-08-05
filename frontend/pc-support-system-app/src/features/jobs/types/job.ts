export interface JobListResponse {

    jobsCount: number;

    jobs: Job[];

}


export interface Job {

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

        brand: string;

        model?: string;

        serialNumber?: string;

    };


    problemDescription: string;


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

}



export interface JobDetails extends Job {

    notes: JobNote[];

}



export interface JobNote {

    id: string;


    author: {

        id: string;

        firstName: string;

        lastName: string;

        role: string;

    };


    note: string;


    createdAt: string;

}