function decodeBase64Payload(payload: string): Uint8Array | null {
  try {
    const binary = atob(payload);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    return bytes;
  } catch {
    return null;
  }
}

/** Decodes a `data:` URI (base64 or percent-encoded) to a UTF-8 string. */
export function decodeDataUri(uri: string): string | null {
  const comma = uri.indexOf(',');
  if (comma < 0) return null;

  const meta = uri.slice(0, comma);
  const payload = uri.slice(comma + 1);

  if (meta.includes(';base64')) {
    const bytes = decodeBase64Payload(payload);
    if (!bytes) return null;
    return new TextDecoder('utf-8').decode(bytes);
  }

  try {
    return decodeURIComponent(payload);
  } catch {
    return payload;
  }
}

/** Decodes a base64 `data:` URI to raw bytes (e.g. for Blob preview without UTF-8 round-trip). */
export function decodeDataUriBytes(uri: string): Uint8Array | null {
  const comma = uri.indexOf(',');
  if (comma < 0) return null;

  const meta = uri.slice(0, comma);
  const payload = uri.slice(comma + 1);
  if (!meta.includes(';base64')) return null;

  return decodeBase64Payload(payload);
}
