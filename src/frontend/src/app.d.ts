// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    interface Locals {
      api: import("#lib/server/api.ts").Api;
      userPromise?: Promise<import("#lib/schemas/shared.ts").User | null>;
    }
    namespace Superforms {
      type Message = {
        text: string;
        tone?: "error" | "success" | "info" | "caution";
        uncertain?: boolean;
      };
    }
    // interface PageData {}
    // interface PageState {}
    // interface Platform {}
  }
}

export {};
