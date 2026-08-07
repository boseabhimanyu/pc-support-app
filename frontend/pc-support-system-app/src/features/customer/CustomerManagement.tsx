import { useState } from "react";
import { useNavigate } from "react-router-dom";

import {
    Alert,
    Button,
    Card,
    Form,
    InputGroup,
    ListGroup,
    Spinner,
} from "react-bootstrap";

import {
    searchCustomers,
} from "./services/customerApi";

import type {
    CustomerSummary,
} from "./services/customerApi";

export default function CustomerManagement() {
    const navigate = useNavigate();

    const [query, setQuery] = useState("");

    const [customers, setCustomers] =
        useState<CustomerSummary[]>([]);

    const [loading, setLoading] =
        useState(false);

    const [error, setError] =
        useState("");

    async function handleSearch() {

        const searchText = query.trim();

        if (searchText.length < 3) {

            setCustomers([]);

            setError(
                "Enter at least 3 characters to search.",
            );

            return;
        }

        try {

            setLoading(true);
            setError("");

            const results =
                await searchCustomers(searchText);

            setCustomers(results);

            if (results.length === 0) {

                setError(
                    "No customers found.",
                );

            }

        } catch (err) {

            console.error(err);

            setCustomers([]);

            setError(
                "Unable to search customers.",
            );

        } finally {

            setLoading(false);

        }
    }

    return (

        <>

            <h2 className="mb-4">
                Customer Management
            </h2>
                <Button
    className="mb-3"
    onClick={() => navigate("create")}
>
    Create Customer
</Button>
            <Card>

                <Card.Body>

                    <Card.Title className="mb-3">
                        Search Customer
                    </Card.Title>

                    <InputGroup className="mb-3">

                        <Form.Control
                            type="text"
                            placeholder="Enter customer name or phone"
                            value={query}
                            onChange={(e) =>
                                setQuery(e.target.value)
                            }
                            onKeyDown={(e) => {

                                if (e.key === "Enter") {
                                    handleSearch();
                                }

                            }}
                        />

                        <Button
                            onClick={handleSearch}
                            disabled={loading}
                        >

                            {loading ? (
                                <>
                                    <Spinner
                                        size="sm"
                                        className="me-2"
                                    />

                                    Searching...
                                </>
                            ) : (
                                "Search"
                            )}

                        </Button>

                    </InputGroup>

                    {error && (
                        <Alert variant="warning">
                            {error}
                        </Alert>
                    )}

                    {customers.length > 0 && (

                        <ListGroup>

                            {customers.map((customer) => (

                                <ListGroup.Item
                                    key={customer.id}
                                    action
                                    onClick={() =>
                                        navigate(
                                            `/receptionist/customers/${customer.id}`,
                                        )
                                    }
                                >

                                    <div>
                                        <strong>
                                            {customer.firstName}{" "}
                                            {customer.lastName}
                                        </strong>
                                    </div>

                                    <div className="text-muted">

                                        Phone:{" "}
                                        {customer.phone}

                                        {customer.email && (
                                            <>
                                                {" • "}
                                                Email:{" "}
                                                {customer.email}
                                            </>
                                        )}

                                    </div>

                                </ListGroup.Item>

                            ))}

                        </ListGroup>

                    )}

                </Card.Body>

            </Card>

        </>

    );
}