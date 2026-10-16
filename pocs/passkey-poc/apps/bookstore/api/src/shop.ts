import { findBook, type Book } from "./catalog.ts";

export type CartLine = { book: Book; qty: number; subtotal: number };
export type Cart = { lines: CartLine[]; total: number };
export type Order = { id: string; createdAt: string; lines: CartLine[]; total: number };

const MAX_QTY = 10;

export class Shop {
  private carts = new Map<string, Map<string, number>>();
  private orders = new Map<string, Order[]>();

  cart(user: string): Cart {
    const items = this.carts.get(user) ?? new Map<string, number>();
    const lines: CartLine[] = [];
    for (const [id, qty] of items) {
      const book = findBook(id);
      if (book) lines.push({ book, qty, subtotal: book.price * qty });
    }
    return { lines, total: lines.reduce((s, l) => s + l.subtotal, 0) };
  }

  setQty(user: string, bookId: string, qty: number): Cart {
    if (!findBook(bookId)) throw new Error("unknown book");
    if (!Number.isInteger(qty) || qty < 0 || qty > MAX_QTY) throw new Error(`quantity must be 0-${MAX_QTY}`);
    const items = this.carts.get(user) ?? new Map<string, number>();
    if (qty === 0) items.delete(bookId);
    else items.set(bookId, qty);
    this.carts.set(user, items);
    return this.cart(user);
  }

  checkout(user: string): Order {
    const cart = this.cart(user);
    if (cart.lines.length === 0) throw new Error("cart is empty");
    const order: Order = { id: crypto.randomUUID(), createdAt: new Date().toISOString(), lines: cart.lines, total: cart.total };
    this.orders.set(user, [order, ...(this.orders.get(user) ?? [])]);
    this.carts.delete(user);
    return order;
  }

  history(user: string): Order[] {
    return this.orders.get(user) ?? [];
  }
}
