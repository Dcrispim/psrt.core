import { visualapp } from '@wails/go/models';
import type {
  TextBlock,
  BlockKind,
  BorderRadius,
  HAlign,
  VAlign,
  Direction,
  WhiteSpace,
  Overflow,
  TextTransform,
} from '../components/editor/types';
import { ZERO_BORDER_RADIUS } from '../components/editor/types';
import { parseGoogleFontsFamilies } from './googleFontsUrl';
import { parseStyle, styleStringValue, toColorInput } from './parseStyle';
import { resolveMaskHeightPercent } from './maskHeight';
import type { PsrtDocument } from '../types/document';
import { isPresentStyleValue, filterPresentStyleProps } from './styleValue';

/** All style keys used for border radius (shorthand + longhands, kebab + camel). */
export const BORDER_RADIUS_STYLE_KEYS = [
  'border-radius',
  'borderRadius',
  'br',
  'border-top-left-radius',
  'borderTopLeftRadius',
  'border-top-right-radius',
  'borderTopRightRadius',
  'border-bottom-right-radius',
  'borderBottomRightRadius',
  'border-bottom-left-radius',
  'borderBottomLeftRadius',
] as const;

const BORDER_RADIUS_STYLE_KEY_SET = new Set<string>(BORDER_RADIUS_STYLE_KEYS);

/** Style keys owned by structured panel fields (not shown in CSS props table). */
export const STRUCTURED_STYLE_KEYS = new Set([
  'color',
  'background',
  'backGround',
  'background-color',
  'font-family',
  'fontFamily',
  'font-weight',
  'fontWeight',
  'fw',
  'font-size',
  'fontSize',
  'line-height',
  'lineHeight',
  'letter-spacing',
  'letterSpacing',
  'font-style',
  'fontStyle',
  'text-decoration',
  'textDecoration',
  'text-transform',
  'textTransform',
  'text-align',
  'textAlign',
  'ta',
  'align-items',
  'alignItems',
  'vertical-align',
  'verticalAlign',
  'direction',
  'white-space',
  'whiteSpace',
  'overflow',
  'text-overflow',
  'textOverflow',
  ...BORDER_RADIUS_STYLE_KEYS,
]);

/** Edited in Posição & Tamanho; hidden from the free-form CSS props list. */
export const POSITION_PANEL_PROP_KEYS = new Set(['height', 'padding']);

export function isUniformBorderRadius(r: BorderRadius): boolean {
  return (
    r.topLeft === r.topRight &&
    r.topRight === r.bottomRight &&
    r.bottomRight === r.bottomLeft
  );
}

const DEFAULT_FONT_FAMILIES = ['Inter', 'Roboto', 'system-ui', 'Georgia', 'monospace'];

function str(style: Record<string, unknown>, ...keys: string[]): string {
  for (const k of keys) {
    const v = style[k];
    if (v !== undefined && v !== null && v !== '') {
      return styleStringValue(v);
    }
  }
  return '';
}

function num(style: Record<string, unknown>, key: string, fallback: number): number {
  const camel = key.replace(/-([a-z])/g, (_, c: string) => c.toUpperCase());
  const raw = style[key] ?? style[camel];
  if (raw === undefined || raw === null || raw === '') return fallback;
  const s = styleStringValue(raw);
  const n = parseFloat(s.replace(/[a-z%]+$/i, ''));
  return Number.isFinite(n) ? n : fallback;
}

function parseWeight(style: Record<string, unknown>): { weight: number; bold: boolean } {
  const raw = str(style, 'font-weight', 'fontWeight', 'fw');
  if (!raw) return { weight: 400, bold: false };
  const n = parseInt(raw, 10);
  if (!Number.isFinite(n)) return { weight: 400, bold: false };
  return { weight: n >= 700 ? n : n, bold: n >= 700 };
}

