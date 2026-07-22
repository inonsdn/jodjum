// Web Push subscription management: browser permission + PushManager on one
// side, the server's /api/v1/push/* endpoints on the other.
import { API_BASE as BASE, VAPID_PUBLIC_KEY } from "./config.js";

async function request(path, { method = "GET", body, token } = {}) {
  const headers = {};
  if (body) headers["Content-Type"] = "application/json";
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  let data = null;
  try {
    data = await res.json();
  } catch {
    data = null;
  }

  if (!res.ok) {
    const message = typeof data === "string" ? data : "Request failed";
    const err = new Error(message);
    err.status = res.status;
    throw err;
  }

  return data;
}

export function isPushSupported() {
  return "serviceWorker" in navigator && "PushManager" in window && "Notification" in window;
}

// "unsupported" | "denied" | "subscribed" | "not-subscribed"
// (Notification.permission === "default" is folded into "not-subscribed" —
// the UI only needs to know whether to offer "enable" or "disable".)
export async function getPushStatus() {
  if (!isPushSupported()) return "unsupported";
  if (Notification.permission === "denied") return "denied";

  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.getSubscription();
  return subscription ? "subscribed" : "not-subscribed";
}

// Ask for permission (if needed) and subscribe this browser to push, then
// register the subscription with the server so it can target this device.
export async function subscribeToPush(token) {
  if (!isPushSupported()) {
    throw new Error("Push notifications aren't supported in this browser");
  }
  if (!VAPID_PUBLIC_KEY) {
    throw new Error("Missing VAPID_PUBLIC_KEY — notifications aren't configured");
  }

  const permission = await Notification.requestPermission();
  if (permission !== "granted") {
    throw new Error("Notification permission was not granted");
  }

  const registration = await navigator.serviceWorker.ready;
  let subscription = await registration.pushManager.getSubscription();
  if (!subscription) {
    subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(VAPID_PUBLIC_KEY),
    });
  }

  const json = subscription.toJSON();
  await request("/push/subscribe", {
    method: "POST",
    token,
    body: { endpoint: json.endpoint, keys: json.keys },
  });

  return subscription;
}

// Unsubscribe this browser from push, both locally and on the server.
export async function unsubscribeFromPush(token) {
  if (!isPushSupported()) return;

  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.getSubscription();
  if (!subscription) return;

  const endpoint = subscription.endpoint;
  await subscription.unsubscribe();
  await request("/push/unsubscribe", { method: "POST", token, body: { endpoint } });
}

// Web Push application server keys are URL-safe base64; PushManager wants raw bytes.
function urlBase64ToUint8Array(base64String) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");
  const raw = atob(base64);
  return Uint8Array.from([...raw].map((c) => c.charCodeAt(0)));
}
