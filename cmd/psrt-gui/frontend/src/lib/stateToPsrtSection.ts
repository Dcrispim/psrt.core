import type { CSSProperties } from 'react';
import { visualapp } from '@wails/go/models';
import { parseStyle } from './parseStyle';
import { filterPresentStyleProps, isPresentStyleValue } from './styleValue';
import type { PSRTEntry, PSRTSection } from '../types/types';

function styleJSONToCSS(raw: string): CSSProperties {
  const parsed = filterPresentStyleProps(parseStyle(raw));
  const out: Record<string, string | number> = {};
  for (const [key, value] of Object.entries(parsed)) {
    const camel = key.replace(/-([a-z])/g, (_, c: string) => c.toUpperCase());
    if (typeof value === 'string' || typeof value === 'number') {
      if (isPresentStyleValue(camel, value)) out[camel] = value;
      if (isPresentStyleValue(key, value)) out[key] = value;
    }
  }
  return out as CSSProperties;
}

export function stateToPsrtSection(
  state: visualapp.UIState,
  pageImageUri: string | null,
): PSRTSection | null {
  if (!state.page) return null;

  const textEntries: PSRTEntry[] = (state.texts ?? []).map((t) => ({
    index: t.index,
    x: t.x,
    y: t.y,
    size: t.textSize,
    width: t.width,
    text: t.content ?? '',
    style: styleJSONToCSS(t.style),
    styleRaw: t.style ?? '{}',
  }));

  const maskEntries: PSRTEntry[] = (state.masks ?? []).map((m) => {
    const style = styleJSONToCSS(m.style);
    return {
      index: m.index,
      x: m.x,
      y: m.y,
      size: 0,
      width: m.width,
      text: '',
      maskHeight: m.height,
      style: { ...style, height: `${m.height}%` },
      styleRaw: m.style ?? '{}',
    };
  });

  const entries = [...textEntries, ...maskEntries].sort(
    (a, b) => (a.index ?? 0) - (b.index ?? 0),
  );

  return {
    title: state.page.name,
    pageLink: pageImageUri ?? undefined,
    pageStyle: styleJSONToCSS(state.page.style),
    entries,
  };
}
