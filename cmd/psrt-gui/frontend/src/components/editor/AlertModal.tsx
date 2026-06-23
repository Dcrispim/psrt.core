import { useEffect, useId, useRef, type FormEvent, type MouseEvent } from 'react';
import s from './alertModal.module.css';

export type AlertModalMode = 'confirm' | 'prompt';

export interface AlertModalProps {
  open: boolean;
  mode: AlertModalMode;
  title: string;
  message?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  defaultValue?: string;
  placeholder?: string;
  danger?: boolean;
  onConfirm: (value?: string) => void;
  onCancel: () => void;
}

export function AlertModal({
  open,
  mode,
  title,
  message,
  confirmLabel,
  cancelLabel,
  defaultValue = '',
  placeholder,
  danger = false,
  onConfirm,
  onCancel,
}: AlertModalProps) {
  const titleId = useId();
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onCancel();
    };
    document.addEventListener('keydown', onKey);
    return () => document.removeEventListener('keydown', onKey);
  }, [open, onCancel]);

  useEffect(() => {
    if (open && mode === 'prompt') {
      const t = window.setTimeout(() => {
        inputRef.current?.focus();
        inputRef.current?.select();
      }, 0);
      return () => window.clearTimeout(t);
    }
  }, [open, mode]);

  if (!open) return null;

  const submitPrompt = (e: FormEvent) => {
    e.preventDefault();
    const v = inputRef.current?.value.trim() ?? '';
    if (v) onConfirm(v);
  };

  const stopDialogClick = (e: MouseEvent) => e.stopPropagation();

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
        <h2 id={titleId} className={s.title}>
          {title}
        </h2>
        {message ? <p className={s.message}>{message}</p> : null}

        {mode === 'prompt' ? (
          <form className={s.field} onSubmit={submitPrompt}>
            <input
              ref={inputRef}
              className={s.input}
              type="text"
              defaultValue={defaultValue}
              placeholder={placeholder}
              autoComplete="off"
              spellCheck={false}
            />
          </form>
        ) : null}

        <div className={s.actions}>
          <button type="button" className={s.btn} onClick={onCancel}>
            {cancelLabel ?? 'Cancelar'}
          </button>
          {mode === 'prompt' ? (
            <button
              type="button"
              className={`${s.btn} ${danger ? s.btnDanger : s.btnPrimary}`}
              onClick={() => {
                const v = inputRef.current?.value.trim() ?? '';
                if (v) onConfirm(v);
              }}
            >
              {confirmLabel ?? 'OK'}
            </button>
          ) : (
            <button
              type="button"
              className={`${s.btn} ${danger ? s.btnDanger : s.btnPrimary}`}
              onClick={() => onConfirm()}
            >
              {confirmLabel ?? 'Confirmar'}
            </button>
          )}
        </div>
        </div>
    </div>
  );
}
