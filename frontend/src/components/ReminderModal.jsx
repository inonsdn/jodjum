import { useEffect, useState } from "react";
import { isoToDateInput } from "../lib/things.js";
import { REMINDER_TYPES } from "../lib/reminders.js";

// Modal for viewing / creating / editing a reminder.
// Callbacks: onSave(fields), onClose(). (The API has no delete endpoint.)
export default function ReminderModal({ mode: initialMode, reminder, onSave, onToggleActive, onClose }) {
  const [mode, setMode] = useState(initialMode);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState(() => ({
    name: reminder?.name ?? "",
    description: reminder?.description ?? "",
    remindDate: reminder ? isoToDateInput(reminder.remind_timestamp) : "",
    reminderType: reminder?.reminder_type ?? "onetime",
    isActive: reminder?.is_active ?? true,
  }));

  useEffect(() => {
    function onKey(e) {
      if (e.key === "Escape") onClose();
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  const isForm = mode === "edit" || mode === "create";

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    if (!form.name.trim()) {
      setError("Name is required");
      return;
    }
    if (!form.remindDate) {
      setError("Pick a date to be reminded on");
      return;
    }
    setBusy(true);
    try {
      await onSave(form);
    } catch (err) {
      setError(err.message || "Could not save");
      setBusy(false);
    }
  }

  async function handleToggle() {
    setError("");
    setBusy(true);
    try {
      await onToggleActive(reminder);
    } catch (err) {
      setError(err.message || "Could not update");
      setBusy(false);
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal" role="dialog" aria-modal="true" onClick={(e) => e.stopPropagation()}>
        <div className="modal-head">
          <h2 className="modal-title">
            {mode === "create" ? "New reminder" : mode === "edit" ? "Edit reminder" : reminder.name}
          </h2>
          <button className="modal-close" onClick={onClose} aria-label="Close">
            ×
          </button>
        </div>

        {error && <div className="auth-error">{error}</div>}

        {isForm ? (
          <form className="modal-body" onSubmit={handleSubmit}>
            <label className="field">
              <span>Name</span>
              <input
                type="text"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="e.g. Pay rent"
                autoFocus
              />
            </label>

            <label className="field">
              <span>Description</span>
              <textarea
                rows={3}
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                placeholder="Optional details"
              />
            </label>

            <div className="field-row">
              <label className="field">
                <span>Remind me on</span>
                <input
                  type="date"
                  className={form.remindDate ? undefined : "date-empty"}
                  value={form.remindDate}
                  onChange={(e) => setForm({ ...form, remindDate: e.target.value })}
                />
              </label>
              <label className="field">
                <span>Repeat</span>
                <select
                  value={form.reminderType}
                  onChange={(e) => setForm({ ...form, reminderType: e.target.value })}
                >
                  {REMINDER_TYPES.map((t) => (
                    <option key={t} value={t}>
                      {typeLabel(t)}
                    </option>
                  ))}
                </select>
              </label>
            </div>

            <div className="modal-actions">
              <button
                type="button"
                className="btn btn-ghost"
                onClick={mode === "edit" ? () => setMode("view") : onClose}
                disabled={busy}
              >
                Cancel
              </button>
              <button type="submit" className="btn" disabled={busy}>
                {busy ? "Saving…" : "Save"}
              </button>
            </div>
          </form>
        ) : (
          <div className="modal-body">
            <dl className="detail-list">
              <div>
                <dt>Description</dt>
                <dd>{reminder.description || <em>None</em>}</dd>
              </div>
              <div>
                <dt>Remind on</dt>
                <dd>{formatDate(reminder.remind_timestamp)}</dd>
              </div>
              <div>
                <dt>Repeat</dt>
                <dd>{typeLabel(reminder.reminder_type)}</dd>
              </div>
              <div>
                <dt>Status</dt>
                <dd>{reminder.is_active ? "Active" : "Paused"}</dd>
              </div>
            </dl>

            <div className="modal-actions">
              <button className="btn btn-ghost" onClick={handleToggle} disabled={busy}>
                {reminder.is_active ? "Set inactive" : "Set active"}
              </button>
              <button className="btn" onClick={() => setMode("edit")} disabled={busy}>
                Edit
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function typeLabel(t) {
  return { onetime: "One-time", daily: "Daily", monthly: "Monthly", yearly: "Yearly" }[t] ?? t;
}

function formatDate(iso) {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" });
}
