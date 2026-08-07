import { Nav } from "react-bootstrap";
import { NavLink } from "react-router-dom";

import type { MenuItem } from "../config/roleMenus";

type SidebarProps = {

    items: MenuItem[];

};

export default function Sidebar({

    items,

}: SidebarProps) {

    return (

        <div
            className="border-end bg-light"
            style={{
                width: "250px",
                minHeight: "100vh",
            }}
        >

            <Nav
                className="flex-column p-3"
            >

                {items.map((item) => (

                    <Nav.Link
                        key={item.path}
                        as={NavLink}
                        to={item.path!}
                    >

                        {item.label}

                    </Nav.Link>

                ))}

            </Nav>

        </div>

    );

}