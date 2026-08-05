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

export interface UserResponse {
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    role: string;
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