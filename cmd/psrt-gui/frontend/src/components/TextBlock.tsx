import React, { useRef, useState, type PointerEvent } from 'react';
import { visualapp } from '@wails/go/models';
import { useEditor } from '../context/useEditor';
import { snapCoord } from '../lib/applyStyle';
import { isMaskTextBlock, parseHeightPercent } from '../lib/textBoxHeight';
import { resolveEntryStyle, useAdaptedEntryStyles, type AdaptedWebStyles } from '../lib/useAdaptedEntryStyles';
import type { PSRTEntry } from '../types/types';
import { FormattedText } from './FormattedText';

interface TextBlockProps {
  text: visualapp.TextDetail;
  displayContent: string;
  isSelected: boolean;
  stageRef: React.RefObject<HTMLDivElement | null>;
  /** When true, text is invisible; ImagePsrt renders styles underneath (canvas WEB). */
  interactionOnly?: boolean;
  imageWidth?: number;
  imageHeight?: number;
  zoom?: number;
  /** Pre-computed styles from parent (avoids N WASM calls in canvas mode). */
  adaptedStyles?: AdaptedWebStyles;
}

type DragMode = 'move' | 'width' | 'textSize';

export function TextBlock({
  text,
  displayContent,
  isSelected,
  stageRef,
  interactionOnly = false,
  imageWidth = 0,
  imageHeight = 0,
  zoom = 1,
  adaptedStyles,
}: TextBlockProps) {
  const {
    patchText,
    beginEdit,
    endEdit,
    selectText,
    toggleMultiSelect,
    state,
  } = useEditor();

  const dragRef = useRef<{
    mode: DragMode;
    startX: number;
    startY: number;
    origX: number;
    origY: number;
    origW: number;
    origTs: number;
    rect: DOMRect;
  } | null>(null);

  const [tooltip, setTooltip] = useState('');

  const onPointerMoveBlock = (e: PointerEvent) => {
    const stage = stageRef.current;
    if (!stage) return;
    const rect = stage.getBoundingClientRect();
    const px = (((e.clientX - rect.left) / rect.width) * 100).toFixed(1);
    const py = (((e.clientY - rect.top) / rect.height) * 100).toFixed(1);
    setTooltip(`#${text.index}  x:${text.x}% y:${text.y}%  (${px}%, ${py}%)`);
  };

  const onPointerDown = (e: PointerEvent, mode: DragMode) => {
    e.stopPropagation();
    e.preventDefault();
    const stage = stageRef.current;
    if (!stage) return;
    void beginEdit();
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
    let origTs = text.textSize;
    if (mode === 'textSize' && isMaskTextBlock(displayContent)) {
      try {
        const parsed = JSON.parse(text.style || '{}') as Record<string, unknown>;
        origTs = parseHeightPercent(String(parsed.height ?? '')) ?? text.textSize;
      } catch {
        origTs = text.textSize;
      }
    }
    dragRef.current = {
      mode,
      startX: e.clientX,
      startY: e.clientY,
      origX: text.x,
      origY: text.y,
      origW: text.width,
      origTs,
      rect: stage.getBoundingClientRect(),
    };
  };

  const onPointerMove = (e: PointerEvent) => {
    const d = dragRef.current;
    if (!d) return;
    const dx = ((e.clientX - d.startX) / d.rect.width) * 100;
    const dy = ((e.clientY - d.startY) / d.rect.height) * 100;
    if (d.mode === 'move') {
      patchText(text.index, { x: snapCoord(d.origX + dx), y: snapCoord(d.origY + dy) });
    } else if (d.mode === 'width') {
      patchText(text.index, { width: snapCoord(Math.max(1, d.origW + dx)) });
    } else if (d.mode === 'textSize') {
      const delta = dy * 0.2;
      if (isMaskTextBlock(displayContent)) {
        patchText(text.index, {
          styleSet: {
            height: `${snapCoord(Math.max(0.5, d.origTs + delta))}%`,
          },
        });
      } else {
        patchText(text.index, {
          textSize: snapCoord(Math.max(0.5, d.origTs + delta)),
        });
      }
    }
  };

  const onPointerUp = (e: PointerEvent) => {
    (e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId);
    dragRef.current = null;
    void endEdit();
  };

  let parsedStyle: Record<string, unknown> = {};
  try {
    parsedStyle = JSON.parse(text.style || '{}') as Record<string, unknown>;
  } catch {
    parsedStyle = {};
  }

  const entry: PSRTEntry = {
    index: text.index,
    x: text.x,
    y: text.y,
    width: text.width,
    size: text.textSize,
    text: displayContent,
    style: parsedStyle,
    styleRaw: text.style ?? '{}',
  };

  const adaptedMap = useAdaptedEntryStyles(
    adaptedStyles !== undefined ? [] : [entry],
    adaptedStyles !== undefined ? 0 : imageWidth,
    adaptedStyles !== undefined ? 0 : imageHeight,
    adaptedStyles !== undefined ? 1 : zoom,
  );
  const metrics =
    imageWidth > 0 && imageHeight > 0
      ? { refWidth: imageWidth, refHeight: imageHeight, zoom }
      : undefined;
  const resolved = resolveEntryStyle(
    entry,
    adaptedStyles ?? adaptedMap.get(text.index),
    metrics,
  );

  const hitStyle: React.CSSProperties = {
    ...resolved.hitArea,
    cursor: 'move',
    background: 'transparent',
    border: 'none',
  };

  return (
    <div
      className={`text-block${isSelected ? ' selected' : ''}${interactionOnly ? ' text-block--interaction-only' : ''}`}
      style={hitStyle}
      title={tooltip}
      data-index={text.index}
      onPointerDown={(e) => onPointerDown(e, 'move')}
      onPointerMove={(e) => {
        onPointerMoveBlock(e);
        onPointerMove(e);
      }}
      onPointerUp={onPointerUp}
      onPointerCancel={onPointerUp}
      onClick={(e) => {
        e.stopPropagation();
        if (e.ctrlKey || e.metaKey) {
          toggleMultiSelect(text.index);
          return;
        }
        selectText(text.index);
      }}
    >
      {!interactionOnly ? (
        <FormattedText content={displayContent} consts={state?.consts} />
      ) : null}
      {isSelected && (
        <>
          <span
            className="resize-handle e"
            onPointerDown={(e) => {
              e.stopPropagation();
              onPointerDown(e, 'width');
            }}
          />
          <span
            className="resize-handle s"
            onPointerDown={(e) => {
              e.stopPropagation();
              onPointerDown(e, 'textSize');
            }}
          />
        </>
      )}
    </div>
  );
}