function parseDecoration(style: Record<string, unknown>): { underline: boolean; strike: boolean } {
  const d = str(style, 'text-decoration', 'textDecoration').toLowerCase();
  return {
    underline: d.includes('underline'),
    strike: d.includes('line-through'),
  };
}

function parseHAlign(style: Record<string, unknown>): HAlign {
  const v = str(style, 'text-align', 'textAlign', 'ta');
  if (v === 'left' || v === 'center' || v === 'right' || v === 'justify') return v;
  return 'center';
}

function parseVAlign(style: Record<string, unknown>): VAlign {
  const v = str(style, 'align-items', 'alignItems', 'vertical-align', 'verticalAlign');
  if (v === 'flex-start' || v === 'start' || v === 'top') return 'flex-start';
  if (v === 'flex-end' || v === 'end' || v === 'bottom') return 'flex-end';
  if (v === 'center' || v === 'middle') return 'center';
  return 'center';
}

function parseDirection(style: Record<string, unknown>): Direction {
  const v = str(style, 'direction');
  return v === 'rtl' ? 'rtl' : 'ltr';
}

function parseWhiteSpace(style: Record<string, unknown>): WhiteSpace {
  const v = str(style, 'white-space', 'whiteSpace');
  if (v === 'nowrap' || v === 'pre' || v === 'pre-wrap') return v;
  return 'normal';
}

function parseOverflow(style: Record<string, unknown>): Overflow {
  const o = str(style, 'overflow');
  const te = str(style, 'text-overflow', 'textOverflow');
  if (te === 'ellipsis' || o === 'ellipsis') return 'ellipsis';
  if (o === 'hidden' || o === 'clip' || o === 'visible') return o;
  return 'visible';
}

function parseTransform(style: Record<string, unknown>): TextTransform {
  const v = str(style, 'text-transform', 'textTransform');
  if (v === 'uppercase' || v === 'lowercase' || v === 'capitalize') return v;
  return 'none';
}

function parsePxToken(token: string): number {
  const n = parseFloat(token.replace(/[a-z%]+$/i, ''));
  return Number.isFinite(n) ? n : 0;
}

function parsePxFromStyle(style: Record<string, unknown>, ...keys: string[]): number | null {
  const raw = str(style, ...keys);
  if (!raw) return null;
  const n = parsePxToken(raw);
  return Number.isFinite(n) ? n : null;
}

function parseBorderRadius(style: Record<string, unknown>): BorderRadius {
  const tl = parsePxFromStyle(
    style,
    'border-top-left-radius',
    'borderTopLeftRadius',
  );
  const tr = parsePxFromStyle(
    style,
    'border-top-right-radius',
    'borderTopRightRadius',
  );
  const br = parsePxFromStyle(
    style,
    'border-bottom-right-radius',
    'borderBottomRightRadius',
  );
  const bl = parsePxFromStyle(
    style,
    'border-bottom-left-radius',
    'borderBottomLeftRadius',
  );
  if (tl !== null || tr !== null || br !== null || bl !== null) {
    return {
      topLeft: tl ?? 0,
      topRight: tr ?? 0,
      bottomRight: br ?? 0,
      bottomLeft: bl ?? 0,
    };
  }

  const raw = str(style, 'border-radius', 'borderRadius', 'br');
  if (!raw) return { ...ZERO_BORDER_RADIUS };
  const parts = raw.split(/\s+/).filter(Boolean).map(parsePxToken);
  if (parts.length === 0) return { ...ZERO_BORDER_RADIUS };
  const [a, b = a, c = a, d = b] = parts;
  if (parts.length === 1) {
    return { topLeft: a, topRight: a, bottomRight: a, bottomLeft: a };
  }
  if (parts.length === 2) {
    return { topLeft: a, topRight: b, bottomRight: a, bottomLeft: b };
  }
  if (parts.length === 3) {
    return { topLeft: a, topRight: b, bottomRight: c, bottomLeft: b };
  }
  return { topLeft: a, topRight: b, bottomRight: c, bottomLeft: d };
}

