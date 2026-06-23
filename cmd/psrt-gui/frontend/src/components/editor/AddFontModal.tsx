import { useEffect, useId, useMemo, useState, type FormEvent, type MouseEvent } from 'react';
import {
  isGoogleFontsCssUrl,
  isLikelyFontUrl,
  normalizeFontUrlInput,
  parseGoogleFontsFamilies,
  primaryGoogleFontFamily,
} from '../../lib/googleFontsUrl';
import s from './assetModal.module.css';

export interface AddFontModalProps {
  open: boolean;
  existingUrls: string[];
  onAdd: (url: string) => void;
  onCancel: () => void;
}

function buildFontPreviewLines(url: string, families: string[]): string {
  const trimmed = url.trim();
  if (!trimmed) return '';

  const lines = ['$FONTS', trimmed, '$ENDFONTS', ''];

  if (families.length > 0) {
    lines.push('font-family no JSON de estilo:');
    for (const family of families) {
      lines.push(`  "font-family": "${family}"`);
    }
  } else if (isGoogleFontsCssUrl(trimmed)) {
    lines.push('Não foi possível extrair o nome da família desta URL.');
  } else {
    lines.push('Use o nome da família tipográfica no campo font-family do JSON.');
  }

  return lines.join('\n');
}

export function AddFontModal({ open, existingUrls, onAdd, onCancel }: AddFontModalProps) {
  const titleId = useId();
  const [url, setUrl] = useState('');
  const [error, setError] = useState<string | null>(null);

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
      setUrl('');
      setError(null);
    }
  }, [open]);

  const normalized = normalizeFontUrlInput(url);
  const googleFamilies = useMemo(
    () => (isGoogleFontsCssUrl(normalized) ? parseGoogleFontsFamilies(normalized) : []),
    [normalized],
  );
  const previewText = useMemo(
    () => buildFontPreviewLines(normalized, googleFamilies),
    [normalized, googleFamilies],
  );
  const primaryFamily = primaryGoogleFontFamily(normalized);

  if (!open) return null;

  const stopDialogClick = (e: MouseEvent) => e.stopPropagation();

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const u = normalized;
    if (!isLikelyFontUrl(u)) {
      setError('Informe uma URL válida (http ou https).');
      return;
    }
    if (existingUrls.includes(u)) {
      setError('Esta fonte já está no documento.');
      return;
    }
    onAdd(u);
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
            Adicionar fonte
          </h2>
          <button type="button" className={s.close} onClick={onCancel} aria-label="Fechar">
            ×
          </button>
        </header>

        <form id="add-font-form" className={s.body} onSubmit={handleSubmit}>
          <div className={s.field}>
            <label className={s.label} htmlFor="font-url">
              URL da fonte
            </label>
            <input
              id="font-url"
              className={s.input}
              type="url"
              value={url}
              placeholder="https://fonts.googleapis.com/css2?family=Roboto&display=swap"
              autoComplete="off"
              onChange={(e) => {
                setUrl(e.target.value);
                setError(null);
              }}
            />
            <p className={s.hint}>
              Cole o link CSS do Google Fonts ou URL de arquivo (.woff2, CDN, etc.).
            </p>
          </div>

          {googleFamilies.length > 0 ? (
            <div className={s.field}>
              <span className={s.label}>Família detectada</span>
              <ul className={s.tagList}>
                {googleFamilies.map((family) => (
                  <li key={family}>
                    <span className={s.tag}>{family}</span>
                  </li>
                ))}
              </ul>
              {primaryFamily ? (
                <p className={s.hint}>
                  Use <code>{primaryFamily}</code> em <code>font-family</code> no estilo do texto.
                </p>
              ) : null}
            </div>
          ) : null}

          {error ? <p className={s.error}>{error}</p> : null}

          {previewText ? (
            <div className={s.previewBlock}>
              <p className={s.previewTitle}>Pré-visualização no arquivo</p>
              <pre className={s.pre}>{previewText}</pre>
            </div>
          ) : (
            <p className={s.hint}>Informe a URL para ver como ficará em $FONTS.</p>
          )}
        </form>

        <footer className={s.footer}>
          <button type="button" className={s.cancel} onClick={onCancel}>
            Cancelar
          </button>
          <button
            type="submit"
            form="add-font-form"
            className={s.confirm}
            disabled={!normalized}
          >
            Adicionar
          </button>
        </footer>
      </div>
    </div>
  );
}
