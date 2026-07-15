import { useEffect, useState } from "react";

type Session = {
  authenticated: boolean;
  sdk: string;
  claims: Record<string, unknown>;
};

const claimOrder = ["name", "email", "nickname", "sub", "iss", "aud", "iat", "exp", "updated_at"];

function sortClaims(claims: Record<string, unknown>): [string, unknown][] {
  const entries = Object.entries(claims);
  return entries.sort(([left], [right]) => {
    const leftRank = claimOrder.indexOf(left);
    const rightRank = claimOrder.indexOf(right);
    if (leftRank !== -1 && rightRank !== -1) return leftRank - rightRank;
    if (leftRank !== -1) return -1;
    if (rightRank !== -1) return 1;
    return left.localeCompare(right);
  });
}

function render(value: unknown): string {
  if (typeof value === "string") return value;
  return JSON.stringify(value);
}

export function App() {
  const [session, setSession] = useState<Session | null>(null);

  useEffect(() => {
    fetch("/api/me")
      .then((response) => response.json() as Promise<Session>)
      .then(setSession)
      .catch(() => setSession({ authenticated: false, sdk: "unknown", claims: {} }));
  }, []);

  if (session === null) {
    return (
      <main className="card">
        <p className="muted">Loading</p>
      </main>
    );
  }

  if (!session.authenticated) {
    return (
      <main className="card">
        <div className="badge">{session.sdk}</div>
        <h1 className="sad">Opps you are not authenticated</h1>
        <p className="muted">You need to sign in with Auth0 to see the protected page.</p>
        <a className="button" href="/login">
          Authenticate with Auth0
        </a>
      </main>
    );
  }

  const claims = sortClaims(session.claims);
  const picture = typeof session.claims.picture === "string" ? session.claims.picture : null;
  const name = typeof session.claims.name === "string" ? session.claims.name : "there";

  return (
    <main className="card wide">
      <div className="badge">{session.sdk}</div>
      <h1 className="happy">Horay you are logged</h1>
      <p className="muted">Welcome {name}, Auth0 authenticated you and issued the claims below.</p>
      {picture && <img className="avatar" src={picture} alt={name} />}
      <table className="claims">
        <thead>
          <tr>
            <th>Claim</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {claims.map(([key, value]) => (
            <tr key={key}>
              <td className="key">{key}</td>
              <td className="value">{render(value)}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <a className="button ghost" href="/logout">
        Sign out
      </a>
    </main>
  );
}
