import { Header } from './components/editor/Header';
import { PageThumbnailsSidebar } from './components/PageThumbnailsSidebar';
import { Canvas } from './components/Canvas';
import { PropertiesPanel } from './components/PropertiesPanel';
import { Toast } from './components/Toast';
import { WebDevBanner } from './dev/WebDevBanner';

const isWebDev = import.meta.env.VITE_WEB_DEV === 'true';

export function App() {
  return (
    <>
      {isWebDev ? <WebDevBanner /> : null}
      <Header />
      <main className="workspace">
        <PageThumbnailsSidebar />
        <Canvas />
        <PropertiesPanel />
      </main>
      <Toast />
    </>
  );
}
