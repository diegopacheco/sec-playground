export type Book = { id: string; title: string; author: string; year: number; price: number; tag: string; color: string };
export type CartLine = { book: Book; qty: number; subtotal: number };
export type Cart = { lines: CartLine[]; total: number };
export type Order = { id: string; createdAt: string; lines: CartLine[]; total: number };

export class ShopApi {
  constructor(private token: string) {}

  private async call<T>(method: string, path: string, body?: unknown): Promise<T> {
    const res = await fetch(path, {
      method,
      headers: { Authorization: `Bearer ${this.token}`, "Content-Type": "application/json" },
      body: body === undefined ? undefined : JSON.stringify(body),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error ?? `request failed with ${res.status}`);
    return data as T;
  }

  cart = () => this.call<Cart>("GET", "/api/cart");
  setQty = (bookId: string, qty: number) => this.call<Cart>("PUT", "/api/cart", { bookId, qty });
  checkout = () => this.call<Order>("POST", "/api/checkout");
  orders = () => this.call<Order[]>("GET", "/api/orders");
}

export async function books(): Promise<Book[]> {
  const res = await fetch("/api/books");
  return res.json();
}

export function money(cents: number): string {
  return `$${(cents / 100).toFixed(2)}`;
}
