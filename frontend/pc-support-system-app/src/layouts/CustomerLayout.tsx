import { Container, Row, Col } from "react-bootstrap";
import TopBar from "../layouts/TopNavbar";
import CustomerSidebar from "./CustomerSidebar";
interface Props { children: React.ReactNode }
export default function CustomerLayout({ children }: Props) {
  return <><TopBar /><Container fluid className="customer-shell"><Row><Col lg={2} md={3} className="customer-sidebar-column"><CustomerSidebar /></Col><Col lg={10} md={9} className="customer-content">{children}</Col></Row></Container></>;
}
