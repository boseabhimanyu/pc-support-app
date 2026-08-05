import { Badge } from "react-bootstrap";


type BadgeVariant =
    | "primary"
    | "warning"
    | "info"
    | "success"
    | "secondary";



function getVariant(
    status: string
): BadgeVariant {

    switch (status.toLowerCase()) {

        case "created":
            return "primary";

        case "assigned":
            return "warning";

        case "in_progress":
            return "info";

        case "closed":
            return "success";

        case "cancelled":
            return "secondary";

        default:
            return "secondary";
    }

}



function formatStatus(
    status: string
) {

    return status
        .replace("_", " ")
        .replace(
            /\b\w/g,
            char => char.toUpperCase()
        );

}



interface Props {

    status: string;

}



export default function JobStatusBadge({
    status,
}: Props) {


    return (

       <Badge
    bg={getVariant(status)}
    className="px-3 py-2 d-inline-flex align-items-center justify-content-center"
>
    {formatStatus(status)}
</Badge>

    );

}