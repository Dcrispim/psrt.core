import { useCallback, useEffect } from 'react';
import { useEditor } from '../context/useEditor';

export function Toolbar() {
  const {
    openFile,
    save,
    saveAs,
    undo,
    redo,
    compileSvg,
    compileHtml,
    setAutoCompile,
    showToast,
  } = useEditor();

  const onKey = useCallback(
    (e: KeyboardEvent) => {
      if (e.ctrlKey && e.key === 'z') {
        e.preventDefault();
        undo();
      }
      if (e.ctrlKey && e.key === 'y') {
        e.preventDefault();
        redo();
      }
    },
    [undo, redo, showToast],
  );

  useEffect(() => {
    document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, [onKey]);

  return (
    <header className="toolbar">
      <button type="button" onClick={() => openFile().catch((e) => showToast(String(e)))}>
        Open
      </button>
      <button type="button" onClick={() => save().catch((e) => showToast(String(e)))}>
        Save
      </button>
      <button type="button" onClick={() => saveAs().catch((e) => showToast(String(e)))}>
        Save As
      </button>
      <button type="button" onClick={undo}>
        Undo
      </button>
      <button type="button" onClick={redo}>
        Redo
      </button>
      <span className="sep" />
      <button type="button" onClick={() => compileSvg().catch((e) => showToast(String(e)))}>
        Compile SVG
      </button>
      <button type="button" onClick={() => compileHtml().catch((e) => showToast(String(e)))}>
        Compile HTML
      </button>
      <label>
        <input
          type="checkbox"
          defaultChecked={false}
          onChange={(e) => setAutoCompile(e.target.checked)}
        />{' '}
        Auto compile
      </label>
    </header>
  );
}
