import { snapCoord } from './applyStyle';

function clamp(value: number, min: number, max: number): number {
  return Math.min(max, Math.max(min, value));
}

/**
 * Maps viewport/client coordinates to PSRT text position (% of image width/height).
 * Uses the rendered image container box (accounts for web zoom and pan transform).
 */
export function clientPointToImagePercent(
  clientX: number,
  clientY: number,
  imageContainer: HTMLElement,
): { x: number; y: number } | null {
  const rect = imageContainer.getBoundingClientRect();
  if (rect.width <= 0 || rect.height <= 0) return null;

  const xPct = ((clientX - rect.left) / rect.width) * 100;
  const yPct = ((clientY - rect.top) / rect.height) * 100;

  return {
    x: snapCoord(clamp(xPct, 0, 100)),
    y: snapCoord(clamp(yPct, 0, 100)),
  };
}