function borderRadiusStyleEntries(r: BorderRadius): Record<string, string> | null {
  const { topLeft, topRight, bottomRight, bottomLeft } = r;
  if (topLeft === 0 && topRight === 0 && bottomRight === 0 && bottomLeft === 0) {
    return null;
  }
  if (
    topLeft === topRight &&
    topRight === bottomRight &&
    bottomRight === bottomLeft
  ) {
    return { 'border-radius': `${topLeft}px` };
  }
  return {
    'border-top-left-radius': `${topLeft}px`,
    'border-top-right-radius': `${topRight}px`,
    'border-bottom-right-radius': `${bottomRight}px`,
    'border-bottom-left-radius': `${bottomLeft}px`,
  };
}

function borderRadiusEqual(a: BorderRadius, b: BorderRadius): boolean {
  return (
    a.topLeft === b.topLeft &&
    a.topRight === b.topRight &&
    a.bottomRight === b.bottomRight &&
    a.bottomLeft === b.bottomLeft
  );
}

function blockName(text: visualapp.TextDetail): string {
  const preview = (text.content ?? '').trim().slice(0, 40);
  return preview || `Texto ${text.index}`;
}

export function fontUrlLabel(url: string): string {
  try {
    const path = new URL(url).pathname;
    const file = path.split('/').pop() ?? url;
    const base = file.replace(/\.(woff2?|ttf|otf)$/i, '');
    const match = base.match(/([a-z]+)-latin/i);
    if (match) return match[1].charAt(0).toUpperCase() + match[1].slice(1);
    return base.slice(0, 24) || url;
  } catch {
    return url.slice(0, 24);
  }
}

export function buildFontSelectOptions(fontUrls: string[]): { value: string; label: string }[] {
  const seen = new Set<string>();
  const options: { value: string; label: string }[] = [];
  for (const f of DEFAULT_FONT_FAMILIES) {
    if (!seen.has(f)) {
      seen.add(f);
      options.push({ value: f, label: f });
    }
  }
  for (const url of fontUrls) {
    const googleFamilies = parseGoogleFontsFamilies(url);
    if (googleFamilies.length > 0) {
      for (const family of googleFamilies) {
        if (!seen.has(family)) {
          seen.add(family);
          options.push({ value: family, label: `${family} (Google)` });
        }
      }
      continue;
    }
    const label = fontUrlLabel(url);
    const value = label;
    if (!seen.has(value)) {
      seen.add(value);
      options.push({ value, label: `${label} (font)` });
    }
  }
  return options;
}

