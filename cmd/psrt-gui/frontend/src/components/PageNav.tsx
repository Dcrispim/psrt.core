import { useEditor } from '../context/useEditor';

export function PageNav() {
  const {
    state,
    thumbs,
    setActivePage,
    addPage,
    removePage,
    movePage,
    pageMoveRef,
    setPageMoveRef,
  } = useEditor();

  const pages = state?.pages ?? [];

  return (
    <nav className="page-nav">
      <label>
        Page{' '}
        <select
          className="page-select"
          value={state?.activePage ?? ''}
          onChange={(e) => setActivePage(e.target.value)}
        >
          {pages.map((p) => (
            <option
              key={p.name}
              value={p.name}
              style={
                thumbs[p.name]
                  ? {
                      backgroundImage: `url(${thumbs[p.name]})`,
                      backgroundSize: '20px auto',
                      backgroundRepeat: 'no-repeat',
                      backgroundPosition: '2px center',
                      paddingLeft: '24px',
                    }
                  : undefined
              }
            >
              {p.name}
            </option>
          ))}
        </select>
      </label>
      <button
        type="button"
        onClick={() => {
          const name = prompt('Page name');
          if (name) addPage(name);
        }}
      >
        Add Page
      </button>
      <button
        type="button"
        onClick={() => {
          if (state?.activePage && confirm(`Remove page ${state.activePage}?`)) {
            removePage();
          }
        }}
      >
        Remove
      </button>
      <input
        placeholder="ref page"
        value={pageMoveRef}
        onChange={(e) => setPageMoveRef(e.target.value)}
      />
      <button
        type="button"
        onClick={() => movePage(pageMoveRef, true)}
      >
        Move Before
      </button>
      <button
        type="button"
        onClick={() => movePage(pageMoveRef, false)}
      >
        Move After
      </button>
    </nav>
  );
}
