import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from 'react';
import { AlertModal, type AlertModalMode } from '../components/editor/AlertModal';

export interface ConfirmOptions {
  title: string;
  message?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  danger?: boolean;
}

export interface PromptOptions {
  title: string;
  message?: string;
  defaultValue?: string;
  placeholder?: string;
  confirmLabel?: string;
  cancelLabel?: string;
}

interface ModalState {
  mode: AlertModalMode;
  title: string;
  message?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  defaultValue?: string;
  placeholder?: string;
  danger?: boolean;
}

interface AlertModalContextValue {
  confirm: (options: ConfirmOptions) => Promise<boolean>;
  prompt: (options: PromptOptions) => Promise<string | null>;
}

const AlertModalContext = createContext<AlertModalContextValue | null>(null);

export function AlertModalProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<ModalState | null>(null);
  const resolverRef = useRef<((value: boolean | string | null) => void) | null>(null);

  const close = useCallback(() => {
    setState(null);
    resolverRef.current = null;
  }, []);

  const confirm = useCallback((options: ConfirmOptions) => {
    return new Promise<boolean>((resolve) => {
      resolverRef.current = (v) => resolve(v === true);
      setState({ mode: 'confirm', ...options });
    });
  }, []);

  const prompt = useCallback((options: PromptOptions) => {
    return new Promise<string | null>((resolve) => {
      resolverRef.current = (v) => resolve(typeof v === 'string' ? v : null);
      setState({ mode: 'prompt', ...options });
    });
  }, []);

  const handleConfirm = useCallback(
    (value?: string) => {
      const resolve = resolverRef.current;
      close();
      if (!resolve) return;
      if (state?.mode === 'prompt') {
        resolve(value ?? null);
      } else {
        resolve(true);
      }
    },
    [close, state?.mode],
  );

  const handleCancel = useCallback(() => {
    const resolve = resolverRef.current;
    close();
    if (!resolve) return;
    if (state?.mode === 'prompt') {
      resolve(null);
    } else {
      resolve(false);
    }
  }, [close, state?.mode]);

  const value = useMemo(() => ({ confirm, prompt }), [confirm, prompt]);

  return (
    <AlertModalContext.Provider value={value}>
      {children}
      {state ? (
        <AlertModal
          open
          mode={state.mode}
          title={state.title}
          message={state.message}
          confirmLabel={state.confirmLabel}
          cancelLabel={state.cancelLabel}
          defaultValue={state.defaultValue}
          placeholder={state.placeholder}
          danger={state.danger}
          onConfirm={handleConfirm}
          onCancel={handleCancel}
        />
      ) : null}
    </AlertModalContext.Provider>
  );
}

export function useAlertModal(): AlertModalContextValue {
  const ctx = useContext(AlertModalContext);
  if (!ctx) {
    throw new Error('useAlertModal must be used within AlertModalProvider');
  }
  return ctx;
}
