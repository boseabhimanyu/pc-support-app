export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  firstName: string;
  lastName: string;
  phone: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export type UserRole =
    | "customer"
    | "receptionist"
    | "technician"
    | "head_technician"
    | "admin"
    | "super_admin";


export interface UserResponse {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    role: UserRole;
    state: string;
}

export interface Device {

    id: string;

    brand: string;

    model: string;

    serialNumber: string;

    category: string;

    type: string;

    notes: string

}