function sharedStyleToBlockFields(
  style: Record<string, unknown>,
  index: number,
  kind: BlockKind,
  x: number,
  y: number,
  width: number,
  extra: Partial<TextBlock>,
): TextBlock {
  const { weight, bold } = parseWeight(style);
  const deco = parseDecoration(style);
  const bgRaw = style.background ?? style.backGround ?? style['background-color'];
  const hasColor =
    style.color !== undefined &&
    style.color !== null &&
    styleStringValue(style.color) !== '';
  const hasBackground =
    bgRaw !== undefined && bgRaw !== null && styleStringValue(bgRaw) !== '';
  const hasFontSize =
    style['font-size'] !== undefined || style.fontSize !== undefined;
  const fontSizePx = hasFontSize ? num(style, 'font-size', 16) : 16;

  const props: { key: string; value: string }[] = [];
  const sourceStyle: { key: string; value: string }[] = [];
  for (const [key, val] of Object.entries(style)) {
    if (val === undefined || val === null || val === '') continue;
    const sv = styleStringValue(val);
    sourceStyle.push({ key, value: sv });
    if (STRUCTURED_STYLE_KEYS.has(key)) continue;
    props.push({ key, value: sv });
  }

  return {
    id: String(index),
    name: kind === 'mask' ? `Cobertura #${index}` : `Texto #${index}`,
    kind,
    x,
    y,
    width,
    size: extra.size ?? 4,
    height: extra.height ?? 5,
    content: extra.content ?? '',
    imageRef: extra.imageRef,
    bgImage: extra.bgImage,
    color: hasColor ? toColorInput(style.color) : '#000000',
    colorSet: hasColor,
    background: hasBackground ? toColorInput(bgRaw) : '#000000',
    backgroundSet: hasBackground,
    font: {
      family: str(style, 'font-family', 'fontFamily') || 'Inter',
      weight: bold ? Math.max(weight, 700) : weight < 700 ? weight : 400,
      size: fontSizePx,
      sizeOverride: hasFontSize,
      lineHeight: num(style, 'line-height', 1.4),
      letterSpacing: num(style, 'letter-spacing', 0),
      italic: str(style, 'font-style', 'fontStyle') === 'italic',
      underline: deco.underline,
      strike: deco.strike,
      bold,
      transform: parseTransform(style),
    },
    align: {
      horizontal: parseHAlign(style),
      vertical: parseVAlign(style),
      direction: parseDirection(style),
      whiteSpace: parseWhiteSpace(style),
      overflow: parseOverflow(style),
    },
    borderRadius: parseBorderRadius(style),
    props,
    sourceStyle,
  };
}

function decorationValue(underline: boolean, strike: boolean): string | null {
  const parts: string[] = [];
  if (underline) parts.push('underline');
  if (strike) parts.push('line-through');
  if (parts.length === 0) return null;
  return parts.join(' ');
}

export function styleEntriesFromBlock(block: TextBlock): Record<string, string> {
  const out: Record<string, string> = {
    'font-family': block.font.family,
    'line-height': String(block.font.lineHeight),
    'letter-spacing': `${block.font.letterSpacing}em`,
    'text-align': block.align.horizontal,
    'align-items': block.align.vertical,
    direction: block.align.direction,
    'white-space': block.align.whiteSpace,
    overflow: block.align.overflow === 'ellipsis' ? 'hidden' : block.align.overflow,
    'text-transform': block.font.transform,
  };

  if (block.colorSet) out.color = block.color;
  if (block.backgroundSet) out.background = block.background;

  const w = block.font.bold ? 700 : block.font.weight;
  out['font-weight'] = String(w);

  if (block.font.sizeOverride && block.font.size > 0) {
    out['font-size'] = `${block.font.size}px`;
  }

  if (block.font.italic) {
    out['font-style'] = 'italic';
  }

  const deco = decorationValue(block.font.underline, block.font.strike);
  if (deco) {
    out['text-decoration'] = deco;
  }

  if (block.align.overflow === 'ellipsis') {
    out['text-overflow'] = 'ellipsis';
  }

  const brEntries = borderRadiusStyleEntries(block.borderRadius);
  if (brEntries) {
    Object.assign(out, brEntries);
  }

  for (const p of block.props) {
    const k = p.key.trim();
    if (!k || BORDER_RADIUS_STYLE_KEY_SET.has(k)) continue;
    out[k] = p.value;
  }

  return out;
}

function normStyleKey(key: string): string {
  return key.trim().toLowerCase();
}

function upsertSourceStyle(
  rows: { key: string; value: string }[],
  key: string,
  value: string,
): { key: string; value: string }[] {
  const nk = normStyleKey(key);
  const next = rows.filter((r) => normStyleKey(r.key) !== nk);
  if (value.trim() !== '') next.push({ key: key.trim(), value });
  return next;
}

