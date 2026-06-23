import { visualapp } from '@wails/go/models';
import type { PsrtDocument, PsrtMask, PsrtPage, PsrtText } from '../types/document';
import { snapCoord } from './applyStyle';
import { resolveMaskHeightPercent, heightPercentFromStyleObject } from './maskHeight';
import { styleStringValue } from './parseStyle';

function applyStyleSet(
  obj: Record<string, unknown>,
  styleSet: Record<string, string>,
): void {
  for (const [k, v] of Object.entries(styleSet)) {
    obj[k] = styleStringValue(v);
  }
}

/** Apply styleRemove before styleSet so replacement keys are not deleted afterward. */
function applyStylePatch(
  obj: Record<string, unknown>,
  styleSet?: Record<string, string>,
  styleRemove?: string[],
): void {
  if (styleRemove?.length) {
    for (const k of styleRemove) {
      delete obj[k];
    }
  }
  if (styleSet && Object.keys(styleSet).length > 0) {
    applyStyleSet(obj, styleSet);
  }
}

export function cloneDocument(doc: PsrtDocument): PsrtDocument {
  return structuredClone(doc);
}

/** Blank PSRT: empty $FONTS / $CONSTS, page "inicio" with centered starter text. */
export function createEmptyDocument(): PsrtDocument {
  return {
    pages: [
      {
        name: 'inicio',
        imageUrl: '',
        style: {},
        texts: [
          {
            index: 0,
            x: 50,
            y: 50,
            width: 80,
            textSize: 3,
            content: 'texto inicio',
            style: {
              color: '#000000',
              'text-align': 'center',
            },
          },
        ],
      },
    ],
    fonts: [],
    consts: {},
  };
}

export function styleToString(style: PsrtText['style']): string {
  if (typeof style === 'string') return style;
  return JSON.stringify(style ?? {});
}

function parseStyleObject(style: PsrtText['style']): Record<string, unknown> {
  if (typeof style === 'object' && style !== null) return { ...style };
  try {
    return JSON.parse(style || '{}') as Record<string, unknown>;
  } catch {
    return {};
  }
}

function pageTexts(page: PsrtPage): PsrtText[] {
  if (page.texts == null) {
    page.texts = [];
  }
  return page.texts;
}

function pageMasks(page: PsrtPage): PsrtMask[] {
  if (page.masks == null) {
    page.masks = [];
  }
  return page.masks;
}

function normalizeDocument(doc: PsrtDocument): PsrtDocument {
  if (!doc.pages) {
    doc.pages = [];
  }
  if (!doc.fonts) {
    doc.fonts = [];
  }
  if (!doc.consts) {
    doc.consts = {};
  }
  for (const page of doc.pages) {
    for (const m of pageMasks(page)) {
      const styleObj = parseStyleObject(m.style);
      m.height = resolveMaskHeightPercent(m.height, styleObj);
    }
  }
  return doc;
}

function maskHeightFromText(t: PsrtText): number {
  const style = parseStyleObject(t.style);
  const fromStyle = heightPercentFromStyleObject(style);
  if (fromStyle !== null) return fromStyle;
  if (t.textSize > 0) return snapCoord(Math.max(0.5, t.textSize));
  return 5;
}

export function nextBlockIndex(page: PsrtPage): number {
  const indices = [
    ...pageTexts(page).map((t) => t.index),
    ...pageMasks(page).map((m) => m.index),
  ];
  if (indices.length === 0) return 0;
  return Math.max(...indices) + 1;
}

