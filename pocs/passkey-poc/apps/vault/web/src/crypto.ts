import { fromBase64url, toBase64url } from "@passkey/passkey-client/base64url";

export type Sealed = { iv: string; ct: string };

const enc = new TextEncoder();
const dec = new TextDecoder();

export const PRF_SALT = enc.encode("passkey-vault/prf/v1");
export const CANARY_AAD = "canary";
const CANARY_TEXT = "passkey-vault-ok";

export async function deriveVaultKey(prfOutput: ArrayBuffer, username: string): Promise<CryptoKey> {
  const ikm = await crypto.subtle.importKey("raw", prfOutput, "HKDF", false, ["deriveKey"]);
  return crypto.subtle.deriveKey(
    { name: "HKDF", hash: "SHA-256", salt: enc.encode("passkey-vault"), info: enc.encode(`vault-key/v1/${username}`) },
    ikm,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"],
  );
}

export async function seal(key: CryptoKey, value: unknown, aad: string): Promise<Sealed> {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const ct = await crypto.subtle.encrypt({ name: "AES-GCM", iv, additionalData: enc.encode(aad) }, key, enc.encode(JSON.stringify(value)));
  return { iv: toBase64url(iv), ct: toBase64url(ct) };
}

export async function open<T>(key: CryptoKey, sealed: Sealed, aad: string): Promise<T> {
  const pt = await crypto.subtle.decrypt({ name: "AES-GCM", iv: fromBase64url(sealed.iv), additionalData: enc.encode(aad) }, key, fromBase64url(sealed.ct));
  return JSON.parse(dec.decode(pt)) as T;
}

export function sealCanary(key: CryptoKey): Promise<Sealed> {
  return seal(key, CANARY_TEXT, CANARY_AAD);
}

export async function checkCanary(key: CryptoKey, canary: Sealed): Promise<boolean> {
  try {
    return (await open<string>(key, canary, CANARY_AAD)) === CANARY_TEXT;
  } catch {
    return false;
  }
}
