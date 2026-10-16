import { remoteIntrospect } from "@passkey/session-guard";
import { createHandler } from "./routes.ts";
import { VaultStore } from "./vault-store.ts";

const port = Number(process.env.VAULT_API_PORT ?? 4002);
const passkeyServer = `http://127.0.0.1:${process.env.PASSKEY_PORT ?? 8090}`;
const store = new VaultStore(process.env.VAULT_DATA ?? "data/vault.json");

Bun.serve({ hostname: "127.0.0.1", port, fetch: createHandler(remoteIntrospect(passkeyServer), store) });
console.log(`vault-api listening on http://127.0.0.1:${port}`);
