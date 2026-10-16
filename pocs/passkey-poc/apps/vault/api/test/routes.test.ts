import { expect, test } from "bun:test";
import { mkdtempSync, readFileSync, statSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import type { Session } from "@passkey/session-guard";
import { createHandler } from "../src/routes.ts";
import { VaultStore } from "../src/vault-store.ts";

const sessions: Record<string, Session> = {
  alice: { app: "vault", username: "alice", userId: "a", expiresAt: "" },
  mallory: { app: "vault", username: "mallory", userId: "m", expiresAt: "" },
  shop: { app: "bookstore", username: "alice", userId: "a", expiresAt: "" },
};
const introspect = async (t: string) => sessions[t] ?? null;
const sealed = { iv: "AAAAAAAAAAAAAAAA", ct: "Y2lwaGVydGV4dC1ieXRlcy1oZXJlLW9ubHk" };
const id = "3f2a6c1e-9b1d-4c33-8f00-1a2b3c4d5e6f";

function client(store = new VaultStore()) {
  const handler = createHandler(introspect, store);
  return (method: string, path: string, token?: string, body?: unknown) =>
    handler(new Request(`http://localhost${path}`, {
      method,
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: body === undefined ? undefined : JSON.stringify(body),
    }));
}

test("a bookstore session cannot open the vault even for the same username", async () => {
  const call = client();
  expect((await call("GET", "/api/vault", "shop")).status).toBe(401);
  expect((await call("GET", "/api/vault")).status).toBe(401);
});

test("the server only accepts ciphertext shaped payloads, so a buggy client cannot store plaintext fields", async () => {
  const call = client();
  expect((await call("PUT", `/api/vault/items/${id}`, "alice", { name: "bank", password: "hunter2" })).status).toBe(400);
  expect((await call("PUT", `/api/vault/items/${id}`, "alice", { iv: "short", ct: sealed.ct })).status).toBe(400);
  expect((await call("PUT", `/api/vault/items/${id}`, "alice", sealed)).status).toBe(200);
});

test("users only ever see their own items", async () => {
  const call = client();
  await call("PUT", `/api/vault/items/${id}`, "alice", sealed);
  const mallory = await (await call("GET", "/api/vault", "mallory")).json();
  expect(mallory.items).toHaveLength(0);
  expect((await call("DELETE", `/api/vault/items/${id}`, "mallory")).status).toBe(404);
});

test("the canary can be set once so a stolen session cannot swap the key check", async () => {
  const call = client();
  expect((await call("PUT", "/api/vault/canary", "alice", sealed)).status).toBe(201);
  expect((await call("PUT", "/api/vault/canary", "alice", sealed)).status).toBe(409);
});

test("item ids must be uuids so paths cannot be abused", async () => {
  const call = client();
  expect((await call("PUT", "/api/vault/items/..%2Fx", "alice", sealed)).status).toBe(400);
});

test("persisted vault file is owner only and holds no plaintext", async () => {
  const path = join(mkdtempSync(join(tmpdir(), "vault-")), "vault.json");
  const call = client(new VaultStore(path));
  await call("PUT", `/api/vault/items/${id}`, "alice", sealed);
  expect(statSync(path).mode & 0o777).toBe(0o600);
  expect(Object.keys(JSON.parse(readFileSync(path, "utf8")).alice.items[0]).sort()).toEqual(["ct", "id", "iv", "updatedAt"]);
});
