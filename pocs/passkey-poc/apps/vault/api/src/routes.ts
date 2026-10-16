import { json, readJSON, requireSession, unauthorized, type Introspect } from "@passkey/session-guard";
import { validId, validSealed, VaultStore } from "./vault-store.ts";

const APP = "vault";

export function createHandler(introspect: Introspect, store = new VaultStore()) {
  return async (req: Request): Promise<Response> => {
    const url = new URL(req.url);
    if (req.method === "GET" && url.pathname === "/api/health") return json({ status: "ok" });

    const session = await requireSession(req, APP, introspect);
    if (!session) return unauthorized();
    const user = session.username;

    if (req.method === "GET" && url.pathname === "/api/vault") return json(store.get(user));

    if (req.method === "PUT" && url.pathname === "/api/vault/canary") {
      const body = await readJSON(req);
      if (!validSealed(body)) return json({ error: "canary must be { iv, ct } in base64url" }, 400);
      return store.setCanary(user, body) ? json({ ok: true }, 201) : json({ error: "canary already set" }, 409);
    }

    const match = url.pathname.match(/^\/api\/vault\/items\/([^/]+)$/);
    if (match) {
      const id = match[1]!;
      if (!validId(id)) return json({ error: "item id must be a uuid" }, 400);
      if (req.method === "PUT") {
        const body = await readJSON(req);
        if (!validSealed(body)) return json({ error: "item must be { iv, ct } in base64url" }, 400);
        try {
          return json(store.put(user, id, body));
        } catch (err) {
          return json({ error: (err as Error).message }, 400);
        }
      }
      if (req.method === "DELETE") return store.remove(user, id) ? new Response(null, { status: 204 }) : json({ error: "not found" }, 404);
    }
    return json({ error: "not found" }, 404);
  };
}
