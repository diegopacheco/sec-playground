import { h, toast } from "@passkey/dom";
import { supported } from "@passkey/passkey-client";
import { createVault, unlock, type Unlocked } from "../unlock.ts";

export function lockView(onUnlock: (u: Unlocked) => void): HTMLElement {
  const input = h("input", { id: "username", placeholder: "username", autocomplete: "username webauthn", maxlength: 32 });
  const status = h("p", { class: "status", id: "status" });
  const run = (fn: (name: string) => Promise<Unlocked>, busy: string) => async () => {
    const name = input.value.trim();
    if (!name) return toast("Type a username first", "error");
    status.textContent = busy;
    try {
      onUnlock(await fn(name));
    } catch (err) {
      status.textContent = "";
      toast((err as Error).message, "error");
    }
  };
  return h("section", { class: "lock" },
    h("div", { class: "lock-card" },
      h("div", { class: "vault-mark" }, h("span", {}, "")),
      h("h1", {}, "PasskeyVault"),
      h("p", { class: "sub" }, "Your passwords, encrypted in this browser with a key that only your passkey can produce. The server stores ciphertext and nothing else."),
      supported() ? null : h("p", { class: "warn" }, "This browser does not support passkeys."),
      h("label", { for: "username" }, "Username"),
      input,
      h("button", { class: "primary", id: "unlock", onclick: run(unlock, "Waiting for your passkey...") }, "Unlock with passkey"),
      h("button", { class: "ghost", id: "create", onclick: run(createVault, "Create the passkey, then touch it again to derive the vault key...") }, "Create a new vault"),
      status,
      h("ul", { class: "facts" },
        h("li", {}, "No master password to forget or phish"),
        h("li", {}, "Vault key = HKDF(passkey PRF output), never sent anywhere"),
        h("li", {}, "AES-256-GCM per entry, bound to its id"),
      ),
    ),
  );
}
