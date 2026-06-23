import { visualapp } from '@wails/go/models';

export interface TextLayoutOverride {
  x?: number;
  y?: number;
  width?: number;
  textSize?: number;
}

export function mergeTextDetail(
  base: visualapp.TextDetail,
  layout: Record<number, TextLayoutOverride>,
  contentOverrides: Record<number, string>,
): visualapp.TextDetail {
  const lo = layout[base.index];
  const content = contentOverrides[base.index];
  const merged = new visualapp.TextDetail({
    index: base.index,
    x: lo?.x ?? base.x,
    y: lo?.y ?? base.y,
    width: lo?.width ?? base.width,
    textSize: lo?.textSize ?? base.textSize,
    content: content ?? base.content,
    imageRef: base.imageRef,
    style: base.style,
  });
  return merged;
}

export function isLayoutPatch(patch: Partial<visualapp.TextPatch>): boolean {
  return (
    patch.x !== undefined ||
    patch.y !== undefined ||
    patch.width !== undefined ||
    patch.textSize !== undefined
  );
}

export function isContentOrStylePatch(patch: Partial<visualapp.TextPatch>): boolean {
  return (
    patch.content !== undefined ||
    patch.imageRef !== undefined ||
    (patch.styleSet !== undefined && Object.keys(patch.styleSet).length > 0) ||
    (patch.styleRemove !== undefined && patch.styleRemove.length > 0)
  );
}
