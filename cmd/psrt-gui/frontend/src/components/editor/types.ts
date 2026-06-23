export type HAlign = "left" | "center" | "right" | "justify";
export type VAlign = "flex-start" | "center" | "flex-end";
export type Direction = "ltr" | "rtl";
export type WhiteSpace = "normal" | "nowrap" | "pre" | "pre-wrap";
export type Overflow = "visible" | "hidden" | "clip" | "ellipsis";
export type TextTransform = "none" | "uppercase" | "lowercase" | "capitalize";

export interface BorderRadius {
  topLeft: number;
  topRight: number;
  bottomRight: number;
  bottomLeft: number;
}

export const ZERO_BORDER_RADIUS: BorderRadius = {
  topLeft: 0,
  topRight: 0,
  bottomRight: 0,
  bottomLeft: 0,
};

export type BlockKind = "text" | "mask";

export interface TextBlock {
  id: string;
  name: string;
  kind: BlockKind;
  x: number;
  y: number;
  width: number;
  /** Text size (% min side); text blocks only. */
  size: number;
  /** Height (% page); mask blocks only. */
  height: number;
  content: string;
  /** Preview URL for mask background (not persisted). */
  bgImage?: string;
  /** Persisted PSRT image ref segment. */
  imageRef?: string;
  color: string;
  /** When false, `color` is omitted from persisted style. */
  colorSet: boolean;
  background: string;
  /** When false, `background` is omitted from persisted style. */
  backgroundSet: boolean;
  font: {
    family: string;
    weight: number;
    size: number;
    /** When true, `font-size` is written to block style JSON. */
    sizeOverride: boolean;
    lineHeight: number;
    letterSpacing: number;
    italic: boolean;
    underline: boolean;
    strike: boolean;
    bold: boolean;
    transform: TextTransform;
  };
  align: {
    horizontal: HAlign;
    vertical: VAlign;
    direction: Direction;
    whiteSpace: WhiteSpace;
    overflow: Overflow;
  };
  borderRadius: BorderRadius;
  props: { key: string; value: string }[];
  /** Raw style entries from PSRT JSON (all keys). */
  sourceStyle: { key: string; value: string }[];
}

export const createBlock = (i: number, kind: BlockKind = "text"): TextBlock => ({
  id: crypto.randomUUID(),
  name: kind === "mask" ? `Cobertura ${i}` : `Texto ${i}`,
  kind,
  x: 50,
  y: 50,
  width: 30,
  size: 4,
  height: 5,
  content: kind === "mask" ? "" : "Novo texto",
  color: "#000000",
  colorSet: true,
  background: "#000000",
  backgroundSet: false,
  font: {
    family: "Inter",
    weight: 400,
    size: 16,
    sizeOverride: false,
    lineHeight: 1.4,
    letterSpacing: 0,
    italic: false,
    underline: false,
    strike: false,
    bold: false,
    transform: "none",
  },
  align: {
    horizontal: "left",
    vertical: "center",
    direction: "ltr",
    whiteSpace: "normal",
    overflow: "visible",
  },
  borderRadius: { ...ZERO_BORDER_RADIUS },
  props: [],
  sourceStyle: [],
});
