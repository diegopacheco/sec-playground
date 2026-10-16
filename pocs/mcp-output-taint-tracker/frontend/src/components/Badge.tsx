import type { Label } from "../types";

export function LabelBadge({ label }: { label: Label }) {
  return <span className={`badge badge-${label}`}>{label}</span>;
}

export function VerdictBadge({ allow }: { allow: boolean }) {
  return <span className={`badge ${allow ? "badge-allow" : "badge-block"}`}>{allow ? "allowed" : "blocked"}</span>;
}
