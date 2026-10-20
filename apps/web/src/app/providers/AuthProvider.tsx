import { type Dispatch, type ReactNode, type SetStateAction, useEffect, useState } from "react";

import { getMe, login as loginApi } from "../../features/auth/repository";
import type { User } from "../../features/auth/types";
import { AuthContext } from "./AuthContext";

type Props = {
  children: ReactNode;
};

const restoreSession = async (): Promise<User | null> => {
  const savedToken = localStorage.getItem("token");
  if (!savedToken) return null;
  try {
    return await getMe(savedToken);
  } catch {
    localStorage.removeItem("token");
    return null;
  }
};

const useUnauthorizedListener = (
  setToken: Dispatch<SetStateAction<string | null>>,
  setUser: Dispatch<SetStateAction<User | null>>,
) => {
  useEffect(() => {
    const clear = () => {
      localStorage.removeItem("token");
      setToken(null);
      setUser(null);
    };
    window.addEventListener("auth:unauthorized", clear);
    return () => window.removeEventListener("auth:unauthorized", clear);
  }, [setToken, setUser]);
};

export function AuthProvider({ children }: Props) {
  const [token, setToken] = useState<string | null>(localStorage.getItem("token"));
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const isAuthenticated = user !== null;
  useUnauthorizedListener(setToken, setUser);

  useEffect(() => {
    void restoreSession().then((restoredUser) => {
      setUser(restoredUser);
      if (!restoredUser) setToken(null);
      setIsLoading(false);
    });
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
