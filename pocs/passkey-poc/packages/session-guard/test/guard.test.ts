import { expect, test } from "bun:test";
import { requireSession, type Session } from "../src/index.ts";

const vaultSession: Session = { app: "vault", username: "alice", userId: "u1", expiresAt: "" };
const introspect = async (token: string) => (token === "vault-token" ? vaultSession : null);

function req(token?: string): Request {
  return new Request("http://localhost/api/x", { headers: token ? { Authorization: `Bearer ${token}` } : {} });
}

test("a token issued for another app is refused, so a bookstore login can never open the vault", async () => {
  expect(await requireSession(req("vault-token"), "bookstore", introspect)).toBeNull();
  expect(await requireSession(req("vault-token"), "vault", introspect)).toEqual(vaultSession);
});

test("missing or unknown tokens are refused", async () => {
  expect(await requireSession(req(), "vault", introspect)).toBeNull();
  expect(await requireSession(req("forged"), "vault", introspect)).toBeNull();
});
