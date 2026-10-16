import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { Policy } from "../types";
import { LabelBadge } from "./Badge";

export function PolicyTab() {
  const [policy, setPolicy] = useState<Policy | null>(null);

  useEffect(() => {
    api.policy().then(setPolicy).catch(() => setPolicy(null));
  }, []);

  if (!policy) return <p className="muted">Loading policy</p>;

  return (
    <div className="policy">
      <h2>Sources: which tool output gets which label</h2>
      <table>
        <thead>
          <tr>
            <th>tool</th>
            <th>condition</th>
            <th>labels</th>
          </tr>
        </thead>
        <tbody>
          {policy.sources.map((s, i) => (
            <tr key={i}>
              <td>
                <code>{s.tool}</code>
              </td>
              <td>{s.pattern ? <code>{`${s.arg} ~ ${s.pattern}`}</code> : <span className="muted">always</span>}</td>
              <td>{s.labels.map((l) => <LabelBadge key={l} label={l} />)}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <p className="muted">Any value matching a secret pattern (AWS keys, GitHub tokens, private keys, KEY=VALUE lines) is labeled secret no matter which tool returned it.</p>
      <h2>Sinks: which argument refuses which label</h2>
      <table>
        <thead>
          <tr>
            <th>tool</th>
            <th>argument</th>
            <th>forbids</th>
            <th>why</th>
          </tr>
        </thead>
        <tbody>
          {policy.sinks.map((s, i) => (
            <tr key={i}>
              <td>
                <code>{s.tool}</code>
              </td>
              <td>
                <code>{s.arg}</code>
              </td>
              <td>{s.forbid.map((l) => <LabelBadge key={l} label={l} />)}</td>
              <td>{s.why}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
