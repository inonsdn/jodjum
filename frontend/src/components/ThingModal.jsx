import { useEffect, useState } from "react";
import { isoToDateInput, daysUntil } from "../lib/things.js";
import DateField from "./DateField.jsx";

// A single modal that handles three modes:
//   "view"   – read-only details with Edit / Delete actions
//   "edit"   – form to update an existing thing
//   "create" – form to add a new thing
// Callbacks: onSave(fields), onDelete(), onClose().
export default function ThingModal({ mode: initialMode, thing, onSave, onDelete, onClose }) {
  const [mode, setMode] = useState(initialMode);
  const [confirmingDelete, setConfirmingDelete] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState(() => ({
    name: thing?.name ?? "",
    description: thing?.description ?? "",
    quantity: thing?.quantity ?? 1,
    expiresDate: thing ? isoToDateInput(thing.expires_at) : "",
    // Expiry reminder (create flow only): notify N days before the expiry date.
    notify: false,
    notifyDaysBefore: 7,
  }));

  // Close on Escape for keyboard users.
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
    setBusy(true);
    try {
      await onSave(form);
    } catch (err) {
      setError(err.message || "Could not save");
      setBusy(false);
    }
  }

  async function handleDelete() {
    setBusy(true);
    try {
      await onDelete();
    } catch (err) {
      setError(err.message || "Could not delete");
      setBusy(false);
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div
        className="modal"
        role="dialog"
        aria-modal="true"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="modal-head">
          <h2 className="modal-title">
            {mode === "create"
              ? "New thing"
              : mode === "edit"
                ? "Edit thing"
                : thing.name}
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
                placeholder="e.g. Milk"
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
                <span>Quantity</span>
                <input
                  type="number"
                  min="0"
                  max="127"
                  value={form.quantity}
                  onChange={(e) => setForm({ ...form, quantity: e.target.value })}
                />
              </label>
              <label className="field">
                <span>Expiry date (optional)</span>
                <DateField
                  value={form.expiresDate}
                  onChange={(e) => setForm({ ...form, expiresDate: e.target.value })}
                />
              </label>
            </div>

            {/* Expiry reminder — only offered on create, and only once a date is set. */}
            {mode === "create" && form.expiresDate && (
              <div className="notify-box">
                <label className="checkbox-row">
                  <input
                    type="checkbox"
                    checked={form.notify}
                    onChange={(e) => setForm({ ...form, notify: e.target.checked })}
                  />
                  <span>Remind me before it expires</span>
                </label>

                {form.notify && (
                  <div className="chip-row">
                    {[3, 7, 30].map((d) => (
                      <button
                        key={d}
                        type="button"
                        className={`chip ${form.notifyDaysBefore === d ? "is-selected" : ""}`}
                        onClick={() => setForm({ ...form, notifyDaysBefore: d })}
                      >
                        {d} days before
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}

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
                <dd>{thing.description || <em>None</em>}</dd>
              </div>
              <div>
                <dt>Quantity</dt>
                <dd>{thing.quantity}</dd>
              </div>
              <div>
                <dt>Expires</dt>
                <dd>
                  {thing.expires_at ? (
                    <>
                      {formatDate(thing.expires_at)}
                      <span className="detail-sub"> · {expiryLabel(thing.expires_at)}</span>
                    </>
                  ) : (
                    <em>No expiry</em>
                  )}
                </dd>
              </div>
              <div>
                <dt>Created</dt>
                <dd>{formatDate(thing.created_at)}</dd>
              </div>
            </dl>

            {confirmingDelete ? (
              <div className="confirm-box">
                <span>Delete “{thing.name}”? This can’t be undone.</span>
                <div className="modal-actions">
                  <button
                    className="btn btn-ghost"
                    onClick={() => setConfirmingDelete(false)}
                    disabled={busy}
                  >
                    Cancel
                  </button>
                  <button className="btn btn-danger" onClick={handleDelete} disabled={busy}>
                    {busy ? "Deleting…" : "Delete"}
                  </button>
                </div>
              </div>
            ) : (
              <div className="modal-actions">
                <button
                  className="btn btn-danger-ghost"
                  onClick={() => setConfirmingDelete(true)}
                >
                  Delete
                </button>
                <button className="btn" onClick={() => setMode("edit")}>
                  Edit
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

function formatDate(iso) {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

function expiryLabel(iso) {
  const days = daysUntil(iso);
  if (days < 0) return `expired ${-days} day${days === -1 ? "" : "s"} ago`;
  if (days === 0) return "expires today";
  return `in ${days} day${days === 1 ? "" : "s"}`;
}
