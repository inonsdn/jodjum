import { useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../lib/auth.jsx";
import { listReminders, createReminder, updateReminder } from "../lib/reminders.js";
import { isoToDateInput } from "../lib/things.js";
import ReminderModal from "../components/ReminderModal.jsx";

export default function RemindersPage() {
  const { token, signOut } = useAuth();
  const navigate = useNavigate();

  const [reminders, setReminders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("all"); // all | active | paused
  const [sort, setSort] = useState("soonest"); // soonest | name
  const [filtersOpen, setFiltersOpen] = useState(false);
  const filterRef = useRef(null);

  const [modal, setModal] = useState(null); // null | { mode, reminder }

  const filtersActive = statusFilter !== "all" || sort !== "soonest";

  function handleAuthError(err) {
    if (err.status === 401) {
      signOut();
      navigate("/login", { replace: true });
      return true;
    }
    return false;
  }

  async function refresh() {
    setLoading(true);
    setError("");
    try {
      const data = await listReminders(token);
      setReminders(Array.isArray(data) ? data : []);
    } catch (err) {
      if (!handleAuthError(err)) setError(err.message || "Could not load reminders");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!filtersOpen) return;
    function onPointer(e) {
      if (filterRef.current && !filterRef.current.contains(e.target)) setFiltersOpen(false);
    }
    function onKey(e) {
      if (e.key === "Escape") setFiltersOpen(false);
    }
    document.addEventListener("mousedown", onPointer);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onPointer);
      document.removeEventListener("keydown", onKey);
    };
  }, [filtersOpen]);

  const visible = useMemo(() => {
    const q = search.trim().toLowerCase();
    let list = reminders.filter((rem) => {
      const matchesSearch =
        !q ||
        rem.name.toLowerCase().includes(q) ||
        (rem.description ?? "").toLowerCase().includes(q);
      const matchesStatus =
        statusFilter === "all" ||
        (statusFilter === "active" && rem.is_active) ||
        (statusFilter === "paused" && !rem.is_active);
      return matchesSearch && matchesStatus;
    });

    list = [...list].sort((a, b) => {
      if (sort === "name") return a.name.localeCompare(b.name);
      return new Date(a.remind_timestamp) - new Date(b.remind_timestamp); // soonest
    });

    return list;
  }, [reminders, search, statusFilter, sort]);

  async function handleSave(fields) {
    if (modal.mode === "create") {
      await createReminder(fields, token);
    } else {
      await updateReminder(modal.reminder.id, fields, token);
    }
    setModal(null);
    await refresh();
  }

  // Flip a reminder's active state, keeping its other fields unchanged.
  async function handleToggleActive(reminder) {
    await updateReminder(
      reminder.id,
      {
        name: reminder.name,
        description: reminder.description,
        remindDate: isoToDateInput(reminder.remind_timestamp),
        reminderType: reminder.reminder_type,
        isActive: !reminder.is_active,
      },
      token,
    );
    setModal(null);
    await refresh();
  }

  return (
    <>
      <main className="things-main">
        <div className="things-head">
          <div>
            <h1 className="things-title">Your reminders</h1>
            <p className="things-count">
              {reminders.length} reminder{reminders.length === 1 ? "" : "s"}
            </p>
          </div>
          <button className="btn" onClick={() => setModal({ mode: "create" })}>
            + New reminder
          </button>
        </div>

        <div className="toolbar">
          <input
            className="toolbar-search"
            type="search"
            placeholder="Search reminders…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <div className="filter-menu" ref={filterRef}>
            <button
              type="button"
              className={`filter-btn ${filtersActive ? "is-active" : ""}`}
              aria-label="Filter and sort"
              aria-expanded={filtersOpen}
              onClick={() => setFiltersOpen((v) => !v)}
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
                <polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3" />
              </svg>
              <span>Filter</span>
              {filtersActive && <span className="filter-dot" />}
            </button>

            {filtersOpen && (
              <div className="filter-popover" role="menu">
                <label className="toolbar-select">
                  <span>Status</span>
                  <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
                    <option value="all">All</option>
                    <option value="active">Active</option>
                    <option value="paused">Paused</option>
                  </select>
                </label>
                <label className="toolbar-select">
                  <span>Sort by</span>
                  <select value={sort} onChange={(e) => setSort(e.target.value)}>
                    <option value="soonest">Soonest</option>
                    <option value="name">Name (A–Z)</option>
                  </select>
                </label>
              </div>
            )}
          </div>
        </div>

        {error && <div className="auth-error">{error}</div>}

        {loading ? (
          <p className="things-empty">Loading…</p>
        ) : visible.length === 0 ? (
          <p className="things-empty">
            {reminders.length === 0
              ? "No reminders yet. Create your first one!"
              : "No reminders match your filters."}
          </p>
        ) : (
          <ul className="things-grid">
            {visible.map((rem) => (
              <li key={rem.id}>
                <button className="thing-card" onClick={() => setModal({ mode: "view", reminder: rem })}>
                  <div className="thing-card-top">
                    <h3 className="thing-name">{rem.name}</h3>
                    <span className={`qty-badge ${rem.is_active ? "" : "qty-zero"}`}>
                      {rem.is_active ? "Active" : "Paused"}
                    </span>
                  </div>
                  {rem.description && <p className="thing-desc">{rem.description}</p>}
                  <p className="thing-meta">
                    {typeLabel(rem.reminder_type)} · {formatDate(rem.remind_timestamp)}
                  </p>
                </button>
              </li>
            ))}
          </ul>
        )}
      </main>

      {modal && (
        <ReminderModal
          mode={modal.mode}
          reminder={modal.reminder}
          onSave={handleSave}
          onToggleActive={handleToggleActive}
          onClose={() => setModal(null)}
        />
      )}
    </>
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
