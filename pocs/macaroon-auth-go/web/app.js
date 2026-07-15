let token = ""
let tokenData = null

const byId = id => document.getElementById(id)

const post = async (path, body) => {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  })
  const data = await response.json()
  if (!response.ok) throw new Error(data.error || "Request failed")
  return data
}

const renderToken = () => {
  byId("empty-chain").hidden = true
  byId("chain").hidden = false
  byId("attenuate-form").hidden = false
  byId("token-box").hidden = false
  byId("chain-count").textContent = `${tokenData.caveats.length} caveats`
  byId("token-value").textContent = token
  byId("chain").innerHTML = tokenData.caveats.map((caveat, index) => `
    <div class="caveat" style="animation-delay:${index * 45}ms">
      <span class="caveat-index">${String(index + 1).padStart(2, "0")}</span>
      <span class="caveat-type">${escapeHTML(caveat.type)}</span>
      <span class="caveat-value" title="${escapeHTML(caveat.value)}">${escapeHTML(caveat.value)}</span>
      <span class="caveat-sign">HMAC ${tokenData.signature.slice(index * 4, index * 4 + 8)}</span>
    </div>
  `).join("")
}

const renderVerdict = result => {
  const verdict = byId("verdict")
  verdict.className = `verdict ${result.allowed ? "allowed" : "denied"}`
  verdict.querySelector(".verdict-word").textContent = result.allowed ? "ALLOWED" : "DENIED"
  verdict.querySelector("p").textContent = result.reason
  byId("checks").innerHTML = result.checks.map(check => `
    <div class="check ${check.passed ? "pass" : "fail"}">
      <span class="check-dot"></span>
      <b>${escapeHTML(check.type)}</b>
      <span>${escapeHTML(check.actual)} → ${escapeHTML(check.required)}</span>
    </div>
  `).join("")
}

const resetVerdict = message => {
  const verdict = byId("verdict")
  verdict.className = "verdict waiting"
  verdict.querySelector(".verdict-word").textContent = "WAITING"
  verdict.querySelector("p").textContent = message
  byId("checks").replaceChildren()
}

const escapeHTML = value => String(value).replace(/[&<>'"]/g, character => ({
  "&": "&amp;", "<": "&lt;", ">": "&gt;", "'": "&#39;", "\"": "&quot;"
})[character])

byId("mint-form").addEventListener("submit", async event => {
  event.preventDefault()
  try {
    tokenData = await post("/api/macaroon", {
      resource: byId("mint-resource").value,
      operation: byId("mint-operation").value,
      location: byId("mint-location").value,
      expires_in_seconds: Number(byId("mint-expiration").value)
    })
    token = tokenData.token
    renderToken()
    resetVerdict("New token ready. Verify access to get a current decision.")
    byId("verdict").scrollIntoView({ behavior: "smooth", block: "center" })
  } catch (error) {
    window.alert(error.message)
  }
})

byId("attenuate-form").addEventListener("submit", async event => {
  event.preventDefault()
  try {
    tokenData = await post("/api/attenuate", {
      token,
      caveat: { type: byId("caveat-type").value, value: byId("caveat-value").value }
    })
    token = tokenData.token
    renderToken()
    resetVerdict("Restriction added. Verify access again against the updated token.")
  } catch (error) {
    window.alert(error.message)
  }
})

byId("verify-form").addEventListener("submit", async event => {
  event.preventDefault()
  if (!token) {
    window.alert("Mint a macaroon first")
    return
  }
  try {
    const result = await post("/api/verify", {
      token,
      resource: byId("verify-resource").value,
      operation: byId("verify-operation").value,
      location: byId("verify-location").value
    })
    renderVerdict(result)
  } catch (error) {
    window.alert(error.message)
  }
})

byId("caveat-type").addEventListener("change", event => {
  const values = {
    resource: "/records/payroll",
    operation: "read",
    location: "us-west",
    expires: new Date(Date.now() + 5 * 60 * 1000).toISOString().replace(".000Z", "Z")
  }
  byId("caveat-value").value = values[event.target.value]
})

byId("copy-token").addEventListener("click", async () => {
  await navigator.clipboard.writeText(token)
  byId("copy-token").textContent = "Copied"
  window.setTimeout(() => { byId("copy-token").textContent = "Copy" }, 1200)
})
