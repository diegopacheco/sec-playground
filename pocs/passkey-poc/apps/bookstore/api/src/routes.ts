import { json, readJSON, requireSession, unauthorized, type Introspect } from "@passkey/session-guard";
import { BOOKS } from "./catalog.ts";
import { Shop } from "./shop.ts";

const APP = "bookstore";

export function createHandler(introspect: Introspect, shop = new Shop()) {
  return async (req: Request): Promise<Response> => {
    const url = new URL(req.url);
    const route = `${req.method} ${url.pathname}`;
    if (route === "GET /api/health") return json({ status: "ok" });
    if (route === "GET /api/books") return json(BOOKS);

    const session = await requireSession(req, APP, introspect);
    if (!session) return unauthorized();
    const user = session.username;

    try {
      if (route === "GET /api/me") return json(session);
      if (route === "GET /api/cart") return json(shop.cart(user));
      if (route === "PUT /api/cart") {
        const body = await readJSON<{ bookId: string; qty: number }>(req);
        if (!body) return json({ error: "invalid body" }, 400);
        return json(shop.setQty(user, body.bookId, body.qty));
      }
      if (route === "POST /api/checkout") return json(shop.checkout(user), 201);
      if (route === "GET /api/orders") return json(shop.history(user));
    } catch (err) {
      return json({ error: (err as Error).message }, 400);
    }
    return json({ error: "not found" }, 404);
  };
}
