// export function mapBackendError(error: string): string {

//     if (error.includes("Phone"))
//         return "Phone Number must contain exactly 10 digits.";

//     if (error.includes("Password"))
//         return "Password must be at least 8 characters.";

//     if (error.includes("Email"))
//         return "Please enter a valid email address.";

//     if (error.includes("First Name"))
//         return "First Name contains invalid characters.";

//     if (error.includes("Last Name"))
//         return "Last Name contains invalid characters.";

//     return error;

// }

export function mapBackendError(error: string): string {
    return error;
}