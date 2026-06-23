import { useEditor } from '../context/useEditor';
import { Sidebar as EditorSidebar } from './editor/Sidebar';
import { useTextEditorSidebar } from '../hooks/useTextEditorSidebar';
import { useResizablePanelWidth } from '../hooks/useResizablePanelWidth';
import { PageProperties } from './PageProperties';
import { PsrtSourceEditor } from './PsrtSourceEditor';
import { useState, type ReactNode } from 'react';

type PanelMode = 'properties' | 'source';

const PANEL_STORAGE_KEY = 'psrt-properties-panel-width';
const PANEL_MIN_WIDTH = 280;
const PANEL_DEFAULT_WIDTH = 360;

function panelMaxWidth(): number {
  return Math.floor(window.innerWidth * 0.55);
}

function PropertiesAside({
  className,
  children,
}: {
  className?: string;
  children: ReactNode;
}) {
  const { width, isResizing, onResizeMouseDown } = useResizablePanelWidth({
    storageKey: PANEL_STORAGE_KEY,
    defaultWidth: PANEL_DEFAULT_WIDTH,
    minWidth: PANEL_MIN_WIDTH,
    maxWidth: panelMaxWidth,
  });

  const asideClass = [
    'properties',
    className,
    isResizing ? 'properties--resizing' : '',
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <aside className={asideClass} style={{ width }}>
      <div
        className="properties-resize-handle"
        role="separator"
        aria-orientation="vertical"
        aria-label="Redimensionar painel"
        onMouseDown={onResizeMouseDown}
      />
      {children}
    </aside>
  );
}

export function PropertiesPanel() {
  const { state, previewTab } = useEditor();
  const [mode, setMode] = useState<PanelMode>('properties');
  const isTextMode = previewTab === 'web';
  const textSidebar = useTextEditorSidebar();

  if (isTextMode) {
    return (
      <PropertiesAside className="properties--editor">
        {textSidebar ? (
          <EditorSidebar {...textSidebar} />
        ) : (
          <p className="properties-empty">Selecione uma página</p>
        )}
      </PropertiesAside>
    );
  }

  return (
    <PropertiesAside>
      <div className="properties-tabs" role="tablist">
        <button
          type="button"
          role="tab"
          className={mode === 'properties' ? 'active' : ''}
          aria-selected={mode === 'properties'}
          onClick={() => setMode('properties')}
        >
          Propriedades
        </button>
        <button
          type="button"
          role="tab"
          className={mode === 'source' ? 'active' : ''}
          aria-selected={mode === 'source'}
          onClick={() => setMode('source')}
        >
          PSRT
        </button>
      </div>
      {mode === 'properties' ? (
        <div className="properties-body">
          {state?.page ? <PageProperties /> : null}
        </div>
      ) : (
        <PsrtSourceEditor active={mode === 'source'} />
      )}
    </PropertiesAside>
  );
}
