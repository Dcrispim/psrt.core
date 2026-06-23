export interface HtmlVariantBody {
  label: string;
  content: string;
}

export interface HtmlVariantPayload {
  paths: string[];
  bodies: HtmlVariantBody[];
}

/** Wails file input may expose `path` on File; otherwise send PSRT text in bodies. */
export async function prepareHtmlVariants(files: File[]): Promise<HtmlVariantPayload> {
  const paths: string[] = [];
  const bodies: HtmlVariantBody[] = [];

  for (const file of files) {
    const path = (file as File & { path?: string }).path?.trim();
    if (path) {
      paths.push(path);
      continue;
    }
    const content = await file.text();
    bodies.push({ label: file.name, content });
  }

  return { paths, bodies };
}
