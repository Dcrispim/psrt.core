import { useEffect, useId, useMemo, useState, type FormEvent, type MouseEvent } from 'react';
import { buildConstPreview, isValidConstName } from '../../lib/psrtConstPreview';
import s from './assetModal.module.css';

export interface AddConstModalProps {
  open: boolean;
  existingNames: string[];
  onAdd: (name: string, value: string) => void;
  onCancel: () => void;
}

export function AddConstModal({ open, existingNames, onAdd, onCancel }: AddConstModalProps) {
  const titleId = useId();
  const [name, setName] = useState('');
  const [value, setValue] = useState('');
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
      setName('');
      setValue('');
      setError(null);
    }
  }, [open]);

  const preview = useMemo(() => buildConstPreview(name, value), [name, value]);

  if (!open) return null;

  const stopDialogClick = (e: MouseEvent) => e.stopPropagation();

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const n = name.trim();
    if (!isValidConstName(n)) {
      setError('Nome inválido (sem espaços ou |).');
      return;
    }
    if (existingNames.includes(n)) {
      setError(`Constante "${n}" já existe.`);
      return;
    }
    onAdd(n, value);
  };

  return (
    <div
      className={s.overlay}
      role="presentation"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onCancel();
      }}
    >
      <div className={s.dialog} role="dialog" aria-modal="true" aria-labelledby={titleId} onMouseDown={stopDialogClick}>
        <header className={s.header}>
          <h2 id={titleId} className={s.title}>Adicionar constante</h2>
          <button type="button" className={s.close} onClick={onCancel} aria-label="Fechar">×</button>
        </header>
        <form id="add-const-form" className={s.body} onSubmit={handleSubmit}>
          <div className={s.field}>
            <label className={s.label} htmlFor="const-name">Nome</label>
            <input id="const-name" className={s.input} value={name} placeholder="baseURL" autoComplete="off" onChange={(e) => { setName(e.target.value); setError(null); }} />
            <p className={s.hint}>Usado como @{name.trim() || 'nome'}@ no arquivo.</p>
          </div>
          <div className={s.field}>
            <label className={s.label} htmlFor="const-value">Valor</label>
            <textarea id="const-value" className={s.textarea} value={value} placeholder="#1DB954 ou https://cdn.example/617/" onChange={(e) => { setValue(e.target.value); setError(null); }} />
          </div>
          {error ? <p className={s.error}>{error}</p> : null}
          {preview ? (
            <>
              <div className={s.previewBlock}>
                <p className={s.previewTitle}>Declaração ($CONSTS)</p>
                <pre className={s.pre}>{preview.declaration}</pre>
              </div>
              <div className={s.previewBlock}>
                <p className={s.previewTitle}>Uso no arquivo (visualização)</p>
                <pre className={s.pre}>{`Token: ${preview.token}\n\nURL da página:\n${preview.examples.pageImageUrl}\n\nCorpo do texto:\n${preview.examples.textBody}\n\nEstilo JSON:\n${preview.examples.styleProperty}\n\n${preview.examples.styleInlineFragment}`}</pre>
              </div>
            </>
          ) : (
            <p className={s.hint}>Informe o nome para ver a pré-visualização.</p>
          )}
        </form>
        <footer className={s.footer}>
          <button type="button" className={s.cancel} onClick={onCancel}>Cancelar</button>
          <button type="submit" form="add-const-form" className={s.confirm} disabled={!name.trim()}>Adicionar</button>
        </footer>
      </div>
    </div>
  );
}
