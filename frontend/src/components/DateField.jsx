import { useState } from "react";

// A date input with NO native placeholder. Browsers (esp. Safari) paint a grey
// "dd/mm/yyyy" or today's date into an empty <input type="date">, which reads
// like pre-filled data. To avoid that entirely, the field renders as a plain,
// placeholder-less text box while it's empty and unfocused; it switches to a
// real date input on focus (for the picker) and stays one once a value is set.
export default function DateField({ value, onChange, className, id }) {
  const [type, setType] = useState(value ? "date" : "text");

  return (
    <input
      id={id}
      className={className}
      type={type}
      value={value}
      placeholder=""
      onFocus={() => setType("date")}
      onBlur={() => {
        if (!value) setType("text");
      }}
      onChange={onChange}
    />
  );
}
