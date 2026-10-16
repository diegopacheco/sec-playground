type Child = Node | string | number | null | undefined | false;
type Attrs = Record<string, string | number | boolean | EventListener | undefined>;

export function h<K extends keyof HTMLElementTagNameMap>(tag: K, attrs: Attrs = {}, ...children: (Child | Child[])[]): HTMLElementTagNameMap[K] {
  const el = document.createElement(tag);
  for (const [k, v] of Object.entries(attrs)) {
    if (v === undefined || v === false) continue;
    if (k.startsWith("on") && typeof v === "function") el.addEventListener(k.slice(2), v);
    else if (v === true) el.setAttribute(k, "");
    else el.setAttribute(k, String(v));
  }
  for (const c of children.flat()) {
    if (c === null || c === undefined || c === false) continue;
    el.append(c instanceof Node ? c : document.createTextNode(String(c)));
  }
  return el;
}

export function mount(root: HTMLElement, ...nodes: Node[]): void {
  root.replaceChildren(...nodes);
}

export function toast(message: string, kind: "ok" | "error" = "ok"): void {
  let stack = document.querySelector<HTMLElement>(".toasts");
  if (!stack) {
    stack = h("div", { class: "toasts" });
    document.body.append(stack);
  }
  const el = h("div", { class: `toast toast-${kind}`, role: "status" }, message);
  stack.append(el);
  setTimeout(() => el.remove(), 3500);
}
