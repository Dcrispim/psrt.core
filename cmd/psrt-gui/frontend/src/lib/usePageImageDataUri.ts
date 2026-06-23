import { useEffect, useState } from 'react';
import * as api from '@wails/go/main/GUIApp';
import { resolveAssetReference } from './expandConsts';

function isDataUri(value: string | null | undefined): value is string {
  return typeof value === 'string' && value.startsWith('data:');
}

function isHttpUrl(value: string | null | undefined): value is string {
  return typeof value === 'string' && /^https?:\/\//i.test(value);
}

/** URL usable as <img src> (data URI from Wails cache or direct http(s) in web dev). */
function isImageSrc(value: string | null | undefined): value is string {
  return isDataUri(value) || isHttpUrl(value);
}

/**
 * Resolves page background for <img src>.
 * - Wails: prefers data URI from GetAssetDataURI
 * - Web dev: uses https URL directly when API is offline
 */
function resolveInitialSrc(
  imageUrl: string | undefined,
  cachedUri: string | null,
): string | null {
  if (!imageUrl) return null;
  const isWebMode =
    import.meta.env.VITE_WEB_DEV === 'true' ||
    import.meta.env.VITE_PRT_WEB === 'true';
  if (isImageSrc(cachedUri)) return cachedUri;
  if (isWebMode && isHttpUrl(imageUrl)) return imageUrl;
  return null;
}

export function usePageImageDataUri(
  imageUrl: string | undefined,
  cachedUri: string | null,
  consts?: Record<string, string>,
): string | null {
  const resolvedUrl = imageUrl
    ? resolveAssetReference(imageUrl, consts)
    : undefined;

  const [src, setSrc] = useState<string | null>(() =>
    resolveInitialSrc(resolvedUrl, cachedUri),
  );

  useEffect(() => {
    if (!resolvedUrl) {
      setSrc(null);
      return;
    }

    let cancelled = false;

    const apply = (value: string | null) => {
      if (!cancelled) setSrc(value);
    };

    const isWebMode =
      import.meta.env.VITE_WEB_DEV === 'true' ||
      import.meta.env.VITE_PRT_WEB === 'true';

    if (isImageSrc(cachedUri)) {
      apply(cachedUri);
    } else if (isWebMode && isHttpUrl(resolvedUrl)) {
      apply(resolvedUrl);
    } else {
      apply(null);
    }

    if (isDataUri(resolvedUrl)) {
      return () => {
        cancelled = true;
      };
    }

    if (isWebMode && isHttpUrl(resolvedUrl)) {
      return () => {
        cancelled = true;
      };
    }

    void (async () => {
      try {
        const uri = await api.GetAssetDataURI(resolvedUrl);
        if (!cancelled && uri) {
          apply(uri);
        }
      } catch {
        if (!cancelled && isHttpUrl(resolvedUrl)) {
          apply(resolvedUrl);
        }
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [resolvedUrl, cachedUri]);

  return src;
}
