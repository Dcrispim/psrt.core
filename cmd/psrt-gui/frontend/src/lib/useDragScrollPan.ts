import { useCallback, useRef, useState } from 'react';

const PAN_SKIP_SELECTOR =
  '.text-block, .resize-handle, button, input, textarea, select, label, a, [role="slider"]';

function shouldSkipPanTarget(target: EventTarget | null): boolean {
  return target instanceof HTMLElement && Boolean(target.closest(PAN_SKIP_SELECTOR));
}

export function useDragScrollPan(enabled: boolean) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const dragRef = useRef<{
    pointerId: number;
    startX: number;
    startY: number;
    scrollLeft: number;
    scrollTop: number;
  } | null>(null);
  const [dragging, setDragging] = useState(false);

  const onPointerDown = useCallback(
    (e: React.PointerEvent) => {
      if (!enabled || e.button !== 0) return;
      if (shouldSkipPanTarget(e.target)) return;

      const el = scrollRef.current;
      if (!el) return;

      e.currentTarget.setPointerCapture(e.pointerId);
      dragRef.current = {
        pointerId: e.pointerId,
        startX: e.clientX,
        startY: e.clientY,
        scrollLeft: el.scrollLeft,
        scrollTop: el.scrollTop,
      };
      setDragging(true);
    },
    [enabled],
  );

  const onPointerMove = useCallback((e: React.PointerEvent) => {
    const d = dragRef.current;
    const el = scrollRef.current;
    if (!d || !el || d.pointerId !== e.pointerId) return;
    el.scrollLeft = d.scrollLeft - (e.clientX - d.startX);
    el.scrollTop = d.scrollTop - (e.clientY - d.startY);
  }, []);

  const endDrag = useCallback((e: React.PointerEvent) => {
    const d = dragRef.current;
    if (!d || d.pointerId !== e.pointerId) return;
    dragRef.current = null;
    setDragging(false);
    if (e.currentTarget.hasPointerCapture(e.pointerId)) {
      e.currentTarget.releasePointerCapture(e.pointerId);
    }
  }, []);

  return {
    scrollRef,
    dragging,
    panHandlers: {
      onPointerDown,
      onPointerMove,
      onPointerUp: endDrag,
      onPointerCancel: endDrag,
    },
  };
}
