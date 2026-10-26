import { beforeNavigate } from "$app/navigation";
import { createContext, onDestroy } from "svelte";
const [getDraft, setDraft] = createContext<{ dirty: boolean }>();
export function provideDraft() {
  const state = $state({ dirty: false });
  setDraft(state);
  return state;
}
export const useDraft = () => getDraft();
export function protectDraft(isDirty: () => boolean) {
  const draft = getDraft();
  $effect(() => {
    draft.dirty = isDirty();
  });
  beforeNavigate(({ from, to, cancel }) => {
    if (from && to && from.url.pathname === to.url.pathname && from.url.search === to.url.search)
      return;
    if (draft.dirty && !window.confirm("Leave this page and discard your unsaved changes?"))
      cancel();
  });
  $effect(() => {
    const warn = (event: BeforeUnloadEvent) => {
      if (draft.dirty) event.preventDefault();
    };
    window.addEventListener("beforeunload", warn);
    return () => window.removeEventListener("beforeunload", warn);
  });
  onDestroy(() => {
    draft.dirty = false;
  });
}
