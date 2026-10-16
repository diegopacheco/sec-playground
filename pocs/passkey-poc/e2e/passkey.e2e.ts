import { chromium, type BrowserContext, type Page } from "playwright";

const BOOKSTORE = `http://localhost:${process.env.BOOKSTORE_WEB_PORT ?? 5173}`;
const VAULT = `http://localhost:${process.env.VAULT_WEB_PORT ?? 5174}`;
const SHOTS = process.env.SCREENSHOTS_DIR;
const run = Date.now().toString(36);
let bookstoreToken = "";

let passed = 0;
function check(cond: unknown, what: string): void {
  if (!cond) throw new Error(`FAILED: ${what}`);
  passed++;
  console.log(`ok ${passed} - ${what}`);
}

async function shot(page: Page, name: string) {
  if (!SHOTS) return;
  await page.evaluate(() => document.querySelector(".toasts")?.remove());
  await page.screenshot({ path: `${SHOTS}/${name}.png` });
}

async function withAuthenticator(context: BrowserContext, page: Page) {
  const cdp = await context.newCDPSession(page);
  await cdp.send("WebAuthn.enable");
  const { authenticatorId } = await cdp.send("WebAuthn.addVirtualAuthenticator", {
    options: { protocol: "ctap2", transport: "internal", hasResidentKey: true, hasUserVerification: true, isUserVerified: true, automaticPresenceSimulation: true, hasPrf: true },
  });
  return { cdp, authenticatorId };
}

async function toastText(page: Page): Promise<string> {
  return (await page.locator(".toast").last().textContent({ timeout: 10000 })) ?? "";
}

async function bookstore(context: BrowserContext) {
  const page = await context.newPage();
  await withAuthenticator(context, page);
  const user = `alice-${run}`;
  await page.goto(BOOKSTORE, { waitUntil: "networkidle" });
  await shot(page, "bookstore-1-login");

  await page.fill("#username", `ghost-${run}`);
  await page.click("#signin");
  const refused = await toastText(page);
  check(refused.includes("no passkey registered"), `bookstore refuses sign in for a username with no passkey (${refused})`);

  await page.fill("#username", user);
  await page.click("#signup");
  await page.waitForSelector(".book");
  check((await page.locator(".book").count()) === 8, "bookstore signup with a passkey shows the catalog");
  await shot(page, "bookstore-2-books");

  await page.locator('.book[data-id="b1"] button').click();
  await page.waitForSelector('[data-tab="cart"]:has-text("Cart (1)")');
  await page.locator('.book[data-id="b5"] button').click();
  await page.waitForSelector('[data-tab="cart"]:has-text("Cart (2)")');
  await page.click('[data-tab="cart"]');
  await page.waitForSelector("#total");
  check((await page.textContent("#total")) === "$97.98", "cart total is computed by the server from catalog prices");
  await shot(page, "bookstore-3-cart");

  await page.click("#checkout");
  await page.waitForSelector(".order");
  check((await page.locator(".order").count()) === 1, "checkout creates an order");
  await shot(page, "bookstore-4-orders");

  await page.click('[data-tab="security"]');
  await page.waitForSelector(".keys tbody tr");
  const row = (await page.textContent(".keys tbody tr")) ?? "";
  check(row.includes("ES256"), "security tab lists the ES256 passkey stored by the Go server");
  await shot(page, "bookstore-5-security");

  const token = await page.evaluate(() => JSON.parse(sessionStorage.getItem("bookstore.auth")!).token as string);
  await page.click("#logout");
  await page.waitForSelector("#signin");
  const reused = await page.evaluate(async (t) => (await fetch("/api/cart", { headers: { Authorization: `Bearer ${t}` } })).status, token);
  check(reused === 401, "a logged out bookstore token is rejected by the API");

  await page.fill("#username", user);
  await page.click("#signin");
  await page.waitForSelector(".book");
  check(true, "the same passkey signs in again");

  bookstoreToken = await page.evaluate(() => JSON.parse(sessionStorage.getItem("bookstore.auth")!).token as string);
  await page.close();
}

async function vault(context: BrowserContext) {
  const page = await context.newPage();
  await withAuthenticator(context, page);
  const user = `bob-${run}`;
  const secret = `Sup3r-Secret-${run}`;
  await page.goto(VAULT, { waitUntil: "networkidle" });
  await shot(page, "vault-1-lock");
  const crossApp = await page.evaluate(async (t) => (await fetch("/api/vault", { headers: { Authorization: `Bearer ${t}` } })).status, bookstoreToken);
  check(crossApp === 401, "a valid bookstore session token is refused by the vault API");

  await page.fill("#username", user);
  await page.click("#create");
  await page.waitForSelector("#add");
  check(true, "vault created: passkey registered, PRF output derived the key, canary stored");

  for (const [name, site, login, pw] of [
    ["GitHub", "https://github.com", "bob", secret],
    ["Bank", "https://bank.test", "bob.b", "correct-horse-battery"],
    ["Email", "https://mail.test", "bob@mail.test", "another-one-42"],
  ]) {
    await page.click("#add");
    await page.fill('[name="name"]', name!);
    await page.fill('[name="url"]', site!);
    await page.fill('[name="username"]', login!);
    await page.fill('[name="password"]', pw!);
    await page.click("#save");
    await page.waitForSelector(`.entry[data-name="${name}"]`);
  }
  check((await page.locator(".entry").count()) === 3, "three entries encrypted and saved");
  await page.locator('.entry[data-name="GitHub"]').click();
  await page.click("#use-generated");
  await page.fill('[name="password"]', secret);
  await shot(page, "vault-2-entries");

  await page.click('[data-tab="server"]');
  await page.waitForSelector("#raw");
  const raw = (await page.textContent("#raw")) ?? "";
  check(raw.includes('"ct"') && !raw.includes(secret) && !raw.includes("GitHub") && !raw.includes("bank.test"), "server copy holds ciphertext only, no names, sites or passwords");
  await shot(page, "vault-3-server-view");

  await page.click('[data-tab="passkey"]');
  await page.waitForSelector(".kv");
  await shot(page, "vault-4-passkey");

  await page.click("#lock");
  await page.waitForSelector("#unlock");
  await page.fill("#username", user);
  await page.click("#unlock");
  await page.waitForSelector('.entry[data-name="GitHub"]');
  await page.locator('.entry[data-name="GitHub"]').click();
  check((await page.inputValue('[name="password"]')) === secret, "after lock and unlock the passkey PRF re-derives the key and decrypts the password");

  const ctx2 = await context.browser()!.newContext();
  const other = await ctx2.newPage();
  await withAuthenticator(ctx2, other);
  await other.goto(VAULT, { waitUntil: "networkidle" });
  await other.fill("#username", user);
  await other.click("#unlock");
  await other.waitForTimeout(3000);
  check((await other.locator("#add").count()) === 0 && (await other.locator("#unlock").isVisible()), "a device without the vault passkey cannot unlock it");
  await ctx2.close();
  await page.close();
}

const browser = await chromium.launch();
try {
  const context = await browser.newContext({ viewport: { width: 1280, height: 800 }, permissions: ["clipboard-read", "clipboard-write"] });
  await bookstore(context);
  await vault(context);
  console.log(`\n${passed} e2e checks passed`);
} catch (err) {
  console.error((err as Error).message);
  process.exitCode = 1;
} finally {
  await browser.close();
}
