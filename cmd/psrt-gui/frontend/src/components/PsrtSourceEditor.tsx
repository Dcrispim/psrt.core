import { useCallback, useEffect, useRef, useState } from 'react';
import * as api from '@wails/go/main/GUIApp';
import { useEditor } from '../context/useEditor';
import { saveLastPsrt } from '../lib/localPsrt';

const APPLY_MS = 400;

export function PsrtSourceEditor({ active }: { active: boolean }) {
  const { state, document, applyPsrtSource, beginEdit, endEdit, showToast } = useEditor();
  const [text, setText] = useState('');
  const dirty = useRef(false);
  const focused = useRef(false);
  const applyTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pending = useRef<string | null>(null);

  const load = useCallback(async () => {
    if (!document || !state?.filePath) {
      setText('');
      return;
    }
    try {
      const src = await api.FormatDocumentJSON(JSON.stringify(document));
      setText(src);
      saveLastPsrt(state.filePath, src);
    } catch (e) {
      showToast(String(e));
    }
  }, [document, state?.filePath, showToast]);

  useEffect(() => {
    if (!active) return;
    if (dirty.current || focused.current) return;
    void load();
  }, [active, document, load]);

  const flush = useCallback(async () => {
    const body = pending.current;
    pending.current = null;
    if (body === null) return;
    try {
      await applyPsrtSource(body);
    } catch (e) {
      showToast(String(e));
    }
  }, [applyPsrtSource, showToast]);

  const scheduleApply = useCallback(
    (value: string) => {
      pending.current = value;
      if (applyTimer.current) clearTimeout(applyTimer.current);
      applyTimer.current = setTimeout(() => {
        applyTimer.current = null;
        void flush();
      }, APPLY_MS);
    },
    [flush],
  );

  const onChange = (value: string) => {
    dirty.current = true;
    setText(value);
    if (state?.filePath) saveLastPsrt(state.filePath, value);
    scheduleApply(value);
  };

  const handleBlur = async () => {
    focused.current = false;
    if (applyTimer.current) {
      clearTimeout(applyTimer.current);
      applyTimer.current = null;
    }
    await flush();
    endEdit();
    dirty.current = false;
    if (active) await load();
  };

  if (!state?.filePath) {
    return <p className="psrt-source-empty">Abra um arquivo .psrt</p>;
  }

  return (
    <div className="psrt-source-wrap">
      <textarea
        className="psrt-source-textarea"
        value={text}
        spellCheck={false}
        onFocus={() => {
          focused.current = true;
          beginEdit();
        }}
        onBlur={() => {
          void handleBlur();
        }}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
