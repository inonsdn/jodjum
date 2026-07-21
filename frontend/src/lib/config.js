// Base URL that all API clients build their request paths on.
//
// VITE_API_BASE_URL is the backend ORIGIN (scheme + host, no trailing path),
// e.g. "https://your-service.run.app". It is a build-time value (Vite inlines
// VITE_* vars when `vite build` runs), not read at runtime.
//
//   • Unset (dev, or Vercel using the /api rewrite in vercel.json):
//       API_BASE = "/api/v1"  → same-origin, no CORS needed.
//   • Set to a cross-origin backend URL:
//       API_BASE = "https://your-service.run.app/api/v1"  → the browser calls
//       the backend directly, so the Go server MUST send CORS headers.
const ORIGIN = (import.meta.env.VITE_API_BASE_URL ?? "").replace(/\/+$/, "");

export const API_BASE = `${ORIGIN}/api/v1`;

// Public half of the server's Web Push VAPID key pair (see server/.env
// VAPID_PUBLIC_KEY — same value, mirrored here since it's safe to ship to the
// browser). Used by src/lib/push.js as the PushManager applicationServerKey.
export const VAPID_PUBLIC_KEY = import.meta.env.VITE_VAPID_PUBLIC_KEY ?? "";
