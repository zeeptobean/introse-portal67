import { readPage } from "#lib/server/forms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, params, depends, locals } = event;
  const documents = locals.api.documents;
  await parent();
  depends("portal:documents");
  return readPage(async () => ({ record: await documents.get(params.request) }));
};
