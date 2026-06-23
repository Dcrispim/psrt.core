import { useEffect, useMemo, useState } from 'react';
import type { CSSProperties } from 'react';
import type { PSRTEntry } from '../types/types';
import {
  adaptEntriesForWeb,
  type AdaptedWebStyles,
} from './adaptWebStyle';
import { isDeclaredInStyleRaw, pickExplicitAdapterCSS } from './explicitWebStyles';
import { estimateTextBoxHeightPct, textFontSizePx } from './psrtGeometry';
import { isPresentStyleValue } from './styleValue';

export type { AdaptedWebStyles };

function adapterContainerCSS(container: CSSProperties | undefined): CSSProperties {
  if (!container) return {};
  const out: Record<string, string> = {};
  for (const [k, v] of Object.entries(container)) {
    if (typeof v !== 'string' && typeof v !== 'number') continue;
    const s = String(v);
    if (s !== '' && isPresentStyleValue(k, s)) out[k] = s;
  }
  return out as CSSProperties;
}

const empty = (): AdaptedWebStyles => ({
  container: {},
  text: {},
  hasStroke: false,
});

export function useAdaptedEntryStyles(
  entries: PSRTEntry[],
  canvasW: number,
  canvasH: number,
  zoom: number,
): Map<number, AdaptedWebStyles> {
  const [map, setMap] = useState<Map<number, AdaptedWebStyles>>(new Map());

  const entriesKey = useMemo(
    () =>
      JSON.stringify(
        entries.map((e) => ({
          i: e.index,
          x: e.x,
          y: e.y,
          w: e.width,
          s: e.size,
          h: e.maskHeight,
          t: e.text,
          styleRaw: e.styleRaw,
        })),
      ),
    [entries],
  );

  useEffect(() => {
    if (canvasW < 1 || canvasH < 1 || entries.length === 0) {
      return;
    }
    let cancelled = false;
    void adaptEntriesForWeb(entries, canvasW, canvasH, zoom).then((next) => {
      if (!cancelled) setMap(next);
    });
    return () => {
      cancelled = true;
    };
  }, [entriesKey, canvasW, canvasH, zoom, entries.length]);

  return map;
}

export function resolveEntryStyle(
  entry: PSRTEntry,
  adapted: AdaptedWebStyles | undefined,
  metrics: { refWidth: number; refHeight: number; zoom: number } | undefined,
): AdaptedWebStyles {
  const base = adapted ?? empty();
  const styleRaw = entry.styleRaw ?? '{}';
  const isMask = entry.maskHeight != null && entry.maskHeight >= 0.5;

  // Text >> : full adapter box (background, padding, flex, computed height).
  // Mask == : explicit style keys only; height from 4th coord (maskHeight).
  const boxStyle = isMask
    ? pickExplicitAdapterCSS(base.container, styleRaw)
    : adapterContainerCSS(base.container);
  const textStyle = pickExplicitAdapterCSS(base.text, styleRaw);

  const fontPx =
    metrics && metrics.refWidth > 0 && metrics.refHeight > 0
      ? textFontSizePx(entry.size, metrics.refWidth, metrics.refHeight, metrics.zoom)
      : 0;

  const fontSize = isDeclaredInStyleRaw('fontSize', styleRaw)
    ? (textStyle.fontSize ?? boxStyle.fontSize)
    : fontPx > 0
      ? `${fontPx}px`
      : undefined;

  const layout: CSSProperties = {
    position: 'absolute',
    left: `${entry.x}%`,
    top: `${entry.y}%`,
    width: `${entry.width}%`,
    boxSizing: 'border-box',
    zIndex: entry.index ?? 0,
    ...(fontSize ? { fontSize } : {}),
    ...(isMask ? { height: `${entry.maskHeight}%` } : {}),
  };

  const rawH =
    typeof base.container?.height === 'string' ? base.container.height : undefined;
  const adapterH =
    rawH && isPresentStyleValue('height', rawH) ? rawH : undefined;
  let hitHeight: string | undefined;
  if (isMask) {
    hitHeight = `${entry.maskHeight}%`;
  } else if (adapterH) {
    hitHeight = adapterH;
  } else if (
    typeof boxStyle.height === 'string' &&
    isPresentStyleValue('height', boxStyle.height)
  ) {
    hitHeight = boxStyle.height;
  } else if (metrics && metrics.refHeight > 0) {
    hitHeight = estimateTextBoxHeightPct(
      {
        width: entry.width,
        size: entry.size,
        text: entry.text,
        styleRaw: entry.styleRaw ?? '{}',
      },
      metrics,
    );
  }

  const hitArea: CSSProperties = {
    position: 'absolute',
    left: `${entry.x}%`,
    top: `${entry.y}%`,
    width: `${entry.width}%`,
    boxSizing: 'border-box',
    zIndex: entry.index ?? 0,
    height: hitHeight,
    minHeight: hitHeight,
  };

  const merged: CSSProperties = {
    ...boxStyle,
    ...textStyle,
    ...layout,
  };

  return {
    container: { ...boxStyle, ...layout },
    text: { ...textStyle, ...(fontSize ? { fontSize } : {}) },
    hasStroke: base.hasStroke,
    merged,
    hitArea,
  };
}
