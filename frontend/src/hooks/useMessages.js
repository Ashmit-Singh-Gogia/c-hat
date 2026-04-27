import { useState, useEffect, useRef } from 'react';
import { chatService } from '../services/chatService';
import { useAuth } from '../contexts/AuthContext';

export function useMessages(chat) {
  const { user } = useAuth();
  const [messages, setMessages] = useState([]);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    if (!chat) return;
    
    const chatId = chat.id || chat.ID;
    
    const fetchMessages = async () => {
      try {
        const res = await chatService.getMessages(chatId);
        // Safely extract array regardless of backend JSON wrapper
        const data = res.data;
        if (Array.isArray(data)) setMessages(data);
        else if (data?.messages && Array.isArray(data.messages)) setMessages(data.messages);
        else if (data?.data && Array.isArray(data.data)) setMessages(data.data);
        else setMessages([]);
      } catch (error) {
        console.error("Failed to fetch messages:", error);
        setMessages([]);
      }
    };

    fetchMessages();
  }, [chat]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

const sendMessage = async (input) => {
    if (!input.trim() || !chat) return;

    const chatId = chat.id || chat.ID;
    const userId = user?.id || user?.ID;

    try {
      const res = await chatService.sendMessage(chatId, input);
      const newMsg = res.data.message || res.data;
      
      // FIX: Explicitly inject 'content: input' so the text we typed 
      // is guaranteed to render, even if the backend doesn't return it.
      setMessages((prev) => [
        ...prev, 
        { 
          ...newMsg, 
          content: input,        // Ensure lowercase 'content' exists
          Content: input,        // Ensure uppercase 'Content' exists (for Go structs)
          sender_id: userId, 
          SenderID: userId 
        }
      ]);
    } catch (error) {
      console.error("Failed to send message", error);
    }
  };

  return { messages, sendMessage, messagesEndRef };
}