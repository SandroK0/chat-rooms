export const config = {
  websocketUrl: import.meta.env.VITE_WEBSOCKET_URL || 'ws://localhost:8080/ws',
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
};