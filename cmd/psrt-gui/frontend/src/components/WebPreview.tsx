import { useCallback, useEffect, useMemo, useRef, useState, type PointerEvent } from 'react';
import ImagePsrt from './ImagePsrt';
import { TextBlock } from './TextBlock';
import { MaskBlock } from './MaskBlock';
import { useEditor } from '../context/useEditor';
import { clientPointToImagePercent } from '../lib/imagePointerCoords';
import { resolveAssetReference } from '../lib/expandConsts';
import { isLocalAssetRef } from '../lib/localAssetRef';
import { stateToPsrtSection } from '../lib/stateToPsrtSection';
import { useAdaptedEntryStyles } from '../lib/useAdaptedEntryStyles';
import { usePageImageDataUri } from '../lib/usePageImageDataUri';
import type { PSRTEntry } from '../types/types';
import { visualapp } from '@wails/go/models';
import '../styles/image-psrt.css';

function textToOverlayEntry(
  text: visualapp.TextDetail,
  content: string,
): PSRTEntry {
  let parsedStyle: Record<string, unknown> = {};
  try {
    parsedStyle = JSON.parse(text.style || '{}') as Record<string, unknown>;
  } catch {
    parsedStyle = {};
  }
  return {
    index: text.index,
    x: text.x,
    y: text.y,
    width: text.width,
    size: text.textSize,
    text: content,
    style: parsedStyle,
    styleRaw: text.style ?? '{}',
  };
}

function maskToOverlayEntry(mask: visualapp.MaskDetail): PSRTEntry {
  let parsedStyle: Record<string, unknown> = {};
  try {
    parsedStyle = JSON.parse(mask.style || '{}') as Record<string, unknown>;
  } catch {
    parsedStyle = {};
  }
  parsedStyle.height = `${mask.height}%`;
  return {
    index: mask.index,
    x: mask.x,
    y: mask.y,
    width: mask.width,
    size: 0,
    maskHeight: mask.height,
    text: '',
    style: parsedStyle,
    styleRaw: mask.style ?? '{}',
  };
}

