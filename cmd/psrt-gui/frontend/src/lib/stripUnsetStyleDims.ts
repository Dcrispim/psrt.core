import type { CSSProperties } from 'react';
import { parseStyle } from './parseStyle';
import { isPresentStyleValue } from './styleValue';

const DIM_KEYS = ['height', 'width', 'minHeight', 'maxHeight', 'minWidth', 'maxWidth'] as const;

/** Removes layout dimensions from adapted CSS unless they exist in the raw PSRT style JSON. */
export function stripUnsetStyleDims(
  style: CSSProperties,
  styleRaw: string | undefined,
): CSSProperties {
  const parsed = parseStyle(styleRaw ?? '{}');
  const out: CSSProperties = { ...style };
  for (const key of DIM_KEYS) {
    if (!isPresentStyleValue(key, parsed[key])) {
      delete (out as Record<string, unknown>)[key];
    }
  }
  return out;
}
