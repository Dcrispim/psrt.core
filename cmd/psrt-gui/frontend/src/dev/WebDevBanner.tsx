import { useEffect, useState } from 'react';
import { devApiHealthy } from './apiClient';

export function WebDevBanner() {
  const [apiOk, setApiOk] = useState<boolean | null>(null);

  useEffect(() => {
    let cancelled = false;
    const check = () => {
      void devApiHealthy().then((ok) => {
        if (!cancelled) setApiOk(ok);
      });
    };
    check();
    const id = window.setInterval(check, 5000);
    return () => {
      cancelled = true;
      window.clearInterval(id);
    };
  }, []);

  return (
    <div
      className={`web-dev-banner${apiOk === false ? ' web-dev-banner--warn' : ''}`}
      role="status"
    >
      <strong>Modo WEB (dev)</strong>
      <span>
        Hot reload no navegador · <code>exemplo.psrt</code> embutido (sem API para abrir).
      </span>
      {apiOk === false ? (
        <span className="web-dev-banner__hint">
          API Go offline — preview/compile precisam de{' '}
          <code>go run ./cmd/psrt-gui-dev</code>
        </span>
      ) : apiOk === true ? (
        <span className="web-dev-banner__ok">API conectada (adapt/compile)</span>
      ) : null}
    </div>
  );
}
