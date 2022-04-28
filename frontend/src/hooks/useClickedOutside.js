import React from "react";

const isTouchCapable = "ontouchstart" in window;

/**
 * Hook that triggers an event when the user clicks outside the specified
 * element ref
 *
 * @param {React.RefObject} elementRef
 * @param {(e: Event) => void} onOutsideClick
 * @param {boolean} includeParent True if include parent as "outside" click
 */
export default function useClickedOutside(
  elementRef,
  onOutsideClick,
  includeParent = false,
) {
  if (!elementRef) {
    throw new Error("ElementRef must be specified");
  }

  React.useEffect(() => {
    const handle = e => {
      const element = includeParent
        ? elementRef.current
        : elementRef.current.parentElement;

      if (element && !element.contains(e.target) && onOutsideClick) {
        onOutsideClick(e);
      }
    };

    const unload = () => {
      document.removeEventListener("touchstart", handle);
      document.removeEventListener("mousedown", handle);
    };

    if (elementRef?.current && onOutsideClick) {
      if (isTouchCapable) {
        document.addEventListener("touchstart", handle);
      } else {
        document.addEventListener("mousedown", handle);
      }
    } else {
      unload(); // Unload if the element isn't a thing anymore //
    }

    return unload;
  }, [elementRef, onOutsideClick, includeParent]);
}
