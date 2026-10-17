import React from "react";
import {
  type ReactNode,
  useState,
  useEffect,
} from "react";
import type { User } from "../../features/auth/types";
import { getMe, login as loginApi } from "../../features/auth/repository"
import { AuthContext } from "./AuthContext";

type Props = {
  children: ReactNode;
};

export function AuthProvider({ children }: Props) {
  const [token, setToken] = useState<string | null>(
    localStorage.getItem("token"),
  );
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const isAuthenticated = user !== null;

  useEffect(() => {
    const restoreUser = async () => {
      const savedToken = localStorage.getItem("token");
      if (savedToken === null) {
        setIsLoading(false);
        return;
      }

      try {
        const restoredUser = await getMe(savedToken);
        setUser(restoredUser);
      } catch {
        localStorage.removeItem("token");
        setToken(null);
      } finally {
        setIsLoading(false);
      }
    };

    restoreUser();
  }, []);

  const login = async (email: string, password: string) => {
    const result = await loginApi({
      email,
      password,
    });

    localStorage.setItem("token", result.token);

    setToken(result.token);
    setUser(result.user);
  };

  const logout = () => {
    localStorage.removeItem("token");

    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        token,
        user,
        isLoading,
        isAuthenticated,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

