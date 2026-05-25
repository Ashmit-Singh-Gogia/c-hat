import { useEffect, useRef, useCallback } from 'react';

export function useWebSocket(chatId, onMessage) {
    const wsRef = useRef(null);
    const reconnectTimer = useRef(null);

    const connect = useCallback(() => {
        if (!chatId) return;

        // Automatically uses the same port as the React app (5173), 
        // which Vite then silently proxies to Go (8082) WITH the cookie!
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/ws/${chatId}`;
        
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
            console.log(`✅ WebSocket Connected to chat ${chatId}`);
        };

        ws.onmessage = (event) => {
            try {
                const msg = JSON.parse(event.data);
                if (msg.type === 'new_message') {
                    onMessage(msg.payload);
                }
            } catch (err) {
                console.error("Failed to parse WebSocket message", err);
            }
        };

        ws.onclose = () => {
            console.warn("⚠️ WebSocket Disconnected. Reconnecting...");
            reconnectTimer.current = setTimeout(connect, 3000);
        };

        ws.onerror = (err) => {
            console.error("❌ WebSocket Error:", err);
            ws.close();
        };
    }, [chatId, onMessage]);

    useEffect(() => {
        connect();
        return () => {
            clearTimeout(reconnectTimer.current);
            wsRef.current?.close();
        };
    }, [connect]);

    const sendMessage = useCallback((content) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({ content }));
        } else {
            console.error("❌ Cannot send message: WebSocket is not open.");
        }
    }, []);

    return { sendMessage };
}