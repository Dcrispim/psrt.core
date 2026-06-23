import type { ReactNode } from "react";
import s from "./sidebar.module.css";
import { SNAP } from "../../lib/applyStyle";

export type PropHandlerName = "color-picker" | "slider" | "select";

export interface PropHandlerContext {
  value: string;
  onChange: (value: string) => void;
}

export type CssPropHandler = (ctx: PropHandlerContext) => ReactNode;
export type CssPropHandlerFactory = () => CssPropHandler;

export interface SliderHandlerOptions {
  min?: number;
  max?: number;
  step?: number;
  unit?: string;
}

const DEFAULT_SLIDER_STEP = SNAP;

function parseCSSValue(value: string): { num: number; unit: string } {
  const trimmed = value.trim();
  if (!trimmed) return { num: 0, unit: "px" };
  const m = trimmed.match(/^(-?[\d.]+)(.*)$/);
  if (m) return { num: Number(m[1]), unit: m[2] || "px" };
  return { num: 0, unit: "px" };
}

function toHexColor(value: string): string {
  const v = value.trim();
  if (/^#[0-9a-f]{6}$/i.test(v)) return v;
  if (/^#[0-9a-f]{3}$/i.test(v)) {
    const [, r, g, b] = v;
    return `#${r}${r}${g}${g}${b}${b}`;
  }
  return "#000000";
}

const colorPickerHandler: CssPropHandler = ({ value, onChange }) => {
  const hex = toHexColor(value);
  return (
    <div className={s.propValueCell}>
      <div className={s.colorRow}>
        <label className={s.swatch} style={{ background: hex }}>
          <input type="color" value={hex} onChange={(e) => onChange(e.target.value)} />
        </label>
        <input
          className={s.hexInput}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          spellCheck={false}
        />
      </div>
    </div>
  );
};

function sliderHandler(opts: SliderHandlerOptions = {}): CssPropHandler {
  const { min = 0, max = 100, step = DEFAULT_SLIDER_STEP, unit } = opts;
  return ({ value, onChange }) => {
    const { num, unit: parsedUnit } = parseCSSValue(value);
    const suffix = unit ?? parsedUnit;
    return (
      <div className={s.propValueCell}>
        <div className={s.propSliderRow}>
          <input
            type="range"
            className={s.slider}
            min={min}
            max={max}
            step={step}
            value={num}
            onChange={(e) => onChange(`${e.target.value}${suffix}`)}
          />
          <input
            type="number"
            className={s.numInput}
            min={min}
            max={max}
            step={step}
            value={num}
            onChange={(e) => onChange(`${e.target.value}${suffix}`)}
          />
        </div>
      </div>
    );
  };
}

function selectHandler(options: string[]): CssPropHandler {
  const opts = options.map((o) => ({ value: o, label: o }));
  const fallback = options[0] ?? "";
  return ({ value, onChange }) => (
    <div className={s.propValueCell}>
      <select
        className={s.select}
        value={options.includes(value) ? value : fallback}
        onChange={(e) => onChange(e.target.value)}
      >
        {opts.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    </div>
  );
}

type PropHandlerMap = {
  "color-picker": () => CssPropHandler;
  slider: (opts?: SliderHandlerOptions) => CssPropHandlerFactory;
  select: (options: string[]) => CssPropHandlerFactory;
};

export const Handlers: PropHandlerMap = {
  "color-picker": () => colorPickerHandler,
  slider: (opts?: SliderHandlerOptions) => () => sliderHandler(opts),
  select: (options: string[]) => () => selectHandler(options),
};

// function pxSlider(max = 128): CssPropHandlerFactory {
//   return Handlers.slider({ min: 0, max, step: DEFAULT_SLIDER_STEP });
// }
function percentSlider(max = 100): CssPropHandlerFactory {
  return Handlers.slider({ min: 0, max, step: DEFAULT_SLIDER_STEP, unit: "%" });
} 

function unitlessSlider(max: number): CssPropHandlerFactory {
  return Handlers.slider({ min: 0, max, step: DEFAULT_SLIDER_STEP, unit: "" });
}

export const CSS_PROP_HANDLER_MAP: Record<string, CssPropHandlerFactory> = {
  padding: percentSlider(20),
  "padding-top": percentSlider(),
  "padding-right": percentSlider(),
  "padding-bottom": percentSlider(),
  "padding-left": percentSlider(),
  margin: percentSlider(),
  "margin-top": percentSlider(),
  "margin-right": percentSlider(),
  "margin-bottom": percentSlider(),
  "margin-left": percentSlider(),
  gap: percentSlider(64),
  "row-gap": percentSlider(64),
  "column-gap": percentSlider(64),
  "font-size": percentSlider(96),
  "border-width": percentSlider(32),
  "border-radius": percentSlider(64),
  width: percentSlider(100),
  height: percentSlider(100),
  opacity: unitlessSlider(1),
  color: Handlers["color-picker"],
  background: Handlers["color-picker"],
  "background-color": Handlers["color-picker"],
  "border-color": Handlers["color-picker"],
  "font-weight": Handlers.slider({ min: 100, max: 900, step: 100, unit: "" }),
  "font-style": Handlers.select(["normal", "italic", "oblique"]),
  "text-decoration": Handlers.select([
    "none",
    "underline",
    "line-through",
    "underline line-through",
  ]),
  "text-transform": Handlers.select(["none", "uppercase", "lowercase", "capitalize"]),
  "text-overflow": Handlers.select(["clip", "ellipsis"]),
  direction: Handlers.select(["ltr", "rtl"]),
  display: Handlers.select(["flex", "block", "inline", "inline-block", "none"]),
  position: Handlers.select(["static", "relative", "absolute", "fixed", "sticky"]),
  overflow: Handlers.select(["visible", "hidden", "scroll", "auto", "clip"]),
  "flex-direction": Handlers.select(["row", "row-reverse", "column", "column-reverse"]),
  "justify-content": Handlers.select([
    "flex-start",
    "flex-end",
    "center",
    "space-between",
    "space-around",
    "space-evenly",
  ]),
  "align-items": Handlers.select([
    "flex-start",
    "flex-end",
    "center",
    "stretch",
    "baseline",
  ]),
  "text-align": Handlers.select(["left", "center", "right", "justify"]),
  "box-sizing": Handlers.select(["content-box", "border-box"]),
  "white-space": Handlers.select([
    "normal",
    "nowrap",
    "pre",
    "pre-wrap",
    "pre-line",
    "break-spaces",
  ]),
};

export const CSS_PROP_KEYS = Object.keys(CSS_PROP_HANDLER_MAP).sort();

export function resolveCssPropHandler(key: string): CssPropHandler | undefined {
  const factory = CSS_PROP_HANDLER_MAP[key.trim().toLowerCase()];
  return factory?.();
}
