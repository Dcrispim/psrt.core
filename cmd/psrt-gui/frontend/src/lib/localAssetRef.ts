/** Normalizes a PSRT local asset reference (file URL or filesystem path). */
export function normalizeFileRef(ref: string): string {
  let r = ref.trim();
  if (!r) return r;

  if (!/^file:/i.test(r)) {
    r = pathToFileUrl(r);
  }

  try {
    if (r.includes('%')) {
      r = decodeURI(r);
    }
  } catch {
    /* keep encoded */
  }

  r = r.replace(/\\/g, '/');

  const winDrive = r.match(/^(file:\/\/\/)([a-zA-Z])(:.*)$/);
  if (winDrive) {
    r = `${winDrive[1]}${winDrive[2].toUpperCase()}${winDrive[3]}`;
  }

  return r;
}

/** Converts an absolute filesystem path to a PSRT file:/// reference. */
export function pathToFileUrl(fsPath: string): string {
  const trimmed = fsPath.trim();
  if (!trimmed) return '';
  if (/^file:/i.test(trimmed)) {
    return normalizeFileRef(trimmed);
  }

  const p = trimmed.replace(/\\/g, '/');
  if (/^[a-zA-Z]:\//.test(p)) {
    return `file:///${p[0].toUpperCase()}${p.slice(1)}`;
  }
  if (p.startsWith('/')) {
    return `file://${p}`;
  }
  return `file:///${p}`;
}

function constPrefixMatches(url: string, prefix: string): boolean {
  if (!prefix) return false;
  if (url === prefix) return true;
  if (!url.startsWith(prefix)) return false;
  const next = url[prefix.length];
  return next === undefined || next === '/';
}

/**
 * Replaces the longest matching $CONSTS value prefix with @name@ when saving a local path.
 */
export function substituteConstPrefixInAssetRef(
  fileUrl: string,
  consts: Record<string, string> | undefined,
): string {
  const normalizedUrl = normalizeFileRef(fileUrl);
  if (!consts || !normalizedUrl) return normalizedUrl;

  const entries = Object.entries(consts).sort((a, b) => {
    const lenDiff = normalizeFileRef(b[1]).length - normalizeFileRef(a[1]).length;
    if (lenDiff !== 0) return lenDiff;
    return a[0].localeCompare(b[0]);
  });

  for (const [name, value] of entries) {
    const prefix = normalizeFileRef(value);
    if (!constPrefixMatches(normalizedUrl, prefix)) continue;
    let suffix = normalizedUrl.slice(prefix.length);
    if (suffix.startsWith('/')) suffix = suffix.slice(1);
    return `@${name}@${suffix}`;
  }

  return normalizedUrl;
}

export function buildPageImageRefFromLocalPath(
  fsPath: string,
  consts: Record<string, string> | undefined,
): string {
  return substituteConstPrefixInAssetRef(pathToFileUrl(fsPath), consts);
}

/** True for file:// URLs and filesystem paths (not http/data/psrt-asset). */
export function isLocalAssetRef(raw: string): boolean {
  const s = raw.trim();
  if (!s) return false;
  if (s.startsWith('data:') || s.startsWith('psrt-asset://')) return false;
  if (/^https?:\/\//i.test(s)) return false;
  if (/^file:/i.test(s)) return true;
  if (/^[a-zA-Z]:[\\/]/.test(s)) return true;
  if (s.includes('\\')) return true;
  if (s.includes('/')) return true;
  return false;
}
