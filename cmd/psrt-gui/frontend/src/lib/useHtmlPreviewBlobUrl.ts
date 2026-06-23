import { useEffect, useState } from 'react';
import { decodeDataUriBytes } from './decodeDataUri';

/**
 * Turns a compiled HTML data URI into a blob: URL for iframe preview.
 * Large documents (embedded JPEG base64) exceed WebView2 limits on data: navigation.
 */
export function useHtmlPreviewBlobUrl(dataUri: string | null): string | null {
  const [blobUrl, setBlobUrl] = useState<string | null>(null);

  useEffect(() => {
    if (!dataUri) {
      setBlobUrl(null);
      return;
    }

    const bytes = decodeDataUriBytes(dataUri);
    if (bytes == null) {
      setBlobUrl(null);
      return;
    }

    const blob = new Blob([bytes], { type: 'text/html;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    setBlobUrl(url);

    return () => {
      URL.revokeObjectURL(url);
    };
  }, [dataUri]);

  return blobUrl;
}
