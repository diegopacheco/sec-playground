import { h, toast } from "@passkey/dom";
import { algName, type PasskeyInfo } from "@passkey/passkey-client";
import { money, type Book, type Cart, type Order, type ShopApi } from "../api.ts";

function cover(b: Book): HTMLElement {
  return h("div", { class: "cover", style: `--c:${b.color}` }, h("span", { class: "cover-title" }, b.title), h("span", { class: "cover-author" }, b.author));
}

export function booksView(list: Book[], api: ShopApi, refresh: (tab?: string) => void): HTMLElement {
  const grid = h("div", { class: "grid" });
  const draw = (q: string) => {
    const term = q.toLowerCase();
    grid.replaceChildren(
      ...list
        .filter((b) => `${b.title} ${b.author} ${b.tag}`.toLowerCase().includes(term))
        .map((b) =>
          h("article", { class: "book", "data-id": b.id },
            cover(b),
            h("div", { class: "book-meta" },
              h("h3", {}, b.title),
              h("p", {}, `${b.author} · ${b.year}`),
              h("div", { class: "book-row" },
                h("span", { class: "price" }, money(b.price)),
                h("button", {
                  class: "small",
                  onclick: async () => {
                    const cart = await api.cart();
                    const qty = (cart.lines.find((l) => l.book.id === b.id)?.qty ?? 0) + 1;
                    try {
                      await api.setQty(b.id, qty);
                      toast(`Added ${b.title}`);
                      refresh();
                    } catch (err) {
                      toast((err as Error).message, "error");
                    }
                  },
                }, "Add to cart"),
              ),
            ),
          ),
        ),
    );
  };
  const search = h("input", { class: "search", placeholder: "Search title, author or topic", oninput: (e: Event) => draw((e.target as HTMLInputElement).value) });
  draw("");
  return h("section", {}, h("div", { class: "section-head" }, h("h2", {}, "The security shelf"), search), grid);
}

export function cartView(cart: Cart, api: ShopApi, refresh: (tab?: string) => void): HTMLElement {
  if (cart.lines.length === 0) return h("section", { class: "empty" }, h("h2", {}, "Your cart is empty"), h("p", {}, "Pick something from the shelf."));
  const set = (id: string, qty: number) => async () => {
    await api.setQty(id, qty);
    refresh();
  };
  return h("section", {},
    h("h2", {}, "Cart"),
    h("div", { class: "lines" },
      ...cart.lines.map((l) =>
        h("div", { class: "line" },
          cover(l.book),
          h("div", { class: "line-info" }, h("strong", {}, l.book.title), h("span", {}, money(l.book.price))),
          h("div", { class: "qty" },
            h("button", { class: "small ghost", onclick: set(l.book.id, l.qty - 1) }, "-"),
            h("span", {}, l.qty),
            h("button", { class: "small ghost", onclick: set(l.book.id, l.qty + 1) }, "+"),
          ),
          h("strong", { class: "subtotal" }, money(l.subtotal)),
        ),
      ),
    ),
    h("div", { class: "checkout" },
      h("span", {}, "Total ", h("strong", { id: "total" }, money(cart.total))),
      h("button", {
        class: "primary",
        id: "checkout",
        onclick: async () => {
          const order = await api.checkout();
          toast(`Order ${order.id.slice(0, 8)} placed`);
          refresh("orders");
        },
      }, "Place order"),
    ),
  );
}

export function ordersView(orders: Order[]): HTMLElement {
  if (orders.length === 0) return h("section", { class: "empty" }, h("h2", {}, "No orders yet"));
  return h("section", {},
    h("h2", {}, "Orders"),
    ...orders.map((o) =>
      h("div", { class: "order" },
        h("div", { class: "order-head" }, h("strong", {}, `#${o.id.slice(0, 8)}`), h("span", {}, new Date(o.createdAt).toLocaleString()), h("strong", {}, money(o.total))),
        h("ul", {}, ...o.lines.map((l) => h("li", {}, `${l.qty} x ${l.book.title}`))),
      ),
    ),
  );
}

export function securityView(keys: PasskeyInfo[], session: { username: string; expiresAt: string }): HTMLElement {
  return h("section", {},
    h("h2", {}, "Security"),
    h("p", { class: "lede" }, `Signed in as ${session.username}. Session ends ${new Date(session.expiresAt).toLocaleTimeString()}.`),
    h("table", { class: "keys" },
      h("thead", {}, h("tr", {}, h("th", {}, "Credential"), h("th", {}, "Algorithm"), h("th", {}, "Sign count"), h("th", {}, "Synced"), h("th", {}, "Transports"))),
      h("tbody", {},
        ...keys.map((k) =>
          h("tr", {},
            h("td", { class: "mono" }, `${k.id.slice(0, 16)}...`),
            h("td", {}, algName(k.alg)),
            h("td", {}, k.signCount),
            h("td", {}, k.backedUp ? "yes" : k.backupEligible ? "eligible" : "device bound"),
            h("td", {}, (k.transports ?? []).join(", ") || "-"),
          ),
        ),
      ),
    ),
  );
}
