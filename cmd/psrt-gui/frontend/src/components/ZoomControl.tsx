const DEFAULT_MIN = 0.25;
const DEFAULT_MAX = 2;
const STEP = 0.05;

interface ZoomControlProps {
  id: string;
  zoom: number;
  onZoomChange: (zoom: number) => void;
  min?: number;
  max?: number;
}

export function ZoomControl({
  id,
  zoom,
  onZoomChange,
  min = DEFAULT_MIN,
  max = DEFAULT_MAX,
}: ZoomControlProps) {
  return (
    <div className="web-zoom-control">
      <label className="web-zoom-label" htmlFor={id}>
        Zoom
      </label>
      <input
        id={id}
        type="range"
        min={min}
        max={max}
        step={STEP}
        value={zoom}
        onChange={(e) => onZoomChange(Number(e.target.value))}
      />
      <span className="web-zoom-value">{Math.round(zoom * 100)}%</span>
    </div>
  );
}
