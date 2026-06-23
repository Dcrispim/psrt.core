/** Minimum base64 payload length before replacing with a placeholder. */
const DEFAULT_MIN_LENGTH = 48;

const DATA_URI_BASE64_RE =
  /data:([a-zA-Z0-9+/.-]+(?:;[a-zA-Z0-9+/.=-]+)*);base64,([A-Za-z0-9+/=]+)/g;

/**
 * Replaces long `data:…;base64,…` payloads in source text with `colapsado`
 * so the code view stays readable.
 */
export function collapseBase64InSource(
  source: string,
  minLength = DEFAULT_MIN_LENGTH,
): string {
  return source.replace(DATA_URI_BASE64_RE, (match, mime: string, payload: string) => {
    if (payload.length < minLength) return match;
    return `data:${mime};base64,colapsado`;
  });
}
