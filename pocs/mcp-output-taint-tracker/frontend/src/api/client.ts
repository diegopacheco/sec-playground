import type { Policy, Session, Summary, TaintEvent } from "../types";

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(path, {
    method,
    headers: body === undefined ? undefined : { "Content-Type": "application/json" },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `${method} ${path} failed with ${res.status}`);
  }
  return res.status === 204 ? (undefined as T) : ((await res.json()) as T);
}

const sessionPath = (id: string) => `/api/sessions/${encodeURIComponent(id)}`;

export const api = {
  policy: () => request<Policy>("GET", "/api/policy"),
  sessions: () => request<Summary[]>("GET", "/api/sessions"),
  session: (id: string) => request<Session>("GET", sessionPath(id)),
  create: (name: string) => request<Session>("POST", "/api/sessions", { name }),
  remove: (id: string) => request<void>("DELETE", sessionPath(id)),
  loadScenarios: () => request<Summary[]>("POST", "/api/scenarios"),
  prompt: (id: string, text: string) => request<TaintEvent>("POST", `${sessionPath(id)}/prompt`, { text }),
  result: (id: string, tool: string, args: Record<string, unknown>, output: string) =>
    request<TaintEvent>("POST", `${sessionPath(id)}/result`, { tool, args, output }),
  call: (id: string, tool: string, args: Record<string, unknown>) =>
    request<TaintEvent>("POST", `${sessionPath(id)}/call`, { tool, args }),
};
