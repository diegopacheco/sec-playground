import { expect, test } from "bun:test";
import type { Session } from "@passkey/session-guard";
import { createHandler } from "../src/routes.ts";

const sessions: Record<string, Session> = {
  alice: { app: "bookstore", username: "alice", userId: "a", expiresAt: "" },
  bob: { app: "bookstore", username: "bob", userId: "b", expiresAt: "" },
  vault: { app: "vault", username: "alice", userId: "a", expiresAt: "" },
};
const handler = createHandler(async (t) => sessions[t] ?? null);

function call(method: string, path: string, token?: string, body?: unknown) {
  const headers: Record<string, string> = token ? { Authorization: `Bearer ${token}` } : {};
  return handler(new Request(`http://localhost${path}`, { method, headers, body: body ? JSON.stringify(body) : undefined }));
}

test("catalog is public but the cart needs a passkey session", async () => {
  expect((await call("GET", "/api/books")).status).toBe(200);
  expect((await call("GET", "/api/cart")).status).toBe(401);
});

test("a vault session cannot shop as the same username", async () => {
  expect((await call("GET", "/api/cart", "vault")).status).toBe(401);
});

test("prices come from the server catalog, the client cannot send its own price", async () => {
  const res = await call("PUT", "/api/cart", "alice", { bookId: "b1", qty: 2, price: 1 });
  const cart = await res.json();
  expect(cart.total).toBe(5499 * 2);
});

test("checkout turns the cart into an order and empties the cart", async () => {
  await call("PUT", "/api/cart", "bob", { bookId: "b2", qty: 1 });
  const order = await (await call("POST", "/api/checkout", "bob")).json();
  expect(order.total).toBe(3999);
  expect((await (await call("GET", "/api/cart", "bob")).json()).lines).toHaveLength(0);
  expect(await (await call("GET", "/api/orders", "bob")).json()).toHaveLength(1);
});

test("carts are isolated between users", async () => {
  await call("PUT", "/api/cart", "alice", { bookId: "b3", qty: 1 });
  const bobCart = await (await call("GET", "/api/cart", "bob")).json();
  expect(bobCart.lines.find((l: { book: { id: string } }) => l.book.id === "b3")).toBeUndefined();
});

test("invalid quantities and unknown books are rejected", async () => {
  expect((await call("PUT", "/api/cart", "alice", { bookId: "b1", qty: -1 })).status).toBe(400);
  expect((await call("PUT", "/api/cart", "alice", { bookId: "b1", qty: 1.5 })).status).toBe(400);
  expect((await call("PUT", "/api/cart", "alice", { bookId: "nope", qty: 1 })).status).toBe(400);
  expect((await call("POST", "/api/checkout", "nobody")).status).toBe(401);
});
