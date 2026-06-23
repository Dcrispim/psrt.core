import { useCallback, useEffect } from 'react';

import { useEditor } from '../context/useEditor';
import { useWebPreviewPan } from '../lib/useWebPreviewPan';

import { ViewTabs } from './ViewTabs';

import { WebPreview } from './WebPreview';

import { WebPreviewToolbar } from './WebPreviewToolbar';

import {

  StaticCompiledPreviewStage,

  StaticCompiledPreviewToolbar,

  useStaticPreviewDisplay,

} from './StaticCompiledPreview';



export function Canvas() {

  const {

    previewTab,

    setPreviewTab,

    getPagePreview,

    compilePreviewSvg,

    compilePreviewHtml,

    state,

    showToast,

    toggleWebTextsVisible,

    webZoom,

    setWebZoom,

  } = useEditor();

  const page = state?.activePage ?? '';

  const previewSvg = page ? getPagePreview(page, 'svg') : null;

  const previewHtml = page ? getPagePreview(page, 'html') : null;



  const svgDisplay = useStaticPreviewDisplay(previewSvg);

  const htmlDisplay = useStaticPreviewDisplay(previewHtml);



  useEffect(() => {

    if (!page) return;

    if (previewTab === 'svg' && !getPagePreview(page, 'svg')) {

      compilePreviewSvg({ notifyGoText: true }).catch((e) => showToast(String(e)));

    } else if (previewTab === 'html' && !getPagePreview(page, 'html')) {

      compilePreviewHtml().catch((e) => showToast(String(e)));

    }

  }, [previewTab, page, getPagePreview, compilePreviewSvg, compilePreviewHtml, showToast]);

  const onWebPreviewKey = useCallback(
    (e: KeyboardEvent) => {
      if (previewTab !== 'web') return;
      if (!(e.ctrlKey || e.metaKey) || e.altKey) return;
      if (e.key.toLowerCase() !== 'h') return;

      const target = e.target;
      if (
        target instanceof HTMLElement &&
        (target.tagName === 'INPUT' ||
          target.tagName === 'TEXTAREA' ||
          target.isContentEditable)
      ) {
        return;
      }

      e.preventDefault();
      toggleWebTextsVisible();
    },
    [previewTab, toggleWebTextsVisible],
  );

  useEffect(() => {
    globalThis.document.addEventListener('keydown', onWebPreviewKey);
    return () => globalThis.document.removeEventListener('keydown', onWebPreviewKey);
  }, [onWebPreviewKey]);

  const switchToWeb = () => {

    if (previewTab !== 'web') setPreviewTab('web');

  };



  const staticDisplay = previewTab === 'svg' ? svgDisplay : previewTab === 'html' ? htmlDisplay : null;

  const staticUri = previewTab === 'svg' ? previewSvg : previewTab === 'html' ? previewHtml : null;

  const webPan = useWebPreviewPan(previewTab === 'web', webZoom, setWebZoom);
  const { resetPan } = webPan;

  useEffect(() => {
    if (previewTab === 'web') resetPan();
  }, [previewTab, page, resetPan]);



  return (

    <section className="canvas-wrap">

      <ViewTabs />

      {previewTab === 'web' && <WebPreviewToolbar />}

      {(previewTab === 'svg' || previewTab === 'html') && staticDisplay && (

        <StaticCompiledPreviewToolbar

          kind={previewTab}

          zoom={staticDisplay.zoom}

          onZoomChange={staticDisplay.setZoom}

          showSource={staticDisplay.showSource}

          onShowSourceChange={staticDisplay.setShowSource}

          disabled={!staticUri}

          minZoom={previewTab === 'svg' || previewTab === 'html' ? 0.25 : undefined}

          maxZoom={previewTab === 'svg' || previewTab === 'html' ? 8 : undefined}

        />

      )}

      <div
        ref={webPan.viewportRef}
        className={`canvas-viewport${previewTab === 'web' ? ' canvas-viewport--web-pan' : ''}${webPan.dragging ? ' canvas-viewport--web-pan-dragging' : ''}`}
        {...(previewTab === 'web' ? webPan.panHandlers : {})}
      >

        {previewTab === 'web' && (

          <div className="canvas-stage canvas-stage-web" style={webPan.stageStyle}>

            <WebPreview variant="canvas" />

          </div>

        )}



        {previewTab === 'svg' && (

          <StaticCompiledPreviewStage

            uri={previewSvg}

            kind="svg"

            loadingMessage="Compilando SVG…"

            onDoubleClick={switchToWeb}

            zoom={svgDisplay.zoom}

            onZoomChange={svgDisplay.setZoom}

            showSource={svgDisplay.showSource}

            source={svgDisplay.source}

            minZoom={0.25}

            maxZoom={8}

          />

        )}



        {previewTab === 'html' && (

          <StaticCompiledPreviewStage

            uri={previewHtml}

            kind="html"

            loadingMessage="Compilando HTML…"

            onDoubleClick={switchToWeb}

            zoom={htmlDisplay.zoom}

            onZoomChange={htmlDisplay.setZoom}

            showSource={htmlDisplay.showSource}

            source={htmlDisplay.source}

            minZoom={0.25}

            maxZoom={8}

          />

        )}

      </div>

      {previewTab !== 'web' && !staticDisplay?.showSource && (

        <p className="canvas-hint">Duplo clique na visualização para editar no modo WEB</p>

      )}

    </section>

  );

}

