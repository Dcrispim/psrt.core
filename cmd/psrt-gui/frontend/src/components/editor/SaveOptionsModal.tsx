import { useEffect, useId, useState, type MouseEvent, type ReactNode } from 'react';
import s from './saveOptionsModal.module.css';

export type SaveOption = 'psrt' | 'svg' | 'html';

const FORMATS: { id: SaveOption; label: string; hint: string }[] = [
  { id: 'psrt', label: 'Salvar PSRT', hint: 'Formato nativo editável (.psrt)' },
  { id: 'svg', label: 'Salvar como SVG', hint: 'Vetor exportado para web' },
  { id: 'html', label: 'Salvar como HTML', hint: 'Variantes embutidas; Ctrl+L no preview' },
];

export interface SaveOptionsModalProps {
  open: boolean;
  onSave: (option: SaveOption, variants: File[]) => void;
  onCancel: () => void;
}

const emptyVariants = (): Record<SaveOption, File[]> => ({
  psrt: [],
  svg: [],
  html: [],
});

export function SaveOptionsModal({ open, onSave, onCancel }: SaveOptionsModalProps) {
  const titleId = useId();
  const [expanded, setExpanded] = useState<SaveOption | null>('psrt');
  const [variants, setVariants] = useState<Record<SaveOption, File[]>>(emptyVariants);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onCancel();
    };
    document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, [open, onCancel]);

  useEffect(() => {
    if (!open) {
      setExpanded('psrt');
      setVariants(emptyVariants());
    }
  }, [open]);

  if (!open) return null;

  const stopDialogClick = (e: MouseEvent) => e.stopPropagation();

  const addFiles = (fmt: SaveOption, files: FileList | null) => {
    if (!files?.length) return;
    const psrtOnly = Array.from(files).filter((f) =>
      f.name.toLowerCase().endsWith('.psrt'),
    );
    if (psrtOnly.length === 0) return;
    setVariants((v) => ({ ...v, [fmt]: [...v[fmt], ...psrtOnly] }));
  };

  const removeFile = (fmt: SaveOption, idx: number) => {
    setVariants((v) => ({ ...v, [fmt]: v[fmt].filter((_, i) => i !== idx) }));
  };

  return (
    <div
      className={s.overlay}
      role="presentation"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onCancel();
      }}
    >
      <div
        className={s.dialog}
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        onMouseDown={stopDialogClick}
      >
        <header className={s.header}>
          <h2 id={titleId} className={s.title}>
            Salvar
          </h2>
          <button type="button" className={s.close} onClick={onCancel} aria-label="Fechar">
            <Icon name="x" />
          </button>
        </header>

        <ul className={s.list}>
          {FORMATS.map((f) => {
            const isOpen = expanded === f.id;
            const count = variants[f.id].length;
            return (
              <li key={f.id} className={`${s.item} ${isOpen ? s.itemOpen : ''}`}>
                <button
                  type="button"
                  className={s.itemHead}
                  onClick={() => setExpanded(isOpen ? null : f.id)}
                  aria-expanded={isOpen}
                >
                  <span className={s.itemMain}>
                    <span className={s.itemLabel}>{f.label}</span>
                    <span className={s.itemHint}>{f.hint}</span>
                  </span>
                  <span className={s.itemRight}>
                    {count > 0 && (
                      <span className={s.badge} title={`${count} variante(s)`}>
                        {count}
                      </span>
                    )}
                    <span className={`${s.chev} ${isOpen ? s.chevOpen : ''}`}>
                      <Icon name="chev" />
                    </span>
                  </span>
                </button>

                {isOpen && (
                  <div className={s.body}>
                    {f.id === 'html' && (
                      <p className={s.variantNote}>
                        O documento atual é a variante principal. Arquivos extras entram no HTML
                        compilado.
                      </p>
                    )}

                    {f.id === 'html' && variants[f.id].length > 0 && (
                      <ul className={s.files}>
                        {variants[f.id].map((file, i) => (
                          <li key={`${file.name}-${i}`} className={s.file}>
                            <Icon name="file" />
                            <span className={s.fileName} title={file.name}>
                              {file.name}
                            </span>
                            <span className={s.fileSize}>{formatSize(file.size)}</span>
                            <button
                              type="button"
                              className={s.fileRemove}
                              onClick={() => removeFile(f.id, i)}
                              aria-label={`Remover ${file.name}`}
                            >
                              <Icon name="x" />
                            </button>
                          </li>
                        ))}
                      </ul>
                    )}

                    <div className={s.actions}>
                      {f.id === 'html' && (
                        <label className={s.variantBtn}>
                          <Icon name="plus" />
                          <span>Adicionar variante</span>
                          <input
                            type="file"
                            accept=".psrt"
                            multiple
                            className={s.fileInput}
                            onChange={(e) => {
                              addFiles(f.id, e.target.files);
                              e.target.value = '';
                            }}
                          />
                        </label>
                      )}
                      <button
                        type="button"
                        className={`${s.confirm} ${f.id !== 'html' ? s.confirmFull : ''}`}
                        onClick={() => {
                          onSave(f.id, variants[f.id]);
                        }}
                      >
                        Confirmar
                      </button>
                    </div>
                  </div>
                )}
              </li>
            );
          })}
        </ul>

        <footer className={s.footer}>
          <button type="button" className={s.cancel} onClick={onCancel}>
            Cancelar
          </button>
        </footer>
      </div>
    </div>
  );
}

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
}

function Icon({ name }: { name: string }) {
  const p: Record<string, ReactNode> = {
    x: <path d="M6 6l12 12M18 6L6 18" />,
    chev: <path d="M6 9l6 6 6-6" />,
    plus: <path d="M12 5v14M5 12h14" />,
    file: (
      <>
        <path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z" />
        <path d="M14 3v6h6" />
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
      {p[name]}
    </svg>
  );
}

