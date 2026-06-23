/**
 * Browser dev replacement for @wails/go/main/GUIApp.
 * Requires the Go dev API: go run ./cmd/psrt-gui-dev
 */
import { styleadapter, visualapp } from '@wails/go/models';
import { devPost } from './apiClient';
import { downloadPsrt, pickPsrtFile } from './fileIO';

let devBaseDir = '';

export function setDevBaseDir(dir: string): void {
  devBaseDir = dir;
}

export async function AdaptEntriesForWeb(
  entriesJSON: string,
  canvasW: number,
  canvasH: number,
  zoom: number,
): Promise<Array<styleadapter.WebPreviewStyle>> {
  const raw = await devPost<unknown[]>('/adapt-entries-for-web', {
    entriesJSON,
    canvasW,
    canvasH,
    zoom,
  });
  return raw.map((r) => styleadapter.WebPreviewStyle.createFrom(r));
}

export async function AdaptTextStyleForWeb(
  styleJSON: string,
  x: number,
  y: number,
  width: number,
  textSize: number,
  canvasW: number,
  canvasH: number,
  zoom: number,
): Promise<styleadapter.WebPreviewStyle> {
  const list = await AdaptEntriesForWeb(
    JSON.stringify([
      { index: 0, style: styleJSON, content: '', x, y, width, textSize },
    ]),
    canvasW,
    canvasH,
    zoom,
  );
  return list[0] ?? styleadapter.WebPreviewStyle.createFrom({});
}

export async function GetAssetDataURI(url: string): Promise<string> {
  const { uri } = await devPost<{ uri: string }>('/get-asset-data-uri', {
    url,
    baseDir: devBaseDir,
  });
  return uri ?? '';
}

export async function FormatDocumentJSON(docJSON: string): Promise<string> {
  const { text } = await devPost<{ text: string }>('/format-document-json', { docJSON });
  return text;
}

export async function FormatPageDocumentJSON(
  docJSON: string,
  pageName: string,
): Promise<string> {
  const { text } = await devPost<{ text: string }>('/format-page-document-json', {
    docJSON,
    pageName,
  });
  return text;
}

export async function MergePageDocumentPSRT(
  fullDocJSON: string,
  pageName: string,
  psrtText: string,
): Promise<string> {
  const { document } = await devPost<{ document: string }>('/merge-page-document-psrt', {
    fullDocJSON,
    pageName,
    psrtText,
  });
  return document;
}

export async function ParseDocumentPSRT(text: string): Promise<string> {
  const { document } = await devPost<{ document: string }>('/parse-document-psrt', { text });
  return document;
}

export async function CompilePageSVGFromDocument(
  docJSON: string,
  page: string,
): Promise<{ uri: string; usedGoTextFallback: boolean }> {
  return devPost<{ uri: string; usedGoTextFallback: boolean }>(
    '/compile-page-svg-from-document',
    {
      docJSON,
      pageName: page,
    },
  );
}

export async function CompilePageHTMLFromDocument(
  docJSON: string,
  page: string,
): Promise<string> {
  const { uri } = await devPost<{ uri: string }>('/compile-page-html-from-document', {
    docJSON,
    pageName: page,
  });
  return uri;
}

export async function OpenImageFileDialog(): Promise<string> {
  return new Promise((resolve) => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/png,image/jpeg,image/gif,image/webp,image/avif,image/svg+xml';
    input.onchange = () => {
      const file = input.files?.[0];
      if (!file) {
        resolve('');
        return;
      }
      const withPath = file as File & { path?: string };
      resolve(withPath.path?.trim() ?? '');
    };
    input.click();
  });
}

export async function OpenFileDialog(): Promise<visualapp.OpenFileResult> {
  const picked = await pickPsrtFile();
  if (!picked) {
    return visualapp.OpenFileResult.createFrom({});
  }
  if (picked.filePath.includes('/') || picked.filePath.includes('\\')) {
    devBaseDir = picked.filePath.replace(/[/\\][^/\\]+$/, '');
  }
  const document = await ParseDocumentPSRT(picked.text);
  return visualapp.OpenFileResult.createFrom({
    filePath: picked.filePath,
    document,
  });
}

export async function SaveDocumentJSON(docJSON: string): Promise<void> {
  const psrt = await FormatDocumentJSON(docJSON);
  downloadPsrt('document.psrt', psrt);
}

export async function SaveAsDocumentJSON(docJSON: string): Promise<string> {
  const psrt = await FormatDocumentJSON(docJSON);
  const name = 'document.psrt';
  downloadPsrt(name, psrt);
  return name;
}

export async function RefreshAssetURL(_url: string): Promise<void> {}

export async function RefreshPageImage(): Promise<void> {}

export async function SetAutoCompile(_on: boolean): Promise<void> {}

export async function BeginEdit(): Promise<void> {}

export async function EndEdit(): Promise<void> {}

export async function Save(): Promise<void> {}

export async function GetDocumentJSON(): Promise<string> {
  return '{}';
}

export async function GetDocumentPSRT(): Promise<string> {
  return '';
}

export async function GetState(): Promise<visualapp.UIState> {
  return visualapp.UIState.createFrom({});
}

export async function SetActivePage(_name: string): Promise<void> {}

export async function SelectText(_index: number): Promise<void> {}

export async function PatchText(
  _pageName: string,
  _index: number,
  _patch: visualapp.TextPatch,
): Promise<void> {}

export async function PatchPage(_patch: visualapp.PagePatch): Promise<void> {}

export async function AddPage(
  _name: string,
  _imageURL: string,
  _styleJSON: string,
): Promise<void> {}

export async function RemovePage(_name: string): Promise<void> {}

export async function MovePage(
  _name: string,
  _ref: string,
  _before: boolean,
): Promise<void> {}

export async function AddTextBlock(
  _index: number,
  _x: number,
  _y: number,
  _width: number,
  _textSize: number,
  _content: string,
  _styleJSON: string,
  _imageRef: string,
): Promise<void> {}

export async function RemoveText(_index: number): Promise<void> {}

export async function ReorderText(
  _index: number,
  _ref: number,
  _before: boolean,
): Promise<void> {}

export async function AddFont(_url: string): Promise<void> {}

export async function RemoveFont(_url: string): Promise<void> {}

export async function AddConst(_name: string, _value: string): Promise<void> {}

export async function RemoveConst(_name: string): Promise<void> {}

export async function SetDocumentFromPSRT(_text: string): Promise<void> {}

export async function Undo(): Promise<void> {}

export async function Redo(): Promise<void> {}

export async function CompileDocumentHTML(): Promise<string> {
  return '';
}

export async function CompilePageHTML(_page: string): Promise<string> {
  return '';
}

export async function CompilePageSVG(_page: string): Promise<{
  uri: string;
  usedGoTextFallback: boolean;
}> {
  return { uri: '', usedGoTextFallback: false };
}

export async function ExportSVG(_dir: string): Promise<{
  uri: string;
  usedGoTextFallback: boolean;
}> {
  return { uri: '', usedGoTextFallback: false };
}

export async function ExportSVGFromDocument(_docJSON: string): Promise<{
  uri: string;
  usedGoTextFallback: boolean;
}> {
  throw new Error('Exportar SVG disponível apenas no aplicativo desktop');
}

export async function ExportHTMLFromDocument(
  _docJSON: string,
  _variantPaths: string[],
  _variantBodies: import('@wails/go/models').visualapp.VariantPSRT[],
): Promise<string> {
  throw new Error('Exportar HTML disponível apenas no aplicativo desktop');
}
