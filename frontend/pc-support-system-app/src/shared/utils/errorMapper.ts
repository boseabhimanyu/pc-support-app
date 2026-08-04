export function mapBackendError(error: string): string {
  if (error.includes("Phone"))
    return "Phone number must contain exactly 10 digits.";

  if (error.includes("Password"))
    return "Password must be at least 8 characters.";

  if (error.includes("Email"))
    return "Please enter a valid email address.";

  if (error.includes("name contains invalid characters"))
    return "Name contains invalid characters.";

  return error;
}
