import { api } from "../../../app/api";
import type {
  LoginRequest,
  RegisterRequest,
  UserResponse,
} from "../types";

export const authApi = {
  register: async (data: RegisterRequest) => {
    const response = await api.post("/auth/register", data);
    return response.data;
  },

  login: async (data: LoginRequest) => {
    const response = await api.post("/auth/login", data);
    return response.data;
  },

  logout: async () => {
    const response = await api.post("/logout");
    return response.data;
},

  getCurrentUser: async (): Promise<UserResponse> => {
    const response = await api.get("/users/me");
    return response.data;
  },

  updateProfile: async (data: Partial<UserResponse>) => {
    const response = await api.patch(
        "/users/me",
        data
    );

    return response.data;
  },
};
