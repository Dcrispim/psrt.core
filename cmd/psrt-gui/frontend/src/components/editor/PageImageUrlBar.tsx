import { useCallback, useEffect, useState } from 'react';
import * as api from '@wails/go/main/GUIApp';
import { buildPageImageRefFromLocalPath } from '../../lib/localAssetRef';
import s from './sidebar.module.css';

interface PageImageUrlBarProps {
  imageUrl: string;
  consts?: Record<string, string>;
  onChange: (url: string) => void;
  onRefresh?: () => void;
  onError?: (message: string) => void;
}

export function PageImageUrlBar({
  imageUrl,
  consts,
  onChange,
  onRefresh,
  onError,
}: PageImageUrlBarProps) {
  const [value, setValue] = useState(imageUrl);
  const [picking, setPicking] = useState(false);

  useEffect(() => {
    setValue(imageUrl);
  }, [imageUrl]);

  const onPickLocalImage = useCallback(async () => {
    setPicking(true);
    try {
      const fsPath = await api.OpenImageFileDialog();
      if (!fsPath) return;
      const ref =
        import.meta.env.VITE_PRT_WEB === 'true'
          ? fsPath
          : buildPageImageRefFromLocalPath(fsPath, consts);
      setValue(ref);
      onChange(ref);
      onRefresh?.();
    } catch (e) {
      onError?.(String(e));
    } finally {
      setPicking(false);
    }
  }, [consts, onChange, onRefresh, onError]);

  return (
    <div className={s.pageImageBar}>
      <label className={s.label} htmlFor="page-image-url-editor">
        Imagem da página
      </label>
      <div className={s.pageImageRow}>
        <input
          id="page-image-url-editor"
          type="text"
          className={s.input}
          value={value}
          spellCheck={false}
          placeholder="URL da imagem…"
          onChange={(e) => {
            const next = e.target.value;
            setValue(next);
            onChange(next);
          }}
        />
        <button
          type="button"
          className={s.pageImageActionBtn}
          title="Usar imagem local"
          disabled={picking}
          onClick={() => {
            onPickLocalImage().catch((e) => onError?.(String(e)));
          }}
        >
          📁
        </button>
        {onRefresh ? (
          <button
            type="button"
            className={s.pageImageActionBtn}
            title="Atualizar imagem do cache"
            onClick={onRefresh}
          >
            ↻
          </button>
        ) : null}
      </div>
    </div>
  );
}
