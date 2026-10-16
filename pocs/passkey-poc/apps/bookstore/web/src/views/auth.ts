import { h, toast } from "@passkey/dom";
import { login, register, supported, type AuthResult } from "@passkey/passkey-client";

export function authView(onAuth: (r: AuthResult) => void): HTMLElement {
  const input = h("input", { id: "username", placeholder: "username", autocomplete: "username webauthn", maxlength: 32 });
  const run = (fn: (app: string, name: string) => Promise<AuthResult>, verb: string) => async () => {
    const name = input.value.trim();
    if (!name) return toast("Type a username first", "error");
    try {
      onAuth(await fn("bookstore", name));
      toast(`${verb} as ${name}`);
    } catch (err) {
      toast((err as Error).message, "error");
    }
  };
  return h(
    "section",
    { class: "auth" },
    h("div", { class: "auth-art" },
      h("p", { class: "eyebrow" }, "No passwords on these shelves"),
      h("h1", {}, "Read more.", h("br"), "Type less."),
      h("p", { class: "lede" }, "Sign in with a passkey: your device keeps a private key and signs a fresh challenge. Nothing reusable ever leaves it, so there is nothing to phish and nothing to leak."),
    ),
    h("div", { class: "auth-card" },
      h("h2", {}, "Your library card"),
      supported() ? null : h("p", { class: "warn" }, "This browser does not support passkeys."),
      h("label", { for: "username" }, "Username"),
      input,
      h("button", { class: "primary", id: "signin", onclick: run(login, "Signed in") }, "Sign in with passkey"),
      h("button", { class: "ghost", id: "signup", onclick: run(register, "Account created") }, "Create account with passkey"),
      h("p", { class: "fine" }, "Touch ID, Windows Hello, a phone or a security key."),
    ),
  );
}
