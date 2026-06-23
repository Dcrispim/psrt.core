import s from "./sidebar.module.css";
import type { ReactNode } from "react";

export function SliderField({
  label, value, min, max, step = 0.001, onChange, disabled = false,
}: {
  label: string;
  value: number;
  min: number;
  max: number;
  step?: number;
  onChange: (v: number) => void;
  disabled?: boolean;
}) {
  return (
    <div className={s.field}>
      <label className={s.label}>{label}</label>
      <div className={s.sliderRow}>
        <input
          type="range"
          className={s.slider}
          min={min} max={max} step={step}
          value={value}
          disabled={disabled}
          onChange={(e) => onChange(Number(e.target.value))}
        />
        <input
          type="number"
          className={s.numInput}
          min={min} max={max} step={step}
          value={value}
          disabled={disabled}
          onChange={(e) => onChange(Number(e.target.value))}
        />
      </div>
    </div>
  );
}

export function ColorField({
  label, value, onChange,
}: { label: string; value: string; onChange: (v: string) => void }) {
  return (
    <div className={s.field}>
      <label className={s.label}>{label}</label>
      <div className={s.colorRow}>
        <label className={s.swatch} style={{ background: value }}>
          <input type="color" value={value} onChange={(e) => onChange(e.target.value)} />
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
}

export function ToggleGroup<T extends string>({
  label, value, options, onChange,
}: {
  label?: string;
  value: T;
  options: { value: T; label: ReactNode; title?: string }[];
  onChange: (v: T) => void;
}) {
  return (
    <div className={s.field}>
      {label && <label className={s.label}>{label}</label>}
      <div className={s.toggleGroup} role="group">
        {options.map((o) => (
          <button
            key={o.value}
            type="button"
            className={`${s.toggle} ${value === o.value ? s.toggleActive : ""}`}
            aria-pressed={value === o.value}
            title={o.title}
            onClick={() => onChange(o.value)}
          >
            {o.label}
          </button>
        ))}
      </div>
    </div>
  );
}

export function MultiToggle({
  label, options,
}: {
  label?: string;
  options: { active: boolean; label: ReactNode; title?: string; onToggle: () => void }[];
}) {
  return (
    <div className={s.field}>
      {label && <label className={s.label}>{label}</label>}
      <div className={s.toggleGroup} role="group">
        {options.map((o, i) => (
          <button
            key={i}
            type="button"
            className={`${s.toggle} ${o.active ? s.toggleActive : ""}`}
            aria-pressed={o.active}
            title={o.title}
            onClick={o.onToggle}
          >
            {o.label}
          </button>
        ))}
      </div>
    </div>
  );
}

export function SelectField<T extends string>({
  label, value, options, onChange,
}: { label: string; value: T; options: { value: T; label: string }[]; onChange: (v: T) => void }) {
  return (
    <div className={s.field}>
      <label className={s.label}>{label}</label>
      <select
        className={s.select}
        value={value}
        onChange={(e) => onChange(e.target.value as T)}
      >
        {options.map((o) => (
          <option key={o.value} value={o.value}>{o.label}</option>
        ))}
      </select>
    </div>
  );
}

export function NumberField({
  label, value, onChange, step = 1, suffix,
}: { label: string; value: number; onChange: (v: number) => void; step?: number; suffix?: string }) {
  return (
    <div className={s.field}>
      <label className={s.label}>{suffix ? `${label} (${suffix})` : label}</label>
      <input
        type="number"
        className={s.input}
        value={value}
        step={step}
        onChange={(e) => onChange(Number(e.target.value))}
      />
    </div>
  );
}
