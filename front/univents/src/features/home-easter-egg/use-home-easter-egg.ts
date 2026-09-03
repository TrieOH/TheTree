import { useEffect, useRef } from "react";

const THEME_CLASS = "univents-easter-egg";

function activateTheme() {
  const applyTheme = () => document.documentElement.classList.add(THEME_CLASS);

  if (document.startViewTransition) {
    document.startViewTransition(applyTheme);
  } else {
    applyTheme();
  }
}

export function useHomeEasterEgg(onComplete: () => void) {
  const clicks = useRef(0);

  useEffect(
    () => () => document.documentElement.classList.remove(THEME_CLASS),
    [],
  );

  return () => {
    window.getSelection()?.removeAllRanges();
    clicks.current += 1;

    if (clicks.current === 10) {
      activateTheme();
    } else if (clicks.current === 20) {
      onComplete();
    }
  };
}
