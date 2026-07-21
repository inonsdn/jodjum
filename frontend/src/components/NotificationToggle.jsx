import { useEffect, useState } from "react";
import { useAuth } from "../lib/auth.jsx";
import { getPushStatus, subscribeToPush, unsubscribeFromPush } from "../lib/push.js";

// Bell button that subscribes/unsubscribes this browser for push notifications.
// Lives in the nav so it's reachable from every page.
export default function NotificationToggle() {
  const { token } = useAuth();
  const [status, setStatus] = useState("checking"); // checking | unsupported | denied | subscribed | not-subscribed
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;
    getPushStatus()
      .then((s) => {
        if (!cancelled) setStatus(s);
      })
      .catch(() => {
        if (!cancelled) setStatus("unsupported");
      });
    return () => {
      cancelled = true;
    };
  }, []);

  if (status === "checking" || status === "unsupported") return null;

  async function handleClick() {
    setError("");
    setBusy(true);
    try {
      if (status === "subscribed") {
        await unsubscribeFromPush(token);
        setStatus("not-subscribed");
      } else {
        await subscribeToPush(token);
        setStatus("subscribed");
      }
    } catch (err) {
      setError(err.message || "Something went wrong");
      setStatus(await getPushStatus());
    } finally {
      setBusy(false);
    }
  }

  const isOn = status === "subscribed";
  const isDenied = status === "denied";

  return (
    <button
      type="button"
      className={`btn btn-ghost nav-bell ${isOn ? "is-on" : ""}`}
      onClick={handleClick}
      disabled={busy || isDenied}
      title={
        isDenied
          ? "Notifications blocked in browser settings"
          : isOn
            ? "Turn off notifications"
            : "Turn on notifications"
      }
    >
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
      >
        <path d="M18 8a6 6 0 0 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
        <path d="M13.73 21a2 2 0 0 1-3.46 0" />
      </svg>
      <span>{isDenied ? "Blocked" : isOn ? "Notifications on" : "Enable notifications"}</span>
      {error && <span className="nav-bell-error">{error}</span>}
    </button>
  );
}
