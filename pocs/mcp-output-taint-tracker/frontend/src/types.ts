export type Label = "untrusted" | "secret";

export type EventKind = "prompt" | "result" | "call";

export interface Entity {
  kind: string;
  value: string;
  labels: Label[];
}

export interface Flow {
  arg: string;
  label: Label;
  sourceId: number;
  sourceTool: string;
  kind: "entity" | "text";
  entityKind?: string;
  evidence: string;
  encoding: string;
  decoded?: string;
}

export interface Decision {
  allow: boolean;
  checked: string[];
  flows: Flow[];
}

export interface TaintEvent {
  seq: number;
  kind: EventKind;
  tool?: string;
  args?: Record<string, unknown>;
  text?: string;
  labels?: Label[];
  entities?: Entity[];
  decision?: Decision;
  at: string;
}

export interface Session {
  id: string;
  name: string;
  created: string;
  events: TaintEvent[];
}

export interface Summary {
  id: string;
  name: string;
  created: string;
  events: number;
  blocked: number;
}

export interface SourceRule {
  tool: string;
  arg?: string;
  pattern?: string;
  labels: Label[];
}

export interface SinkRule {
  tool: string;
  arg: string;
  forbid: Label[];
  why: string;
}

export interface Policy {
  sources: SourceRule[];
  sinks: SinkRule[];
}
