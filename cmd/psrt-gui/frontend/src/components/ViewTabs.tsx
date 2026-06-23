import { useEditor } from '../context/useEditor';
import type { PreviewTab } from '../context/EditorContext';

const TABS: { id: PreviewTab; label: string }[] = [
  { id: 'svg', label: 'SVG' },
  { id: 'web', label: 'WEB' },
  { id: 'html', label: 'HTML' },
];

export function ViewTabs() {
  const { previewTab, setPreviewTab } = useEditor();

  return (
    <div className="view-tabs">
      {TABS.map(({ id, label }) => (
        <button
          key={id}
          type="button"
          className={previewTab === id ? 'active' : ''}
          onClick={() => setPreviewTab(id)}
        >
          {label}
        </button>
      ))}
    </div>
  );
}
