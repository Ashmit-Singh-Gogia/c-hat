import axios from 'axios';

export const apiClient = axios.create({
  baseURL: 'http://localhost:8082/api',
  withCredentials: true, // Crucial for HTTP-Only Cookies
});