export function buildUIState(
  doc: PsrtDocument,
  filePath: string,
  activePage: string,
  selectedIndex: number,
  autoCompile: boolean,
): visualapp.UIState {
  const pages = (doc.pages ?? []).map(
    (p) =>
      new visualapp.PageSummary({
        name: p.name,
        imageUrl: p.imageUrl,
      }),
  );

  const page = doc.pages.find((p) => p.name === activePage);
  let pageDetail: visualapp.PageDetail | undefined;
  let texts: visualapp.TextDetail[] | undefined;
  let text: visualapp.TextDetail | undefined;

  if (page) {
    pageDetail = new visualapp.PageDetail({
      name: page.name,
      imageUrl: page.imageUrl,
      style: styleToString(page.style),
    });
    texts = pageTexts(page).map(
      (t) =>
        new visualapp.TextDetail({
          index: t.index,
          x: t.x,
          y: t.y,
          width: t.width,
          textSize: t.textSize,
          content: t.content,
          imageRef: t.imageRef ?? '',
          style: styleToString(t.style),
        }),
    );
    const masks = pageMasks(page).map(
      (m) =>
        new visualapp.MaskDetail({
          index: m.index,
          x: m.x,
          y: m.y,
          width: m.width,
          height: m.height,
          imageRef: m.imageRef ?? '',
          style: styleToString(m.style),
        }),
    );
    if (selectedIndex >= 0) {
      text = texts.find((t) => t.index === selectedIndex);
      if (!text) {
        const mask = masks.find((m) => m.index === selectedIndex);
        if (mask) {
          return new visualapp.UIState({
            filePath,
            activePage,
            selectedIndex,
            pages,
            page: pageDetail,
            texts,
            masks,
            mask,
            fonts: [...(doc.fonts ?? [])],
            consts: { ...(doc.consts ?? {}) },
            autoCompile,
          });
        }
      }
    }

    return new visualapp.UIState({
      filePath,
      activePage,
      selectedIndex,
      pages,
      page: pageDetail,
      texts,
      masks,
      text,
      fonts: [...(doc.fonts ?? [])],
      consts: { ...(doc.consts ?? {}) },
      autoCompile,
    });
  }

  return new visualapp.UIState({
    filePath,
    activePage,
    selectedIndex,
    pages,
    page: pageDetail,
    texts,
    text,
    fonts: [...(doc.fonts ?? [])],
    consts: { ...(doc.consts ?? {}) },
    autoCompile,
  });
}

/** Promotes an empty >> block to == when patching mask fields. */
export function ensureMaskBlockInDocument(
  doc: PsrtDocument,
  pageName: string,
  index: number,
): PsrtDocument {
  const page = doc.pages.find((p) => p.name === pageName);
  if (!page) return doc;
  if (pageMasks(page).some((m) => m.index === index)) return doc;
  const text = pageTexts(page).find((t) => t.index === index);
  if (!text || text.content.trim() !== '') return doc;
  return convertBlockKindInDocument(doc, pageName, index, 'mask');
}

export function patchMaskInDocument(
  doc: PsrtDocument,
  pageName: string,
  index: number,
  patch: Partial<visualapp.MaskPatch>,
): PsrtDocument {
  let next = ensureMaskBlockInDocument(doc, pageName, index);
  next = cloneDocument(next);
  const page = next.pages.find((p) => p.name === pageName);
  if (!page) return next;
  const m = pageMasks(page).find((x) => x.index === index);
  if (!m) return next;

  if (patch.x !== undefined) m.x = snapCoord(patch.x);
  if (patch.y !== undefined) m.y = snapCoord(patch.y);
  if (patch.width !== undefined) m.width = snapCoord(Math.max(1, patch.width));
  if (patch.height !== undefined) m.height = snapCoord(Math.max(0.5, patch.height));
  if (patch.imageRef !== undefined) m.imageRef = patch.imageRef;

  if (
    (patch.styleSet && Object.keys(patch.styleSet).length > 0) ||
    patch.styleRemove?.length
  ) {
    const obj = parseStyleObject(m.style);
    applyStylePatch(obj, patch.styleSet, patch.styleRemove);
    m.style = obj;
  }

  return next;
}

export function patchTextInDocument(
  doc: PsrtDocument,
  pageName: string,
  index: number,
  patch: Partial<visualapp.TextPatch>,
): PsrtDocument {
  const next = cloneDocument(doc);
  const page = next.pages.find((p) => p.name === pageName);
  if (!page) return next;
  const t = pageTexts(page).find((x) => x.index === index);
  if (!t) return next;

  if (patch.content !== undefined) {
    t.content = patch.append ? `${t.content}${patch.content}` : patch.content;
  }
  if (patch.x !== undefined) t.x = snapCoord(patch.x);
  if (patch.y !== undefined) t.y = snapCoord(patch.y);
  if (patch.width !== undefined) t.width = snapCoord(Math.max(1, patch.width));
  if (patch.textSize !== undefined) t.textSize = snapCoord(Math.max(0.5, patch.textSize));
  if (patch.imageRef !== undefined) t.imageRef = patch.imageRef;

  if (
    (patch.styleSet && Object.keys(patch.styleSet).length > 0) ||
    patch.styleRemove?.length
  ) {
    const obj = parseStyleObject(t.style);
    applyStylePatch(obj, patch.styleSet, patch.styleRemove);
    t.style = obj;
  }

  return next;
}

