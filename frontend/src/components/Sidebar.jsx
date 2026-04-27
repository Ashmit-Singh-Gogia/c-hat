import { Search, UserPlus, LogOut } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { useChats } from '../hooks/useChats';
import { authService } from '../services/authService';

export default function Sidebar({ onSelectChat, activeChatId }) {
  const { user } = useAuth();
  
  const { chats, loading, createChat } = useChats();

  // Handle the Add Chat button click
  const handleAddChat = () => {
    const uid = prompt("Enter the User ID you want to chat with:");
    if (uid && !isNaN(uid)) {
      createChat(uid);
    }
  };

  // Helper function to extract the other user's name
const getChatName = (chat) => {
    if (chat.name) return chat.name; 

    const users = chat.Users || chat.users || chat.Participants || chat.participants;
    
    if (Array.isArray(users)) {
      const currentUserId = user?.id || user?.ID;
      const otherUser = users.find((u) => (u.id || u.ID) !== currentUserId && (u.user_id || u.UserID) !== currentUserId);
      
      if (otherUser) {
        // We now check if the nested 'user' object exists!
        const nestedUser = otherUser.user || otherUser.User;
        
        if (nestedUser) {
          return nestedUser.username || nestedUser.Username || nestedUser.name || nestedUser.Name;
        }
        
        // Fallback if nested user is missing
        return `User ${otherUser.user_id || otherUser.UserID || otherUser.id}`;
      }
    }
    return `Chat ${chat.id || chat.ID}`;
  };

  return (
    <div className="w-80 bg-white border-r border-slate-200 flex flex-col h-full">
      {/* Header */}
      <div className="p-4 border-b border-slate-200 flex items-center justify-between">
        <h2 className="text-xl font-bold text-slate-800">Messages</h2>
        {/* Wired up the button! */}
        <button 
          onClick={handleAddChat}
          className="p-2 bg-indigo-50 text-indigo-600 rounded-full hover:bg-indigo-100 transition"
        >
          <UserPlus size={20} />
        </button>
      </div>

      {/* Search */}
      <div className="p-4">
        <div className="relative">
          <Search className="absolute left-3 top-2.5 text-slate-400" size={18} />
          <input 
            type="text" 
            placeholder="Search chats..." 
            className="w-full bg-slate-100 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
      </div>

      {/* Chat List */}
      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <p className="text-center text-slate-400 mt-10 text-sm animate-pulse">Loading chats...</p>
        ) : chats.length === 0 ? (
          <p className="text-center text-slate-400 mt-10 text-sm">No chats yet.</p>
        ) : (
          chats.map((chat) => {
            const chatName = getChatName(chat); 
            
            return (
              <div 
                key={chat.id || chat.ID} 
                onClick={() => onSelectChat(chat)}
                className={`p-4 flex items-center gap-3 cursor-pointer hover:bg-slate-50 transition-colors ${activeChatId === (chat.id || chat.ID) ? 'bg-indigo-50 border-r-4 border-indigo-500' : ''}`}
              >
                <div className="w-12 h-12 bg-indigo-200 rounded-full flex items-center justify-center text-indigo-700 font-bold uppercase">
                  {chatName.charAt(0)}
                </div>
                <div>
                  <h3 className="font-semibold text-slate-800">{chatName}</h3>
                  <p className="text-sm text-slate-500 truncate">Click to view messages</p>
                </div>
              </div>
            );
          })
        )}
      </div>

      {/* Profile Footer */}
      <div className="p-4 border-t border-slate-200 bg-slate-50 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-indigo-600 rounded-full flex items-center justify-center text-white font-bold shadow-sm uppercase">
            {user?.username ? user.username.charAt(0) : 'U'}
          </div>
          <div>
            <p className="text-sm font-semibold text-slate-800">{user?.username || 'My Profile'}</p>
            <div className="flex items-center gap-1">
              <span className="text-xs text-slate-500">My ID:</span>
              <span className="text-xs font-mono bg-slate-200 text-slate-700 px-1.5 py-0.5 rounded">
                {user?.id || user?.ID || 'Unknown'}
              </span>
            </div>
          </div>
        </div>
        <button 
          onClick={authService.logout}
          className="p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-full transition-colors"
        >
          <LogOut size={18} />
        </button>
      </div>
    </div>
  );
}