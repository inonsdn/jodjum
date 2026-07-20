// A native date input that opens the calendar on a single click/tap anywhere on
// the field (not just the little icon), and shows no placeholder text when empty.
//
// It stays a real <input type="date"> so mobile taps open the OS date picker
// natively; onClick also calls showPicker() so a desktop click anywhere opens it.
export default function DateField({ value, onChange }) {
  function openPicker(e) {
    const el = e.currentTarget;
    if (el && typeof el.showPicker === "function") {
      try {
        el.showPicker();
      } catch {
        // showPicker() can throw if already open / no activation; the native
        // click behaviour still opens the picker, so ignore.
      }
    }
  }

  return (
    <input
      className={`field-date ${value ? "" : "is-empty"}`}
      type="date"
      value={value}
      onChange={onChange}
      onClick={openPicker}
    />
  );
}
