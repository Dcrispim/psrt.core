import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { EditorProvider } from './context/EditorContext';
import { AlertModalProvider } from './context/AlertModalContext';
import { App } from './App';
import { AppErrorBoundary } from './components/AppErrorBoundary';
import './styles/global.css';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AppErrorBoundary>
      <AlertModalProvider>
        <EditorProvider>
          <App />
        </EditorProvider>
      </AlertModalProvider>
    </AppErrorBoundary>
  </StrictMode>,
);
