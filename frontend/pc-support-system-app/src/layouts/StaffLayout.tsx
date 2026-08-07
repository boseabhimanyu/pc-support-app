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
            style={{
                minHeight: "100vh",
                width: "100%",
            }}
        >

            {/* Top bar */}

            <Header />

            {/* Sidebar + page content */}

            <div
                className="d-flex"
                style={{
                    minHeight: "calc(100vh - 56px)",
                }}
            >

                <Sidebar items={menuItems} />

                <main
                style={{
                    flex: 1,
                    minWidth: 0,
                    width: 0,
                    padding: "24px",
                }}
>

                    <Outlet />

                </main>

            </div>

        </div>

    );

}