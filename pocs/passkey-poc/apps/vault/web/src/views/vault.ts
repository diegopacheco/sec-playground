import { h, toast } from "@passkey/dom";
import type { StoredVault } from "../api.ts";
import { open, seal } from "../crypto.ts";
import { DEFAULT_OPTIONS, entropyBits, generatePassword, type GeneratorOptions } from "../generator.ts";
import type { Unlocked } from "../unlock.ts";

export type Entry = { name: string; url: string; username: string; password: string; notes: string };
export type Decrypted = { id: string; updatedAt: string; entry: Entry };

const EMPTY: Entry = { name: "", url: "", username: "", password: "", notes: "" };

export async function decryptAll(u: Unlocked, vault: StoredVault): Promise<Decrypted[]> {
  const out = await Promise.all(vault.items.map(async (i) => ({ id: i.id, updatedAt: i.updatedAt, entry: await open<Entry>(u.key, i, i.id) })));
  return out.sort((a, b) => a.entry.name.localeCompare(b.entry.name));
}

function field(label: string, name: keyof Entry, value: string, type = "text"): HTMLElement {
  const input = name === "notes" ? h("textarea", { name, rows: 4 }) : h("input", { name, type, value, autocomplete: "off" });
  if (input instanceof HTMLTextAreaElement) input.value = value;
  return h("label", { class: "field" }, h("span", {}, label), input);
}

function generatorPanel(onUse: (pw: string) => void): HTMLElement {
  const opts: GeneratorOptions = { ...DEFAULT_OPTIONS };
  const out = h("code", { class: "gen-out", id: "gen-out" });
  const bits = h("span", { class: "bits" });
  const draw = () => {
    try {
      out.textContent = generatePassword(opts);
      bits.textContent = `${entropyBits(opts)} bits`;
    } catch (err) {
      out.textContent = (err as Error).message;
      bits.textContent = "";
    }
  };
  const toggle = (k: "lower" | "upper" | "digits" | "symbols", label: string) =>
    h("label", { class: "check" }, h("input", { type: "checkbox", checked: opts[k], onchange: (e: Event) => { opts[k] = (e.target as HTMLInputElement).checked; draw(); } }), label);
  const length = h("input", { type: "range", min: 8, max: 64, value: opts.length, oninput: (e: Event) => { opts.length = Number((e.target as HTMLInputElement).value); lenLabel.textContent = String(opts.length); draw(); } });
  const lenLabel = h("strong", {}, opts.length);
  draw();
  return h("div", { class: "generator" },
    h("div", { class: "gen-head" }, h("h3", {}, "Generator"), bits),
    out,
    h("div", { class: "gen-row" }, h("span", {}, "Length"), length, lenLabel),
    h("div", { class: "gen-row" }, toggle("lower", "a-z"), toggle("upper", "A-Z"), toggle("digits", "2-9"), toggle("symbols", "!@#")),
    h("div", { class: "gen-row" },
      h("button", { class: "small ghost", onclick: draw }, "Regenerate"),
      h("button", { class: "small", id: "use-generated", onclick: () => onUse(out.textContent ?? "") }, "Use password"),
    ),
  );
}

export function vaultView(u: Unlocked, items: Decrypted[], reload: () => void): HTMLElement {
  let selected: string | null = items[0]?.id ?? null;
  let query = "";
  const list = h("div", { class: "list", id: "entries" });
  const detail = h("div", { class: "detail" });

  const drawList = () => {
    const q = query.toLowerCase();
    const visible = items.filter((i) => `${i.entry.name} ${i.entry.url} ${i.entry.username}`.toLowerCase().includes(q));
    list.replaceChildren(
      ...visible.map((i) =>
        h("button", { class: `entry ${i.id === selected ? "active" : ""}`, "data-name": i.entry.name, onclick: () => { selected = i.id; drawList(); drawDetail(); } },
          h("span", { class: "avatar" }, (i.entry.name[0] ?? "?").toUpperCase()),
          h("span", { class: "entry-text" }, h("strong", {}, i.entry.name), h("small", {}, i.entry.username || i.entry.url)),
        ),
      ),
    );
    if (!visible.length) list.append(h("p", { class: "muted pad" }, items.length ? "No match." : "Your vault is empty. Add your first login."));
  };

  const drawDetail = (editing: Decrypted | null = items.find((i) => i.id === selected) ?? null, isNew = false) => {
    if (!editing && !isNew) {
      detail.replaceChildren(h("div", { class: "placeholder" }, h("h2", {}, "Nothing selected"), h("p", {}, "Pick an entry or add a new one.")));
      return;
    }
    const entry = editing?.entry ?? EMPTY;
    const id = editing?.id ?? crypto.randomUUID();
    const form = h("form", { class: "editor", id: "editor" },
      field("Name", "name", entry.name),
      field("Website", "url", entry.url),
      field("Username", "username", entry.username),
      h("div", { class: "pw-row" }, field("Password", "password", entry.password, "password"),
        h("button", { type: "button", class: "small ghost", id: "reveal", onclick: () => { const p = form.querySelector<HTMLInputElement>("[name=password]")!; p.type = p.type === "password" ? "text" : "password"; } }, "Show"),
        h("button", { type: "button", class: "small ghost", id: "copy", onclick: async () => { await navigator.clipboard.writeText(form.querySelector<HTMLInputElement>("[name=password]")!.value); toast("Password copied"); } }, "Copy"),
      ),
      field("Notes", "notes", entry.notes),
      h("div", { class: "actions" },
        isNew ? null : h("button", { type: "button", class: "danger ghost", id: "delete", onclick: async () => { await u.api.remove(id); toast("Entry deleted"); reload(); } }, "Delete"),
        h("button", { type: "submit", class: "primary", id: "save" }, isNew ? "Encrypt and add" : "Encrypt and save"),
      ),
    );
    form.addEventListener("submit", async (e) => {
      e.preventDefault();
      const data = Object.fromEntries(new FormData(form)) as unknown as Entry;
      if (!data.name.trim()) return toast("Name is required", "error");
      await u.api.put(id, await seal(u.key, data, id));
      selected = id;
      toast("Encrypted and saved");
      reload();
    });
    const gen = generatorPanel((pw) => {
      const p = form.querySelector<HTMLInputElement>("[name=password]")!;
      p.value = pw;
      p.type = "text";
    });
    detail.replaceChildren(h("div", { class: "detail-head" }, h("h2", {}, isNew ? "New login" : entry.name), editing ? h("small", { class: "muted" }, `updated ${new Date(editing.updatedAt).toLocaleString()}`) : null), form, gen);
  };

  drawList();
  drawDetail();
  return h("div", { class: "vault" },
    h("aside", {},
      h("input", { class: "search", id: "search", placeholder: `Search ${items.length} entries`, oninput: (e: Event) => { query = (e.target as HTMLInputElement).value; drawList(); } }),
      h("button", { class: "primary wide", id: "add", onclick: () => { selected = null; drawList(); drawDetail(null, true); } }, "Add login"),
      list,
    ),
    detail,
  );
}

export function serverView(vault: StoredVault, username: string): HTMLElement {
  return h("section", { class: "server" },
    h("h2", {}, "What the server stores"),
    h("p", { class: "muted" }, `This is the exact JSON the vault API returns for ${username}. Names, sites, usernames, passwords and notes are all inside ct. The key to open it never left your browser.`),
    h("pre", { id: "raw" }, JSON.stringify(vault, null, 2)),
  );
}
