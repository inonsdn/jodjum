import NavBar from "./NavBar.jsx";

// Shell for authenticated pages: the shared top nav plus the page content.
export default function Layout({ children }) {
  return (
    <div className="app-shell">
      <NavBar />
      {children}
    </div>
  );
}
