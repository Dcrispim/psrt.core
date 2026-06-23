import { styleStringValue } from './parseStyle';

/** Height % from style object (`height` / `Height`), or null if unset. */
export function heightPercentFromStyleObject(
  style: Record<string, unknown>,
): number | null {
  const hRaw = style.height ?? style.Height;
  if (hRaw === undefined || hRaw === null) return null;
  const s = styleStringValue(hRaw).replace(/%$/, '');
  const n = parseFloat(s);
  return Number.isFinite(n) && n > 0 ? n : null;
}

/** Resolve == block height: explicit field, style %, legacy >> textSize, or default. */
export function resolveMaskHeightPercent(
  height: number | undefined,
  style: Record<string, unknown>,
  textSizeFallback?: number,
): number {
  if (typeof height === 'number' && Number.isFinite(height) && height >= 0.5) {
    return height;
  }
  const fromStyle = heightPercentFromStyleObject(style);
  if (fromStyle !== null) return fromStyle;
  if (
    textSizeFallback !== undefined &&
    Number.isFinite(textSizeFallback) &&
    textSizeFallback >= 0.5
  ) {
    return textSizeFallback;
  }
  return 5;
}