export function patchPageInDocument(
  doc: PsrtDocument,
  activePage: string,
  patch: Partial<visualapp.PagePatch> & { styleRemove?: string[] },
): PsrtDocument {
  const next = cloneDocument(doc);
  const page = next.pages.find((p) => p.name === activePage);
  if (!page) return next;

  if (patch.name !== undefined && patch.name !== page.name) {
    page.name = patch.name;
  }
  if (patch.imageUrl !== undefined) page.imageUrl = patch.imageUrl;
  if (
    (patch.styleSet && Object.keys(patch.styleSet).length > 0) ||
    patch.styleRemove?.length
  ) {
    const obj = parseStyleObject(page.style);
    applyStylePatch(obj, patch.styleSet, patch.styleRemove);
    page.style = obj;
  }
  return next;
}

export function setActivePageInDoc(_doc: PsrtDocument, _name: string): PsrtDocument {
  return _doc;
}

export function addPageToDocument(
  doc: PsrtDocument,
  name: string,
  imageURL: string,
  styleJSON = '{}',
): PsrtDocument {
  const next = cloneDocument(doc);
  let style: string | Record<string, unknown> = styleJSON;
  try {
    style = JSON.parse(styleJSON) as Record<string, unknown>;
  } catch {
    /* keep string */
  }
  const page: PsrtPage = { name, imageUrl: imageURL, style, texts: [] };
  next.pages.push(page);
  return next;
}

export function removePageFromDocument(doc: PsrtDocument, name: string): PsrtDocument {
  const next = cloneDocument(doc);
  next.pages = next.pages.filter((p) => p.name !== name);
  return next;
}

export function movePageInDocument(
  doc: PsrtDocument,
  name: string,
  ref: string,
  before: boolean,
): PsrtDocument {
  const next = cloneDocument(doc);
  const from = next.pages.findIndex((p) => p.name === name);
  const refIdx = next.pages.findIndex((p) => p.name === ref);
  if (from < 0 || refIdx < 0) return next;
  const [item] = next.pages.splice(from, 1);
  let to = refIdx;
  if (from < refIdx) to--;
  if (!before) to++;
  next.pages.splice(to, 0, item);
  return next;
}

export function addFontToDocument(doc: PsrtDocument, url: string): PsrtDocument {
  const next = cloneDocument(doc);
  if (!next.fonts.includes(url)) next.fonts.push(url);
  return next;
}

export function removeFontFromDocument(doc: PsrtDocument, url: string): PsrtDocument {
  const next = cloneDocument(doc);
  next.fonts = next.fonts.filter((f) => f !== url);
  return next;
}

export function addConstToDocument(
  doc: PsrtDocument,
  name: string,
  value: string,
): PsrtDocument {
  const next = cloneDocument(doc);
  if (!next.consts) next.consts = {};
  next.consts[name] = value;
  return next;
}

export function removeConstFromDocument(doc: PsrtDocument, name: string): PsrtDocument {
  const next = cloneDocument(doc);
  if (next.consts) delete next.consts[name];
  return next;
}

function defaultStyleForNewText(existing: PsrtText[]): Record<string, unknown> {
  if (existing.length === 0) {
    return {
      color: '#000000',
      padding: '8px',
      'text-align': 'center',
    };
  }
  const prev = existing[existing.length - 1];
  const prevStyle = parseStyleObject(prev.style);
  const style: Record<string, unknown> = {
    padding: prevStyle.padding ?? '8px',
    'text-align': 'center',
  };
  const align = prevStyle['text-align'] ?? prevStyle.textAlign;
  if (align !== undefined && align !== null && styleStringValue(align) !== '') {
    style['text-align'] = styleStringValue(align);
  }
  if (prevStyle.color !== undefined && prevStyle.color !== null) {
    style.color = styleStringValue(prevStyle.color);
  }
  const bg =
    prevStyle.background ??
    prevStyle.backGround ??
    prevStyle['background-color'];
  if (bg !== undefined && bg !== null && styleStringValue(bg) !== '') {
    style.background = styleStringValue(bg);
  }
  return style;
}

