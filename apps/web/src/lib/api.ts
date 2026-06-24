import axios from "axios";

const API_URL = (import.meta.env.VITE_API_URL as string | undefined) ?? "";

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Auto-attach Idempotency-Key on unsafe methods so any mutation gets
// safe-retry semantics for free.
api.interceptors.request.use((config) => {
  const method = (config.method || "get").toUpperCase();
  const unsafe = method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE";
  if (unsafe && config.headers && !config.headers["Idempotency-Key"]) {
    config.headers["Idempotency-Key"] = crypto.randomUUID();
  }
  return config;
});
