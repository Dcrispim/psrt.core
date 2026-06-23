import { useEditor } from '../context/useEditor';

export function PageThumbnailsSidebar() {
  const { state, thumbs, setActivePage } = useEditor();
  const pages = state?.pages ?? [];
  const activePage = state?.activePage ?? '';

  if (pages.length === 0) {
    return <aside className="page-thumbs-sidebar" aria-label="Páginas" />;
  }

  return (
    <aside className="page-thumbs-sidebar" aria-label="Páginas">
      <div className="page-thumbs-grid">
        {pages.map((page) => {
          const url = page.imageUrl ?? '';
          const thumb =
            thumbs[page.name] ?? (/^https?:\/\//i.test(url) ? url : null);
          const isActive = page.name === activePage;

          return (
            <button
              key={page.name}
              type="button"
              className={`page-thumb${isActive ? ' page-thumb--active' : ''}`}
              title={page.name}
              aria-current={isActive ? 'page' : undefined}
              onClick={() => setActivePage(page.name)}
            >
              {thumb ? (
                <img
                  className="page-thumb__img"
                  src={thumb}
                  alt=""
                  draggable={false}
                />
              ) : (
                <span className="page-thumb__placeholder" aria-hidden />
              )}
              <span className="page-thumb__label">{page.name}</span>
            </button>
          );
        })}
      </div>
    </aside>
  );
}
