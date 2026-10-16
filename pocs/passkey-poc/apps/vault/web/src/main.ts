import "./style.css";
import { h, mount, toast } from "@passkey/dom";
import { algName, logout, passkeys } from "@passkey/passkey-client";
import type { Unlocked } from "./unlock.ts";
import { lockView } from "./views/lock.ts";
import { decryptAll, serverView, vaultView } from "./views/vault.ts";

type Tab = "vault" | "server" | "passkey";
const IDLE_LOCK_MS = 5 * 60 * 1000;
const root = document.getElementById("app")!;

let unlocked: Unlocked | null = null;
let tab: Tab = "vault";
let idleTimer: ReturnType<typeof setTimeout> | undefined;

async function lock(reason?: string) {
  if (unlocked) await logout(unlocked.auth.token).catch(() => undefined);
  unlocked = null;
  clearTimeout(idleTimer);
  if (reason) toast(reason);
  render();
}

function touch() {
  clearTimeout(idleTimer);
  if (unlocked) idleTimer = setTimeout(() => lock("Vault locked after 5 minutes idle"), IDLE_LOCK_MS);
}

function header(): HTMLElement {
  const nav = (t: Tab, label: string) => h("button", { class: `tab ${tab === t ? "active" : ""}`, "data-tab": t, onclick: () => { tab = t; render(); } }, label);
  return h("header", {},
    h("div", { class: "brand" }, h("span", { class: "vault-dot" }), "PasskeyVault"),
    unlocked ? h("nav", {}, nav("vault", "Vault"), nav("server", "Server view"), nav("passkey", "Passkey")) : null,
    unlocked ? h("div", { class: "who" }, h("span", { class: "chip" }, unlocked.auth.session.username), h("button", { class: "small ghost", id: "lock", onclick: () => lock("Vault locked") }, "Lock")) : null,
  );
}

async function passkeyView(u: Unlocked): Promise<HTMLElement> {
  const keys = await passkeys("vault", u.auth.token);
  return h("section", { class: "server" },
    h("h2", {}, "Your vault passkey"),
    h("p", { class: "muted" }, "The passkey both authenticates you to the server and, through the PRF extension, produces the secret that derives your vault key. Lose every passkey and the vault cannot be decrypted by anyone."),
    ...keys.map((k) => h("dl", { class: "kv" },
      h("dt", {}, "Credential"), h("dd", { class: "mono" }, k.id),
      h("dt", {}, "Algorithm"), h("dd", {}, algName(k.alg)),
      h("dt", {}, "Sign count"), h("dd", {}, k.signCount),
      h("dt", {}, "Synced"), h("dd", {}, k.backedUp ? "yes" : k.backupEligible ? "eligible" : "device bound"),
    )),
  );
}

async function render() {
  if (!unlocked) {
    mount(root, header(), h("main", {}, lockView((u) => { unlocked = u; tab = "vault"; toast("Vault unlocked"); touch(); render(); })));
    return;
  }
  const u = unlocked;
  try {
    const stored = await u.api.load();
    let body: HTMLElement;
    if (tab === "vault") body = vaultView(u, await decryptAll(u, stored), render);
    else if (tab === "server") body = serverView(stored, u.auth.session.username);
    else body = await passkeyView(u);
    mount(root, header(), h("main", {}, body));
  } catch (err) {
    await lock((err as Error).message);
  }
}

for (const ev of ["click", "keydown"]) document.addEventListener(ev, touch);
render();
