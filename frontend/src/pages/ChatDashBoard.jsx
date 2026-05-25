import { useState, useEffect } from "react";
import {
  Search, Plus, Settings, LogOut, Send, Paperclip,
  Smile, MoreVertical, Phone, Video, Circle,
} from "lucide-react";
import { useAuth } from "../contexts/AuthContext.jsx";
import { chatService } from "../services/chatService";

// IMPORT YOUR NEW WEBSOCKET HOOK HERE
import { useMessages } from "../hooks/useMessages";

function StatusDot({ status }) {
  const colors = { online: "bg-emerald-400", away: "bg-amber-400", offline: "bg-slate-600", group: "bg-violet-400" };
  return <span className={`absolute bottom-0 right-0 w-2.5 h-2.5 rounded-full border-2 border-[#0d0d18] ${colors[status] ?? colors.offline}`} />;
}

function Avatar({ initials, picture, color = "#22d3ee", size = "md", status }) {
  const sz = size === "sm" ? "w-8 h-8 text-xs" : size === "lg" ? "w-11 h-11 text-sm" : "w-10 h-10 text-sm";
  return (
    <div className={`relative shrink-0 ${sz} rounded-full flex items-center justify-center font-semibold text-white overflow-hidden`}
      style={{ background: `${color}22`, border: `1.5px solid ${color}55`, color }}>
      {picture ? <img src={picture} alt="" className="w-full h-full object-cover" /> : initials}
      {status && <StatusDot status={status} />}
    </div>
  );
}

function ContactRow({ chat, active, onClick }) {
  return (
    <button onClick={onClick}
      className={`w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-left transition-all duration-150 ${active ? "bg-white/[0.07]" : "hover:bg-white/[0.04]"}`}>
      <Avatar initials={chat.initials} color={chat.color} status={chat.status} />
      <div className="flex-1 min-w-0">
        <div className="flex items-baseline justify-between gap-2">
          <span className={`text-sm font-medium truncate ${active ? "text-white" : "text-slate-200"}`}>{chat.name}</span>
          <span className="text-[11px] text-slate-600 shrink-0">{chat.time}</span>
        </div>
        <p className="text-xs text-slate-500 truncate mt-0.5">{chat.lastMsg}</p>
      </div>
      {chat.unread > 0 && (
        <span className="shrink-0 w-5 h-5 rounded-full bg-cyan-500 flex items-center justify-center text-[10px] font-bold text-white">{chat.unread}</span>
      )}
    </button>
  );
}

function MessageBubble({ msg, currentUserId }) {
  const isMe = msg.sender_id === currentUserId;
  return (
    <div className={`flex items-end gap-2.5 ${isMe ? "flex-row-reverse" : "flex-row"}`}>
      {!isMe && <Avatar initials={msg.sender_initials} color={msg.sender_color ?? "#22d3ee"} size="sm" />}
      <div className={`max-w-[68%] flex flex-col gap-1 ${isMe ? "items-end" : "items-start"}`}>
        <div className={`px-4 py-2.5 rounded-2xl text-sm leading-relaxed ${isMe ? "rounded-br-sm text-white" : "rounded-bl-sm text-slate-200"}`}
          style={isMe
            ? { background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)", boxShadow: "0 4px 16px rgba(34,211,238,0.2)" }
            : { background: "rgba(255,255,255,0.06)", border: "1px solid rgba(255,255,255,0.08)" }}>
          {msg.content}
        </div>
        <span className="text-[11px] text-slate-600 px-1">
          {msg.created_at ? new Date(msg.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }) : ""}
        </span>
      </div>
    </div>
  );
}

const COLORS = ["#22d3ee","#a78bfa","#34d399","#f472b6","#fb923c","#60a5fa"];
function chatColor(id) { return COLORS[id % COLORS.length]; }
function toInitials(name = "") { return (name || "").split(" ").map(w => w[0]).join("").slice(0, 2).toUpperCase(); }

