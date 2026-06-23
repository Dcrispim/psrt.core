const GOOGLE_FONTS_HOSTS = new Set(['fonts.googleapis.com', 'fonts.google.com']);

export function isGoogleFontsCssUrl(url: string): boolean {
  try {
    const u = new URL(url.trim());
    const host = u.hostname.replace(/^www\./, '');
    if (!GOOGLE_FONTS_HOSTS.has(host)) return false;
    return u.pathname.includes('/css');
  } catch {
    return false;
  }
}

/** Decode `Roboto+Mono` → `Roboto Mono`. */
function decodeFamilyName(encoded: string): string {
  return decodeURIComponent(encoded.replace(/\+/g, ' ')).trim();
}

/**
 * Extract font-family names from Google Fonts CSS API URLs.
 * Supports css2 (`family=Roboto:wght@400`) and legacy (`family=Roboto:400,700|Open+Sans`).
 */
export function parseGoogleFontsFamilies(url: string): string[] {
  try {
    const u = new URL(url.trim());
    const families: string[] = [];

    const familyParams = u.searchParams.getAll('family');
    for (const param of familyParams) {
      const segments = param.includes('|') ? param.split('|') : [param];
      for (const segment of segments) {
        const raw = segment.split(':')[0]?.trim();
        if (raw) families.push(decodeFamilyName(raw));
      }
    }

    return [...new Set(families.filter((f) => f.length > 0))];
  } catch {
    return [];
  }
}

/** Primary family for `font-family` in PSRT JSON (first listed in URL). */
export function primaryGoogleFontFamily(url: string): string | null {
  const families = parseGoogleFontsFamilies(url);
  return families[0] ?? null;
}

export function normalizeFontUrlInput(input: string): string {
  return input.trim();
}

export function isLikelyFontUrl(url: string): boolean {
  const u = normalizeFontUrlInput(url);
  if (!u) return false;
  return /^https?:\/\//i.test(u);
}
