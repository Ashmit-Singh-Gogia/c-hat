import { useEffect, useState, useRef } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { MessageSquare, AlertCircle, X, Zap, Mail, ArrowLeft, Loader2 } from "lucide-react";
import { useAuth } from "../contexts/AuthContext.jsx";

const OAUTH_ERRORS = {
  auth_failed: "Authentication failed. Please try again.",
  access_denied: "Access was denied. Did you cancel the sign-in?",
  server_error: "A server error occurred. Please try again shortly.",
  session_expired: "Your session expired. Please sign in again.",
  invalid_provider: "Invalid provider. Please try again.",
  user_creation_failed: "Failed to create account. Please try again.",
  token_generation_failed: "A server error occurred. Please try again.",
};

function ErrorToast({ message, onDismiss }) {
  return (
    <div className="flex items-start gap-3 bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 backdrop-blur-sm animate-fade-in">
      <AlertCircle size={18} className="mt-0.5 shrink-0 text-red-400" />
      <p className="text-sm leading-snug flex-1">{message}</p>
      <button onClick={onDismiss} className="text-red-400/60 hover:text-red-300 transition-colors shrink-0">
        <X size={14} />
      </button>
    </div>
  );
}

function Field({ label, type = "text", value, onChange, placeholder, disabled }) {
  const [focused, setFocused] = useState(false);
  return (
    <div className="flex flex-col gap-1.5">
      <label className="text-xs text-slate-500 font-medium uppercase tracking-wider">{label}</label>
      <input
        type={type}
        value={value}
        onChange={onChange}
        placeholder={placeholder}
        disabled={disabled}
        className="w-full rounded-xl px-4 py-3 text-sm text-slate-200 placeholder-slate-600 outline-none transition-all duration-200 disabled:opacity-40"
        style={{
          background: "rgba(255,255,255,0.05)",
          border: focused ? "1px solid rgba(34,211,238,0.4)" : "1px solid rgba(255,255,255,0.08)",
          boxShadow: focused ? "0 0 0 3px rgba(34,211,238,0.06)" : "none",
        }}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
      />
    </div>
  );
}

