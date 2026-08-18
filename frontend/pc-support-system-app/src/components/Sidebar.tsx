import { Nav } from "react-bootstrap";
import { NavLink } from "react-router";
import type { MenuItem } from "../config/roleMenus";

type SidebarProps = { items: MenuItem[] };
export default function Sidebar({ items }: SidebarProps) {
  return <aside className="app-sidebar"><Nav className="flex-column"><span className="sidebar-label">Workspace</span>{items.map((item) => <Nav.Link key={item.path} as={NavLink} to={item.path!}>{item.label}</Nav.Link>)}</Nav></aside>;
}
