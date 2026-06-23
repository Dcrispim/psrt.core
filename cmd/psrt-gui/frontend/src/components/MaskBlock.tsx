import React, { useRef, useState, type PointerEvent } from 'react';
import { visualapp } from '@wails/go/models';
import { useEditor } from '../context/useEditor';
import { snapCoord } from '../lib/applyStyle';
import type { PSRTEntry } from '../types/types';
import { resolveEntryStyle, useAdaptedEntryStyles, type AdaptedWebStyles } from '../lib/useAdaptedEntryStyles';

interface MaskBlockProps {
  mask: visualapp.MaskDetail;
  isSelected: boolean;
  stageRef: React.RefObject<HTMLDivElement | null>;
  interactionOnly?: boolean;
  imageWidth?: number;
  imageHeight?: number;
  zoom?: number;
  adaptedStyles?: AdaptedWebStyles;
}

type DragMode = 'move' | 'width' | 'height';

export function MaskBlock({
  mask,
  isSelected,
  stageRef,
  interactionOnly = false,
  imageWidth = 0,
  imageHeight = 0,
  zoom = 1,
  adaptedStyles,
}: MaskBlockProps) {
  const { patchMask, beginEdit, endEdit, selectText, toggleMultiSelect } = useEditor();

  const dragRef = useRef<{
    mode: DragMode;
    startX: number;
    startY: number;
    origX: number;
    origY: number;
    origW: number;
    origH: number;
    rect: DOMRect;
  } | null>(null);

  const [tooltip, setTooltip] = useState('');

  const onPointerMoveBlock = (e: PointerEvent) => {
    const stage = stageRef.current;
    if (!stage) return;
    const rect = stage.getBoundingClientRect();
    const px = (((e.clientX - rect.left) / rect.width) * 100).toFixed(1);
    const py = (((e.clientY - rect.top) / rect.height) * 100).toFixed(1);
    setTooltip(`#${mask.index}  x:${mask.x}% y:${mask.y}%  (${px}%, ${py}%)`);
  };

  const onPointerDown = (e: PointerEvent, mode: DragMode) => {
    e.stopPropagation();
    e.preventDefault();
    const stage = stageRef.current;
    if (!stage) return;
    void beginEdit();
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
    dragRef.current = {
      mode,
      startX: e.clientX,
      startY: e.clientY,
      origX: mask.x,
      origY: mask.y,
      origW: mask.width,
      origH: mask.height,
      rect: stage.getBoundingClientRect(),
    };
  };

  const onPointerMove = (e: PointerEvent) => {
    const d = dragRef.current;
    if (!d) return;
    const dx = ((e.clientX - d.startX) / d.rect.width) * 100;
    const dy = ((e.clientY - d.startY) / d.rect.height) * 100;
    if (d.mode === 'move') {
      patchMask(mask.index, { x: snapCoord(d.origX + dx), y: snapCoord(d.origY + dy) });
    } else if (d.mode === 'width') {
      patchMask(mask.index, { width: snapCoord(Math.max(1, d.origW + dx)) });
    } else if (d.mode === 'height') {
      const delta = dy * 0.2;
      patchMask(mask.index, { height: snapCoord(Math.max(0.5, d.origH + delta)) });
    }
  };

  const onPointerUp = (e: PointerEvent) => {
    (e.currentTarget as HTMLElement).releasePointerCapture(e.pointerId);
    dragRef.current = null;
    void endEdit();
  };

  let parsedStyle: Record<string, unknown> = {};
  try {
    parsedStyle = JSON.parse(mask.style || '{}') as Record<string, unknown>;
  } catch {
    parsedStyle = {};
  }
  parsedStyle.height = `${mask.height}%`;

  const entry: PSRTEntry = {
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
    adaptedStyles ?? adaptedMap.get(mask.index),
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
      className={`text-block mask-block${isSelected ? ' selected' : ''}${interactionOnly ? ' text-block--interaction-only' : ''}`}
      style={hitStyle}
      title={tooltip}
      data-index={mask.index}
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
          toggleMultiSelect(mask.index);
          return;
        }
        selectText(mask.index);
      }}
    >
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
              onPointerDown(e, 'height');
            }}
          />
        </>
      )}
    </div>
  );
}
