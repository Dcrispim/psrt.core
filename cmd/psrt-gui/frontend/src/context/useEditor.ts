import { useContext } from 'react';
import { EditorContext, type EditorContextValue } from './EditorContext';

export function useEditor(): EditorContextValue {
  const ctx = useContext(EditorContext);
  if (!ctx) {
    throw new Error('useEditor must be used within EditorProvider');
  }
  return ctx;
}
