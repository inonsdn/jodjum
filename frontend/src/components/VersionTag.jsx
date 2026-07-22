import { useEffect, useState } from "react";
import { getAppVersion } from "../lib/version.js";

// Small, unobtrusive build version tag for debugging — e.g. "v1.0.0-202607231844".
export default function VersionTag() {
  const [version, setVersion] = useState(null);

  useEffect(() => {
    let cancelled = false;
    getAppVersion().then((v) => {
      if (!cancelled) setVersion(v);
    });
    return () => {
      cancelled = true;
    };
  }, []);

  if (!version) return null;

  return <span className="nav-version">{version}</span>;
}