export default function ChatDashboard() {
  const { user, logout } = useAuth();
  const [chats, setChats] = useState([]);
  const [activeChat, setActiveChat] = useState(null);
  const [input, setInput] = useState("");
  const [query, setQuery] = useState("");

  // THIS IS THE BRAIN OF YOUR COMPONENT NOW. 
  // It handles state, fetching, scrolling, and WebSockets automatically.
  const { messages, sendMessage, messagesEndRef } = useMessages(activeChat);

  // Fetch sidebar chats
  useEffect(() => {
    chatService.getChats()
      .then(res => setChats(res.data.map(c => ({
        ...c,
        initials: toInitials(c.name),
        color: chatColor(c.id),
        time: c.last_message_at ? new Date(c.last_message_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }) : "",
        lastMsg: c.last_message ?? "",
        unread: c.unread_count ?? 0,
        status: c.is_group ? "group" : "online",
      }))))
      .catch(console.error);
  }, []);

  const handleSend = () => {
    if (!input.trim() || !activeChat) return;
    const content = input.trim();
    setInput("");
    sendMessage(content); 
  };

  const userInitials = toInitials(user?.name);
  const filtered = chats.filter(c => (c.name ?? "").toLowerCase().includes(query.toLowerCase()));

  return (
    <div className="flex h-screen w-screen overflow-hidden" style={{ background: "#07070f", fontFamily: "'DM Sans', sans-serif" }}>

      {/* SIDEBAR */}
      <aside className="flex flex-col w-72 shrink-0 h-full border-r" style={{ borderColor: "rgba(255,255,255,0.06)", background: "rgba(255,255,255,0.02)" }}>
        {/* User profile */}
        <div className="flex items-center gap-3 px-4 py-4 border-b" style={{ borderColor: "rgba(255,255,255,0.06)" }}>
          <div className="w-9 h-9 rounded-full flex items-center justify-center text-sm font-bold text-white shrink-0 overflow-hidden"
            style={{ background: "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)" }}>
            {user?.picture ? <img src={user.picture} alt="" className="w-full h-full object-cover" /> : userInitials}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-semibold text-white truncate">{user?.name ?? "You"}</p>
            <div className="flex items-center gap-1">
              <Circle size={7} className="text-emerald-400 fill-emerald-400" />
              <span className="text-[11px] text-slate-500">ID: {user?.id ?? "—"}</span>
            </div>
          </div>
          <div className="flex items-center gap-1">
            <button onClick={() => logout("google")} title="Sign out"
              className="p-1.5 rounded-lg text-slate-600 hover:text-red-400 hover:bg-red-500/10 transition-all duration-150">
              <LogOut size={15} />
            </button>
            <button className="p-1.5 rounded-lg text-slate-600 hover:text-slate-300 hover:bg-white/[0.06] transition-all duration-150" title="Settings">
              <Settings size={15} />
            </button>
          </div>
        </div>

        {/* Search + New Chat */}
        <div className="px-3 pt-3 pb-2 flex items-center gap-2">
          <div className="flex-1 flex items-center gap-2 px-3 py-2 rounded-xl" style={{ background: "rgba(255,255,255,0.04)", border: "1px solid rgba(255,255,255,0.07)" }}>
            <Search size={13} className="text-slate-600 shrink-0" />
            <input value={query} onChange={e => setQuery(e.target.value)}
              placeholder="Search conversations…"
              className="flex-1 bg-transparent text-xs text-slate-300 placeholder-slate-600 outline-none" />
          </div>
          <button
            onClick={() => {
              const uid = prompt("Enter the User ID to start a chat:");
              if (!uid || isNaN(uid)) return;
              chatService.createDirectChat(uid)
                .then(res => {
                  const chat = res.data;
                  const formatted = {
                    ...chat,
                    initials: toInitials(chat.name),
                    color: chatColor(chat.id),
                    time: "",
                    lastMsg: "",
                    unread: 0,
                    status: chat.is_group ? "group" : "online",
                  };
                  setChats(prev => [formatted, ...prev.filter(x => x.id !== chat.id)]);
                  setActiveChat(formatted);
                })
                .catch(err => alert(err.response?.data?.error ?? "Failed to create chat"));
            }}
            className="w-8 h-8 rounded-xl flex items-center justify-center text-slate-400 hover:text-white transition-all duration-150 shrink-0"
            style={{ background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.08)" }}
            title="New chat">
            <Plus size={15} />
          </button>
        </div>

        {/* Contact list */}
        <div className="flex-1 overflow-y-auto px-2 pb-3 scrollbar-hide">
          <p className="text-[10px] uppercase tracking-widest text-slate-700 font-semibold px-2 pt-3 pb-2">Messages</p>
          <div className="flex flex-col gap-0.5">
            {filtered.map(c => (
              <ContactRow key={c.id} chat={c} active={activeChat?.id === c.id} onClick={() => setActiveChat(c)} />
            ))}
          </div>
        </div>
      </aside>

      {/* MAIN CHAT AREA */}
      <main className="flex flex-col flex-1 min-w-0">
        {activeChat ? (
          <>
            {/* Header */}
            <header className="flex items-center justify-between px-5 py-3.5 border-b shrink-0"
              style={{ borderColor: "rgba(255,255,255,0.06)", background: "rgba(255,255,255,0.02)" }}>
              <div className="flex items-center gap-3">
                <Avatar initials={activeChat.initials} color={activeChat.color} size="md" status={activeChat.status} />
                <div>
                  <h2 className="text-sm font-semibold text-white">{activeChat.name}</h2>
                  <p className="text-[11px] text-slate-500">
                    {activeChat.status === "online" ? "Active now" : activeChat.status === "away" ? "Away" : activeChat.status === "group" ? "Group" : "Last seen recently"}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-1">
                {[{ icon: Phone, label: "Call" }, { icon: Video, label: "Video" }, { icon: MoreVertical, label: "More" }].map(({ icon: Icon, label }) => (
                  <button key={label} title={label} className="p-2 rounded-xl text-slate-500 hover:text-slate-200 hover:bg-white/[0.06] transition-all duration-150">
                    <Icon size={17} />
                  </button>
                ))}
              </div>
            </header>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto px-5 py-5 flex flex-col gap-4 scrollbar-hide">
              {messages.map((msg, idx) => (
                <MessageBubble key={msg.id || idx} msg={msg} currentUserId={user?.id} />
              ))}
              {/* This connects to the hook to auto-scroll */}
              <div ref={messagesEndRef} />
            </div>

            {/* Input bar */}
            <div className="shrink-0 px-4 py-3 border-t flex items-end gap-3"
              style={{ borderColor: "rgba(255,255,255,0.06)", background: "rgba(255,255,255,0.02)" }}>
              <button className="p-2 text-slate-600 hover:text-slate-300 transition-colors mb-0.5"><Paperclip size={18} /></button>
              <div className="flex-1 flex items-end gap-2 px-4 py-2.5 rounded-2xl"
                style={{ background: "rgba(255,255,255,0.05)", border: "1px solid rgba(255,255,255,0.09)" }}>
                <textarea value={input} onChange={e => setInput(e.target.value)}
                  onKeyDown={e => { if (e.key === "Enter" && !e.shiftKey) { e.preventDefault(); handleSend(); } }}
                  placeholder={`Message ${activeChat.name}…`} rows={1}
                  className="flex-1 bg-transparent text-sm text-slate-200 placeholder-slate-600 outline-none resize-none max-h-32 leading-relaxed"
                  style={{ scrollbarWidth: "none" }} />
                <button className="text-slate-600 hover:text-slate-300 transition-colors shrink-0 mb-0.5"><Smile size={17} /></button>
              </div>
              <button onClick={handleSend} disabled={!input.trim()}
                className="p-2.5 rounded-xl transition-all duration-150 shrink-0 mb-0.5 disabled:opacity-40 disabled:cursor-not-allowed"
                style={{
                  background: input.trim() ? "linear-gradient(135deg, #22d3ee 0%, #6366f1 100%)" : "rgba(255,255,255,0.06)",
                  boxShadow: input.trim() ? "0 4px 16px rgba(34,211,238,0.25)" : "none",
                  color: "white",
                }}>
                <Send size={17} />
              </button>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <p className="text-slate-600 text-sm">Select a conversation</p>
          </div>
        )}
      </main>

      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600&display=swap');
        .scrollbar-hide::-webkit-scrollbar { display: none; }
        .scrollbar-hide { -ms-overflow-style: none; scrollbar-width: none; }
      `}</style>
    </div>
  );
}