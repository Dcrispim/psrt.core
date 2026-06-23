export const GO_TEXT_SVG_NOTICE =
  'Compilação SVG com motor go-text (Chromium não encontrado). Pode haver pequenas diferenças em relação ao preview.';

export interface CompileSVGResult {
  uri: string;
  usedGoTextFallback: boolean;
}

export function notifySvgGoTextFallback(
  showToast: (message: string) => void,
  usedGoTextFallback: boolean,
): void {
  if (usedGoTextFallback) {
    showToast(GO_TEXT_SVG_NOTICE);
  }
}
