import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import * as api from '@wails/go/main/GUIApp';
import { useEditor } from '../context/useEditor';
import { useAlertModal } from '../context/AlertModalContext';
import type { EditorSidebarProps } from '../components/editor/Sidebar';
import {
  textDetailToBlock,
  maskDetailToBlock,
  blockPatchToEntryPatch,
  buildFontSelectOptions,
  singleStylePropPatch,
} from '../lib/textBlockAdapter';
import { toggleTextSelectionMarkup, type InlineWrapKind } from '../lib/inlineMarkup';

const CONTENT_DEBOUNCE_MS = 400;
const PSRT_APPLY_MS = 400;

export function useTextEditorSidebar(): EditorSidebarProps | null {
  const {
    state,
    document,
    patchText,
    patchMask,
    convertBlockKind,
    patchPage,
    selectText,
    addText,
    duplicateText,
    removeText,
    applyPagePsrtSource,
    beginEdit,
    endEdit,
    refreshPageImage,
    showToast,
    setTextContentDraft,
    clearTextContentDraft,
    getTextDisplayContent,
  } = useEditor();
  const { confirm } = useAlertModal();

  const pageName = state?.activePage ?? '';
  const texts = state?.texts ?? [];
  const masks = state?.masks ?? [];
  const selectedIndex = state?.selectedIndex ?? -1;

  const blocks = useMemo(() => {
    const byIndex = new Map<string, ReturnType<typeof textDetailToBlock>>();
    for (const t of texts) {
      byIndex.set(String(t.index), textDetailToBlock(t));
    }
    for (const m of masks) {
      byIndex.set(String(m.index), maskDetailToBlock(m));
    }
    return [...byIndex.values()].sort((a, b) => Number(a.id) - Number(b.id));
  }, [texts, masks]);

  const activeId = selectedIndex >= 0 ? String(selectedIndex) : '';
  const activeBlock = blocks.find((b) => b.id === activeId);

  const fontOptions = useMemo(
    () => buildFontSelectOptions(state?.fonts ?? []),
    [state?.fonts],
  );

  const contentTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pendingContent = useRef<{ index: number; content: string } | null>(null);

  const [psrtValue, setPsrtValue] = useState('');
  const psrtDirty = useRef(false);
  const psrtFocused = useRef(false);
  const psrtApplyTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const psrtPending = useRef<string | null>(null);

  const activeTextContent =
    selectedIndex >= 0
      ? getTextDisplayContent(selectedIndex, activeBlock?.content ?? '')
      : '';

  const loadPsrt = useCallback(async () => {
    if (!document || !pageName) {
      setPsrtValue('');
      return;
    }
    try {
      const src = await api.FormatPageDocumentJSON(
        JSON.stringify(document),
        pageName,
      );
      setPsrtValue(src);
    } catch (e) {
      showToast(String(e));
    }
  }, [document, pageName, showToast]);

  useEffect(() => {
    if (psrtDirty.current || psrtFocused.current) return;
    void loadPsrt();
  }, [loadPsrt, document, pageName]);

  const flushPsrt = useCallback(async () => {
    const body = psrtPending.current;
    psrtPending.current = null;
    if (body === null || !pageName) return;
    try {
      await applyPagePsrtSource(body);
    } catch (e) {
      showToast(String(e));
    }
  }, [applyPagePsrtSource, pageName, showToast]);

  const schedulePsrtApply = useCallback(
    (value: string) => {
      psrtPending.current = value;
      if (psrtApplyTimer.current) clearTimeout(psrtApplyTimer.current);
      psrtApplyTimer.current = setTimeout(() => {
        psrtApplyTimer.current = null;
        void flushPsrt();
      }, PSRT_APPLY_MS);
    },
    [flushPsrt],
  );

  const flushContent = useCallback(() => {
    if (contentTimer.current) {
      clearTimeout(contentTimer.current);
      contentTimer.current = null;
    }
    const pending = pendingContent.current;
    pendingContent.current = null;
    if (!pending) return;
    patchText(pending.index, { content: pending.content, append: false });
    clearTextContentDraft(pending.index);
  }, [patchText, clearTextContentDraft]);

  useEffect(() => {
    return () => {
      flushContent();
    };
  }, [selectedIndex, flushContent]);

  const onActiveTextContentChange = useCallback(
    (value: string) => {
      const idx = selectedIndex;
      if (idx < 0) return;
      setTextContentDraft(idx, value);
      pendingContent.current = { index: idx, content: value };
      if (contentTimer.current) clearTimeout(contentTimer.current);
      contentTimer.current = setTimeout(() => {
        contentTimer.current = null;
        flushContent();
      }, CONTENT_DEBOUNCE_MS);
    },
    [selectedIndex, setTextContentDraft, flushContent],
  );

  const onTypographyWrap = useCallback(
    (kind: InlineWrapKind, textarea: HTMLTextAreaElement): boolean => {
      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;
      if (start === end) return false;
      const { value, selectionStart, selectionEnd } = toggleTextSelectionMarkup(
        activeTextContent,
        start,
        end,
        kind,
      );
      onActiveTextContentChange(value);
      requestAnimationFrame(() => {
        textarea.focus();
        textarea.setSelectionRange(selectionStart, selectionEnd);
      });
      return true;
    },
    [activeTextContent, onActiveTextContentChange],
  );

  const onPatchStyleProp: EditorSidebarProps['onPatchStyleProp'] = useCallback(
    (key: string, value: string | null) => {
      const idx = selectedIndex;
      if (idx < 0) return;
      const stylePatch = singleStylePropPatch(key, value);
      if (!stylePatch.styleSet && !stylePatch.styleRemove) return;
      flushContent();
      void beginEdit();
      const block = blocks.find((b) => b.id === String(idx));
      if (block?.kind === 'mask') {
        patchMask(idx, stylePatch);
      } else {
        patchText(idx, stylePatch);
      }
    },
    [
      selectedIndex,
      blocks,
      patchText,
      patchMask,
      beginEdit,
      flushContent,
    ],
  );

  const onBlockChange: EditorSidebarProps['onBlockChange'] = useCallback(
    (updater) => {
      const idx = selectedIndex;
      if (idx < 0 || !document) return;
      const prev = blocks.find((b) => b.id === String(idx));
      if (!prev) return;
      const next = updater(prev);
      const entry = blockPatchToEntryPatch(prev, next);

      flushContent();
      void beginEdit();

      if (entry.kindChange) {
        convertBlockKind(idx, entry.kindChange, {
          mask: entry.mask,
          text: entry.text,
        });
        return;
      }
      if (entry.mask) {
        patchMask(idx, entry.mask);
        return;
      }
      if (entry.text) {
        patchText(idx, entry.text);
      }
    },
    [
      selectedIndex,
      document,
      blocks,
      patchText,
      patchMask,
      convertBlockKind,
      beginEdit,
      flushContent,
    ],
  );

  const onSelect = useCallback(
    (id: string) => {
      flushContent();
      const n = Number(id);
      if (Number.isFinite(n)) selectText(n);
    },
    [selectText, flushContent],
  );

  const onAdd = useCallback(() => {
    flushContent();
    addText();
  }, [addText, flushContent]);

  const onDuplicate = useCallback(() => {
    if (selectedIndex < 0) return;
    flushContent();
    duplicateText();
  }, [selectedIndex, duplicateText, flushContent]);

  const onRemove = useCallback(async () => {
    if (selectedIndex < 0) return;
    const ok = await confirm({
      title: 'Remover texto',
      message: `Remover o texto #${selectedIndex}?`,
      confirmLabel: 'Remover',
      danger: true,
    });
    if (!ok) return;
    flushContent();
    removeText(selectedIndex);
  }, [selectedIndex, removeText, confirm, flushContent]);

  const onPsrtChange = useCallback(
    (value: string) => {
      psrtDirty.current = true;
      setPsrtValue(value);
      schedulePsrtApply(value);
    },
    [schedulePsrtApply],
  );

  const onPsrtFocus = useCallback(() => {
    psrtFocused.current = true;
    beginEdit();
  }, [beginEdit]);

  const pageImageUrl = state?.page?.imageUrl ?? '';

  const onPageImageUrlChange = useCallback(
    (url: string) => {
      patchPage({ imageUrl: url });
    },
    [patchPage],
  );

  const onPageImageError = useCallback(
    (message: string) => {
      showToast(message);
    },
    [showToast],
  );

  const onRefreshPageImage = useCallback(() => {
    refreshPageImage().catch((e) => showToast(String(e)));
  }, [refreshPageImage, showToast]);

  const onPsrtBlur = useCallback(async () => {
    psrtFocused.current = false;
    if (psrtApplyTimer.current) {
      clearTimeout(psrtApplyTimer.current);
      psrtApplyTimer.current = null;
    }
    await flushPsrt();
    endEdit();
    psrtDirty.current = false;
    await loadPsrt();
  }, [flushPsrt, endEdit, loadPsrt]);

  if (!state?.page) {
    return null;
  }

  return {
    blocks,
    activeId,
    activeTextContent,
    onActiveTextContentChange,
    pageImageUrl,
    pageImageConsts: state?.consts,
    onPageImageUrlChange,
    onRefreshPageImage,
    onPageImageError,
    psrtValue,
    onSelect,
    onAdd,
    onDuplicate,
    onRemove,
    onBlockChange,
    onPatchStyleProp,
    onTypographyWrap,
    onPsrtChange,
    onPsrtFocus,
    onPsrtBlur,
    fontOptions,
    emptyHint: 'Selecione um texto na lista ou clique no bloco na imagem.',
  };
}
