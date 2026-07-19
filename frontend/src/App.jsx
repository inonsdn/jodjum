import { Navigate, Route, Routes } from "react-router-dom";
import { useAuth } from "./lib/auth.jsx";
import LoginPage from "./pages/LoginPage.jsx";
import RegisterPage from "./pages/RegisterPage.jsx";
import ThingsPage from "./pages/ThingsPage.jsx";
import Layout from "./components/Layout.jsx";

// Redirect to the login page when there is no active session.
function RequireAuth({ children }) {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? children : <Navigate to="/login" replace />;
}

// Keep signed-in users out of the auth pages.
function RedirectIfAuthed({ children }) {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <Navigate to="/things" replace /> : children;
}

export default function App() {
  return (
    <Routes>
      <Route
        path="/login"
        element={
          <RedirectIfAuthed>
            <LoginPage />
          </RedirectIfAuthed>
        }
      />
      <Route
        path="/register"
        element={
          <RedirectIfAuthed>
            <RegisterPage />
          </RedirectIfAuthed>
        }
      />
      <Route
        path="/things"
        element={
          <RequireAuth>
            <Layout>
              <ThingsPage />
            </Layout>
          </RequireAuth>
        }
      />
      <Route path="*" element={<Navigate to="/things" replace />} />
    </Routes>
  );
}
