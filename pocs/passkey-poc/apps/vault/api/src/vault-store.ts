import { existsSync, mkdirSync, readFileSync, renameSync, writeFileSync } from "node:fs";
import { dirname } from "node:path";

export type Sealed = { iv: string; ct: string };
export type Item = Sealed & { id: string; updatedAt: string };
export type Vault = { canary: Sealed | null; items: Item[] };

const B64URL = /^[A-Za-z0-9_-]+$/;
const UUID = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/;
const IV_CHARS = 16;
const MAX_CT_CHARS = 32 * 1024;
const MIN_CT_CHARS = 22;
const MAX_ITEMS = 500;

export function validSealed(v: unknown): v is Sealed {
  if (typeof v !== "object" || v === null) return false;
  const { iv, ct } = v as Record<string, unknown>;
  return (
    typeof iv === "string" && iv.length === IV_CHARS && B64URL.test(iv) &&
    typeof ct === "string" && ct.length >= MIN_CT_CHARS && ct.length <= MAX_CT_CHARS && B64URL.test(ct)
  );
}

export function validId(id: string): boolean {
  return UUID.test(id);
}

export class VaultStore {
  private vaults: Record<string, Vault> = {};

  constructor(private path = "") {
    if (path && existsSync(path)) this.vaults = JSON.parse(readFileSync(path, "utf8"));
  }

  private save() {
    if (!this.path) return;
    mkdirSync(dirname(this.path), { recursive: true, mode: 0o700 });
    writeFileSync(`${this.path}.tmp`, JSON.stringify(this.vaults, null, 2), { mode: 0o600 });
    renameSync(`${this.path}.tmp`, this.path);
  }

  get(user: string): Vault {
    return this.vaults[user] ?? { canary: null, items: [] };
  }

  setCanary(user: string, canary: Sealed): boolean {
    const v = this.get(user);
    if (v.canary) return false;
    this.vaults[user] = { ...v, canary: { iv: canary.iv, ct: canary.ct } };
    this.save();
    return true;
  }

  put(user: string, id: string, sealed: Sealed): Item {
    const v = this.get(user);
    const exists = v.items.some((i) => i.id === id);
    if (!exists && v.items.length >= MAX_ITEMS) throw new Error("vault is full");
    const item: Item = { id, iv: sealed.iv, ct: sealed.ct, updatedAt: new Date().toISOString() };
    this.vaults[user] = { ...v, items: [...v.items.filter((i) => i.id !== id), item] };
    this.save();
    return item;
  }

  remove(user: string, id: string): boolean {
    const v = this.get(user);
    const items = v.items.filter((i) => i.id !== id);
    if (items.length === v.items.length) return false;
    this.vaults[user] = { ...v, items };
    this.save();
    return true;
  }
}
