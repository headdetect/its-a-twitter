import React from "react";

import "./FloatingLabelInput.css";

export default function FloatingLabelInput({ label, ...inputProps }) {
  const [isExpanded, setIsExpanded] = React.useState(
    Boolean(inputProps?.value),
  );

  const handleFocus = e => {
    setIsExpanded(true);
    inputProps?.onFocus?.(e);
  };

  const handleBlur = e => {
    if (!e.target.value) {
      setIsExpanded(false);
    }
    inputProps?.onBlur?.(e);
  };

  return (
    <div className={`input-group ${isExpanded ? "input-group-focused" : ""}`}>
      <label htmlFor={label.toLowerCase()}>{label}</label>
      <input
        id={label.toLowerCase()}
        {...inputProps}
        onFocus={handleFocus}
        onBlur={handleBlur}
      />
    </div>
  );
}
