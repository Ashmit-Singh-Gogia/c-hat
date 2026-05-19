import { createContext, useContext, useState, useEffect, useCallback } from "react";
import { apiClient as api } from "../api/client";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = useCallback(async () => {
    try {
      const { data } = await api.get("/users/me");
      setUser(data);
    } catch {
      // 401 → not authenticated; any other error → treat as unauthenticated
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  // Hard redirect — browser carries the cookie; no JS token handling needed
  const login = (provider) => {
    // Navigate to the backend auth route via the dev-server proxy
    window.location.href = `/api/auth/${provider}`;
  };

  const logout = (provider) => {
    window.location.href = `/api/auth/${provider}/logout`;
  };
  return (
    <AuthContext.Provider value={{ user, isLoading, login, logout, refetchUser: checkAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}