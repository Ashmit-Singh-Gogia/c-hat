import { useEffect, useState } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { MessageSquare, AlertCircle, X, Zap } from "lucide-react";
import { useAuth } from "../contexts/AuthContext.jsx";

const ERROR_MESSAGES = {
  auth_failed: "Authentication failed. Please try again.",
  access_denied: "Access was denied. Did you cancel the sign-in?",
  server_error: "A server error occurred. Please try again shortly.",
  session_expired: "Your session expired. Please sign in again.",
};

function ErrorToast({ message, onDismiss }) {
  return (
    <div className="flex items-start gap-3 bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 backdrop-blur-sm animate-fade-in">
      <AlertCircle size={18} className="mt-0.5 shrink-0 text-red-400" />
      <p className="text-sm leading-snug flex-1">{message}</p>
      <button
        onClick={onDismiss}
        className="text-red-400/60 hover:text-red-300 transition-colors shrink-0"
        aria-label="Dismiss error"
      >
        <X size={14} />
      </button>
    </div>
  );
}

export default function Login() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { user, isLoading, login } = useAuth();
  const [error, setError] = useState(null);
  const [hovering, setHovering] = useState(false);

  useEffect(() => {
    // If already authenticated, bounce to dashboard
    if (!isLoading && user) navigate("/", { replace: true });
  }, [user, isLoading, navigate]);

  useEffect(() => {
    const code = searchParams.get("error");
    if (code) {
      setError(ERROR_MESSAGES[code] ?? "An unexpected error occurred.");
      // Clean the URL so a refresh doesn't re-display the error
      window.history.replaceState({}, "", window.location.pathname);
    }
  }, [searchParams]);

  return (
    <div className="relative min-h-screen w-full flex items-center justify-center overflow-hidden bg-[#07070f]">

      {/* ── Ambient background glow ── */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          background: `
            radial-gradient(ellipse 70% 55% at 50% -10%, rgba(56,189,248,0.12) 0%, transparent 60%),
            radial-gradient(ellipse 50% 40% at 80% 80%,  rgba(139,92,246,0.10) 0%, transparent 55%),
            radial-gradient(ellipse 40% 30% at 10% 90%,  rgba(34,211,238,0.07) 0%, transparent 50%)
          `,
        }}
      />

      {/* ── Subtle grid overlay ── */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 opacity-[0.025]"
        style={{
          backgroundImage: `linear-gradient(rgba(148,163,184,1) 1px, transparent 1px),
                            linear-gradient(90deg, rgba(148,163,184,1) 1px, transparent 1px)`,
          backgroundSize: "48px 48px",
        }}
      />

      {/* ── Card ── */}
      <div className="relative z-10 w-full max-w-sm mx-4">
        <div
          className="rounded-2xl p-8 flex flex-col gap-7"
          style={{
            background: "rgba(255,255,255,0.04)",
            backdropFilter: "blur(24px)",
            WebkitBackdropFilter: "blur(24px)",
            border: "1px solid rgba(255,255,255,0.08)",
            boxShadow: "0 0 0 1px rgba(0,0,0,0.3), 0 32px 64px rgba(0,0,0,0.5), inset 0 1px 0 rgba(255,255,255,0.07)",
          }}
        >
          {/* Logo */}
          <div className="flex flex-col items-center gap-4 text-center">
            <div className="relative">
              <div
                className="w-14 h-14 rounded-2xl flex items-center justify-center"
                style={{
                  background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)",
                  boxShadow: "0 0 32px rgba(34,211,238,0.3), 0 8px 24px rgba(0,0,0,0.4)",
                }}
              >
                <MessageSquare size={26} className="text-white" strokeWidth={1.8} />
              </div>
              <span
                className="absolute -top-1 -right-1 w-4 h-4 rounded-full flex items-center justify-center"
                style={{ background: "linear-gradient(135deg, #f59e0b, #ef4444)" }}
              >
                <Zap size={9} className="text-white" fill="white" />
              </span>
            </div>

            <div>
              <h1
                className="text-2xl font-bold tracking-tight text-white"
                style={{ fontFamily: "'Syne', 'DM Sans', sans-serif", letterSpacing: "-0.02em" }}
              >
                c‑hat
              </h1>
              <p className="mt-1 text-sm text-slate-400">
                Real-time conversations, beautifully simple.
              </p>
            </div>
          </div>

          {/* Error toast */}
          {error && (
            <ErrorToast message={error} onDismiss={() => setError(null)} />
          )}

          {/* Divider */}
          <div className="flex items-center gap-3">
            <div className="flex-1 h-px bg-white/[0.07]" />
            <span className="text-xs text-slate-600 font-medium uppercase tracking-widest">sign in</span>
            <div className="flex-1 h-px bg-white/[0.07]" />
          </div>

          {/* Google Button */}
          <button
            onClick={() => login("google")}
            onMouseEnter={() => setHovering(true)}
            onMouseLeave={() => setHovering(false)}
            disabled={isLoading}
            className="group relative w-full flex items-center justify-center gap-3 py-3.5 rounded-xl font-medium text-sm transition-all duration-200 select-none disabled:opacity-50 disabled:cursor-not-allowed"
            style={{
              background: hovering
                ? "rgba(255,255,255,0.1)"
                : "rgba(255,255,255,0.06)",
              border: hovering
                ? "1px solid rgba(255,255,255,0.2)"
                : "1px solid rgba(255,255,255,0.1)",
              color: "#e2e8f0",
              boxShadow: hovering ? "0 4px 24px rgba(34,211,238,0.08)" : "none",
              transform: hovering ? "translateY(-1px)" : "translateY(0)",
            }}
          >
            {/* Google SVG icon */}
            <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden>
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" />
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
            </svg>
            Continue with Google
          </button>

          {/* Footer note */}
          <p className="text-center text-xs text-slate-600 leading-relaxed">
            By signing in you agree to our{" "}
            <a href="#" className="text-slate-500 hover:text-cyan-400 transition-colors underline underline-offset-2">
              Terms
            </a>{" "}
            and{" "}
            <a href="#" className="text-slate-500 hover:text-cyan-400 transition-colors underline underline-offset-2">
              Privacy Policy
            </a>.
          </p>
        </div>

        {/* Below-card glow line */}
        <div
          aria-hidden
          className="absolute -bottom-px left-1/2 -translate-x-1/2 h-px w-2/3"
          style={{ background: "linear-gradient(90deg, transparent, rgba(34,211,238,0.4), transparent)" }}
        />
      </div>

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Syne:wght@700;800&family=DM+Sans:wght@400;500&display=swap');
        @keyframes fade-in { from { opacity:0; transform:translateY(-6px); } to { opacity:1; transform:translateY(0); } }
        .animate-fade-in { animation: fade-in 0.25s ease forwards; }
      `}</style>
    </div>
  );
}