import { Navigate } from "react-router-dom";
import { useAuth } from "../features/auth/hooks/useAuth";

interface Props {
    children: React.ReactNode;
}

export default function ProtectedRoute({
    children,
}: Props) {

    const {
        user,
        loading,
    } = useAuth();



    if (loading) {
        return <div>Loading...</div>;
    }


    if (!user) {
        return <Navigate to="/" replace />;
    }


    return <>{children}</>;
}