const state = {
  tab: "approvals",
  token: "",
  intents: [],
  credentials: [],
  audit: { verified: true, entries: [] },
  scenarios: [],
  selected: "",
  transcript: [],
  running: false,
};

function h(tag, attrs, ...children) {
  const el = document.createElement(tag);
  for (const [k, v] of Object.entries(attrs || {})) {
    if (v === false || v === null || v === undefined) continue;
    if (k === "class") el.className = v;
    else if (k === "style") el.style.cssText = v;
    else if (k.startsWith("on")) el.addEventListener(k.slice(2), v);
    else el.setAttribute(k, v === true ? "" : v);
  }
  for (const c of children.flat()) {
    if (c === null || c === undefined || c === false) continue;
    el.append(c instanceof Node ? c : document.createTextNode(String(c)));
  }
  return el;
}

function badge(text, kind) {
  return h("span", { class: "badge b-" + (kind || text) }, text);
}

function short(hash) {
  return hash ? hash.slice(0, 10) : "";
}

function ago(iso) {
  if (!iso) return "";
  const s = Math.round((Date.now() - new Date(iso).getTime()) / 1000);
  if (s < 60) return s + "s ago";
  if (s < 3600) return Math.round(s / 60) + "m ago";
  return Math.round(s / 3600) + "h ago";
}

function remaining(iso) {
  const s = Math.round((new Date(iso).getTime() - Date.now()) / 1000);
  if (s <= 0) return "expired";
  return Math.floor(s / 60) + "m " + String(s % 60).padStart(2, "0") + "s left";
}

async function api(method, path) {
  const res = await fetch(path, { method, headers: { "X-Admin-Token": state.token } });
  if (res.status === 401) {
    lock(true);
    throw new Error("unauthorized");
  }
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || res.statusText);
  return data;
}

function lock(on) {
  document.getElementById("lock").hidden = !on;
  document.querySelectorAll(".panel").forEach((p) => (p.hidden = on || p.id !== "panel-" + state.tab));
}

function readToken() {
  const m = location.hash.match(/token=([a-f0-9]+)/);
  if (m) {
    sessionStorage.setItem("broker-admin-token", m[1]);
    history.replaceState(null, "", location.pathname + location.search);
  }
  state.token = sessionStorage.getItem("broker-admin-token") || "";
  const tab = new URLSearchParams(location.search).get("tab");
  if (tab) state.tab = tab;
}

async function refresh() {
  if (!state.token) {
    lock(true);
    return;
  }
  try {
    const [intents, credentials, auditData] = await Promise.all([
      api("GET", "/api/intents"),
      api("GET", "/api/credentials"),
      api("GET", "/api/audit"),
    ]);
    if (!state.scenarios.length) state.scenarios = await api("GET", "/api/lab");
    state.intents = intents;
    state.credentials = credentials;
    state.audit = auditData;
    lock(false);
    render();
  } catch (e) {
    if (e.message !== "unauthorized") console.error(e);
  }
}

function renderStats() {
  const count = (s) => state.intents.filter((i) => i.status === s).length;
  const denied = state.audit.entries.filter((e) => e.outcome === "denied").length;
  const items = [
    ["Pending", count("pending")],
    ["Tokens used", count("used")],
    ["Blocked", denied],
    ["Audit events", state.audit.entries.length],
  ];
  document.getElementById("stats").replaceChildren(...items.map(([l, v]) => h("div", { class: "stat" }, h("b", {}, v), h("span", {}, l))));
  const pc = document.getElementById("pending-count");
  pc.textContent = count("pending");
  pc.className = "count" + (count("pending") ? " hot" : "");
}

async function decide(id, action) {
  try {
    await api("POST", "/api/intents/" + id + "/" + action);
  } catch (e) {
    alert(e.message);
  }
  refresh();
}

