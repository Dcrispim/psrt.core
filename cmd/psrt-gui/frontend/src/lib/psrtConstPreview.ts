const PIPE = ' | ';

/** One line inside `$CONSTS` … `$ENDCONSTS` (matches Go `FormatConstPSRT`). */
export function formatConstDeclarationLine(name: string, value: string): string {
  const n = name.trim();
  return `@ ${n}${PIPE}${value}`;
}

/** Placeholder replaced at compile time. */
export function formatConstUsageToken(name: string): string {
  return `@${name.trim()}@`;
}

export interface ConstPreviewExamples {
  pageImageUrl: string;
  textBody: string;
  styleProperty: string;
  styleInlineFragment: string;
}

export interface ConstPreview {
  declaration: string;
  token: string;
  expandedValue: string;
  examples: ConstPreviewExamples;
}

export function buildConstPreview(name: string, value: string): ConstPreview | null {
  const n = name.trim();
  if (!n) return null;

  const token = formatConstUsageToken(n);
  const declaration = formatConstDeclarationLine(n, value);
  const expandedValue = value;

  const looksLikeJsonFragment =
    value.includes(':') && (value.trimStart().startsWith('"') || value.includes('":'));

  return {
    declaration,
    token,
    expandedValue,
    examples: {
      pageImageUrl: `$START pagina | {} | ${token}01.webp`,
      textBody: `Visite ${token} para mais informações`,
      styleProperty: looksLikeJsonFragment
        ? `{"backGround":"#fff",${token}}`
        : `{"color":"${token}","font-weight":"700"}`,
      styleInlineFragment: looksLikeJsonFragment
        ? `… após expansão, o JSON inclui: ${value || '—'}`
        : `… após expansão: "color": "${expandedValue || '—'}"`,
    },
  };
}

export function isValidConstName(name: string): boolean {
  const n = name.trim();
  return n.length > 0 && !/\s/.test(n) && !n.includes('|');
}