/** Apply a single CSS key/value to panel block state (structured fields + props). */
export function applyStyleEntryToBlock(
  block: TextBlock,
  key: string,
  value: string,
): TextBlock {
  const k = normStyleKey(key);
  const v = value.trim();
  let next: TextBlock = {
    ...block,
    sourceStyle: upsertSourceStyle(block.sourceStyle ?? [], key, value),
  };

  switch (k) {
    case 'color':
      return { ...next, color: toColorInput(v || '#000000'), colorSet: v !== '' };
    case 'background':
    case 'background-color':
      return { ...next, background: toColorInput(v || '#000000'), backgroundSet: v !== '' };
    case 'font-family':
    case 'fontfamily':
      return { ...next, font: { ...next.font, family: v || 'Inter' } };
    case 'font-weight':
    case 'fontweight':
    case 'fw': {
      const n = parseInt(v, 10);
      const bold = Number.isFinite(n) && n >= 700;
      return {
        ...next,
        font: {
          ...next.font,
          bold,
          weight: Number.isFinite(n) ? n : next.font.weight,
        },
      };
    }
    case 'font-size':
    case 'fontsize': {
      const n = parseFloat(v.replace(/px$/i, ''));
      return {
        ...next,
        font: {
          ...next.font,
          size: Number.isFinite(n) ? n : next.font.size,
          sizeOverride: v !== '',
        },
      };
    }
    case 'line-height':
    case 'lineheight': {
      const lh = parseFloat(v);
      return {
        ...next,
        font: { ...next.font, lineHeight: Number.isFinite(lh) ? lh : next.font.lineHeight },
      };
    }
    case 'letter-spacing':
    case 'letterspacing':
      return {
        ...next,
        font: { ...next.font, letterSpacing: parseFloat(v.replace(/em$/i, '')) || 0 },
      };
    case 'font-style':
    case 'fontstyle':
      return { ...next, font: { ...next.font, italic: v === 'italic' } };
    case 'text-decoration':
    case 'textdecoration': {
      const d = v.toLowerCase();
      return {
        ...next,
        font: {
          ...next.font,
          underline: d.includes('underline'),
          strike: d.includes('line-through'),
        },
      };
    }
    case 'text-transform':
    case 'texttransform':
      return {
        ...next,
        font: {
          ...next.font,
          transform: (v as TextTransform) || 'none',
        },
      };
    case 'text-align':
    case 'textalign':
    case 'ta':
      return {
        ...next,
        align: {
          ...next.align,
          horizontal: (v as HAlign) || next.align.horizontal,
        },
      };
    case 'align-items':
    case 'alignitems':
    case 'vertical-align':
    case 'verticalalign':
      return { ...next, align: { ...next.align, vertical: parseVAlign({ [k]: v }) } };
    case 'direction':
      return { ...next, align: { ...next.align, direction: v === 'rtl' ? 'rtl' : 'ltr' } };
    case 'white-space':
    case 'whitespace':
      return { ...next, align: { ...next.align, whiteSpace: parseWhiteSpace({ [k]: v }) } };
    case 'overflow':
      return { ...next, align: { ...next.align, overflow: parseOverflow({ overflow: v }) } };
    case 'text-overflow':
    case 'textoverflow':
      return {
        ...next,
        align: {
          ...next.align,
          overflow: v === 'ellipsis' ? 'ellipsis' : next.align.overflow,
        },
      };
    default: {
      const props = [...next.props];
      const idx = props.findIndex((p) => normStyleKey(p.key) === k);
      if (idx >= 0) {
        if (v === '') props.splice(idx, 1);
        else props[idx] = { key: key.trim(), value };
      } else if (v !== '') {
        props.push({ key: key.trim(), value });
      }
      return { ...next, props };
    }
  }
}

/** All effective style entries (panel fields + custom props), unfiltered. */
export function allStyleRowsFromBlock(block: TextBlock): { key: string; value: string }[] {
  return Object.entries(styleEntriesFromBlock(block))
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([key, value]) => ({ key, value }));
}

export function maskDetailToBlock(mask: visualapp.MaskDetail): TextBlock {
  const style = parseStyle(mask.style ?? '{}');
  return sharedStyleToBlockFields(style, mask.index, 'mask', mask.x, mask.y, mask.width, {
    height: resolveMaskHeightPercent(mask.height, style),
    content: '',
    imageRef: mask.imageRef || undefined,
  });
}

