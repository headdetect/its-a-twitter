import React from "react";

/**
 * Is essentially a hook that assets/throws it should be ran once.
 *
 * @param {React.EffectCallback} effect The effect to run
 * @param {React.DependencyList} deps The dependencies to pass into the hook
 */
export default function useEffectOnce(
  effect,
  deps = null,
  throwOnRerun = false,
) {
  const trackCountRef = React.useRef(false);

  React.useEffect(() => {
    if (trackCountRef.current) {
      if (throwOnRerun) {
        throw Error("Hook ran multiple times");
      } else {
        console.warn("Hook ran multiple times");
        console.trace();
      }
    }

    trackCountRef.current = true;

    return effect();

    // Disabled because it should be updating at the same rate
    // as the normal useEffect should be
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, deps);
}
