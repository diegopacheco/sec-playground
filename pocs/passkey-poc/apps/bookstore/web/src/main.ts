import "./style.css";
import { h, mount } from "@passkey/dom";
import { logout, passkeys, type AuthResult } from "@passkey/passkey-client";
import { books, ShopApi } from "./api.ts";
import { authView } from "./views/auth.ts";
import { booksView, cartView, ordersView, securityView } from "./views/shop.ts";

type Tab = "books" | "cart" | "orders" | "security";
const STORAGE_KEY = "bookstore.auth";
const root = document.getElementById("app")!;

function loadAuth(): AuthResult | null {
  const raw = sessionStorage.getItem(STORAGE_KEY);
  if (!raw) return null;
  const auth = JSON.parse(raw) as AuthResult;
  return new Date(auth.session.expiresAt) > new Date() ? auth : null;
}

let auth = loadAuth();
let tab: Tab = "books";

function signedIn(r: AuthResult) {
  auth = r;
  sessionStorage.setItem(STORAGE_KEY, JSON.stringify(r));
  tab = "books";
  render();
}

async function signOut() {
  if (auth) await logout(auth.token).catch(() => undefined);
  auth = null;
  sessionStorage.removeItem(STORAGE_KEY);
  render();
}

function header(cartCount: number): HTMLElement {
  const nav = (t: Tab, label: string) => h("button", { class: `tab ${tab === t ? "active" : ""}`, "data-tab": t, onclick: () => { tab = t; render(); } }, label);
  return h("header", {},
    h("div", { class: "brand" }, h("span", { class: "mark" }, "PB"), h("span", {}, "Passkey Books")),
    auth ? h("nav", {}, nav("books", "Books"), nav("cart", `Cart (${cartCount})`), nav("orders", "Orders"), nav("security", "Security")) : null,
    auth ? h("div", { class: "who" }, h("span", { class: "chip" }, auth.session.username), h("button", { class: "small ghost", id: "logout", onclick: signOut }, "Sign out")) : null,
  );
}

async function render(next?: string) {
  if (next) tab = next as Tab;
  if (!auth) {
    mount(root, header(0), h("main", {}, authView(signedIn)));
    return;
  }
  const api = new ShopApi(auth.token);
  try {
    const cart = await api.cart();
    const count = cart.lines.reduce((s, l) => s + l.qty, 0);
    let body: HTMLElement;
    if (tab === "books") body = booksView(await books(), api, render);
    else if (tab === "cart") body = cartView(cart, api, render);
    else if (tab === "orders") body = ordersView(await api.orders());
    else body = securityView(await passkeys("bookstore", auth.token), auth.session);
    mount(root, header(count), h("main", {}, body));
  } catch {
    await signOut();
  }
}

render();
