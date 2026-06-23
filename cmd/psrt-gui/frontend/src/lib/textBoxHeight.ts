import { snapCoord } from './applyStyle';

export function isMaskTextBlock(content: string): boolean {
  return content.trim() === '';
}

function stylePropIndex(props: { key: string; value: string }[], key: string): number {
  const want = key.trim().toLowerCase();
  return props.findIndex((p) => p.key.trim().toLowerCase() === want);
}

function heightPropIndex(props: { key: string; value: string }[]): number {
  return stylePropIndex(props, 'height');
}

function paddingPropIndex(props: { key: string; value: string }[]): number {
  return stylePropIndex(props, 'padding');
}

export function parseHeightPercent(value: string): number | null {
  const v = value.trim();
  if (!v.endsWith('%')) return null;
  const n = parseFloat(v.slice(0, -1));
  return Number.isFinite(n) && n >= 0 ? n : null;
}

/** Height % from style props (`height: N%`), or null if unset. */
export function getStyleHeightPercent(
  props: { key: string; value: string }[],
): number | null {
  const row = props[heightPropIndex(props)];
  if (!row) return null;
  return parseHeightPercent(row.value);
}

export function setStyleHeightPercent(
  props: { key: string; value: string }[],
  percent: number,
): { key: string; value: string }[] {
  const pct = `${snapCoord(Math.max(0, percent))}%`;
  const i = heightPropIndex(props);
  if (i >= 0) {
    const next = [...props];
    next[i] = { ...next[i], value: pct };
    return next;
  }
  return [...props, { key: 'height', value: pct }];
}

export function removeStyleHeight(
  props: { key: string; value: string }[],
): { key: string; value: string }[] {
  const i = heightPropIndex(props);
  if (i < 0) return props;
  return props.filter((_, idx) => idx !== i);
}

export function parsePaddingPercent(value: string): number | null {
  const v = value.trim();
  if (!v.endsWith('%')) return null;
  const n = parseFloat(v.slice(0, -1));
  return Number.isFinite(n) && n >= 0 ? n : null;
}

export function getStylePaddingPercent(
  props: { key: string; value: string }[],
): number | null {
  const row = props[paddingPropIndex(props)];
  if (!row) return null;
  return parsePaddingPercent(row.value);
}

export function setStylePaddingPercent(
  props: { key: string; value: string }[],
  percent: number,
): { key: string; value: string }[] {
  const pct = `${snapCoord(Math.max(0, percent))}%`;
  const i = paddingPropIndex(props);
  if (i >= 0) {
    const next = [...props];
    next[i] = { ...next[i], value: pct };
    return next;
  }
  return [...props, { key: 'padding', value: pct }];
}
