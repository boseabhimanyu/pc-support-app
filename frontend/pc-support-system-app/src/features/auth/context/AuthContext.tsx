import {
  createContext,
  useState,
  useEffect,
} from "react";

import type { ReactNode } from "react";
import type { UserResponse } from "../types";

import { authApi } from "../services/authApi";

export interface AuthContextType {
  user: UserResponse | null;

  loading: boolean;

  setUser: (user: UserResponse | null) => void;

  logout: () => Promise<void>;

  refreshUser: () => Promise<void>;
}

export const AuthContext =
  createContext<AuthContextType | null>(null);

interface Props {
  children: ReactNode;
}

export function AuthProvider({
  children,
}: Props) {

  const [user, setUserState] =
    useState<UserResponse | null>(null);

  const [loading, setLoading] =
    useState(true);

    useEffect(() => {
    refreshUser();
    }, []);

  function setUser(user: UserResponse | null) {
    setUserState(user);
  }

 async function refreshUser() {
    try {

        setLoading(true);

        const currentUser =
            await authApi.getCurrentUser();

        setUser(currentUser);

    } catch {

        setUser(null);

    } finally {

        setLoading(false);

    }
}

  async function logout() {
    try {
        await authApi.logout();
    } finally {
        setUser(null);
    }
}

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        setUser,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
    
  );
}