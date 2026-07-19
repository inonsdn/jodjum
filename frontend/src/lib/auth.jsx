import { createContext, useContext, useMemo, useState } from "react";

// Simple auth state held in React + localStorage so a refresh keeps you signed in.
const AuthContext = createContext(null);

const STORAGE_KEY = "jodjum.auth";

function readStored() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function AuthProvider({ children }) {
  const [auth, setAuth] = useState(readStored);

  const value = useMemo(
    () => ({
      // The Go server returns { AccessToken, UserId } (Go field names, no json tags).
      token: auth?.token ?? null,
      userId: auth?.userId ?? null,
      isAuthenticated: Boolean(auth?.token),
      signIn(result) {
        const next = { token: result.AccessToken, userId: result.UserId };
        localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
        setAuth(next);
      },
      signOut() {
        localStorage.removeItem(STORAGE_KEY);
        setAuth(null);
      },
    }),
    [auth],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
