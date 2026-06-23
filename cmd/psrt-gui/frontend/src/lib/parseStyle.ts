export function parseStyle(raw: string): Record<string, unknown> {
  try {
    return JSON.parse(raw || '{}') as Record<string, unknown>;
  } catch {
    return {};
  }
}

/** Strips JSON string quotes from style values stored incorrectly. */
export function styleStringValue(v: unknown): string {
  if (v === null || v === undefined) return '';
  if (typeof v === 'number' || typeof v === 'boolean') return String(v);
  if (typeof v !== 'string') return '';
  let s = v.trim();
  if (
    (s.startsWith('"') && s.endsWith('"')) ||
    (s.startsWith("'") && s.endsWith("'"))
  ) {
    try {
      const parsed: unknown = JSON.parse(s);
      if (typeof parsed === 'string') return parsed;
    } catch {
      return s.slice(1, -1);
    }
  }
  return s;
}

export function toColorInput(v: unknown): string {
  const s = styleStringValue(v);
  if (s.startsWith('#')) {
    if (s.length >= 7) return s.slice(0, 7);
    if (s.length === 4) {
      return `#${s[1]}${s[1]}${s[2]}${s[2]}${s[3]}${s[3]}`;
    }
  }
  return '#ffffff';
}
