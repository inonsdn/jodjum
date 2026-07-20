// Client for the Go server's "things" API. Every endpoint requires a Bearer token.
// Requests hit /api/... and are proxied to the Go server in dev (see vite.config.js).

import { API_BASE as BASE } from "./config.js";

// The server stores expires_at / created_at as timestamptz, serialized as RFC3339
// strings (e.g. "2026-07-26T00:00:00Z"). The UI works with <input type="date">
// values ("YYYY-MM-DD"), so convert at the boundary.

// ISO timestamp -> "YYYY-MM-DD" for a date input.
export function isoToDateInput(iso) {
  if (!iso) return "";
  const d = new Date(iso);
  return Number.isNaN(d.getTime()) ? "" : d.toISOString().slice(0, 10);
}

// "YYYY-MM-DD" -> RFC3339 timestamp (UTC midnight) for the API.
export function dateInputToIso(date) {
  return new Date(`${date}T00:00:00Z`).toISOString();
}

// Whole days from now until the given ISO timestamp (negative if past).
// Returns null when there is no expiry date.
export function daysUntil(iso) {
  if (!iso) return null;
  const ms = new Date(iso).getTime() - Date.now();
  return Math.ceil(ms / (24 * 60 * 60 * 1000));
}

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

export function listThings(token) {
  return request("/myThings", { token });
}

export function getThing(id, token) {
  return request(`/things/${id}`, { token });
}

// fields: { name, description, quantity, expiresDate }  (expiresDate = "YYYY-MM-DD")
export function createThing(fields, token) {
  return request("/things", {
    method: "POST",
    token,
    body: toApiBody(fields),
  });
}

export function updateThing(id, fields, token) {
  return request(`/things/${id}`, {
    method: "PUT",
    token,
    body: toApiBody(fields),
  });
}

export function deleteThing(id, token) {
  return request(`/things/${id}`, { method: "DELETE", token });
}

function toApiBody({ name, description, quantity, expiresDate, notify, notifyDaysBefore }) {
  const body = {
    name,
    description,
    quantity: Number(quantity),
    // Expiry is optional: send null when the user left the date blank.
    expires_at: expiresDate ? dateInputToIso(expiresDate) : null,
  };

  // Seconds from now until the reminder should fire; the server does
  // `now + value·seconds` and creates a one-time reminder. Only include the
  // field when a reminder was actually requested.
  const offset = remindOffsetSeconds(expiresDate, notify, notifyDaysBefore);
  if (offset !== null) body.next_remind_timestamp = offset;

  return body;
}

// Seconds from now until "notifyDaysBefore days before the expiry date", or null
// when the user didn't ask to be reminded.
function remindOffsetSeconds(expiresDate, notify, notifyDaysBefore) {
  if (!notify || !expiresDate || !notifyDaysBefore) return null;
  const dayMs = 24 * 60 * 60 * 1000;
  const expiryMs = new Date(`${expiresDate}T00:00:00Z`).getTime();
  const remindMs = expiryMs - Number(notifyDaysBefore) * dayMs;
  return Math.round((remindMs - Date.now()) / 1000);
}
