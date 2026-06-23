/** Keys where 0 / 0px / 0% mean "not set" rather than an explicit zero. */
const OMIT_ZERO_KEYS = new Set([
  'height',
  'width',
  'padding',
  'paddingTop',
  'paddingRight',
  'paddingBottom',
  'paddingLeft',
  'borderWidth',
  'strokeWidth',
  'borderRadius',
  'letterSpacing',
  'wordSpacing',
  'lineHeight',
  'textIndent',
]);

function isZeroLikeCSSValue(value: string): boolean {
  const s = value.trim().toLowerCase();
  if (
    s === '0' ||
    s === '0%' ||
    s === '0px' ||
    s === '0em' ||
    s === '0rem' ||
    s === '0pt' ||
    s === '0cqh'
  ) {
    return true;
  }
  const n = Number.parseFloat(s);
  return !Number.isNaN(n) && n === 0;
}

/** Whether a parsed style value was effectively passed (not null/empty/zero). */
export function isPresentStyleValue(key: string, value: unknown): boolean {
  if (value === null || value === undefined) return false;
  if (typeof value === 'boolean') return value;
  if (typeof value === 'number') {
    if (OMIT_ZERO_KEYS.has(key)) return value !== 0;
    return true;
  }
  if (typeof value !== 'string') return false;
  const s = value.trim();
  if (s === '' || s === 'null') return false;
  if (OMIT_ZERO_KEYS.has(key) && isZeroLikeCSSValue(s)) return false;
  return true;
}

/** Strip absent / zero-like entries from a style object. */
export function filterPresentStyleProps(
  style: Record<string, unknown>,
): Record<string, unknown> {
  const out: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(style)) {
    if (isPresentStyleValue(key, value)) {
      out[key] = value;
    }
  }
  return out;
}
