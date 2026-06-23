import { useEditor } from '../context/useEditor';
import { parseStyle, styleStringValue, toColorInput } from '../lib/parseStyle';

export type StylePatch = {
  styleSet?: Record<string, string>;
  styleRemove?: string[];
};

export function StyleFields({
  styleRaw,
  onPatch,
}: {
  styleRaw: string;
  onPatch: (patch: StylePatch) => void;
}) {
  const { propertyMap } = useEditor();
  const style = parseStyle(styleRaw);
  const known = new Set<string>();

  return (
    <>
      {(propertyMap['color-picker'] ?? []).map((key) => {
        if (!(key in style)) return null;
        known.add(key);
        return (
          <ColorField
            key={key}
            label={key}
            value={toColorInput(style[key])}
            onChange={(hex) => onPatch({ styleSet: { [key]: hex } })}
          />
        );
      })}
      {Object.entries(propertyMap.dropdown ?? {}).map(([key, opts]) => {
        if (!(key in style)) return null;
        known.add(key);
        const val = styleStringValue(style[key]);
        return (
          <div key={key} className="field">
            <label>{key}</label>
            <select
              value={val}
              onChange={(e) => onPatch({ styleSet: { [key]: e.target.value } })}
            >
              {opts.map((opt) => (
                <option key={opt} value={opt}>
                  {opt}
                </option>
              ))}
            </select>
          </div>
        );
      })}
      {Object.entries(style).map(([key, val]) => {
        if (known.has(key)) return null;
        return (
          <div key={key} className="field custom-style">
            <label>{key}</label>
            <input
              type="text"
              defaultValue={typeof val === 'string' ? styleStringValue(val) : JSON.stringify(val)}
              onChange={(e) => onPatch({ styleSet: { [key]: e.target.value } })}
            />
            <button type="button" onClick={() => onPatch({ styleRemove: [key] })}>
              Remove
            </button>
          </div>
        );
      })}
      <button
        type="button"
        onClick={() => {
          const key = prompt('Property name');
          if (key) onPatch({ styleSet: { [key]: '#ffffff' } });
        }}
      >
        Add property
      </button>
    </>
  );
}

function ColorField({
  label,
  value,
  onChange,
}: {
  label: string;
  value: string;
  onChange: (hex: string) => void;
}) {
  return (
    <div className="field">
      <label>{label}</label>
      <input type="color" value={value} onChange={(e) => onChange(e.target.value)} />
    </div>
  );
}
