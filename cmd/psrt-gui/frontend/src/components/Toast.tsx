import { useEditor } from '../context/useEditor';

export function Toast() {
  const { toast } = useEditor();
  return <div className={`toast${toast ? '' : ' hidden'}`}>{toast ?? ''}</div>;
}
