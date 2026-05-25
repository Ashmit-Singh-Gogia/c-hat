import { useState, useEffect, useRef, useCallback } from 'react';
import { chatService } from '../services/chatService';
import { useAuth } from '../contexts/AuthContext';
import { useWebSocket } from './useWebSocket';

export function useMessages(chat) {
    const { user } = useAuth();
    const [messages, setMessages] = useState([]);
    const messagesEndRef = useRef(null);

    // Called by WebSocket when a new message arrives from server
    const handleIncoming = useCallback((wsEventData) => {
        // 1. Unwrap the Go envelope
        const newMsg = wsEventData.payload ? wsEventData.payload : wsEventData;

        // 2. Ignore non-message events
        if (wsEventData.type && wsEventData.type !== 'new_message') return;

        setMessages(prev => {
            const isFromMe = newMsg.sender_id === (user?.id || user?.ID);

            if (isFromMe) {
                const optimisticIndex = prev.findIndex(m => 
                    String(m.id).startsWith('temp-') && m.content === newMsg.content
                );

                if (optimisticIndex !== -1) {
                    const updated = [...prev];
                    updated[optimisticIndex] = newMsg;
                    return updated;
                }
            }
            
            return [...prev, newMsg];
        });
    }, [user]);

    const { sendMessage: wsSend } = useWebSocket(
        chat?.id || chat?.ID,
        handleIncoming
    );

    // Fetch message history via HTTP on chat open
    useEffect(() => {
        if (!chat) return;
        const chatId = chat.id || chat.ID;

        chatService.getMessages(chatId)
            .then(res => {
                const data = res.data;
                if (Array.isArray(data)) setMessages(data);
                else if (data?.messages) setMessages(data.messages);
                else setMessages([]);
            })
            .catch(() => setMessages([]));
    }, [chat]);

    // Auto-scroll to bottom on new messages
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    const sendMessage = (input) => {
        if (!input.trim() || !chat) return;
        const userId = user?.id || user?.ID;

        // Optimistic update
        const optimistic = {
            id: `temp-${Date.now()}`,
            sender_id: userId,
            content: input,
            created_at: new Date().toISOString(),
        };
        setMessages(prev => [...prev, optimistic]);

        // Send over WebSocket
        wsSend(input);
    };

    return { messages, sendMessage, messagesEndRef };
}