import { defineEnvVars } from "@sveltejs/kit/env";
import { z } from "zod";
export const variables = defineEnvVars({
  BACKEND_URL: {
    schema: z
      .url()
      .default("http://127.0.0.1:8080")
      .refine((value) => {
        const url = new URL(value);
        return (
          ["http:", "https:"].includes(url.protocol) &&
          !url.username &&
          !url.password &&
          url.pathname === "/" &&
          !url.search &&
          !url.hash
        );
      }, "BACKEND_URL must be an HTTP(S) origin without credentials, path or query."),
    description: "Private Go API origin; never sent to the browser.",
  },
});
