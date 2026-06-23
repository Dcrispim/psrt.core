import { useCallback, useEffect, useRef, useState, type MouseEvent } from 'react';

function clamp(n: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, n));
}

function readStoredWidth(
  key: string,
  fallback: number,
  min: number,
  max: number,
): number {
  try {
    const raw = localStorage.getItem(key);
    if (raw == null) return fallback;
    const n = Number.parseInt(raw, 10);
    if (!Number.isFinite(n)) return fallback;
    return clamp(n, min, max);
  } catch {
    return fallback;
  }
}

export interface UseResizablePanelWidthOptions {
  storageKey: string;
  defaultWidth: number;
  minWidth: number;
  maxWidth: number | (() => number);
}

export function useResizablePanelWidth({
  storageKey,
  defaultWidth,
  minWidth,
  maxWidth,
}: UseResizablePanelWidthOptions) {
  const getMax = useCallback(
    () => (typeof maxWidth === 'function' ? maxWidth() : maxWidth),
    [maxWidth],
  );

  const [width, setWidth] = useState(() =>
    readStoredWidth(storageKey, defaultWidth, minWidth, getMax()),
  );
  const [isResizing, setIsResizing] = useState(false);
  const widthRef = useRef(width);
  widthRef.current = width;

  const onResizeMouseDown = useCallback(
    (e: MouseEvent<HTMLDivElement>) => {
      if (e.button !== 0) return;
      e.preventDefault();
      const startX = e.clientX;
      const startWidth = widthRef.current;
      setIsResizing(true);

      const onMove = (ev: globalThis.MouseEvent) => {
        const delta = startX - ev.clientX;
        const next = clamp(startWidth + delta, minWidth, getMax());
        setWidth(next);
      };

      const onUp = () => {
        setIsResizing(false);
        document.removeEventListener('mousemove', onMove);
        document.removeEventListener('mouseup', onUp);
        document.body.style.cursor = '';
        document.body.style.userSelect = '';
        try {
          localStorage.setItem(storageKey, String(widthRef.current));
        } catch {
          /* ignore quota / private mode */
        }
      };

      document.body.style.cursor = 'ew-resize';
      document.body.style.userSelect = 'none';
      document.addEventListener('mousemove', onMove);
      document.addEventListener('mouseup', onUp);
    },
    [storageKey, minWidth, getMax],
  );

  useEffect(() => {
    const onWinResize = () => {
      setWidth((w) => clamp(w, minWidth, getMax()));
    };
    window.addEventListener('resize', onWinResize);
    return () => window.removeEventListener('resize', onWinResize);
  }, [minWidth, getMax]);

  return { width, isResizing, onResizeMouseDown };
}