export function textDetailToBlock(text: visualapp.TextDetail): TextBlock {
  const style = parseStyle(text.style ?? '{}');
  const block = sharedStyleToBlockFields(
    style,
    text.index,
    'text',
    text.x,
    text.y,
    text.width,
    {
      size: text.textSize,
      content: text.content ?? '',
      imageRef: text.imageRef || undefined,
    },
  );
  block.name = blockName(text);
  return block;
}

/** Patch one CSS key in the document — does not touch other style keys. */
export function singleStylePropPatch(
  key: string,
  value: string | null,
): Pick<Partial<visualapp.TextPatch>, 'styleSet' | 'styleRemove'> {
  const k = key.trim();
  if (!k) return {};
  if (value != null && value !== '' && isPresentStyleValue(k, value)) {
    return { styleSet: { [k]: value } };
  }
  return { styleRemove: [k] };
}

function propsByKey(rows: { key: string; value: string }[]): Map<string, string> {
  const m = new Map<string, string>();
  for (const p of rows) {
    const k = p.key.trim();
    if (k) m.set(k, p.value);
  }
  return m;
}

function sanitizeStyleField(style: string | Record<string, unknown>): string {
  const parsed =
    typeof style === 'string' ? parseStyle(style) : (style as Record<string, unknown>);
  return JSON.stringify(filterPresentStyleProps(parsed));
}

/** Remove absent / zero-like style keys from the document before persisting to disk. */
export function sanitizeDocumentStylesForSave(doc: PsrtDocument): PsrtDocument {
  return {
    ...doc,
    pages: doc.pages.map((page) => ({
      ...page,
      style: sanitizeStyleField(page.style),
      texts: page.texts.map((t) => ({ ...t, style: sanitizeStyleField(t.style) })),
      masks: page.masks?.map((m) => ({ ...m, style: sanitizeStyleField(m.style) })),
    })),
  };
}

/** Diff only free-form CSS props rows (sidebar list), never panel defaults. */
function propsStylePatch(
  prevProps: { key: string; value: string }[],
  nextProps: { key: string; value: string }[],
): Pick<Partial<visualapp.TextPatch>, 'styleSet' | 'styleRemove'> {
  const styleSet: Record<string, string> = {};
  const styleRemove: string[] = [];
  const prev = propsByKey(prevProps);
  const next = propsByKey(nextProps);
  const keys = new Set([...prev.keys(), ...next.keys()]);
  for (const key of keys) {
    const pv = prev.get(key);
    const nv = next.get(key);
    if (pv === nv) continue;
    if (nv === undefined || nv === '') {
      styleRemove.push(key);
    } else {
      styleSet[key] = nv;
    }
  }
  const out: Pick<Partial<visualapp.TextPatch>, 'styleSet' | 'styleRemove'> = {};
  if (Object.keys(styleSet).length > 0) out.styleSet = styleSet;
  if (styleRemove.length > 0) out.styleRemove = [...new Set(styleRemove)];
  return out;
}

