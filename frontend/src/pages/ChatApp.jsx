import { useState } from 'react';
import Sidebar from '../components/Sidebar';
import ChatWindow from '../components/ChatWindow';

export default function ChatApp() {
  const [activeChat, setActiveChat] = useState(null);

  return (
    <div className="flex h-screen w-full overflow-hidden">
      <Sidebar onSelectChat={setActiveChat} activeChatId={activeChat?.id} />
      <ChatWindow chat={activeChat} />
    </div>
  );
}