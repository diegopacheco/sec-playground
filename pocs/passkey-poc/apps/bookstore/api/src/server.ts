import { remoteIntrospect } from "@passkey/session-guard";
import { createHandler } from "./routes.ts";

const port = Number(process.env.BOOKSTORE_API_PORT ?? 4001);
const passkeyServer = `http://127.0.0.1:${process.env.PASSKEY_PORT ?? 8090}`;

Bun.serve({ hostname: "127.0.0.1", port, fetch: createHandler(remoteIntrospect(passkeyServer)) });
console.log(`bookstore-api listening on http://127.0.0.1:${port}`);
