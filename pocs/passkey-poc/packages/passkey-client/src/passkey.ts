import { fromBase64url, toBase64url } from "./base64url.ts";

export type Session = { app: string; username: string; userId: string; expiresAt: string };
export type AuthResult = { token: string; session: Session; prf?: ArrayBuffer };
export type PasskeyInfo = {
  id: string;
  alg: number;
  signCount: number;
  attestationFmt: string;
  backupEligible: boolean;
  backedUp: boolean;
  transports: string[] | null;
};

type Descriptor = { type: "public-key"; id: string; transports?: string[] };

type CreationJSON = {
  rp: { id: string; name: string };
  user: { id: string; name: string; displayName: string };
  challenge: string;
  pubKeyCredParams: { type: "public-key"; alg: number }[];
  timeout: number;
  excludeCredentials: Descriptor[];
  authenticatorSelection: AuthenticatorSelectionCriteria;
  attestation: AttestationConveyancePreference;
};

type RequestJSON = {
  challenge: string;
  rpId: string;
  timeout: number;
  allowCredentials: Descriptor[];
  userVerification: UserVerificationRequirement;
};

type Begin<T> = { ceremonyId: string; publicKey: T };

const BASE = "/passkey";

async function post<T>(path: string, body: unknown, token?: string): Promise<T> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token) headers.Authorization = `Bearer ${token}`;
  const res = await fetch(`${BASE}${path}`, { method: "POST", headers, body: JSON.stringify(body) });
  const data = res.status === 204 ? {} : await res.json();
  if (!res.ok) throw new Error(data.error ?? `request failed with ${res.status}`);
  return data as T;
}

function toDescriptors(list: Descriptor[]): PublicKeyCredentialDescriptor[] {
  return list.map((d) => ({ type: d.type, id: fromBase64url(d.id), transports: d.transports as AuthenticatorTransport[] | undefined }));
}

export function supported(): boolean {
  return typeof window !== "undefined" && "PublicKeyCredential" in window;
}

export async function register(app: string, username: string, withPrf = false): Promise<AuthResult> {
  const begin = await post<Begin<CreationJSON>>(`/${app}/register/begin`, { username });
  const o = begin.publicKey;
  const cred = (await navigator.credentials.create({
    publicKey: {
      rp: o.rp,
      user: { ...o.user, id: fromBase64url(o.user.id) },
      challenge: fromBase64url(o.challenge),
      pubKeyCredParams: o.pubKeyCredParams,
      timeout: o.timeout,
      excludeCredentials: toDescriptors(o.excludeCredentials),
      authenticatorSelection: o.authenticatorSelection,
      attestation: o.attestation,
      extensions: withPrf ? { prf: {} } : {},
    },
  })) as PublicKeyCredential | null;
  if (!cred) throw new Error("passkey creation was cancelled");
  const r = cred.response as AuthenticatorAttestationResponse;
  return post<AuthResult>(`/${app}/register/finish`, {
    ceremonyId: begin.ceremonyId,
    credential: {
      id: cred.id,
      rawId: toBase64url(cred.rawId),
      type: cred.type,
      response: {
        clientDataJSON: toBase64url(r.clientDataJSON),
        attestationObject: toBase64url(r.attestationObject),
        transports: r.getTransports?.() ?? [],
      },
    },
  });
}

export async function login(app: string, username: string, prfSalt?: Uint8Array<ArrayBuffer>): Promise<AuthResult> {
  const begin = await post<Begin<RequestJSON>>(`/${app}/login/begin`, { username });
  const o = begin.publicKey;
  const extensions: AuthenticationExtensionsClientInputs = prfSalt ? { prf: { eval: { first: prfSalt } } } : {};
  const cred = (await navigator.credentials.get({
    publicKey: {
      challenge: fromBase64url(o.challenge),
      rpId: o.rpId,
      timeout: o.timeout,
      allowCredentials: toDescriptors(o.allowCredentials),
      userVerification: o.userVerification,
      extensions,
    },
  })) as PublicKeyCredential | null;
  if (!cred) throw new Error("passkey login was cancelled");
  const r = cred.response as AuthenticatorAssertionResponse;
  const result = await post<AuthResult>(`/${app}/login/finish`, {
    ceremonyId: begin.ceremonyId,
    credential: {
      id: cred.id,
      rawId: toBase64url(cred.rawId),
      type: cred.type,
      response: {
        clientDataJSON: toBase64url(r.clientDataJSON),
        authenticatorData: toBase64url(r.authenticatorData),
        signature: toBase64url(r.signature),
        userHandle: r.userHandle ? toBase64url(r.userHandle) : undefined,
      },
    },
  });
  const prf = cred.getClientExtensionResults().prf?.results?.first;
  if (prf) result.prf = prf instanceof ArrayBuffer ? prf : new Uint8Array(new Uint8Array(prf.buffer, prf.byteOffset, prf.byteLength)).buffer;
  return result;
}

export async function logout(token: string): Promise<void> {
  await post(`/logout`, {}, token);
}

export async function passkeys(app: string, token: string): Promise<PasskeyInfo[]> {
  const res = await fetch(`${BASE}/${app}/passkeys`, { headers: { Authorization: `Bearer ${token}` } });
  if (!res.ok) throw new Error("could not load passkeys");
  return res.json();
}

export function algName(alg: number): string {
  return ({ [-7]: "ES256 (P-256)", [-8]: "EdDSA (Ed25519)", [-257]: "RS256 (RSA)" } as Record<number, string>)[alg] ?? `alg ${alg}`;
}
