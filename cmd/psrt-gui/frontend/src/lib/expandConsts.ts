/** Mirrors psrt.ExpandConsts — longer keys first to avoid prefix clashes. */
export function expandConsts(
  content: string,
  consts: Record<string, string> | undefined,
): string {
  if (!consts || !content) return content;
  const keys = Object.keys(consts).sort((a, b) => {
    if (b.length !== a.length) return b.length - a.length;
    return a.localeCompare(b);
  });
  let out = content;
  for (const k of keys) {
    out = out.split(`@${k}@`).join(consts[k]);
  }
  return out;
}

export function resolveAssetReference(
  raw: string,
  consts: Record<string, string> | undefined,
): string {
  return expandConsts(raw.trim(), consts);
}