export function WebPreview({ variant = 'footer' }: { variant?: 'footer' | 'canvas' }) {
  const {
    state,
    pageImageUri,
    selectText,
    clearMultiSelect,
    multiSelected,
    webZoom,
    webTextsVisible,
    getTextDisplayContent,
    patchText,
    beginEdit,
    endEdit,
    showToast,
  } = useEditor();
  const stageRef = useRef<HTMLDivElement>(null);
  const imageContainerRef = useRef<HTMLDivElement>(null);
  const middleClickRef = useRef({ time: 0, x: 0, y: 0 });
  const [imageSize, setImageSize] = useState<{ w: number; h: number } | null>(null);
  const isCanvas = variant === 'canvas';

  const placeSelectedTextAtClient = useCallback(
    (clientX: number, clientY: number) => {
      const idx = state?.selectedIndex ?? -1;
      if (idx < 0) {
        showToast('Selecione um texto para reposicionar (duplo clique no botão do scroll).');
        return;
      }

      const container = imageContainerRef.current;
      if (!container) return;

      const pt = clientPointToImagePercent(clientX, clientY, container);
      if (!pt) return;

      void beginEdit();
      patchText(idx, { x: pt.x, y: pt.y });
      void endEdit();
    },
    [state?.selectedIndex, beginEdit, patchText, endEdit, showToast],
  );

  const onMiddleDoubleClickCapture = useCallback(
    (e: PointerEvent<HTMLDivElement>) => {
      if (!isCanvas || e.button !== 1) return;
      e.preventDefault();

      const isDouble =
        e.detail === 2 ||
        (() => {
          const now = Date.now();
          const prev = middleClickRef.current;
          const close =
            now - prev.time < 450 &&
            Math.hypot(e.clientX - prev.x, e.clientY - prev.y) < 16;
          middleClickRef.current = close
            ? { time: 0, x: 0, y: 0 }
            : { time: now, x: e.clientX, y: e.clientY };
          return close;
        })();

      if (!isDouble) return;

      e.stopPropagation();
      placeSelectedTextAtClient(e.clientX, e.clientY);
    },
    [isCanvas, placeSelectedTextAtClient],
  );

  const imageUrl = state?.page?.imageUrl;
  const docConsts = state?.consts;
  const resolvedImageSrc = usePageImageDataUri(imageUrl, pageImageUri, docConsts);
  const imageSrc =
    resolvedImageSrc ??
    pageImageUri ??
    (imageUrl && /^https?:\/\//i.test(imageUrl) ? imageUrl : null);

  const pageKey = state?.activePage ?? '';
  const imageKey = state?.page?.imageUrl ?? '';

  const section = useMemo(() => {
    if (!state) return null;
    const sec = stateToPsrtSection(state, imageSrc);
    if (!sec) return null;
    if (!isCanvas || !sec.entries?.length) return sec;
    return {
      ...sec,
      entries: sec.entries.map((entry) => ({
        ...entry,
        text: getTextDisplayContent(entry.index ?? 0, entry.text ?? ''),
      })),
    };
  }, [state, imageSrc, isCanvas, getTextDisplayContent]);

  useEffect(() => {
    setImageSize(null);
  }, [pageKey, imageKey]);

  const texts = state?.texts ?? [];
  const masks = state?.masks ?? [];
  const hasOverlay = texts.length > 0 || masks.length > 0;

  /** Higher index = painted later = on top for hit-test and z-order. */
  const overlayBlocks = useMemo(() => {
    const items: Array<
      | { kind: 'text'; data: (typeof texts)[number] }
      | { kind: 'mask'; data: (typeof masks)[number] }
    > = [];
    for (const t of texts) items.push({ kind: 'text', data: t });
    for (const m of masks) items.push({ kind: 'mask', data: m });
    items.sort((a, b) => a.data.index - b.data.index);
    return items;
  }, [texts, masks]);

  const overlayEntries = useMemo(() => {
    if (!isCanvas) return [];
    return overlayBlocks.map((item) =>
      item.kind === 'text'
        ? textToOverlayEntry(
            item.data,
            getTextDisplayContent(item.data.index, item.data.content ?? ''),
          )
        : maskToOverlayEntry(item.data),
    );
  }, [isCanvas, overlayBlocks, getTextDisplayContent]);

  const adaptedByIndex = useAdaptedEntryStyles(
    overlayEntries,
    imageSize?.w ?? 0,
    imageSize?.h ?? 0,
    webZoom,
  );

  if (!section) {
    return <div className="preview-web-empty">Open a page to preview</div>;
  }

  const resolvedImageUrl = imageUrl
    ? resolveAssetReference(imageUrl, docConsts)
    : undefined;

  const needsAssetResolve =
    resolvedImageUrl &&
    !imageSrc &&
    !isLocalAssetRef(resolvedImageUrl) &&
    !/^https?:\/\//i.test(resolvedImageUrl) &&
    !resolvedImageUrl.startsWith('data:') &&
    !resolvedImageUrl.startsWith('psrt-asset://');

  if (needsAssetResolve) {
    return <div className="preview-web-empty">Carregando imagem…</div>;
  }

  const showLocalAssetsHint =
    import.meta.env.VITE_PRT_WEB === 'true' &&
    resolvedImageUrl !== undefined &&
    isLocalAssetRef(resolvedImageUrl) &&
    !imageSrc;

  const selectedIndex = state?.selectedIndex ?? -1;

  const rootClass =
    variant === 'canvas' ? 'preview-web preview-web-canvas' : 'preview-web';

  return (
    <div className={rootClass}>
      {showLocalAssetsHint ? (
        <p className="preview-web-local-hint">
          Vincule a pasta de imagens com o botão <strong>Pasta assets</strong> (topo) ou 📁
          ao lado da URL — o browser não acessa paths do disco diretamente.
        </p>
      ) : null}
      <div
        className={isCanvas ? 'canvas-editor-stage' : undefined}
        onPointerDownCapture={isCanvas ? onMiddleDoubleClickCapture : undefined}
      >
        <ImagePsrt
          key={`${pageKey}::${imageKey}::${imageSrc ?? 'loading'}`}
          pageData={section}
          alt={section.title}
          pageLink={imageSrc ?? undefined}
          editor={isCanvas}
          hideSub={isCanvas && !webTextsVisible}
          parentZoom={webZoom}
          fixedReferenceSize={isCanvas}
          getSize={({ w, h }) => {
            if (w && h) setImageSize({ w, h });
          }}
          selectedIndex={
            isCanvas ? undefined : selectedIndex >= 0 ? selectedIndex : undefined
          }
          onClickEntry={isCanvas ? undefined : (index) => selectText(index)}
          consts={docConsts}
          imageContainerRef={imageContainerRef}
        />
        {isCanvas && webTextsVisible && hasOverlay ? (
          <div
            ref={stageRef}
            className="canvas-text-overlay"
            onPointerDown={(e) => {
              if (e.target === e.currentTarget) clearMultiSelect();
            }}
          >
            {overlayBlocks.map((item) =>
              item.kind === 'text' ? (
                <TextBlock
                  key={`t-${item.data.index}`}
                  text={item.data}
                  displayContent={getTextDisplayContent(
                    item.data.index,
                    item.data.content ?? '',
                  )}
                  isSelected={
                    selectedIndex === item.data.index ||
                    multiSelected.has(item.data.index)
                  }
                  stageRef={stageRef}
                  interactionOnly
                  imageWidth={imageSize?.w ?? 0}
                  imageHeight={imageSize?.h ?? 0}
                  zoom={webZoom}
                  adaptedStyles={adaptedByIndex.get(item.data.index)}
                />
              ) : (
                <MaskBlock
                  key={`m-${item.data.index}`}
                  mask={item.data}
                  isSelected={
                    selectedIndex === item.data.index ||
                    multiSelected.has(item.data.index)
                  }
                  stageRef={stageRef}
                  interactionOnly
                  imageWidth={imageSize?.w ?? 0}
                  imageHeight={imageSize?.h ?? 0}
                  zoom={webZoom}
                  adaptedStyles={adaptedByIndex.get(item.data.index)}
                />
              ),
            )}
          </div>
        ) : null}
      </div>
    </div>
  );
}
