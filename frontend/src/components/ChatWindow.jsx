import { useState } from 'react';
import { Send } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { useMessages } from '../hooks/useMessages';

export default function ChatWindow({ chat }) {
  const { user } = useAuth();
  const [input, setInput] = useState('');
  
  // All fetching and sending logic abstracted into this hook!
  const { messages, sendMessage, messagesEndRef } = useMessages(chat);

  const handleSend = (e) => {
    e.preventDefault();
    sendMessage(input);
    setInput('');
  };

  // --- NEW HELPER FUNCTION GOES RIGHT HERE ---
  // We put it inside the component so it has access to the 'user' variable above
  const getChatName = (c) => {
    if (c.name) return c.name;
    const users = c.Users || c.users || c.Participants || c.participants;
    if (Array.isArray(users)) {
      const currentUserId = user?.id || user?.ID;
      const otherUser = users.find((u) => (u.id || u.ID) !== currentUserId);
      if (otherUser) return otherUser.username || otherUser.Username || otherUser.name || otherUser.Name;
    }
    return `Chat ${c.id || c.ID}`;
  };

  // If no chat is selected, show the empty state screen
  if (!chat) {
    return (
      <div className="flex-1 flex items-center justify-center bg-slate-50 h-full">
        <div className="text-center text-slate-400">
          <div className="bg-slate-200 p-4 rounded-full inline-block mb-4">
            <Send size={32} className="text-slate-400" />
          </div>
          <h2 className="text-xl font-medium">Select a chat to start messaging</h2>
        </div>
      </div>
    );
  }

  // If a chat IS selected, render the actual chat window
  return (
    <div className="flex-1 flex flex-col h-full bg-white">
      
      {/* Header */}
      <div className="h-16 border-b border-slate-200 flex items-center justify-between px-6 shadow-sm z-10 bg-white">
        {/* Left Side: Chat Name & Chat ID */}
        <div className="flex items-center gap-3">
          <h2 className="text-lg font-semibold text-slate-800">{getChatName(chat)}</h2>
          <span className="text-xs font-mono bg-slate-100 text-slate-500 px-2 py-1 rounded border border-slate-200">
            Chat ID: {chat.id || chat.ID}
          </span>
        </div>

        {/* Right Side: Logged-in User ID */}
        <div className="flex items-center gap-2">
          <span className="text-xs text-slate-400 font-medium uppercase tracking-wider">Session</span>
          <span className="text-sm font-mono font-bold text-indigo-600 bg-indigo-50 border border-indigo-100 px-2.5 py-1 rounded-md">
            User ID: {user?.id || user?.ID || 'Unknown'}
          </span>
        </div>
      </div>
      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto p-6 bg-slate-50 space-y-4">
        {Array.isArray(messages) && messages.length === 0 && (
          <div className="text-center text-slate-400 mt-10 text-sm">No messages yet. Say hi!</div>
        )}
        
        {Array.isArray(messages) && messages.map((msg, idx) => {
          const userId = user?.id || user?.ID;
          const senderId = msg.sender_id || msg.SenderID;
          const isMine = senderId === userId;
          
          return (
            <div key={idx} className={`flex ${isMine ? 'justify-end' : 'justify-start'}`}>
              <div 
                className={`max-w-[70%] px-4 py-2.5 rounded-2xl shadow-sm ${
                  isMine 
                    ? 'bg-indigo-600 text-white rounded-br-none' 
                    : 'bg-white border border-slate-200 text-slate-800 rounded-bl-none'
                }`}
              >
                <p>{msg.content || msg.Content}</p>
              </div>
            </div>
          );
        })}
        <div ref={messagesEndRef} />
      </div>

      {/* Input Area */}
      <div className="p-4 bg-white border-t border-slate-200">
        <form onSubmit={handleSend} className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a message..."
            className="flex-1 bg-slate-100 rounded-full px-6 py-3 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all"
          />
          <button 
            type="submit"
            disabled={!input.trim()}
            className="p-3 bg-indigo-600 text-white rounded-full hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-md"
          >
            <Send size={20} className="ml-1" />
          </button>
        </form>
      </div>
    </div>
  );
}