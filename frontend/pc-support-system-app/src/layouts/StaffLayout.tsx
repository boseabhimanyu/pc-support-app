import { Outlet } from "react-router";
import { useMemo } from "react";
import Header from "../components/Header";
import Sidebar from "../components/Sidebar";
import { useAuth } from "../features/auth/hooks/useAuth";
import { roleMenus } from "../config/roleMenus";

export default function StaffLayout() {
  const { user } = useAuth();
  const menuItems = useMemo(() => user ? roleMenus[user.role] ?? [] : [], [user]);
  return <div className="app-shell"><Header /><div className="app-body"><Sidebar items={menuItems} /><main className="app-main"><Outlet /></main></div></div>;
}
