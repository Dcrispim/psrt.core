import type { CSSProperties } from 'react';
import { parseStyle } from './parseStyle';

export function styleFromJSON(raw: string): CSSProperties {
  const s = parseStyle(raw);
  const out: CSSProperties = {};
  if (typeof s.color === 'string') out.color = s.color;
  const bg = s.background ?? s.backGround;
  if (typeof bg === 'string') out.background = bg;
  if (typeof s['text-align'] === 'string') out.textAlign = s['text-align'] as CSSProperties['textAlign'];
  const fw = s.fontWeight ?? s['font-weight'];
  if (fw != null) out.fontWeight = String(fw) as CSSProperties['fontWeight'];
  if (typeof s['font-style'] === 'string') out.fontStyle = s['font-style'] as CSSProperties['fontStyle'];
  return out;
}

export const SNAP = 0.01;

export function snapCoord(v: number): number {
  return Math.round(v / SNAP) * SNAP;
}
