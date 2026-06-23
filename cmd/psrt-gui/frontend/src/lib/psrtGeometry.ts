/**
 * Mirrors psrt/coords.go and compilesvg (PSRT % semantics).
 * - x, width: % of canvas width
 * - y: % of canvas height
 * - textSize: % of min(canvas width, canvas height)
 */

export function textSizeBasisPx(canvasW: number, canvasH: number): number {
  const w = canvasW;
  const h = canvasH;
  if (w <= 0 && h <= 0) return 1;
  if (w <= 0) return h;
  if (h <= 0) return w;
  return Math.min(w, h);
}

export function textFontSizePx(
  textSizePct: number,
  canvasW: number,
  canvasH: number,
  zoom = 1,
): number {
  const px = (textSizePct / 100) * textSizeBasisPx(canvasW, canvasH) * zoom;
  return px < 1 ? 1 : px;
}

export function textBlockWidthPx(widthPct: number, canvasW: number): number {
  const w = Math.round((canvasW * widthPct) / 100);
  return w < 1 ? 1 : w;
}

const DEFAULT_LINE_HEIGHT = 1.2;
const CHARS_PER_EM = 0.48;

function estimateLineCount(content: string, widthPx: number, fontPx: number): number {
  const plain = content.replace(/\*\*|\*|_|~|\\/g, '');
  if (!plain.trim()) return 1;
  const charsPerLine = Math.max(1, Math.floor(widthPx / (fontPx * CHARS_PER_EM)));
  let total = 0;
  for (const part of plain.split('\n')) {
    const n = [...part.trim()].length;
    if (n === 0) continue;
    total += Math.ceil(n / charsPerLine);
  }
  return Math.max(1, total);
}

function parsePaddingPx(styleRaw: string, canvasW: number, canvasH: number, fontPx: number): number {
  try {
    const style = JSON.parse(styleRaw || '{}') as Record<string, unknown>;
    const pad = style.padding;
    if (typeof pad === 'number') return pad * 2;
    if (typeof pad !== 'string' || !pad.trim()) return 0;
    const s = pad.trim();
    if (s.endsWith('px')) {
      const n = Number.parseFloat(s);
      return Number.isNaN(n) ? 0 : n * 2;
    }
    if (s.endsWith('%')) {
      const n = Number.parseFloat(s);
      if (Number.isNaN(n)) return 0;
      const basis = Math.max(canvasW, canvasH);
      return (n / 100) * basis * 2;
    }
    const n = Number.parseFloat(s);
    if (!Number.isNaN(n)) return n * fontPx * 2;
  } catch {
    /* ignore */
  }
  return 0;
}

/** Mirrors compilesvg.TextBlockGeometry height → canvas % (for hit targets). */
export function estimateTextBoxHeightPct(
  entry: { width: number; size: number; text: string; styleRaw: string },
  metrics: { refWidth: number; refHeight: number; zoom: number },
): string {
  const { refWidth, refHeight, zoom } = metrics;
  if (refWidth < 1 || refHeight < 1) return '1%';

  const fontPx = textFontSizePx(entry.size, refWidth, refHeight, zoom);
  const linePx = fontPx * DEFAULT_LINE_HEIGHT;
  const widthPx = textBlockWidthPx(entry.width, refWidth) * zoom;
  const lines = estimateLineCount(entry.text, widthPx, fontPx);
  const padPx = parsePaddingPx(entry.styleRaw, refWidth, refHeight, fontPx);
  const heightPx = linePx * lines + padPx;
  const pct = (heightPx / refHeight) * 100;
  return `${Math.max(pct, 0.1)}%`;
}

export function pctToPx(valuePct: number, axisLength: number, zoom = 1): string {
  return `${(valuePct / 100) * axisLength * zoom}px`;
}
