import type { CSSProperties } from 'react';
import { parseStyle } from './parseStyle';
import { isPresentStyleValue } from './styleValue';

/** PSRT entry geometry — not style JSON; always from entry fields in web preview. */
const ENTRY_LAYOUT_KEYS = new Set([
  'position',
  'left',
  'top',
  'width',
]);

/** Flex layout from styleadapter when text-align / align-items are set (not in raw JSON). */
const FLEX_LAYOUT_KEYS = new Set([
  'display',
  'flexDirection',
  'justifyContent',
  'alignItems',
]);

/** Span layout from styleadapter for flex text blocks (display:block; width:100%). */
const SPAN_ALIGN_LAYOUT_KEYS = new Set(['display', 'width']);

const RAW_TO_ADAPTER_KEY: Record<string, string> = {
  background: 'backgroundColor',
  backGround: 'backgroundColor',
  'background-color': 'backgroundColor',
  'font-weight': 'fontWeight',
  fw: 'fontWeight',
  'text-align': 'textAlign',
  ta: 'textAlign',
  'font-size': 'fontSize',
  pd: 'padding',
  'border-width': 'borderWidth',
  bw: 'borderWidth',
  'border-color': 'borderColor',
  bc: 'borderColor',
  'border-radius': 'borderRadius',
  br: 'borderRadius',
  'text-shadow': 'textShadow',
  ts: 'textShadow',
  'box-shadow': 'boxShadow',
  bsh: 'boxShadow',
  'line-height': 'lineHeight',
  'letter-spacing': 'letterSpacing',
  'word-spacing': 'wordSpacing',
  'stroke-width': 'strokeWidth',
  sw: 'strokeWidth',
  'stroke-color': 'strokeColor',
  sc: 'strokeColor',
};

function toKebab(key: string): string {
  return key.replace(/([A-Z])/g, '-$1').toLowerCase().replace(/^-/, '');
}

/** True if this CSS property was declared in the PSRT style JSON string. */
export function isDeclaredInStyleRaw(cssKey: string, styleRaw: string): boolean {
  const raw = parseStyle(styleRaw);
  if (isPresentStyleValue(cssKey, raw[cssKey])) return true;
  const kebab = toKebab(cssKey);
  if (isPresentStyleValue(kebab, raw[kebab])) return true;
  for (const [rawKey, adapterKey] of Object.entries(RAW_TO_ADAPTER_KEY)) {
    if (adapterKey === cssKey && isPresentStyleValue(rawKey, raw[rawKey])) {
      return true;
    }
  }
  return false;
}

/** Keep only adapter CSS keys that exist in style JSON (no inferred layout). */
function hasAlignInStyleRaw(styleRaw: string): boolean {
  return (
    isDeclaredInStyleRaw('textAlign', styleRaw) ||
    isDeclaredInStyleRaw('alignItems', styleRaw)
  );
}

export function pickExplicitAdapterCSS(
  props: CSSProperties,
  styleRaw: string,
): CSSProperties {
  const out: CSSProperties = {};
  const keepFlexLayout = hasAlignInStyleRaw(styleRaw);
  for (const [key, value] of Object.entries(props)) {
    if (value === undefined || value === null || value === '') continue;
    if (ENTRY_LAYOUT_KEYS.has(key)) continue;
    if (
      keepFlexLayout &&
      (FLEX_LAYOUT_KEYS.has(key) || SPAN_ALIGN_LAYOUT_KEYS.has(key))
    ) {
      (out as Record<string, unknown>)[key] = value;
      continue;
    }
    if (!isDeclaredInStyleRaw(key, styleRaw)) continue;
    (out as Record<string, unknown>)[key] = value;
  }
  return out;
}
