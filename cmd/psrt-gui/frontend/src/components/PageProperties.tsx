import { useEditor } from '../context/useEditor';
import { StyleFields } from './StyleFields';

export function PageProperties() {
  const { state, patchPage, refreshPageImage, showToast } = useEditor();
  const page = state?.page;
  if (!page) return null;

  return (
    <>
      <h3>Page</h3>
      <div className="field">
        <label htmlFor="page-name">Name</label>
        <input
          id="page-name"
          type="text"
          defaultValue={page.name}
          onChange={(e) => patchPage({ name: e.target.value })}
        />
      </div>
      <div className="field">
        <label htmlFor="page-image-url">Image URL</label>
        <input
          id="page-image-url"
          type="text"
          defaultValue={page.imageUrl}
          onChange={(e) => patchPage({ imageUrl: e.target.value })}
        />
        <button
          type="button"
          onClick={() => refreshPageImage().catch((e) => showToast(String(e)))}
        >
          Refresh image
        </button>
      </div>

      <h3>Page style</h3>
      <StyleFields styleRaw={page.style} onPatch={(patch) => patchPage(patch)} />
    </>
  );
}
