import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { MessageSquare } from "lucide-react";
import { AuthProvider, useAuth } from "./contexts/AuthContext.jsx";
import Login from "./pages/Login";
import ChatDashboard from "./pages/ChatDashBoard";
import VerifyEmail from "./pages/VerifyEmail";

function GlobalLoader() {
  return (
    <div className="fixed inset-0 flex flex-col items-center justify-center gap-5" style={{ background: "#07070f" }}>
      <div className="relative">
        <div className="w-16 h-16 rounded-2xl flex items-center justify-center animate-pulse" style={{
          background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)",
          boxShadow: "0 0 48px rgba(34,211,238,0.25)",
        }}>
          <MessageSquare size={30} className="text-white" strokeWidth={1.8} />
        </div>
        <div className="absolute inset-0 -m-2 rounded-2xl border border-cyan-500/20 animate-ping" style={{ animationDuration: "1.8s" }} />
      </div>
      <p className="text-slate-600 text-sm tracking-widest uppercase font-medium animate-pulse">c‑hat</p>
    </div>
  );
}

function ProtectedRoute({ children }) {
  const { user, isLoading } = useAuth();
  if (isLoading) return <GlobalLoader />;
  if (!user) return <Navigate to="/login" replace />;
  return children;
}

function PublicRoute({ children }) {
  const { user, isLoading } = useAuth();
  if (isLoading) return <GlobalLoader />;
  if (user) return <Navigate to="/" replace />;
  return children;
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<PublicRoute><Login /></PublicRoute>} />
      <Route path="/" element={<ProtectedRoute><ChatDashboard /></ProtectedRoute>} />
      <Route path="/verify-email" element={<VerifyEmail />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}