function intentCard(i) {
  return h("div", { class: "card intent" },
    h("div", { class: "spread" },
      h("div", { class: "row" }, h("span", { class: "method " + i.method }, i.method), h("h3", {}, i.credential)),
      h("span", { class: "muted" }, remaining(i.expires_at))),
    h("div", { class: "url", style: "margin-top:8px" }, i.url),
    h("div", { class: "reason" }, "“" + (i.reason || "no reason given") + "”"),
    h("dl", { class: "kv" },
      h("dt", {}, "Intent"), h("dd", { class: "mono" }, i.id),
      h("dt", {}, "Body sha256"), h("dd", { class: "mono" }, short(i.body_sha256) + "…"),
      h("dt", {}, "Requested"), h("dd", {}, ago(i.created_at))),
    i.body ? h("pre", { class: "body" }, i.body) : null,
    h("div", { class: "row", style: "margin-top:14px" },
      h("button", { class: "btn primary", onclick: () => decide(i.id, "approve") }, "Approve this exact request"),
      h("button", { class: "btn danger", onclick: () => decide(i.id, "deny") }, "Deny")));
}

function renderApprovals() {
  const pending = state.intents.filter((i) => i.status === "pending");
  const history = state.intents.filter((i) => i.status !== "pending");
  const panel = document.getElementById("panel-approvals");
  panel.replaceChildren(
    h("div", { class: "section-title" }, h("h2", {}, "Waiting for your approval"), h("span", { class: "muted" }, "Approving mints one token bound to method, URL and body hash")),
    pending.length ? h("div", { class: "grid" }, pending.map(intentCard)) : h("div", { class: "empty" }, h("b", {}, "No pending requests"), "When the agent calls request_access for a human approved credential it shows up here."),
    h("div", { class: "section-title", style: "margin-top:28px" }, h("h2", {}, "Intent history"), h("span", { class: "muted" }, history.length + " intent(s)")),
    history.length
      ? h("div", { class: "card table-card" }, h("table", {},
          h("thead", {}, h("tr", {}, ["Status", "Credential", "Request", "Reason", "Approval", "Created"].map((t) => h("th", {}, t)))),
          h("tbody", {}, history.map((i) => h("tr", {},
            h("td", {}, badge(i.status)),
            h("td", { class: "mono" }, i.credential),
            h("td", {}, h("span", { class: "method " + i.method }, i.method), " ", h("span", { class: "url" }, i.url)),
            h("td", {}, i.reason),
            h("td", {}, badge(i.approval)),
            h("td", { class: "muted" }, ago(i.created_at)))))))
      : h("div", { class: "empty" }, h("b", {}, "Nothing yet"), "Run a scenario in the Attack lab tab or connect an agent over MCP."));
}

function renderCredentials() {
  const panel = document.getElementById("panel-credentials");
  panel.replaceChildren(
    h("div", { class: "section-title" }, h("h2", {}, "Credentials in policy"), h("span", { class: "muted" }, "Values stay in Vault or Keychain and are fetched only at call time")),
    h("div", { class: "grid" }, state.credentials.map((c) => h("div", { class: "card cred" },
      h("div", { class: "spread" }, h("h3", {}, c.name), h("div", { class: "row" }, badge(c.provider), badge(c.approval === "auto" ? "auto approve" : "human approval", c.approval))),
      h("div", { class: "muted" }, c.description),
      h("div", { class: "secret-line" }, c.inject_header + ": ••••••••  injected by broker"),
      h("div", { class: "label" }, "Allowed egress"),
      c.allow.map((r) => h("div", { class: "rule-line" },
        h("span", { class: "row" }, r.methods.map((m) => h("span", { class: "method " + m }, m))),
        h("span", { class: "rule" }, r.scheme + "://" + r.host + r.path_prefix + "*")))))));
}

