const INDENT = '  ';

const VOID_TAGS = new Set([
  'area',
  'base',
  'br',
  'col',
  'embed',
  'hr',
  'img',
  'input',
  'link',
  'meta',
  'param',
  'source',
  'track',
  'wbr',
  'circle',
  'ellipse',
  'line',
  'path',
  'polygon',
  'polyline',
  'rect',
  'use',
  'image',
]);

const RAW_CONTENT_TAGS = new Set(['script', 'style']);

/** Pretty-prints HTML/SVG for the compiled preview source panel. */
export function formatMarkupSource(source: string): string {
  const trimmed = source.trim();
  if (!trimmed) return source;

  const compact = trimmed.replace(/>\s+</g, '><');
  const lines: string[] = [];
  let depth = 0;
  let i = 0;

  const pushLine = (text: string) => {
    if (text.length > 0) lines.push(text);
  };

  const pad = (level: number) => INDENT.repeat(Math.max(0, level));

  while (i < compact.length) {
    if (compact[i] !== '<') {
      const next = compact.indexOf('<', i);
      const text = (next < 0 ? compact.slice(i) : compact.slice(i, next)).trim();
      if (text) pushLine(pad(depth) + text);
      i = next < 0 ? compact.length : next;
      continue;
    }

    const tagRead = readTag(compact, i);
    if (!tagRead) break;
    const { tag, end } = tagRead;
    i = end;

    const info = parseTagInfo(tag);
    if (!info) {
      pushLine(pad(depth) + tag);
      continue;
    }

    if (info.kind === 'close') {
      depth = Math.max(0, depth - 1);
      pushLine(pad(depth) + tag);
      continue;
    }

    if (info.kind === 'self' || VOID_TAGS.has(info.name)) {
      pushLine(pad(depth) + formatOpeningTag(tag, pad(depth)));
      continue;
    }

    if (info.kind === 'open' && RAW_CONTENT_TAGS.has(info.name)) {
      const closeName = info.name;
      const closeRe = new RegExp(`</${closeName}\\s*>`, 'i');
      const closeMatch = closeRe.exec(compact.slice(i));
      const innerEnd = closeMatch ? i + closeMatch.index : compact.length;
      const inner = compact.slice(i, innerEnd);
      const closeTag = closeMatch ? closeMatch[0] : '';

      pushLine(pad(depth) + formatOpeningTag(tag, pad(depth)));
      const innerTrimmed = inner.trim();
      if (innerTrimmed) {
        for (const innerLine of innerTrimmed.split('\n')) {
          const line = innerLine.trimEnd();
          if (line) pushLine(pad(depth + 1) + line);
        }
      }
      if (closeTag) {
        pushLine(pad(depth) + closeTag);
        i = innerEnd + closeTag.length;
      }
      continue;
    }

    if (info.kind === 'open') {
      pushLine(pad(depth) + formatOpeningTag(tag, pad(depth)));
      depth += 1;
      continue;
    }

    pushLine(pad(depth) + tag);
  }

  return lines.join('\n');
}

function readTag(source: string, start: number): { tag: string; end: number } | null {
  if (source[start] !== '<') return null;

  let inQuote: '"' | "'" | null = null;
  for (let i = start + 1; i < source.length; i++) {
    const c = source[i];
    if (inQuote) {
      if (c === inQuote) inQuote = null;
    } else if (c === '"' || c === "'") {
      inQuote = c;
    } else if (c === '>') {
      return { tag: source.slice(start, i + 1), end: i + 1 };
    }
  }

  return { tag: source.slice(start), end: source.length };
}

type TagKind = 'open' | 'close' | 'self' | 'other';

function parseTagInfo(tag: string): { kind: TagKind; name: string } | null {
  if (tag.startsWith('<!--') || tag.startsWith('<!') || tag.startsWith('<?')) {
    return { kind: 'other', name: '' };
  }

  const match = tag.match(/^<\/?([a-zA-Z][\w:.-]*)/);
  if (!match) return null;

  const name = match[1].toLowerCase();
  if (tag.startsWith('</')) return { kind: 'close', name };
  if (tag.endsWith('/>') || tag.endsWith(' />')) return { kind: 'self', name };
  return { kind: 'open', name };
}

function formatOpeningTag(tag: string, baseIndent: string): string {
  const openMatch = tag.match(/^<([a-zA-Z][\w:.-]*)([\s\S]*?)(\s*\/?>)$/);
  if (!openMatch) return tag;

  const [, name, attrPart, closing] = openMatch;
  const attrs = splitAttributes(attrPart);
  const selfClose = closing.trimStart().startsWith('/');
  const oneLine = tag.length <= 100 || attrs.length <= 1;

  if (oneLine) return tag;

  const attrIndent = baseIndent + INDENT;
  const lines = [`<${name}`];
  for (const attr of attrs) {
    lines.push(`${attrIndent}${attr}`);
  }
  lines.push(`${baseIndent}${selfClose ? '/>' : '>'}`);
  return lines.join('\n');
}

function splitAttributes(attrPart: string): string[] {
  const attrs: string[] = [];
  let current = '';
  let inQuote: '"' | "'" | null = null;

  for (let i = 0; i < attrPart.length; i++) {
    const c = attrPart[i];
    if (inQuote) {
      current += c;
      if (c === inQuote) inQuote = null;
    } else if (c === '"' || c === "'") {
      current += c;
      inQuote = c;
    } else if (/\s/.test(c)) {
      if (current.trim()) {
        attrs.push(current.trim());
        current = '';
      }
    } else {
      current += c;
    }
  }

  if (current.trim()) attrs.push(current.trim());
  return attrs;
}
