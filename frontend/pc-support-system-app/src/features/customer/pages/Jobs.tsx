import { useEffect, useState } from "react";

import {
    Alert,
    Container,
    Spinner,
} from "react-bootstrap";

import CustomerLayout from "../../../layouts/CustomerLayout";

import JobCard from "../../jobs/components/JobCard";

import {
    jobApi
} from "../../jobs/services/jobApi";

import type {
    Job
} from "../../jobs/types/job";


export default function Jobs() {


    const [jobs, setJobs] =
        useState<Job[]>([]);


    const [jobsCount, setJobsCount] =
        useState(0);


    const [loading, setLoading] =
        useState(true);


    const [error, setError] =
        useState("");



    useEffect(() => {

        getJobs();

    }, []);



    async function getJobs() {

        try {

            setLoading(true);

            setError("");


            const data =
                await jobApi.getMyJobs();


            console.log(
                "Jobs:",
                data
            );


            setJobs(
                data.jobs
            );


            setJobsCount(
                data.jobsCount
            );


        } catch (err) {

            console.log(err);

            setError(
                "Unable to load jobs."
            );


        } finally {

            setLoading(false);

        }

    }



    return (

        <CustomerLayout>

            <Container className="py-4">


                <div className="d-flex justify-content-between align-items-center mb-4">

                    <h3>
                        My Jobs
                    </h3>


                    <h5>
                        Total Jobs:
                        {" "}
                        {jobsCount}
                    </h5>


                </div>



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
                    !loading && jobs.length === 0 &&

                    <Alert variant="info">

                        No jobs found.

                    </Alert>

                }



                {
                    jobs.map(job => (

                        <JobCard
                            key={job.id}
                            job={job}
                        />

                    ))
                }



            </Container>


        </CustomerLayout>

    );

}