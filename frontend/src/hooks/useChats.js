import { useState, useEffect } from 'react';
import { chatService } from '../services/chatService';

export function useChats() {
  const [chats, setChats] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchChats = async () => {
    try {
      const res = await chatService.getChats();
      setChats(res.data || []);
    } catch (error) {
      console.error("Failed to fetch chats", error);
      setChats([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchChats();
  }, []);

  // NEW: Function to create a chat and refresh the list
  const createChat = async (uid) => {
    try {
      await chatService.createDirectChat(uid);
      await fetchChats(); // Refresh the list so the new chat appears!
    } catch (error) {
      console.error("Failed to create chat", error);
      alert("Failed to create chat. Make sure the user ID exists.");
    }
  };

  return { chats, loading, createChat };
}