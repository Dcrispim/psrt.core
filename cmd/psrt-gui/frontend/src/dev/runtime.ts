/** Browser dev shim for @wails/runtime (no-op events). */

export function EventsOnMultiple(
  _eventName: string,
  _callback: (...data: unknown[]) => void,
  _maxCallbacks?: number,
): () => void {
  return () => {};
}

export function EventsOn(_eventName: string, callback: (...data: unknown[]) => void): () => void {
  return EventsOnMultiple(_eventName, callback, -1);
}

export function EventsOff(_eventName: string, ..._additionalEventNames: string[]): void {}

export function EventsEmit(_eventName: string, ..._data: unknown[]): void {}
