export type Session = { app: string; username: string; userId: string; expiresAt: string };

export type Introspect = (token: string) => Promise<Session | null>;

export function remoteIntrospect(passkeyServer: string): Introspect {
  return async (token) => {
    const res = await fetch(`${passkeyServer}/api/session`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!res.ok) return null;
    return (await res.json()) as Session;
  };
}

export function bearer(req: Request): string | null {
  const h = req.headers.get("Authorization") ?? "";
  return h.startsWith("Bearer ") ? h.slice(7) : null;
}

export async function requireSession(req: Request, app: string, introspect: Introspect): Promise<Session | null> {
  const token = bearer(req);
  if (!token) return null;
  const session = await introspect(token);
  if (!session || session.app !== app) return null;
  return session;
}

export function json(data: unknown, status = 200): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: { "Content-Type": "application/json", "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" },
  });
}

export function unauthorized(): Response {
  return json({ error: "login with your passkey first" }, 401);
}

export async function readJSON<T>(req: Request, maxBytes = 64 * 1024): Promise<T | null> {
  const text = await req.text();
  if (text.length > maxBytes) return null;
  try {
    return JSON.parse(text) as T;
  } catch {
    return null;
  }
}
