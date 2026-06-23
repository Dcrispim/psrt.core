import { useCallback, useEffect, useRef, useState, type CSSProperties } from 'react';

const WHEEL_ZOOM_FACTOR = 1.08;
const DEFAULT_MIN_ZOOM = 0.25;
const DEFAULT_MAX_ZOOM = 2;

const PAN_SKIP_SELECTOR =
  '.text-block, .resize-handle, button, input, textarea, select, label, a, [role="slider"]';

function clamp(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value));
}

function shouldSkipPanTarget(target: EventTarget | null): boolean {
  return target instanceof HTMLElement && Boolean(target.closest(PAN_SKIP_SELECTOR));
}

export function useWebPreviewPan(
  enabled: boolean,
  zoom: number,
  onZoomChange: (zoom: number) => void,
  minZoom = DEFAULT_MIN_ZOOM,
  maxZoom = DEFAULT_MAX_ZOOM,
) {
  const viewportRef = useRef<HTMLDivElement>(null);
  const panRef = useRef({ x: 0, y: 0 });
  const zoomRef = useRef(zoom);
  const dragRef = useRef<{
    pointerId: number;
    startX: number;
    startY: number;
    panX: number;
    panY: number;
  } | null>(null);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [dragging, setDragging] = useState(false);

  panRef.current = pan;
  zoomRef.current = zoom;

  const resetPan = useCallback(() => {
    panRef.current = { x: 0, y: 0 };
    setPan({ x: 0, y: 0 });
  }, []);

  const zoomAtPoint = useCallback(
    (clientX: number, clientY: number, nextZoom: number) => {
      const el = viewportRef.current;
      if (!el) return;

      const rect = el.getBoundingClientRect();
      const ratio = nextZoom / zoomRef.current;
      const cx = clientX - rect.left - panRef.current.x;
      const cy = clientY - rect.top - panRef.current.y;
      const next = {
        x: panRef.current.x - cx * (ratio - 1),
        y: panRef.current.y - cy * (ratio - 1),
      };
      panRef.current = next;
      setPan(next);
      zoomRef.current = nextZoom;
      onZoomChange(nextZoom);
    },
    [onZoomChange],
  );

  useEffect(() => {
    const el = viewportRef.current;
    if (!el || !enabled) return;

    const blockNativeDrag = (e: DragEvent) => {
      e.preventDefault();
    };
    el.addEventListener('dragstart', blockNativeDrag);

    const onWheel = (e: WheelEvent) => {
      if (!(e.ctrlKey || e.metaKey)) return;
      e.preventDefault();
      const delta = e.deltaY < 0 ? WHEEL_ZOOM_FACTOR : 1 / WHEEL_ZOOM_FACTOR;
      const next = clamp(zoomRef.current * delta, minZoom, maxZoom);
      zoomAtPoint(e.clientX, e.clientY, next);
    };

    el.addEventListener('wheel', onWheel, { passive: false });
    return () => {
      el.removeEventListener('dragstart', blockNativeDrag);
      el.removeEventListener('wheel', onWheel);
    };
  }, [enabled, minZoom, maxZoom, zoomAtPoint]);

  const onPointerDown = useCallback(
    (e: React.PointerEvent) => {
      if (!enabled || e.button !== 0) return;
      if (shouldSkipPanTarget(e.target)) return;

      e.preventDefault();
      e.currentTarget.setPointerCapture(e.pointerId);
      dragRef.current = {
        pointerId: e.pointerId,
        startX: e.clientX,
        startY: e.clientY,
        panX: panRef.current.x,
        panY: panRef.current.y,
      };
      setDragging(true);
    },
    [enabled],
  );

  const onPointerMove = useCallback((e: React.PointerEvent) => {
    const d = dragRef.current;
    if (!d || d.pointerId !== e.pointerId) return;
    const next = {
      x: d.panX + (e.clientX - d.startX),
      y: d.panY + (e.clientY - d.startY),
    };
    panRef.current = next;
    setPan(next);
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

  const stageStyle: CSSProperties = {
    transform: `translate(${pan.x}px, ${pan.y}px)`,
  };

  return {
    viewportRef,
    dragging,
    resetPan,
    stageStyle,
    panHandlers: {
      onPointerDown,
      onPointerMove,
      onPointerUp: endDrag,
      onPointerCancel: endDrag,
    },
  };
}
