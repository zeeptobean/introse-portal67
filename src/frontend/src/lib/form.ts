import { superForm, type FormOptions, type SuperValidated } from "sveltekit-superforms";
export type FormMessage = App.Superforms.Message;
export const fieldId = (path: string) => path.replace(/[.[\]]+/g, "-").replace(/-$/, "");
export function actionUrl(url: URL | { href: string }, action: string) {
  const target = new URL(url.href);
  // Copy the keys first: deleting while iterating the live params skips entries.
  for (const key of Array.from(target.searchParams.keys()))
    if (key.startsWith("/")) target.searchParams.delete(key);
  target.searchParams.set(`/${action}`, "");
  target.hash = "";
  return target.pathname + target.search;
}
// Superforms owns validation, fields and submission state. These options encode
// Portal67's no-replay contract, without hiding domain workflows in a form engine.
export function portalForm<T extends Record<string, unknown>>(
  initial: SuperValidated<T>,
  options: FormOptions<NoInfer<T>, FormMessage> = {},
) {
  const controller = superForm<T, FormMessage>(initial, {
    resetForm: false,
    taintedMessage: false,
    multipleSubmits: "prevent",
    ...options,
    onError: ({ result }) => {
      controller.message.set({
        text:
          result.error.message ??
          "We couldn't confirm this action. Check current records before trying again.",
        tone: "error",
        uncertain: true,
      });
    },
  });
  return controller;
}
