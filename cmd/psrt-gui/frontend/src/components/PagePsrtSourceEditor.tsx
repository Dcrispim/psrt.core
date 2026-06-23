import { useCallback, useEffect, useRef, useState } from 'react';
import * as api from '@wails/go/main/GUIApp';
import { useEditor } from '../context/useEditor';

const APPLY_MS = 400;

export function PagePsrtSourceEditor({ active }: { active: boolean }) {
  const { state, document, applyPagePsrtSource, beginEdit, endEdit, showToast } =
    useEditor();
  const pageName = state?.activePage ?? '';
  const [text, setText] = useState('');
  const dirty = useRef(false);
  const focused = useRef(false);
  const applyTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pending = useRef<string | null>(null);

  const load = useCallback(async () => {
    if (!document || !pageName) {
      setText('');
      return;
    }
    try {
      const src = await api.FormatPageDocumentJSON(
        JSON.stringify(document),
        pageName,
      );
      setText(src);
    } catch (e) {
      showToast(String(e));
    }
  }, [document, pageName, showToast]);

  useEffect(() => {
    if (!active) return;
    if (dirty.current || focused.current) return;
    void load();
  }, [active, load]);

  const flush = useCallback(async () => {
    const body = pending.current;
    pending.current = null;
    if (body === null || !pageName) return;
    try {
      await applyPagePsrtSource(body);
    } catch (e) {
      showToast(String(e));
    }
  }, [applyPagePsrtSource, pageName, showToast]);

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

  if (!state?.page) {
    return <p className="psrt-source-empty">Selecione uma página</p>;
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
