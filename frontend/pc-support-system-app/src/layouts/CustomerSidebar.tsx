import { Nav } from "react-bootstrap";
import { NavLink } from "react-router-dom";

export default function CustomerSidebar() {
    return (
        <div className="border-end vh-100 bg-light">

            <Nav className="flex-column p-3">

                <Nav.Link
                    as={NavLink}
                    to="/customer/devices"
                >
                    My Devices
                </Nav.Link>

                <Nav.Link
                    as={NavLink}
                    to="/customer/jobs"
                >
                    My Jobs
                </Nav.Link>

                <Nav.Link
                    as={NavLink}
                    to="/customer/profile"
                >
                    Profile
                </Nav.Link>

            </Nav>

        </div>
    );
}