export default function Login() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { user, isLoading, login, localLogin, localRegister, resendVerification } = useAuth();

  // view: 'signin' | 'signup' | 'verify_required'
  const [view, setView] = useState("signin");
  const [error, setError] = useState(null);
  const [successMsg, setSuccessMsg] = useState(null);
  const [submitting, setSubmitting] = useState(false);
  const [hovering, setHovering] = useState(false);

  // Sign in
  const [siEmail, setSiEmail] = useState("");
  const [siPassword, setSiPassword] = useState("");

  // Sign up
  const [suUsername, setSuUsername] = useState("");
  const [suEmail, setSuEmail] = useState("");
  const [suPassword, setSuPassword] = useState("");

  // Verify required
  const [verifyEmail, setVerifyEmail] = useState("");
  const [countdown, setCountdown] = useState(0);
  const timerRef = useRef(null);

  useEffect(() => {
    if (!isLoading && user) navigate("/", { replace: true });
  }, [user, isLoading, navigate]);

  useEffect(() => {
    const code = searchParams.get("error");
    if (code) {
      setError(OAUTH_ERRORS[code] ?? "An unexpected error occurred.");
      window.history.replaceState({}, "", window.location.pathname);
    }
  }, [searchParams]);

  useEffect(() => () => clearInterval(timerRef.current), []);

  const startCountdown = () => {
    setCountdown(60);
    timerRef.current = setInterval(() => {
      setCountdown(prev => {
        if (prev <= 1) { clearInterval(timerRef.current); return 0; }
        return prev - 1;
      });
    }, 1000);
  };

  const switchView = (v) => { setError(null); setSuccessMsg(null); setView(v); };

  const handleSignIn = async () => {
    if (!siEmail.trim() || !siPassword.trim()) { setError("Please fill in all fields."); return; }
    setError(null); setSuccessMsg(null); setSubmitting(true);
    const res = await localLogin(siEmail, siPassword);
    setSubmitting(false);
    if (res.success) {
      navigate("/", { replace: true });
    } else if (res.type === "unverified") {
      setVerifyEmail(siEmail);
      setView("verify_required");
    } else {
      setError("Invalid email or password.");
    }
  };

  const handleSignUp = async () => {
    if (!suUsername.trim() || !suEmail.trim() || !suPassword.trim()) { setError("Please fill in all fields."); return; }
    setError(null); setSuccessMsg(null); setSubmitting(true);
    const res = await localRegister(suUsername, suEmail, suPassword);
    setSubmitting(false);
    if (res.success) {
      setSuccessMsg("Account created! Check your email to verify, then sign in.");
      setSiEmail(suEmail);
      setSuUsername(""); setSuEmail(""); setSuPassword("");
      setView("signin");
    } else {
      setError(res.message);
    }
  };

  const handleResend = async () => {
    if (countdown > 0) return;
    setError(null);
    const res = await resendVerification(verifyEmail);
    if (res.success) {
      startCountdown();
      setSuccessMsg("Verification email sent! Check your inbox.");
    } else {
      setError(res.message);
    }
  };

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
        <div className="rounded-2xl p-8 flex flex-col gap-6" style={cardStyle}>

          {/* Logo — unchanged */}
          <div className="flex flex-col items-center gap-4 text-center">
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
            <div>
              <h1 className="text-2xl font-bold tracking-tight text-white" style={{ fontFamily: "'Syne', 'DM Sans', sans-serif", letterSpacing: "-0.02em" }}>
                c‑hat
              </h1>
              <p className="mt-1 text-sm text-slate-400">Real-time conversations, beautifully simple.</p>
            </div>
          </div>

          {/* Toasts */}
          {error && <ErrorToast message={error} onDismiss={() => setError(null)} />}
          {successMsg && (
            <div className="bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 rounded-xl px-4 py-3 text-sm animate-fade-in">
              {successMsg}
            </div>
          )}

          {/* VERIFY REQUIRED */}
          {view === "verify_required" && (
            <div className="flex flex-col gap-5">
              <div className="flex flex-col items-center gap-3 text-center">
                <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ background: "rgba(234,179,8,0.15)", border: "1px solid rgba(234,179,8,0.3)" }}>
                  <Mail size={22} className="text-yellow-400" />
                </div>
                <div>
                  <p className="text-sm font-semibold text-slate-200">Verify your email</p>
                  <p className="text-xs text-slate-500 mt-1">
                    We sent a link to <span className="text-cyan-400">{verifyEmail}</span>
                  </p>
                </div>
              </div>
              <button
                onClick={handleResend}
                disabled={countdown > 0}
                className="w-full py-3 rounded-xl text-sm font-medium transition-all duration-200 disabled:cursor-not-allowed"
                style={{
                  background: countdown > 0 ? "rgba(255,255,255,0.03)" : "rgba(34,211,238,0.1)",
                  border: countdown > 0 ? "1px solid rgba(255,255,255,0.06)" : "1px solid rgba(34,211,238,0.25)",
                  color: countdown > 0 ? "#475569" : "#22d3ee",
                }}
              >
                {countdown > 0 ? `Resend available in ${countdown}s` : "Resend Verification Email"}
              </button>
              <button onClick={() => switchView("signin")} className="flex items-center justify-center gap-2 text-xs text-slate-500 hover:text-slate-300 transition-colors">
                <ArrowLeft size={12} /> Back to Sign In
              </button>
            </div>
          )}

          {/* TABS + FORMS */}
          {view !== "verify_required" && (
            <>
              {/* Tab switcher */}
              <div className="flex rounded-xl p-1" style={{ background: "rgba(255,255,255,0.04)", border: "1px solid rgba(255,255,255,0.06)" }}>
                {["signin", "signup"].map(tab => (
                  <button key={tab} onClick={() => switchView(tab)}
                    className="flex-1 py-2 rounded-lg text-xs font-medium transition-all duration-200"
                    style={{
                      background: view === tab ? "rgba(255,255,255,0.09)" : "transparent",
                      color: view === tab ? "#e2e8f0" : "#475569",
                      boxShadow: view === tab ? "0 1px 4px rgba(0,0,0,0.3)" : "none",
                    }}
                  >
                    {tab === "signin" ? "Sign In" : "Sign Up"}
                  </button>
                ))}
              </div>

              {/* SIGN IN */}
              {view === "signin" && (
                <div className="flex flex-col gap-4">
                  <Field label="Email" type="email" value={siEmail} onChange={e => setSiEmail(e.target.value)} placeholder="you@example.com" disabled={submitting} />
                  <Field label="Password" type="password" value={siPassword} onChange={e => setSiPassword(e.target.value)} placeholder="••••••••" disabled={submitting} />
                  <button
                    onClick={handleSignIn}
                    disabled={submitting || !siEmail.trim() || !siPassword.trim()}
                    className="w-full py-3 rounded-xl text-sm font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    style={{ background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)", color: "white", boxShadow: "0 4px 16px rgba(34,211,238,0.2)" }}
                  >
                    {submitting && <Loader2 size={15} className="animate-spin" />}
                    {submitting ? "Signing in…" : "Sign In"}
                  </button>
                  <div className="flex items-center gap-3">
                    <div className="flex-1 h-px bg-white/[0.07]" />
                    <span className="text-xs text-slate-600 font-medium">or</span>
                    <div className="flex-1 h-px bg-white/[0.07]" />
                  </div>
                  {/* Google button — identical to original */}
                  <button
                    onClick={() => login("google")}
                    onMouseEnter={() => setHovering(true)}
                    onMouseLeave={() => setHovering(false)}
                    disabled={isLoading || submitting}
                    className="group relative w-full flex items-center justify-center gap-3 py-3.5 rounded-xl font-medium text-sm transition-all duration-200 select-none disabled:opacity-50 disabled:cursor-not-allowed"
                    style={{
                      background: hovering ? "rgba(255,255,255,0.1)" : "rgba(255,255,255,0.06)",
                      border: hovering ? "1px solid rgba(255,255,255,0.2)" : "1px solid rgba(255,255,255,0.1)",
                      color: "#e2e8f0",
                      boxShadow: hovering ? "0 4px 24px rgba(34,211,238,0.08)" : "none",
                      transform: hovering ? "translateY(-1px)" : "translateY(0)",
                    }}
                  >
                    <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden>
                      <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                      <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                      <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" />
                      <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
                    </svg>
                    Continue with Google
                  </button>
                </div>
              )}

              {/* SIGN UP */}
              {view === "signup" && (
                <div className="flex flex-col gap-4">
                  <Field label="Username" value={suUsername} onChange={e => setSuUsername(e.target.value)} placeholder="yourname" disabled={submitting} />
                  <Field label="Email" type="email" value={suEmail} onChange={e => setSuEmail(e.target.value)} placeholder="you@example.com" disabled={submitting} />
                  <Field label="Password" type="password" value={suPassword} onChange={e => setSuPassword(e.target.value)} placeholder="••••••••" disabled={submitting} />
                  <button
                    onClick={handleSignUp}
                    disabled={submitting || !suUsername.trim() || !suEmail.trim() || !suPassword.trim()}
                    className="w-full py-3 rounded-xl text-sm font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    style={{ background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)", color: "white", boxShadow: "0 4px 16px rgba(34,211,238,0.2)" }}
                  >
                    {submitting && <Loader2 size={15} className="animate-spin" />}
                    {submitting ? "Creating account…" : "Create Account"}
                  </button>
                </div>
              )}
            </>
          )}

          {/* Footer */}
          <p className="text-center text-xs text-slate-600 leading-relaxed">
            By signing in you agree to our{" "}
            <a href="#" className="text-slate-500 hover:text-cyan-400 transition-colors underline underline-offset-2">Terms</a>{" "}
            and{" "}
            <a href="#" className="text-slate-500 hover:text-cyan-400 transition-colors underline underline-offset-2">Privacy Policy</a>.
          </p>
        </div>

        <div aria-hidden className="absolute -bottom-px left-1/2 -translate-x-1/2 h-px w-2/3" style={{ background: "linear-gradient(90deg, transparent, rgba(34,211,238,0.4), transparent)" }} />
      </div>

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Syne:wght@700;800&family=DM+Sans:wght@400;500&display=swap');
        @keyframes fade-in { from { opacity:0; transform:translateY(-6px); } to { opacity:1; transform:translateY(0); } }
        .animate-fade-in { animation: fade-in 0.25s ease forwards; }
      `}</style>
    </div>
  );
}