export function nextTextIndex(page: PsrtPage): number {
  return nextBlockIndex(page);
}

export function addTextToDocument(
  doc: PsrtDocument,
  pageName: string,
  partial?: Partial<Omit<PsrtText, 'index'>>,
): { doc: PsrtDocument; index: number } {
  const next = cloneDocument(doc);
  const page = next.pages.find((p) => p.name === pageName);
  if (!page) return { doc: next, index: -1 };

  const index = nextTextIndex(page);
  const style = defaultStyleForNewText(pageTexts(page));
  const text: PsrtText = {
    index,
    x: 10,
    y: 10,
    width: 40,
    textSize: 3,
    content: 'Novo texto',
    style,
    ...partial,
  };
  pageTexts(page).push(text);
  return { doc: next, index };
}

export function duplicateTextInDocument(
  doc: PsrtDocument,
  pageName: string,
  sourceIndex: number,
): { doc: PsrtDocument; index: number } {
  const next = cloneDocument(doc);
  const page = next.pages.find((p) => p.name === pageName);
  if (!page) return { doc: next, index: -1 };

  const sourceMask = pageMasks(page).find((m) => m.index === sourceIndex);
  if (sourceMask) {
    const index = nextBlockIndex(page);
    const mask = structuredClone(sourceMask);
    mask.index = index;
    pageMasks(page).push(mask);
    return { doc: next, index };
  }

  const texts = pageTexts(page);
  const source = texts.find((t) => t.index === sourceIndex);
  if (!source) return { doc: next, index: -1 };

  const index = nextBlockIndex(page);
  const text = structuredClone(source);
  text.index = index;
  texts.push(text);
  return { doc: next, index };
}

export function removeTextFromDocument(
  doc: PsrtDocument,
  pageName: string,
  index: number,
): PsrtDocument {
  const next = cloneDocument(doc);
  const page = next.pages.find((p) => p.name === pageName);
  if (!page) return next;
  page.texts = pageTexts(page).filter((t) => t.index !== index);
  page.masks = pageMasks(page).filter((m) => m.index !== index);
  return next;
}

export function addMaskToDocument(
  doc: PsrtDocument,
  pageName: string,
  partial?: Partial<Omit<PsrtMask, 'index'>>,
): { doc: PsrtDocument; index: number } {
  const next = cloneDocument(doc);
  const page = next.pages.find((p) => p.name === pageName);
  if (!page) return { doc: next, index: -1 };

  const index = nextBlockIndex(page);
  const mask: PsrtMask = {
    index,
    x: 10,
    y: 10,
    width: 20,
    height: 5,
    style: { background: '#eee9b2' },
    ...partial,
  };
  pageMasks(page).push(mask);
  return { doc: next, index };
}

export function convertBlockKindInDocument(
  doc: PsrtDocument,
  pageName: string,
  index: number,
  toKind: 'text' | 'mask',
): PsrtDocument {
  const next = cloneDocument(doc);
  const page = next.pages.find((p) => p.name === pageName);
  if (!page) return next;

  const textIdx = pageTexts(page).findIndex((t) => t.index === index);
  const maskIdx = pageMasks(page).findIndex((m) => m.index === index);

  if (toKind === 'mask' && textIdx >= 0) {
    const t = pageTexts(page)[textIdx];
    const style = parseStyleObject(t.style);
    const height = maskHeightFromText(t);
    delete style.height;
    delete style.Height;
    const mask: PsrtMask = {
      index: t.index,
      x: t.x,
      y: t.y,
      width: t.width,
      height,
      imageRef: t.imageRef,
      style,
    };
    page.texts = pageTexts(page).filter((_, i) => i !== textIdx);
    pageMasks(page).push(mask);
    return next;
  }

  if (toKind === 'text' && maskIdx >= 0) {
    const m = pageMasks(page)[maskIdx];
    const text: PsrtText = {
      index: m.index,
      x: m.x,
      y: m.y,
      width: m.width,
      textSize: m.height > 0 ? m.height : 3,
      content: '',
      imageRef: m.imageRef,
      style: m.style,
    };
    page.masks = pageMasks(page).filter((_, i) => i !== maskIdx);
    pageTexts(page).push(text);
    return next;
  }

  return next;
}

export function parseDocumentJson(json: string): PsrtDocument {
  return normalizeDocument(JSON.parse(json) as PsrtDocument);
}
