import { useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../lib/auth.jsx";
import {
  listThings,
  createThing,
  updateThing,
  deleteThing,
  daysUntil,
} from "../lib/things.js";
import { createReminder, expiryReminderFields } from "../lib/reminders.js";
import ThingModal from "../components/ThingModal.jsx";

export default function ThingsPage() {
  const { token, signOut } = useAuth();
  const navigate = useNavigate();

  const [things, setThings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Toolbar state.
  const [search, setSearch] = useState("");
  const [expiryFilter, setExpiryFilter] = useState("all"); // all | valid | expired
  const [sort, setSort] = useState("recent"); // recent | name | quantity | expiry
  const [filtersOpen, setFiltersOpen] = useState(false);
  const filterRef = useRef(null);

  // Modal state: null, or { mode, thing }.
  const [modal, setModal] = useState(null);

  // Highlight the filter icon when anything is set away from its default.
  const filtersActive = expiryFilter !== "all" || sort !== "recent";

  // Close the filter popover on outside click or Escape.
  useEffect(() => {
    if (!filtersOpen) return;
    function onPointer(e) {
      if (filterRef.current && !filterRef.current.contains(e.target)) {
        setFiltersOpen(false);
      }
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

  // If a request comes back 401 the session is dead — sign out and bail.
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
      const data = await listThings(token);
      setThings(Array.isArray(data) ? data : []);
    } catch (err) {
      if (!handleAuthError(err)) setError(err.message || "Could not load things");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const visible = useMemo(() => {
    const q = search.trim().toLowerCase();
    let list = things.filter((t) => {
      const matchesSearch =
        !q ||
        t.name.toLowerCase().includes(q) ||
        (t.description ?? "").toLowerCase().includes(q);
      const expired = daysUntil(t.expires_at) < 0;
      const matchesExpiry =
        expiryFilter === "all" ||
        (expiryFilter === "valid" && !expired) ||
        (expiryFilter === "expired" && expired);
      return matchesSearch && matchesExpiry;
    });

    list = [...list].sort((a, b) => {
      if (sort === "name") return a.name.localeCompare(b.name);
      if (sort === "quantity") return b.quantity - a.quantity;
      if (sort === "expiry") {
        // Things with no expiry sort last (they never expire).
        const ax = a.expires_at ? new Date(a.expires_at).getTime() : Infinity;
        const bx = b.expires_at ? new Date(b.expires_at).getTime() : Infinity;
        return ax - bx;
      }
      return new Date(b.created_at) - new Date(a.created_at); // recent
    });

    return list;
  }, [things, search, expiryFilter, sort]);

  async function handleSave(fields) {
    if (modal.mode === "create") {
      await createThing(fields, token);
      // If the user asked to be reminded before expiry, create a one-time reminder.
      if (fields.notify && fields.expiresDate) {
        await createReminder(
          expiryReminderFields(fields.name, fields.expiresDate, fields.notifyDaysBefore),
          token,
        );
      }
    } else {
      await updateThing(modal.thing.id, fields, token);
    }
    setModal(null);
    await refresh();
  }

  async function handleDelete() {
    await deleteThing(modal.thing.id, token);
    setModal(null);
    await refresh();
  }

  return (
    <>
      <main className="things-main">
        <div className="things-head">
          <div>
            <h1 className="things-title">Your things</h1>
            <p className="things-count">
              {things.length} item{things.length === 1 ? "" : "s"}
            </p>
          </div>
          <button className="btn" onClick={() => setModal({ mode: "create" })}>
            + New thing
          </button>
        </div>

        <div className="toolbar">
          <input
            className="toolbar-search"
            type="search"
            placeholder="Search by name or description…"
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
                  <span>Expiry</span>
                  <select value={expiryFilter} onChange={(e) => setExpiryFilter(e.target.value)}>
                    <option value="all">All</option>
                    <option value="valid">Valid</option>
                    <option value="expired">Expired</option>
                  </select>
                </label>
                <label className="toolbar-select">
                  <span>Sort by</span>
                  <select value={sort} onChange={(e) => setSort(e.target.value)}>
                    <option value="recent">Newest</option>
                    <option value="name">Name (A–Z)</option>
                    <option value="quantity">Quantity</option>
                    <option value="expiry">Expiry soonest</option>
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
            {things.length === 0
              ? "No things yet. Add your first one!"
              : "No things match your filters."}
          </p>
        ) : (
          <ul className="things-grid">
            {visible.map((t) => (
              <li key={t.id}>
                <button
                  className="thing-card"
                  onClick={() => setModal({ mode: "view", thing: t })}
                >
                  <div className="thing-card-top">
                    <h3 className="thing-name">{t.name}</h3>
                    <span
                      className={`qty-badge ${t.quantity > 0 ? "" : "qty-zero"}`}
                    >
                      ×{t.quantity}
                    </span>
                  </div>
                  {t.description && (
                    <p className="thing-desc">{t.description}</p>
                  )}
                  <p className={`thing-meta ${daysUntil(t.expires_at) < 0 ? "is-expired" : ""}`}>
                    {expiryText(t.expires_at)}
                  </p>
                </button>
              </li>
            ))}
          </ul>
        )}
      </main>

      {modal && (
        <ThingModal
          mode={modal.mode}
          thing={modal.thing}
          onSave={handleSave}
          onDelete={handleDelete}
          onClose={() => setModal(null)}
        />
      )}
    </>
  );
}

function expiryText(iso) {
  if (!iso) return "No expiry";
  const days = daysUntil(iso);
  const date = new Date(iso).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
  if (days < 0) return `Expired · ${date}`;
  if (days === 0) return `Expires today · ${date}`;
  return `Expires in ${days}d · ${date}`;
}
