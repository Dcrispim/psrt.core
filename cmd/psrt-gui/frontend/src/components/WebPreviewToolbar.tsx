import { useEditor } from '../context/useEditor';
import { WebZoomControl } from './WebZoomControl';

export function WebPreviewToolbar() {
  const { webTextsVisible, toggleWebTextsVisible } = useEditor();

  return (
    <div className="web-preview-toolbar">
      <WebZoomControl />
      <button
        type="button"
        className={`web-texts-toggle-btn${webTextsVisible ? '' : ' active'}`}
        aria-pressed={!webTextsVisible}
        title={webTextsVisible ? 'Ocultar textos (Ctrl+H)' : 'Mostrar textos (Ctrl+H)'}
        onClick={toggleWebTextsVisible}
      >
        {webTextsVisible ? 'Ocultar textos' : 'Mostrar textos'}
      </button>
    </div>
  );
}
