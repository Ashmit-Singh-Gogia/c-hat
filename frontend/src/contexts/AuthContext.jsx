import { createContext, useContext, useState, useEffect, useCallback } from "react";
import axios from "axios";

// Axios instance with HttpOnly cookie support baked in
export const api = axios.create({
  withCredentials: true,
});

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = useCallback(async () => {
    try {
      const { data } = await api.get("/api/users/me");
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
    window.location.href = `http://localhost:8082/api/auth/${provider}`;
  };

  const logout = (provider) => {
    window.location.href = `http://localhost:8082/api/auth/${provider}/logout`;
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