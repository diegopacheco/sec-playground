import { login, register, type AuthResult } from "@passkey/passkey-client";
import { VaultApi } from "./api.ts";
import { checkCanary, deriveVaultKey, PRF_SALT, sealCanary } from "./crypto.ts";

export type Unlocked = { auth: AuthResult; key: CryptoKey; api: VaultApi };

const APP = "vault";

async function unlockWith(auth: AuthResult, username: string): Promise<Unlocked> {
  if (!auth.prf) throw new Error("This passkey does not support the PRF extension, so it cannot derive a vault key");
  const key = await deriveVaultKey(auth.prf, username);
  const api = new VaultApi(auth.token);
  const vault = await api.load();
  if (!vault.canary) await api.setCanary(await sealCanary(key));
  else if (!(await checkCanary(key, vault.canary))) throw new Error("This passkey cannot decrypt this vault");
  return { auth, key, api };
}

export async function unlock(username: string): Promise<Unlocked> {
  return unlockWith(await login(APP, username, PRF_SALT), username);
}

export async function createVault(username: string): Promise<Unlocked> {
  await register(APP, username, true);
  return unlock(username);
}
