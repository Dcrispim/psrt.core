import { Component, type ErrorInfo, type ReactNode } from 'react';

interface Props {
  children: ReactNode;
  filePath?: string;
  onError?: (message: string) => void;
}

interface State {
  error: Error | null;
}

function reportToBackend(operation: string, filePath: string, error: Error) {
  const message = `${error.name}: ${error.message}\n${error.stack ?? ''}`;
  const wails = (
    window as unknown as {
      go?: { main?: { GUIApp?: { ReportClientError?: (op: string, fp: string, msg: string) => void } } };
    }
  ).go?.main?.GUIApp;
  if (wails?.ReportClientError) {
    try {
      wails.ReportClientError(operation, filePath, message);
    } catch {
      /* backend unavailable (dev web) */
    }
  }
}

export class AppErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    const detail = info.componentStack
      ? `${error.message}\n${info.componentStack}`
      : error.message;
    reportToBackend('react', this.props.filePath ?? '', new Error(detail));
    this.props.onError?.(error.message);
  }

  render() {
    if (this.state.error) {
      return (
        <div className="app-error-boundary" role="alert">
          <p>Algo deu errado na interface.</p>
          <p className="app-error-boundary__detail">{this.state.error.message}</p>
          <p className="app-error-boundary__hint">
            Detalhes foram gravados em <strong>psrt-gui.log</strong> (pasta do executável).
          </p>
        </div>
      );
    }
    return this.props.children;
  }
}
