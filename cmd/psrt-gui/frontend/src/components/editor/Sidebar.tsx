import { useState } from "react";
import s from "./sidebar.module.css";
import type { InlineWrapKind } from "../../lib/inlineMarkup";
import type { TextBlock } from "./types";
import { PropertiesPanel } from "./PropertiesPanel.new";
import { PsrtPanel } from "./PsrtPanel";
import { PageImageUrlBar } from "./PageImageUrlBar";

type Tab = "props" | "psrt";

export interface EditorSidebarProps {
  blocks: TextBlock[];
  activeId: string;
  activeTextContent: string;
  onActiveTextContentChange: (value: string) => void;
  pageImageUrl: string;
  pageImageConsts?: Record<string, string>;
  onPageImageUrlChange: (url: string) => void;
  onRefreshPageImage?: () => void;
  onPageImageError?: (message: string) => void;
  psrtValue: string;
  onSelect: (id: string) => void;
  onAdd: () => void;
  onDuplicate: () => void;
  onRemove: () => void;
  onBlockChange: (updater: (b: TextBlock) => TextBlock) => void;
  /** Writes one CSS key to the block style JSON — does not rewrite other keys. */
  onPatchStyleProp: (key: string, value: string | null) => void;
  onTypographyWrap?: (kind: InlineWrapKind, textarea: HTMLTextAreaElement) => boolean;
  onPsrtChange: (value: string) => void;
  onPsrtFocus?: () => void;
  onPsrtBlur?: () => void;
  fontOptions?: { value: string; label: string }[];
  emptyHint?: string;
}

export function Sidebar({
  blocks,
  activeId,
  activeTextContent,
  onActiveTextContentChange,
  pageImageUrl,
  pageImageConsts,
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
  emptyHint,
}: EditorSidebarProps) {
  const [tab, setTab] = useState<Tab>("props");

  return (
    <aside className={`${s.root} ${s.embedded}`}>
      <PageImageUrlBar
        imageUrl={pageImageUrl}
        consts={pageImageConsts}
        onChange={onPageImageUrlChange}
        onRefresh={onRefreshPageImage}
        onError={onPageImageError}
      />
      <div className={s.tabs}>
        <button
          type="button"
          className={`${s.tab} ${tab === "props" ? s.tabActive : ""}`}
          onClick={() => setTab("props")}
        >
          Propriedades
        </button>
        <button
          type="button"
          className={`${s.tab} ${tab === "psrt" ? s.tabActive : ""}`}
          onClick={() => setTab("psrt")}
        >
          PSRT
        </button>
      </div>
      <div className={s.scroll}>
        {tab === "props" ? (
          <PropertiesPanel
            blocks={blocks}
            activeId={activeId}
            textContent={activeTextContent}
            onTextContentChange={onActiveTextContentChange}
            onSelect={onSelect}
            onAdd={onAdd}
            onDuplicate={onDuplicate}
            onRemove={onRemove}
            onChange={onBlockChange}
            onPatchStyleProp={onPatchStyleProp}
            onTypographyWrap={onTypographyWrap}
            fontOptions={fontOptions}
            emptyHint={emptyHint}
          />
        ) : (
          <PsrtPanel
            value={psrtValue}
            onChange={onPsrtChange}
            onFocus={onPsrtFocus}
            onBlur={onPsrtBlur}
          />
        )}
      </div>
    </aside>
  );
}
