import { Outlet } from "react-router-dom";
import { useMemo } from "react";

import Header from "../components/Header";
import Sidebar from "../components/Sidebar";

import { useAuth } from "../features/auth/hooks/useAuth";
import { roleMenus } from "../config/roleMenus";

export default function StaffLayout() {

    const { user } = useAuth();

    const menuItems = useMemo(() => {

        if (!user) return [];

        return roleMenus[user.role] ?? [];

    }, [user]);

    return (

         <div
        className="d-flex"
        style={{
            minHeight: "100vh",
            width: "100%",
        }}
    >

        <Sidebar items={menuItems} />

        <div
            className="d-flex flex-column"
            style={{
                flex: 1,
                minWidth: 0,
            }}
        >

            <Header />

            <main
                style={{
                    flex: 1,
                    width: "100%",
                    padding: "24px",
                }}
            >

                <Outlet />

            </main>

        </div>

    </div>

    );
}