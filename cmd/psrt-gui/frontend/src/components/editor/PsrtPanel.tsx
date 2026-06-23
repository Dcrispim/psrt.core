import { useMemo, useRef } from "react";
import s from "./sidebar.module.css";

function escapeHtml(src: string) {
  return src.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function span(className: string, text: string): string {
  return `<span class=${className}>${text}</span>`;
}

/** Highlight without inserting quoted class attrs (breaks later string regex). */
function highlightLine(line: string): string {
  let l = escapeHtml(line);
  if (!l) return "&nbsp;";

  const slots: string[] = [];
  const stash = (html: string) => {
    const id = slots.length;
    slots.push(html);
    return `\x01${id}\x01`;
  };

  l = l.replace(/("[^"\\]*(?:\\.[^"\\]*)*")/g, (m) => stash(span(s.tokString, m)));
  l = l.replace(/(https?:\/\/[^\s|]+)/g, (m) => stash(span(s.tokUrl, m)));
  l = l.replace(/^(@\s+\w+)/, (m) => stash(span(s.tokKey, m)));
  l = l.replace(/^(&gt;&gt;[\d.\-]+)/, (m) => stash(span(s.tokNumber, m)));
  l = l.replace(/^(\$[A-Z]+)/, (m) => stash(span(s.tokDirective, m)));

  return l.replace(/\x01(\d+)\x01/g, (_, id) => slots[Number(id)] ?? "");
}

function highlight(src: string): string {
  return src.split("\n").map(highlightLine).join("\n");
}

export interface PsrtPanelProps {
  value: string;
  onChange: (value: string) => void;
  onFocus?: () => void;
  onBlur?: () => void;
}

export function PsrtPanel({ value, onChange, onFocus, onBlur }: PsrtPanelProps) {
  const preRef = useRef<HTMLPreElement>(null);
  const taRef = useRef<HTMLTextAreaElement>(null);
  const html = useMemo(() => highlight(value), [value]);
  const lines = value.split("\n").length;

  const onScroll = () => {
    if (preRef.current && taRef.current) {
      preRef.current.scrollTop = taRef.current.scrollTop;
      preRef.current.scrollLeft = taRef.current.scrollLeft;
    }
  };

  return (
    <div className={s.psrt}>
      <div className={s.psrtGutter} aria-hidden>
        {Array.from({ length: lines }, (_, i) => (
          <span key={i}>{i + 1}</span>
        ))}
      </div>
      <div className={s.psrtCode}>
        <pre
          ref={preRef}
          className={s.psrtPre}
          dangerouslySetInnerHTML={{ __html: html }}
        />
        <textarea
          ref={taRef}
          className={s.psrtTextarea}
          value={value}
          spellCheck={false}
          onFocus={onFocus}
          onBlur={onBlur}
          onChange={(e) => onChange(e.target.value)}
          onScroll={onScroll}
        />
      </div>
    </div>
  );
}
