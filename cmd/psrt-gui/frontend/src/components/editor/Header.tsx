import { useCallback, useEffect, useState, type ReactNode } from 'react';
import { useEditor } from '../../context/useEditor';
import { useAlertModal } from '../../context/AlertModalContext';
import { AddConstModal } from './AddConstModal';
import { AddFontModal } from './AddFontModal';
import { SaveOptionsModal, type SaveOption } from './SaveOptionsModal';
import s from './header.module.css';

export function Header() {
  const {
    document: editorDoc,
    state,
    openFile,
    newFile,
    save,
    saveAs,
    saveAsSvg,
    saveAsHtml,
    undo,
    redo,
    compileSvg,
    compileHtml,
    setAutoCompile,
    setActivePage,
    addPage,
    removePage,
    movePage,
    pageMoveRef,
    setPageMoveRef,
    addConst,
    addFont,
    importAssetsFolder,
    showToast,
  } = useEditor();
  const { confirm, prompt } = useAlertModal();
  const [saveModalOpen, setSaveModalOpen] = useState(false);
  const [constModalOpen, setConstModalOpen] = useState(false);
  const [fontModalOpen, setFontModalOpen] = useState(false);

  const pages = state?.pages ?? [];
  const activePage = state?.activePage ?? '';
  const autoCompile = state?.autoCompile ?? false;
  const hasDoc = Boolean(editorDoc);
  const existingConstNames = Object.keys(state?.consts ?? {});
  const existingFontUrls = state?.fonts ?? [];

  const onKey = useCallback(
    (e: KeyboardEvent) => {
      if (!e.ctrlKey) return;
      const key = e.key.toLowerCase();
      if (key === 'n') {
        e.preventDefault();
        newFile();
      }
      if (key === 'o') {
        e.preventDefault();
        openFile().catch((err) => showToast(String(err)));
      }
      if (key === 's' && e.shiftKey) {
        e.preventDefault();
        setSaveModalOpen(true);
      } else if (key === 's') {
        e.preventDefault();
        save().catch((err) => showToast(String(err)));
      }
      if (key === 'z') {
        e.preventDefault();
        undo();
      }
      if (key === 'y') {
        e.preventDefault();
        redo();
      }
    },
    [newFile, openFile, save, undo, redo, showToast],
  );

  const handleSaveOption = useCallback(
    (option: SaveOption, variantFiles: File[]) => {
      setSaveModalOpen(false);
      const run =
        option === 'psrt'
          ? saveAs()
          : option === 'svg'
            ? saveAsSvg()
            : saveAsHtml(variantFiles);
      run.catch((err) => showToast(String(err)));
    },
    [saveAs, saveAsSvg, saveAsHtml, showToast],
  );
  useEffect(() => {
    globalThis.document.addEventListener('keydown', onKey);
    return () => globalThis.document.removeEventListener('keydown', onKey);
  }, [onKey]);

  const handleAddPage = async () => {
    const previousPageName = pages.find((p) => p.name === activePage)?.name ?? '';
    const name = await prompt({
      title: 'Nova página',
      message: 'Informe o nome da página.',
      placeholder: 'nome-da-pagina',
      confirmLabel: 'Criar',
      defaultValue: previousPageName,
    });
    if (name?.trim()) addPage(name.trim());
  };

  const handleRemovePage = async () => {
    if (!activePage) return;
    const ok = await confirm({
      title: 'Remover página',
      message: `Remover a página "${activePage}"? Esta ação não pode ser desfeita facilmente.`,
      confirmLabel: 'Remover',
      danger: true,
    });
    if (ok) removePage();
  };

  return (
    <header className={s.root} role="toolbar" aria-label="Barra do editor">
      <div className={s.row}>
        <div className={s.brand}>
          <span className={s.logo} aria-hidden>
            ◆
          </span>
          <div className={s.brandText}>
            <strong>Block Editor</strong>
            <span className={s.subtitle}>
              {activePage || '—'}
              {activePage ? <em className={s.dot} title="Página ativa" /> : null}
            </span>
          </div>
        </div>

        <div className={s.group} aria-label="Arquivo">
          <button
            type="button"
            className={s.btn}
            title="Novo (Ctrl+N)"
            onClick={newFile}
          >
            <Icon name="new" /> <span>Novo</span>
          </button>
          <button
            type="button"
            className={s.btn}
            title="Abrir (Ctrl+O)"
            onClick={() => openFile().catch((err) => showToast(String(err)))}
          >
            <Icon name="folder" /> <span>Abrir</span>
          </button>
          <button
            type="button"
            className={`${s.btn} ${s.primary}`}
            title="Salvar (Ctrl+S)"
            disabled={!hasDoc}
            onClick={() => save().catch((err) => showToast(String(err)))}
          >
            <Icon name="save" /> <span>Salvar</span>
          </button>
          {import.meta.env.VITE_PRT_WEB === 'true' ? (
            <button
              type="button"
              className={s.btn}
              title="Vincular pasta de imagens locais (001.jpg, …)"
              disabled={!hasDoc}
              onClick={() =>
                importAssetsFolder().catch((err) => showToast(String(err)))
              }
            >
              <Icon name="folder" /> <span>Pasta assets</span>
            </button>
          ) : null}
          <button
            type="button"
            className={s.iconBtn}
            title="Salvar / exportar (Ctrl+Shift+S)"
            aria-label="Salvar como"
            disabled={!hasDoc}
            onClick={() => setSaveModalOpen(true)}
          >
            <Icon name="saveAs" />
          </button>
        </div>

        <div className={s.group} aria-label="Histórico">
          <button
            type="button"
            className={s.iconBtn}
            title="Desfazer (Ctrl+Z)"
            aria-label="Desfazer"
            disabled={!hasDoc}
            onClick={undo}
          >
            <Icon name="undo" />
          </button>
          <button
            type="button"
            className={s.iconBtn}
            title="Refazer (Ctrl+Y)"
            aria-label="Refazer"
            disabled={!hasDoc}
            onClick={redo}
          >
            <Icon name="redo" />
          </button>
        </div>

        <div className={s.spacer} />

        <div className={`${s.group} ${s.segmented}`} aria-label="Compilar">
          <button
            type="button"
            className={s.segBtn}
            title="Compilar SVG"
            disabled={!hasDoc}
            onClick={() => compileSvg().catch((err) => showToast(String(err)))}
          >
            <Icon name="svg" /> SVG
          </button>
          <button
            type="button"
            className={s.segBtn}
            title="Compilar HTML"
            disabled={!hasDoc}
            onClick={() => compileHtml().catch((err) => showToast(String(err)))}
          >
            <Icon name="html" /> HTML
          </button>
        </div>

        <label className={s.switch} title="Recompila automaticamente ao editar">
          <input
            type="checkbox"
            checked={autoCompile}
            disabled={!hasDoc}
            onChange={(e) => setAutoCompile(e.target.checked)}
          />
          <span className={s.track}>
            <span className={s.thumb} />
          </span>
          <span className={s.switchLabel}>Auto</span>
        </label>
      </div>

      <div className={`${s.row} ${s.rowSub}`}>
        <div className={s.field}>
          <label className={s.fieldLabel} htmlFor="page-select">
            Página
          </label>
          <select
            id="page-select"
            className={s.select}
            value={activePage}
            disabled={pages.length === 0}
            onChange={(e) => setActivePage(e.target.value)}
          >
            {pages.length === 0 ? (
              <option value="">—</option>
            ) : (
              pages.map((p) => (
                <option key={p.name} value={p.name}>
                  {p.name}
                </option>
              ))
            )}
          </select>
        </div>

        <div className={s.group} aria-label="Gerenciar páginas">
          <button
            type="button"
            className={s.btn}
            title="Adicionar página"
            disabled={!state}
            onClick={handleAddPage}
          >
            <Icon name="plus" /> <span>Nova</span>
          </button>
          <button
            type="button"
            className={s.btn}
            title="Adicionar constante ($CONSTS)"
            disabled={!hasDoc}
            onClick={() => setConstModalOpen(true)}
          >
            <Icon name="const" /> <span>Const</span>
          </button>
          <button
            type="button"
            className={s.btn}
            title="Adicionar fonte ($FONTS)"
            disabled={!hasDoc}
            onClick={() => setFontModalOpen(true)}
          >
            <Icon name="font" /> <span>Font</span>
          </button>
          <button
            type="button"
            className={`${s.iconBtn} ${s.danger}`}
            title="Remover página atual"
            aria-label="Remover página"
            disabled={!activePage}
            onClick={handleRemovePage}
          >
            <Icon name="trash" />
          </button>
        </div>

        <div className={s.spacer} />

        <div className={s.field}>
          <label className={s.fieldLabel} htmlFor="ref-page">
            Referência
          </label>
          <div className={s.inputGroup}>
            <Icon name="link" />
            <input
              id="ref-page"
              className={s.input}
              placeholder="página de referência"
              value={pageMoveRef}
              disabled={!hasDoc}
              onChange={(e) => setPageMoveRef(e.target.value)}
            />
          </div>
        </div>

        <div className={`${s.group} ${s.segmented}`} aria-label="Reordenar">
          <button
            type="button"
            className={s.segBtn}
            title="Mover antes da referência"
            disabled={!activePage || !pageMoveRef}
            onClick={() => movePage(pageMoveRef, true)}
          >
            <Icon name="arrowUp" /> Antes
          </button>
          <button
            type="button"
            className={s.segBtn}
            title="Mover depois da referência"
            disabled={!activePage || !pageMoveRef}
            onClick={() => movePage(pageMoveRef, false)}
          >
            <Icon name="arrowDown" /> Depois
          </button>
        </div>
      </div>

      <SaveOptionsModal
        open={saveModalOpen}
        onSave={handleSaveOption}
        onCancel={() => setSaveModalOpen(false)}
      />
      <AddConstModal
        open={constModalOpen}
        existingNames={existingConstNames}
        onAdd={(name, value) => {
          addConst(name, value);
          setConstModalOpen(false);
          showToast(`Constante @${name}@ adicionada`);
        }}
        onCancel={() => setConstModalOpen(false)}
      />
      <AddFontModal
        open={fontModalOpen}
        existingUrls={existingFontUrls}
        onAdd={(url) => {
          addFont(url);
          setFontModalOpen(false);
          showToast('Fonte adicionada');
        }}
        onCancel={() => setFontModalOpen(false)}
      />
    </header>
  );
}

