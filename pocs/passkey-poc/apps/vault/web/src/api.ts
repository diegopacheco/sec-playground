import type { Sealed } from "./crypto.ts";

export type StoredItem = Sealed & { id: string; updatedAt: string };
export type StoredVault = { canary: Sealed | null; items: StoredItem[] };

export class VaultApi {
  constructor(private token: string) {}

  private async call<T>(method: string, path: string, body?: unknown): Promise<T> {
    const res = await fetch(path, {
      method,
      headers: { Authorization: `Bearer ${this.token}`, "Content-Type": "application/json" },
      body: body === undefined ? undefined : JSON.stringify(body),
    });
    if (res.status === 204) return undefined as T;
    const data = await res.json();
    if (!res.ok) throw new Error(data.error ?? `request failed with ${res.status}`);
    return data as T;
  }

  load = () => this.call<StoredVault>("GET", "/api/vault");
  setCanary = (canary: Sealed) => this.call<void>("PUT", "/api/vault/canary", canary);
  put = (id: string, sealed: Sealed) => this.call<StoredItem>("PUT", `/api/vault/items/${id}`, sealed);
  remove = (id: string) => this.call<void>("DELETE", `/api/vault/items/${id}`);
}
