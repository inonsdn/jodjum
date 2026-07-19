// Thin wrapper around the Go server's auth API.
// Requests go to /api/... and are proxied to the Go server in dev (see vite.config.js).

const BASE = "/api/v1";

async function request(path, { method = "GET", body, token } = {}) {
  const headers = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  // The server always responds with JSON (a value or an error string).
  let data = null;
  try {
    data = await res.json();
  } catch {
    data = null;
  }

  if (!res.ok) {
    // Error responses are a plain JSON string message.
    const message = typeof data === "string" ? data : "Request failed";
    throw new Error(message);
  }

  return data;
}

export function login(email, password) {
  return request("/login", { method: "POST", body: { email, password } });
}

export function register(username, email, password) {
  return request("/register", {
    method: "POST",
    body: { username, email, password },
  });
}

export function logout(token) {
  return request("/logout", { method: "POST", token });
}
