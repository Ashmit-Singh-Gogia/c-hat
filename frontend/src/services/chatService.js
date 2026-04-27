import { apiClient } from '../api/client';

export const chatService = {
  getChats: () => apiClient.get('/chats/'),
  getMessages: (chatId) => apiClient.get(`/chats/${chatId}/messages`),
  sendMessage: (chatId, content) => apiClient.post('/messages/', { chat_id: chatId, content }),
  createDirectChat: (uid) => apiClient.post('/chats/direct', { uid: parseInt(uid) })
};