import { exemploDocument } from './exemploDocument';
import type { PsrtDocument } from '../types/document';

const DEV_FILE_PATH = 'exemplo.psrt';

/** Fallback image when PSRT references a local file:// URL (browser cannot load it). */
const WEB_DEV_CAPA_IMAGE = 'https://picsum.photos/seed/psrt-capa/1080/1920';

function documentForBrowser(doc: PsrtDocument): PsrtDocument {
  const next = structuredClone(doc);
  for (const page of next.pages) {
    const url = page.imageUrl ?? '';
    if (url.startsWith('file:') || url.startsWith('file://')) {
      page.imageUrl = WEB_DEV_CAPA_IMAGE;
    }
  }
  return next;
}

/** Loads exemplo.psrt embedded in the bundle — no Go dev API required. */
export function loadDevFixture(): {
  filePath: string;
  document: PsrtDocument;
} {
  const document = documentForBrowser(structuredClone(exemploDocument));
  return {
    filePath: DEV_FILE_PATH,
    document,
  };
}
