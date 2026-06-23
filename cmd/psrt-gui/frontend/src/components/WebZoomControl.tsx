import { useEditor } from '../context/useEditor';
import { ZoomControl } from './ZoomControl';

export function WebZoomControl() {
  const { webZoom, setWebZoom } = useEditor();

  return (
    <ZoomControl id="web-zoom-slider" zoom={webZoom} onZoomChange={setWebZoom} />
  );
}
