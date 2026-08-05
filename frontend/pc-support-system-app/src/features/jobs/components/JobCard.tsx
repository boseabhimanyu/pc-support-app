import { Card } from "react-bootstrap";
import { Link } from "react-router-dom";

import JobStatusBadge from "./JobStatusBadge";

import type { Job } from "../types/job";


interface JobCardProps {

    job: Job;

}


function formatDeviceType(
    value: string
) {

    return value
        .charAt(0)
        .toUpperCase() +
        value.slice(1).toLowerCase();

}



function formatBrand(
    value: string
) {

    if (
        value.toLowerCase() === "assembled"
    ) {

        return "Assembled";

    }

    return value.toUpperCase();

}



export default function JobCard({
    job,
}: JobCardProps) {


    return (

        <Card className="mb-3 shadow-sm">

            <Card.Body>


                <div className="d-flex justify-content-between align-items-center">


                    <Link
                        to={`/customer/jobs/${job.id}`}
                        className="fw-bold text-decoration-none"
                    >

                        {job.jobNumber}

                    </Link>


<div className="d-flex align-items-center">
    <JobStatusBadge status={job.status} />
</div>


                </div>



                <hr />



                <p className="mb-1">

                    <strong>
                        Device Type:
                    </strong>
                    {" "}
                    {formatDeviceType(job.device.type)}

                </p>



                <p className="mb-1">

                    <strong>
                        Brand:
                    </strong>
                    {" "}
                    {formatBrand(job.device.brand)}

                </p>



                {
                    job.device.model &&

                    <p className="mb-1">

                        <strong>
                            Model:
                        </strong>
                        {" "}
                        {job.device.model}

                    </p>

                }



                {
                    job.device.serialNumber &&

                    <p className="mb-1">

                        <strong>
                            Serial No:
                        </strong>
                        {" "}
                        {job.device.serialNumber}

                    </p>

                }



                <p className="mt-3">

                    {job.problemDescription}

                </p>



                <small className="text-muted">

                    Created:
                    {" "}
                    {new Date(
                        job.createdAt
                    ).toLocaleString()}

                </small>


            </Card.Body>

        </Card>

    );

}