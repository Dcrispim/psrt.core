import type { CSSProperties } from 'react';
import { AdaptEntriesForWeb } from '@wails/go/main/GUIApp';
import { styleadapter } from '@wails/go/models';
import type { PSRTEntry } from '../types/types';
import { isPresentStyleValue } from './styleValue';

export type AdaptedWebStyles = {
  container: CSSProperties;
  text: CSSProperties;
  hasStroke: boolean;
  merged?: CSSProperties;
  hitArea?: CSSProperties;
};

function mapToCSSProperties(m: Record<string, string> | undefined): CSSProperties {
  if (!m) return {};
  const out: Record<string, string> = {};
  for (const [k, v] of Object.entries(m)) {
    if (v !== '' && isPresentStyleValue(k, v)) out[k] = v;
  }
  return out as CSSProperties;
}

export function adaptedToCSS(adapted: styleadapter.WebPreviewStyle): AdaptedWebStyles {
  return {
    container: mapToCSSProperties(adapted.container),
    text: mapToCSSProperties(adapted.text),
    hasStroke: adapted.hasStroke,
  };
}

export type WebEntryInput = {
  index: number;
  style: string;
  content: string;
  x: number;
  y: number;
  width: number;
  textSize: number;
  height?: number;
  isMask?: boolean;
};

export function entriesToAdaptInput(entries: PSRTEntry[]): WebEntryInput[] {
  return entries?.map((e) => {
    const isMask = e.maskHeight !== undefined && e.maskHeight > 0;
    return {
      index: e.index ?? 0,
      style:
        e.styleRaw && e.styleRaw.trim() !== ''
          ? e.styleRaw
          : JSON.stringify(e.style ?? {}),
      content: e.text ?? '',
      x: e.x,
      y: e.y,
      width: e.width,
      textSize: e.size,
      height: isMask ? e.maskHeight : undefined,
      isMask,
    };
  });
}

/** Adapts all page entries via Go styleadapter (percent → px, stroke → WebKit, etc.). */
export async function adaptEntriesForWeb(
  entries: PSRTEntry[],
  canvasW: number,
  canvasH: number,
  zoom: number,
): Promise<Map<number, AdaptedWebStyles>> {
  const map = new Map<number, AdaptedWebStyles>();
  if (entries.length === 0 || canvasW < 1 || canvasH < 1) return map;

  try {
    const inputs = entriesToAdaptInput(entries);
    const raw = await AdaptEntriesForWeb(
      JSON.stringify(inputs),
      canvasW,
      canvasH,
      zoom,
    );
    if (!Array.isArray(raw)) {
      return map;
    }
    raw.forEach((r, i) => {
      const item = styleadapter.WebPreviewStyle.createFrom(r);
      const idx = inputs[i]?.index ?? i;
      map.set(idx, adaptedToCSS(item));
    });
  } catch (err) {
    console.error('AdaptEntriesForWeb failed:', err);
  }
  return map;
}
