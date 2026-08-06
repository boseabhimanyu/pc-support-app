import { useEffect, useState } from "react";

import JobStatusBadge from "../../jobs/components/JobStatusBadge";
import {
    Alert,
    Card,
    Container,
    Spinner,
} from "react-bootstrap";

import { useParams } from "react-router-dom";

import CustomerLayout from "../../../layouts/CustomerLayout";

import {
    jobApi
} from "../../jobs/services/jobApi";

import type {
    JobDetails as JobDetailsType
} from "../../jobs/types/job";



export default function JobDetails() {


    function formatDateTime(date: string) {

    return new Date(date).toLocaleString("en-IN", {
        day: "2-digit",
        month: "short",
        year: "numeric",
        hour: "numeric",
        minute: "2-digit",
    });

}

    const { jobNumber } = useParams();


    const [job, setJob] =
        useState<JobDetailsType | null>(null);


    const [loading, setLoading] =
        useState(true);


    const [error, setError] =
        useState("");



    useEffect(() => {

        if (jobNumber) {

            loadJob();

        }

    }, [jobNumber]);



    async function loadJob() {

        try {

            setLoading(true);

            const data =
                await jobApi.getJob(jobNumber!);

            setJob(data);


        } catch {

            setError(
                "Unable to load job details."
            );

        } finally {

            setLoading(false);

        }

    }



    return (

        <CustomerLayout>

            <Container className="py-4">


                {
                    loading &&

                    <div className="text-center">

                        <Spinner animation="border"/>

                    </div>
                }



                {
                    error &&

                    <Alert variant="danger">

                        {error}

                    </Alert>

                }



                {
                    job &&

                    <>


                            <div className="d-flex justify-content-between align-items-start mb-4">

                                <div>

                                    <h3 className="mb-1">
                                        {job.jobNumber}
                                    </h3>

                                    <small className="text-muted">
                                        Created on{" "}
                                        {formatDateTime(job.createdAt)}
                                    </small>

                                </div>

                                <JobStatusBadge
                                    status={job.status}
                                />

                            </div>



                    <Card className="mb-3">

                        <Card.Body>

                            <h5>
                                Device
                            </h5>


                            <p>
                                Type: {job.device.type}
                            </p>


                            <p>
                                Brand: {job.device.brand}
                            </p>


                            {
                                job.device.model &&

                                <p>
                                    Model: {job.device.model}
                                </p>

                            }


                            {
                                job.device.serialNumber &&

                                <p>
                                    Serial No:
                                    {" "}
                                    {job.device.serialNumber}
                                </p>

                            }


                        </Card.Body>

                    </Card>





                    <Card className="mb-3">

                        <Card.Body>

                            <h5>
                                Problem
                            </h5>


                            <p>
                                {job.problemDescription}
                            </p>


                        </Card.Body>

                    </Card>





                    <Card>

                        <Card.Body>

                            <h5>
                                Notes
                            </h5>


                            {job.notes.map((note, index) => (
    <div
        key={note.id}
        className={`py-2 ${
            index !== job.notes.length - 1 ? "border-bottom" : ""
        }`}
    >
        <div className="d-flex justify-content-between align-items-start">

    <div>

        <strong>
            {note.author.firstName} {note.author.lastName}
        </strong>

        {" - "}

        <small className="text-muted">
            {note.author.role
                .replace("_", " ")
                .replace(/\b\w/g, c => c.toUpperCase())}
        </small>

    </div>

    <small className="text-muted">
        {formatDateTime(note.createdAt)}
    </small>

</div>

<p className="mb-0 mt-2">
    {note.note}
</p>
    </div>
))}


                        </Card.Body>

                    </Card>



                    </>

                }


            </Container>

        </CustomerLayout>

    );

}