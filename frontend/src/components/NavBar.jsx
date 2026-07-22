import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../lib/auth.jsx";
import NotificationToggle from "./NotificationToggle.jsx";
import VersionTag from "./VersionTag.jsx";

// Top navigation: horizontal menu on desktop, hamburger-toggled panel on mobile.
const LINKS = [
  { to: "/things", label: "My Things" },
  { to: "/reminders", label: "Reminders" },
];

export default function NavBar() {
  const { signOut } = useAuth();
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);

  function handleLogout() {
    setOpen(false);
    signOut();
    navigate("/login", { replace: true });
  }

  return (
    <header className="nav">
      <div className="nav-inner">
        <div className="nav-brand-group">
          <span className="nav-brand">Jodjum</span>
          <VersionTag />
        </div>

        <button
          className="nav-toggle"
          aria-label="Toggle menu"
          aria-expanded={open}
          onClick={() => setOpen((v) => !v)}
        >
          <span className="nav-toggle-bar" />
          <span className="nav-toggle-bar" />
          <span className="nav-toggle-bar" />
        </button>

        <nav className={`nav-menu ${open ? "is-open" : ""}`}>
          {LINKS.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              className={({ isActive }) => `nav-link ${isActive ? "active" : ""}`}
              onClick={() => setOpen(false)}
            >
              {link.label}
            </NavLink>
          ))}
          <NotificationToggle />
          <button className="btn btn-ghost nav-logout" onClick={handleLogout}>
            Log out
          </button>
        </nav>
      </div>
    </header>
  );
}