function Icon({ name }: { name: string }) {
  const paths: Record<string, ReactNode> = {
    new: (
      <>
        <path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
        <path d="M14 3v6h6M12 11v6M9 14h6" />
      </>
    ),
    folder: (
      <path d="M3 5a2 2 0 0 1 2-2h3l2 2h7a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
    ),
    save: (
      <>
        <path d="M5 3h11l3 3v13a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2z" />
        <path d="M7 3v6h9V3M7 21v-7h10v7" />
      </>
    ),
    saveAs: (
      <>
        <path d="M5 3h9l3 3v8" />
        <path d="M3 19l4-1 9-9-3-3-9 9z" />
      </>
    ),
    undo: <path d="M9 14l-4-4 4-4M5 10h9a5 5 0 0 1 0 10h-3" />,
    redo: <path d="M15 14l4-4-4-4M19 10h-9a5 5 0 0 0 0 10h3" />,
    svg: (
      <>
        <path d="M4 4h16v16H4z" />
        <path d="M7 14c1 1 2 1 3 0M14 10l3 4-3 4" />
      </>
    ),
    html: (
      <>
        <path d="M4 4l1.5 16L12 22l6.5-2L20 4z" />
        <path d="M8 8h8l-.5 4H9l.3 3 2.7.8 2.7-.8.2-2" />
      </>
    ),
    plus: <path d="M12 5v14M5 12h14" />,
    trash: (
      <>
        <path d="M4 7h16M9 7V4h6v3M6 7l1 13h10l1-13" />
      </>
    ),
    link: (
      <>
        <path d="M10 14a4 4 0 0 0 5.7 0l3-3a4 4 0 0 0-5.7-5.7l-1 1" />
        <path d="M14 10a4 4 0 0 0-5.7 0l-3 3a4 4 0 0 0 5.7 5.7l1-1" />
      </>
    ),
    arrowUp: <path d="M12 19V5M6 11l6-6 6 6" />,
    arrowDown: <path d="M12 5v14M6 13l6 6 6-6" />,
    const: (
      <>
        <path d="M4 7h16M4 12h10M4 17h6" />
        <path d="M16 12h4M18 10v4" />
      </>
    ),
    font: (
      <>
        <path d="M4 20h16M6 4v8l6-4 6 4V4" />
      </>
    ),
  };

  return (
    <svg
      width="14"
      height="14"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden
    >
      {paths[name]}
    </svg>
  );
}