function renderAudit() {
  const a = state.audit;
  const entries = [...a.entries].reverse();
  const panel = document.getElementById("panel-audit");
  panel.replaceChildren(
    a.verified
      ? h("div", { class: "banner ok" }, "Hash chain verified", h("small", {}, a.entries.length + " entries, each hash covers the previous one"))
      : h("div", { class: "banner bad" }, "Hash chain broken at entry " + a.broken_at, h("small", {}, "the audit file was modified")),
    entries.length
      ? h("div", { class: "card table-card" }, h("table", {},
          h("thead", {}, h("tr", {}, ["#", "When", "Event", "Outcome", "Credential", "Request", "Detail", "Hash / prev"].map((t) => h("th", {}, t)))),
          h("tbody", {}, entries.map((e) => h("tr", {},
            h("td", { class: "mono" }, e.seq),
            h("td", { class: "muted" }, ago(e.time)),
            h("td", { class: "mono" }, e.event),
            h("td", {}, badge(e.outcome)),
            h("td", { class: "mono" }, e.credential || ""),
            h("td", {}, e.method ? h("span", { class: "method " + e.method }, e.method) : "", " ", h("span", { class: "url" }, e.url || "")),
            h("td", {}, e.detail || ""),
            h("td", {}, h("div", { class: "hash" }, short(e.hash)), h("div", { class: "hash" }, "↳ " + short(e.prev))))))))
      : h("div", { class: "empty" }, h("b", {}, "Audit log is empty"), "Every request, decision and call is appended here."));
}

async function runScenario(id) {
  state.selected = id;
  state.running = true;
  state.transcript = [];
  renderLab();
  try {
    state.transcript = await api("POST", "/api/lab/" + id);
  } catch (e) {
    state.transcript = [{ actor: "broker", action: "run " + id, outcome: "error", detail: e.message }];
  }
  state.running = false;
  refresh();
}

function renderLab() {
  const panel = document.getElementById("panel-lab");
  const current = state.scenarios.find((s) => s.id === state.selected);
  const blocked = state.transcript.filter((s) => s.outcome === "blocked").length;
  panel.replaceChildren(h("div", { class: "lab" },
    h("div", { class: "stack" },
      h("div", { class: "section-title" }, h("h2", {}, "Scenarios"), h("span", { class: "muted" }, "run against the live broker")),
      state.scenarios.map((s) => h("div", { class: "card scenario" + (s.id === state.selected ? " selected" : ""), onclick: () => runScenario(s.id) },
        h("h3", {}, s.title),
        h("p", {}, h("span", { class: "label" }, "Threat"), s.threat),
        h("p", {}, h("span", { class: "label" }, "Control"), s.control)))),
    h("div", { class: "card" },
      current
        ? [
            h("div", { class: "spread" }, h("h2", { style: "margin:0;font-size:16px" }, current.title), h("button", { class: "btn small", onclick: () => runScenario(current.id), disabled: state.running }, state.running ? "Running…" : "Run again")),
            h("p", { class: "muted", style: "margin:4px 0 10px" }, current.control),
            h("ol", { class: "timeline" }, state.transcript.map((s) => h("li", { class: "step" },
              h("div", { class: "actor " + s.actor }, s.actor),
              h("div", {},
                h("div", { class: "spread" }, h("span", { class: "step-action" }, s.action), badge(s.outcome)),
                h("div", { class: "step-detail" }, s.detail))))),
            state.transcript.length ? h("div", { class: "verdict ok" }, blocked ? blocked + " unsafe step(s) stopped by the broker, credential never reached the agent" : "Request completed, the agent only saw the redacted response") : null,
          ]
        : h("div", { class: "empty" }, h("b", {}, "Pick a scenario"), "Each one drives the real broker code path: request_access, human approval, call_api."))));
}

function render() {
  renderStats();
  if (state.tab === "approvals") renderApprovals();
  if (state.tab === "credentials") renderCredentials();
  if (state.tab === "audit") renderAudit();
  if (state.tab === "lab") renderLab();
  document.querySelectorAll(".tab").forEach((t) => t.classList.toggle("active", t.dataset.tab === state.tab));
  document.querySelectorAll(".panel").forEach((p) => (p.hidden = p.id !== "panel-" + state.tab));
}

document.getElementById("tabs").addEventListener("click", (e) => {
  const tab = e.target.closest(".tab");
  if (!tab) return;
  state.tab = tab.dataset.tab;
  render();
});

document.getElementById("lock-form").addEventListener("submit", (e) => {
  e.preventDefault();
  state.token = document.getElementById("lock-input").value.trim();
  sessionStorage.setItem("broker-admin-token", state.token);
  refresh();
});

readToken();
refresh();
setInterval(() => {
  if (!state.running) refresh();
}, 2000);
