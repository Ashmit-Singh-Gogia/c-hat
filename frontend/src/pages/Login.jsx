import { MessageSquare } from 'lucide-react';
import { authService } from '../services/authService';

export default function Login() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-slate-100">
      <div className="w-full max-w-md p-8 bg-white rounded-2xl shadow-xl border border-slate-100 text-center">
        <div className="flex justify-center mb-6">
          <div className="p-3 bg-indigo-100 rounded-full text-indigo-600">
            <MessageSquare size={40} />
          </div>
        </div>
        <h1 className="text-3xl font-bold text-slate-800 mb-2">Welcome to C-Hat</h1>
        <p className="text-slate-500 mb-8">Sign in to connect with your friends.</p>
        
        <button 
          onClick={() => window.location.href = authService.loginUrl}
          className="w-full flex items-center justify-center gap-3 px-4 py-3 bg-white border border-slate-300 rounded-xl hover:bg-slate-50 transition-all font-medium text-slate-700 shadow-sm"
        >
          <img src="https://www.svgrepo.com/show/475656/google-color.svg" alt="Google" className="w-6 h-6" />
          Continue with Google
        </button>
      </div>
    </div>
  );
}