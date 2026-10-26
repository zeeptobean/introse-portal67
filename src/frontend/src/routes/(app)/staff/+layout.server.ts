import { requireUser } from "#lib/server/auth.ts";
import type { LayoutServerLoad } from "./$types";
export const load: LayoutServerLoad = async (event) => {
  await requireUser(event, "staff");
  return {};
};
