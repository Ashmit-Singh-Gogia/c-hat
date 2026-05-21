import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { MessageSquare, Zap, CheckCircle, XCircle, Loader2 } from "lucide-react";
import { apiClient as api } from "../api/client";

export default function VerifyEmail() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState("loading"); // 'loading' | 'success' | 'error'
  const [errorMsg, setErrorMsg] = useState("");

  useEffect(() => {
    const token = searchParams.get("token");
    if (!token) {
      setStatus("error");
      setErrorMsg("No verification token found in the link.");
      return;
    }
    api.get(`/auth/local/verify-email?token=${token}`)
      .then(() => setStatus("success"))
      .catch(err => {
        setStatus("error");
        setErrorMsg(err.response?.data?.error || "Invalid or expired verification link.");
      });
  }, []);

  const cardStyle = {
    background: "rgba(255,255,255,0.04)",
    backdropFilter: "blur(24px)",
    WebkitBackdropFilter: "blur(24px)",
    border: "1px solid rgba(255,255,255,0.08)",
    boxShadow: "0 0 0 1px rgba(0,0,0,0.3), 0 32px 64px rgba(0,0,0,0.5), inset 0 1px 0 rgba(255,255,255,0.07)",
  };

  return (
    <div className="relative min-h-screen w-full flex items-center justify-center overflow-hidden bg-[#07070f]">
      <div aria-hidden className="pointer-events-none absolute inset-0" style={{
        background: `
          radial-gradient(ellipse 70% 55% at 50% -10%, rgba(56,189,248,0.12) 0%, transparent 60%),
          radial-gradient(ellipse 50% 40% at 80% 80%, rgba(139,92,246,0.10) 0%, transparent 55%),
          radial-gradient(ellipse 40% 30% at 10% 90%, rgba(34,211,238,0.07) 0%, transparent 50%)
        `,
      }} />
      <div aria-hidden className="pointer-events-none absolute inset-0 opacity-[0.025]" style={{
        backgroundImage: `linear-gradient(rgba(148,163,184,1) 1px, transparent 1px), linear-gradient(90deg, rgba(148,163,184,1) 1px, transparent 1px)`,
        backgroundSize: "48px 48px",
      }} />

      <div className="relative z-10 w-full max-w-sm mx-4">
        <div className="rounded-2xl p-8 flex flex-col gap-6 items-center text-center" style={cardStyle}>
          {/* Logo */}
          <div className="relative">
            <div className="w-14 h-14 rounded-2xl flex items-center justify-center" style={{
              background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)",
              boxShadow: "0 0 32px rgba(34,211,238,0.3), 0 8px 24px rgba(0,0,0,0.4)",
            }}>
              <MessageSquare size={26} className="text-white" strokeWidth={1.8} />
            </div>
            <span className="absolute -top-1 -right-1 w-4 h-4 rounded-full flex items-center justify-center" style={{ background: "linear-gradient(135deg, #f59e0b, #ef4444)" }}>
              <Zap size={9} className="text-white" fill="white" />
            </span>
          </div>

          {status === "loading" && (
            <>
              <Loader2 size={32} className="text-cyan-400 animate-spin" />
              <p className="text-sm text-slate-400">Verifying your email…</p>
            </>
          )}

          {status === "success" && (
            <>
              <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ background: "rgba(16,185,129,0.15)", border: "1px solid rgba(16,185,129,0.3)" }}>
                <CheckCircle size={24} className="text-emerald-400" />
              </div>
              <div>
                <p className="text-base font-semibold text-slate-200">Email Verified!</p>
                <p className="text-sm text-slate-500 mt-1">Your account is now active. You can sign in.</p>
              </div>
              <button onClick={() => navigate("/login")} className="w-full py-3 rounded-xl text-sm font-medium transition-all duration-200"
                style={{ background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)", color: "white", boxShadow: "0 4px 16px rgba(34,211,238,0.2)" }}>
                Go to Sign In
              </button>
            </>
          )}

          {status === "error" && (
            <>
              <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ background: "rgba(239,68,68,0.12)", border: "1px solid rgba(239,68,68,0.25)" }}>
                <XCircle size={24} className="text-red-400" />
              </div>
              <div>
                <p className="text-base font-semibold text-slate-200">Verification Failed</p>
                <p className="text-sm text-slate-500 mt-1">{errorMsg}</p>
              </div>
              <button onClick={() => navigate("/login")} className="w-full py-3 rounded-xl text-sm font-medium transition-all duration-200"
                style={{ background: "rgba(255,255,255,0.06)", border: "1px solid rgba(255,255,255,0.1)", color: "#e2e8f0" }}>
                Back to Sign In
              </button>
            </>
          )}
        </div>

        <div aria-hidden className="absolute -bottom-px left-1/2 -translate-x-1/2 h-px w-2/3" style={{ background: "linear-gradient(90deg, transparent, rgba(34,211,238,0.4), transparent)" }} />
      </div>

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Syne:wght@700;800&family=DM+Sans:wght@400;500&display=swap');
      `}</style>
    </div>
  );
}