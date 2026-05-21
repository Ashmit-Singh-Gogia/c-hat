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
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  // OAuth — hard redirect through Vite proxy
  const login = (provider) => {
    window.location.href = `/api/auth/oauth/${provider}`;
  };

  // Logout — axios clears cookie, then we navigate manually
  // (backend returns 200 JSON, not a redirect, so window.location alone won't work)
  const logout = async () => {
    try { await api.get("/auth/logout"); } catch {}
    window.location.href = "/login";
  };

  const localLogin = async (email, password) => {
    try {
      await api.post("/auth/local/login", { email, password });
      await checkAuth();
      return { success: true };
    } catch (err) {
      const errorMsg = err.response?.data?.error;
      if (errorMsg === "unverified_account") return { success: false, type: "unverified" };
      return { success: false, type: "invalid_credentials" };
    }
  };

  const localRegister = async (username, email, password) => {
    try {
      await api.post("/auth/local/register", { username, email, password });
      return { success: true };
    } catch (err) {
      return { success: false, message: err.response?.data?.error || "Registration failed. Please try again." };
    }
  };

  const resendVerification = async (email) => {
    try {
      await api.post("/auth/local/resend-verification", { email });
      return { success: true };
    } catch (err) {
      return { success: false, message: err.response?.data?.error || "Failed to resend. Please try again." };
    }
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, login, logout, localLogin, localRegister, resendVerification, refetchUser: checkAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}