function appendStructuralStylePatch(
  prev: TextBlock,
  next: TextBlock,
  styleSet: Record<string, string>,
  styleRemove: string[],
): void {
  if (prev.font.family !== next.font.family) {
    styleSet['font-family'] = next.font.family;
  }
  if (prev.font.lineHeight !== next.font.lineHeight) {
    styleSet['line-height'] = String(next.font.lineHeight);
  }
  if (prev.font.letterSpacing !== next.font.letterSpacing) {
    styleSet['letter-spacing'] = `${next.font.letterSpacing}em`;
  }
  if (prev.font.sizeOverride !== next.font.sizeOverride) {
    if (next.font.sizeOverride && next.font.size > 0) {
      styleSet['font-size'] = `${next.font.size}px`;
    } else {
      styleRemove.push('font-size', 'fontSize');
    }
  } else if (next.font.sizeOverride && prev.font.size !== next.font.size) {
    styleSet['font-size'] = `${next.font.size}px`;
  }

  if (prev.font.bold !== next.font.bold || prev.font.weight !== next.font.weight) {
    if (next.font.bold) {
      styleSet['font-weight'] = '700';
    } else {
      styleRemove.push('font-weight', 'fontWeight', 'fw');
      if (next.font.weight !== 400) {
        styleSet['font-weight'] = String(next.font.weight);
      }
    }
  }

  if (prev.font.italic !== next.font.italic) {
    if (next.font.italic) styleSet['font-style'] = 'italic';
    else styleSet['font-style'] = 'normal';
  }

  const prevDeco = decorationValue(prev.font.underline, prev.font.strike);
  const nextDeco = decorationValue(next.font.underline, next.font.strike);
  if (prevDeco !== nextDeco) {
    if (nextDeco) styleSet['text-decoration'] = nextDeco;
    else styleRemove.push('text-decoration', 'textDecoration');
  }

  if (prev.font.transform !== next.font.transform) {
    styleSet['text-transform'] = next.font.transform;
  }

  if (prev.align.horizontal !== next.align.horizontal) {
    styleSet['text-align'] = next.align.horizontal;
  }
  if (prev.align.vertical !== next.align.vertical) {
    styleSet['align-items'] = next.align.vertical;
  }
  if (prev.align.direction !== next.align.direction) {
    styleSet.direction = next.align.direction;
  }
  if (prev.align.whiteSpace !== next.align.whiteSpace) {
    styleSet['white-space'] = next.align.whiteSpace;
  }
  if (prev.align.overflow !== next.align.overflow) {
    styleSet.overflow = next.align.overflow === 'ellipsis' ? 'hidden' : next.align.overflow;
    if (next.align.overflow === 'ellipsis') {
      styleSet['text-overflow'] = 'ellipsis';
    } else {
      styleRemove.push('text-overflow', 'textOverflow');
    }
  }

  if (prev.colorSet !== next.colorSet) {
    if (next.colorSet) styleSet.color = next.color;
    else styleRemove.push('color');
  } else if (next.colorSet && prev.color !== next.color) {
    styleSet.color = next.color;
  }

  if (prev.backgroundSet !== next.backgroundSet) {
    if (next.backgroundSet) styleSet.background = next.background;
    else styleRemove.push('background', 'backGround', 'background-color');
  } else if (next.backgroundSet && prev.background !== next.background) {
    styleSet.background = next.background;
  }

  if (!borderRadiusEqual(prev.borderRadius, next.borderRadius)) {
    const brEntries = borderRadiusStyleEntries(next.borderRadius);
    if (brEntries) {
      for (const k of BORDER_RADIUS_STYLE_KEYS) styleRemove.push(k);
      Object.assign(styleSet, brEntries);
    } else {
      styleRemove.push(...BORDER_RADIUS_STYLE_KEYS);
    }
  }
}

export function blockPatchToTextPatch(
  prev: TextBlock,
  next: TextBlock,
): Partial<visualapp.TextPatch> {
  const patch: Partial<visualapp.TextPatch> = {};
  const styleSet: Record<string, string> = {};
  const styleRemove: string[] = [];

  if (prev.x !== next.x) patch.x = next.x;
  if (prev.y !== next.y) patch.y = next.y;
  if (prev.width !== next.width) patch.width = next.width;
  if (prev.size !== next.size) patch.textSize = next.size;
  if (prev.content !== next.content) patch.content = next.content;

  const propsDelta = propsStylePatch(prev.props, next.props);
  if (propsDelta.styleSet) Object.assign(styleSet, propsDelta.styleSet);
  if (propsDelta.styleRemove) styleRemove.push(...propsDelta.styleRemove);

  appendStructuralStylePatch(prev, next, styleSet, styleRemove);

  if (Object.keys(styleSet).length > 0) patch.styleSet = styleSet;
  if (styleRemove.length > 0) {
    patch.styleRemove = [...new Set(styleRemove)];
  }

  return patch;
}

