import type { PsrtDocument } from '../types/document';

export function parseDocumentJson(json: string): PsrtDocument {
  return JSON.parse(json) as PsrtDocument;
}
