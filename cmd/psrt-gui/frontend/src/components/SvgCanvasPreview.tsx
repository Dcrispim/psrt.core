import { useCallback, useEffect, useRef, useState } from 'react';

const WHEEL_ZOOM_FACTOR = 1.08;

interface SvgCanvasPreviewProps {
  uri: string;
  zoom: number;
  onZoomChange: (zoom: number) => void;
  minZoom?: number;
  maxZoom?: number;
}

interface Pan {
  x: number;
  y: number;
}

function clamp(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value));
}

export function SvgCanvasPreview({
  uri,
  zoom,
  onZoomChange,
  minZoom = 0.25,
  maxZoom = 8,
}: SvgCanvasPreviewProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const imgRef = useRef<HTMLImageElement | null>(null);
  const fitScaleRef = useRef(1);
  const panRef = useRef<Pan>({ x: 0, y: 0 });
  const dragRef = useRef<{ pointerId: number; startX: number; startY: number; panX: number; panY: number } | null>(
    null,
  );
  const zoomRef = useRef(zoom);
  const skipSliderCenterRef = useRef(false);

  const [pan, setPan] = useState<Pan>({ x: 0, y: 0 });
  const [dragging, setDragging] = useState(false);
  const [ready, setReady] = useState(false);

  panRef.current = pan;
  zoomRef.current = zoom;

  const centerOnImage = useCallback((fit: number, z: number, img: HTMLImageElement, w: number, h: number) => {
    const scale = fit * z;
    const drawW = img.naturalWidth * scale;
    const drawH = img.naturalHeight * scale;
    const next = { x: (w - drawW) / 2, y: (h - drawH) / 2 };
    panRef.current = next;
    setPan(next);
  }, []);

  const draw = useCallback(() => {
    const container = containerRef.current;
    const canvas = canvasRef.current;
    const img = imgRef.current;
    if (!container || !canvas || !img?.naturalWidth) return;

    const cssW = container.clientWidth;
    const cssH = container.clientHeight;
    if (cssW <= 0 || cssH <= 0) return;

    const dpr = window.devicePixelRatio || 1;
    canvas.width = Math.round(cssW * dpr);
    canvas.height = Math.round(cssH * dpr);
    canvas.style.width = `${cssW}px`;
    canvas.style.height = `${cssH}px`;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, cssW, cssH);
    ctx.fillStyle = '#111';
    ctx.fillRect(0, 0, cssW, cssH);

    const scale = fitScaleRef.current * zoomRef.current;
    const { x, y } = panRef.current;
    ctx.drawImage(img, x, y, img.naturalWidth * scale, img.naturalHeight * scale);
  }, []);

  const fitToContainer = useCallback(() => {
    const container = containerRef.current;
    const img = imgRef.current;
    if (!container || !img?.naturalWidth) return;

    const w = container.clientWidth;
    const h = container.clientHeight;
    if (w <= 0 || h <= 0) return;

    const fit = Math.min(w / img.naturalWidth, h / img.naturalHeight) * 0.92;
    fitScaleRef.current = fit;
    centerOnImage(fit, zoomRef.current, img, w, h);
    draw();
  }, [centerOnImage, draw]);

  useEffect(() => {
    setReady(false);
    imgRef.current = null;
    panRef.current = { x: 0, y: 0 };
    setPan({ x: 0, y: 0 });

    const img = new Image();
    img.onload = () => {
      imgRef.current = img;
      setReady(true);
      fitToContainer();
    };
    img.onerror = () => {
      imgRef.current = null;
      setReady(false);
    };
    img.src = uri;

    return () => {
      img.onload = null;
      img.onerror = null;
    };
  }, [uri, fitToContainer]);

  useEffect(() => {
    if (!ready) return;

    if (skipSliderCenterRef.current) {
      skipSliderCenterRef.current = false;
      zoomRef.current = zoom;
      draw();
      return;
    }

    const img = imgRef.current;
    const container = containerRef.current;
    if (!img?.naturalWidth || !container) return;

    const w = container.clientWidth;
    const h = container.clientHeight;
    if (w <= 0 || h <= 0) return;

    const fit = fitScaleRef.current;
    const cx = w / 2;
    const cy = h / 2;
    const oldScale = fit * zoomRef.current;
    const newScale = fit * zoom;
    const worldX = (cx - panRef.current.x) / oldScale;
    const worldY = (cy - panRef.current.y) / oldScale;
    const next = {
      x: cx - worldX * newScale,
      y: cy - worldY * newScale,
    };
    panRef.current = next;
    setPan(next);
    zoomRef.current = zoom;
    draw();
  }, [zoom, ready, draw]);

  useEffect(() => {
    if (!ready) return;
    draw();
  }, [pan, ready, draw]);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    const ro = new ResizeObserver(() => {
      if (imgRef.current?.naturalWidth) fitToContainer();
    });
    ro.observe(container);
    return () => ro.disconnect();
  }, [fitToContainer, ready]);

  const zoomAtPoint = useCallback(
    (clientX: number, clientY: number, nextZoom: number) => {
      const container = containerRef.current;
      const img = imgRef.current;
      if (!container || !img?.naturalWidth) return;

      const rect = container.getBoundingClientRect();
      const px = clientX - rect.left;
      const py = clientY - rect.top;
      const fit = fitScaleRef.current;
      const oldScale = fit * zoomRef.current;
      const newScale = fit * nextZoom;

      const worldX = (px - panRef.current.x) / oldScale;
      const worldY = (py - panRef.current.y) / oldScale;

      const nextPan = {
        x: px - worldX * newScale,
        y: py - worldY * newScale,
      };
      panRef.current = nextPan;
      setPan(nextPan);
      skipSliderCenterRef.current = true;
      zoomRef.current = nextZoom;
      onZoomChange(nextZoom);
      draw();
    },
    [draw, onZoomChange],
  );

  const onWheel = useCallback(
    (e: React.WheelEvent) => {
      e.preventDefault();
      const delta = e.deltaY < 0 ? WHEEL_ZOOM_FACTOR : 1 / WHEEL_ZOOM_FACTOR;
      const next = clamp(zoomRef.current * delta, minZoom, maxZoom);
      zoomAtPoint(e.clientX, e.clientY, next);
    },
    [minZoom, maxZoom, zoomAtPoint],
  );

  const onPointerDown = useCallback((e: React.PointerEvent) => {
    if (e.button !== 0) return;
    e.currentTarget.setPointerCapture(e.pointerId);
    dragRef.current = {
      pointerId: e.pointerId,
      startX: e.clientX,
      startY: e.clientY,
      panX: panRef.current.x,
      panY: panRef.current.y,
    };
    setDragging(true);
  }, []);

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

  return (
    <div
      ref={containerRef}
      className={`svg-canvas-preview${dragging ? ' svg-canvas-preview--dragging' : ''}`}
      onWheel={onWheel}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={endDrag}
      onPointerCancel={endDrag}
      role="img"
      aria-label="Visualização SVG"
    >
      <canvas ref={canvasRef} className="svg-canvas-preview__canvas" />
      {!ready && <p className="svg-canvas-preview__loading">Carregando SVG…</p>}
    </div>
  );
}
