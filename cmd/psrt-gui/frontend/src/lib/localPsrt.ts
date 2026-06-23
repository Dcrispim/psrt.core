const PSRT_KEY = 'psrt-gui:last-psrt';
const PATH_KEY = 'psrt-gui:last-psrt-path';
const UPDATED_KEY = 'psrt-gui:last-psrt-updated';

export interface StoredPsrt {
  filePath: string;
  content: string;
  updatedAt: number;
}

export function saveLastPsrt(filePath: string, content: string, _documentJson?: string): void {
  if (!filePath || !content) return;
  try {
    localStorage.setItem(PSRT_KEY, content);
    localStorage.setItem(PATH_KEY, filePath);
    localStorage.setItem(UPDATED_KEY, String(Date.now()));
  } catch {
    /* quota or private mode */
  }
}

export function loadLastPsrt(): StoredPsrt | null {
  try {
    const content = localStorage.getItem(PSRT_KEY);
    const filePath = localStorage.getItem(PATH_KEY);
    const updatedAt = Number(localStorage.getItem(UPDATED_KEY) ?? 0);
    if (!content || !filePath) return null;
    return { filePath, content, updatedAt };
  } catch {
    return null;
  }
}
