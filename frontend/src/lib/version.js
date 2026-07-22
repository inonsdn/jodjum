// App version for debugging — read from a static file (public/version.json)
// rather than baked in at build time, so it can be updated at deploy without
// needing a rebuild for the value itself to be current.
export async function getAppVersion() {
  try {
    const res = await fetch("/version.json", { cache: "no-store" });
    if (!res.ok) return null;
    const data = await res.json();
    return typeof data.version === "string" ? data.version : null;
  } catch {
    return null;
  }
}
