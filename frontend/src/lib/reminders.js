// Client for the Go server's "reminders" API. Every endpoint requires a Bearer token.
import { API_BASE as BASE } from "./config.js";
import { dateInputToIso } from "./things.js";

export const REMINDER_TYPES = ["onetime", "daily", "monthly", "yearly"];

// The backend represents reminder_type as an int (1-4); the UI uses strings.
// Translate at the API boundary so components stay string-based.
const TYPE_TO_INT = { onetime: 1, daily: 2, monthly: 3, yearly: 4 };
const INT_TO_TYPE = { 1: "onetime", 2: "daily", 3: "monthly", 4: "yearly" };

// Normalize a reminder coming back from the API: int reminder_type -> string.
function fromApi(reminder) {
  if (!reminder) return reminder;
  return { ...reminder, reminder_type: INT_TO_TYPE[reminder.reminder_type] ?? reminder.reminder_type };
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

export async function listReminders(token) {
  const list = await request("/myReminders", { token });
  return Array.isArray(list) ? list.map(fromApi) : list;
}

export async function getReminder(id, token) {
  return fromApi(await request(`/reminders/${id}`, { token }));
}

// fields: { name, description, remindDate ("YYYY-MM-DD"), reminderType, isActive }
export async function createReminder(fields, token) {
  return fromApi(await request("/reminders", { method: "POST", token, body: toApiBody(fields) }));
}

export async function updateReminder(id, fields, token) {
  return fromApi(await request(`/reminders/${id}`, { method: "PUT", token, body: toApiBody(fields) }));
}

function toApiBody({ name, description, remindDate, reminderType, isActive }) {
  return {
    name,
    description,
    remind_timestamp: dateInputToIso(remindDate),
    reminder_type: TYPE_TO_INT[reminderType] ?? 1,
    is_active: Boolean(isActive),
  };
}

// "YYYY-MM-DD" minus N days, still as "YYYY-MM-DD".
export function minusDays(dateStr, days) {
  const d = new Date(`${dateStr}T00:00:00Z`);
  d.setUTCDate(d.getUTCDate() - Number(days));
  return d.toISOString().slice(0, 10);
}

// Build a one-time reminder that fires `daysBefore` days before a thing expires.
export function expiryReminderFields(thingName, expiryDateStr, daysBefore) {
  return {
    name: `${thingName} expires soon`,
    description: `"${thingName}" expires on ${expiryDateStr} — reminder set ${daysBefore} day(s) before.`,
    remindDate: minusDays(expiryDateStr, daysBefore),
    reminderType: "onetime",
    isActive: true,
  };
}
