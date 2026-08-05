import { Container, Row, Col } from "react-bootstrap";

import TopBar from "../layouts/TopNavbar";
import CustomerSidebar from "./CustomerSidebar";

interface Props {
    children: React.ReactNode;
}

export default function CustomerLayout({
    children,
}: Props) {

    return (

        <>

            <TopBar />

            <Container fluid>

                <Row>

                    <Col md={2}>

                        <CustomerSidebar />

                    </Col>

                    <Col md={10}>

                        {children}

                    </Col>

                </Row>

            </Container>

        </>

    );

}