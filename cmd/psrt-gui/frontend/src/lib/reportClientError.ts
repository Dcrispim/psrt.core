export function reportClientError(
  operation: string,
  filePath: string,
  error: unknown,
): void {
  const err = error instanceof Error ? error : new Error(String(error));
  const message = `${err.name}: ${err.message}\n${err.stack ?? ''}`;
  const wails = (
    window as unknown as {
      go?: { main?: { GUIApp?: { ReportClientError?: (op: string, fp: string, msg: string) => void } } };
    }
  ).go?.main?.GUIApp;
  if (wails?.ReportClientError) {
    try {
      wails.ReportClientError(operation, filePath, message);
    } catch {
      /* dev web / bindings missing */
    }
  }
}