export function isEmptyPatch(patch: Partial<visualapp.TextPatch>): boolean {
  return (
    patch.content === undefined &&
    patch.x === undefined &&
    patch.y === undefined &&
    patch.width === undefined &&
    patch.textSize === undefined &&
    patch.imageRef === undefined &&
    (!patch.styleSet || Object.keys(patch.styleSet).length === 0) &&
    (!patch.styleRemove || patch.styleRemove.length === 0)
  );
}

export function blockPatchToMaskPatch(
  prev: TextBlock,
  next: TextBlock,
): Partial<visualapp.MaskPatch> {
  const patch: Partial<visualapp.MaskPatch> = {};
  const styleSet: Record<string, string> = {};
  const styleRemove: string[] = [];

  if (prev.x !== next.x) patch.x = next.x;
  if (prev.y !== next.y) patch.y = next.y;
  if (prev.width !== next.width) patch.width = next.width;
  if (prev.height !== next.height) patch.height = next.height;
  if (prev.imageRef !== next.imageRef) patch.imageRef = next.imageRef ?? '';

  const propsDelta = propsStylePatch(prev.props, next.props);
  if (propsDelta.styleSet) Object.assign(styleSet, propsDelta.styleSet);
  if (propsDelta.styleRemove) styleRemove.push(...propsDelta.styleRemove);

  if (prev.colorSet !== next.colorSet) {
    if (next.colorSet) styleSet.color = next.color;
    else styleRemove.push('color');
  } else if (next.colorSet && prev.color !== next.color) {
    styleSet.color = next.color;
  }
  if (prev.backgroundSet !== next.backgroundSet) {
    if (next.backgroundSet) styleSet.background = next.background;
    else styleRemove.push('background', 'backGround', 'background-color');
  } else if (next.backgroundSet && prev.background !== next.background) {
    styleSet.background = next.background;
  }

  if (Object.keys(styleSet).length > 0) patch.styleSet = styleSet;
  if (styleRemove.length > 0) patch.styleRemove = [...new Set(styleRemove)];
  return patch;
}

export function isEmptyMaskPatch(patch: Partial<visualapp.MaskPatch>): boolean {
  return (
    patch.x === undefined &&
    patch.y === undefined &&
    patch.width === undefined &&
    patch.height === undefined &&
    patch.imageRef === undefined &&
    (!patch.styleSet || Object.keys(patch.styleSet).length === 0) &&
    (!patch.styleRemove || patch.styleRemove.length === 0)
  );
}

export type BlockEntryPatch = {
  text?: Partial<visualapp.TextPatch>;
  mask?: Partial<visualapp.MaskPatch>;
  kindChange?: BlockKind;
};

export function blockPatchToEntryPatch(prev: TextBlock, next: TextBlock): BlockEntryPatch {
  if (prev.kind !== next.kind) {
    const out: BlockEntryPatch = { kindChange: next.kind };
    if (next.kind === 'mask') {
      const mask = blockPatchToMaskPatch({ ...prev, kind: 'mask' }, next);
      if (!isEmptyMaskPatch(mask)) out.mask = mask;
    } else {
      const text = blockPatchToTextPatch({ ...prev, kind: 'text' }, next);
      if (!isEmptyPatch(text)) out.text = text;
    }
    return out;
  }
  if (next.kind === 'mask') {
    const mask = blockPatchToMaskPatch(prev, next);
    return isEmptyMaskPatch(mask) ? {} : { mask };
  }
  const text = blockPatchToTextPatch(prev, next);
  return isEmptyPatch(text) ? {} : { text };
}
