const statusText = document.querySelector("#status")
const signal = document.querySelector("#signal-light")
const registerForm = document.querySelector("#register-form")
const loginForm = document.querySelector("#login-form")

function decode(value) {
  const base64 = value.replace(/-/g, "+").replace(/_/g, "/")
  const padded = base64.padEnd(Math.ceil(base64.length / 4) * 4, "=")
  const binary = atob(padded)
  return Uint8Array.from(binary, character => character.charCodeAt(0))
}

function encode(value) {
  const bytes = new Uint8Array(value)
  let binary = ""
  bytes.forEach(byte => binary += String.fromCharCode(byte))
  return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "")
}

function creationOptions(options) {
  const publicKey = options.publicKey
  publicKey.challenge = decode(publicKey.challenge)
  publicKey.user.id = decode(publicKey.user.id)
  publicKey.excludeCredentials = (publicKey.excludeCredentials || []).map(item => ({...item, id: decode(item.id)}))
  return {publicKey}
}

function requestOptions(options) {
  const publicKey = options.publicKey
  publicKey.challenge = decode(publicKey.challenge)
  publicKey.allowCredentials = (publicKey.allowCredentials || []).map(item => ({...item, id: decode(item.id)}))
  return {publicKey}
}

function registrationCredential(credential) {
  return {
    id: credential.id,
    rawId: encode(credential.rawId),
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment,
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      attestationObject: encode(credential.response.attestationObject),
      clientDataJSON: encode(credential.response.clientDataJSON),
      transports: credential.response.getTransports ? credential.response.getTransports() : []
    }
  }
}

function loginCredential(credential) {
  return {
    id: credential.id,
    rawId: encode(credential.rawId),
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment,
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      authenticatorData: encode(credential.response.authenticatorData),
      clientDataJSON: encode(credential.response.clientDataJSON),
      signature: encode(credential.response.signature),
      userHandle: credential.response.userHandle ? encode(credential.response.userHandle) : null
    }
  }
}

async function call(path, body) {
  const response = await fetch(path, {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify(body)
  })
  const data = await response.json()
  if (!response.ok) throw new Error(data.error || "request failed")
  return data
}

function setStatus(message, state = "ready") {
  statusText.textContent = message
  signal.className = state === "ready" ? "" : state
}

function lock(form, locked) {
  form.querySelectorAll("button, input").forEach(element => element.disabled = locked)
}

async function refresh() {
  const response = await fetch("/api/status")
  const data = await response.json()
  document.querySelector("#user-count").textContent = data.users
  document.querySelector("#credential-count").textContent = data.credentials
}

registerForm.addEventListener("submit", async event => {
  event.preventDefault()
  lock(registerForm, true)
  try {
    setStatus("Waiting for your authenticator", "busy")
    const started = await call("/api/register/begin", {
      username: document.querySelector("#register-username").value,
      displayName: document.querySelector("#display-name").value
    })
    const credential = await navigator.credentials.create(creationOptions(started.options))
    await call(`/api/register/finish?session=${encodeURIComponent(started.sessionId)}`, registrationCredential(credential))
    setStatus("Passkey registered", "ready")
    document.querySelector("#login-username").value = document.querySelector("#register-username").value
    await refresh()
  } catch (error) {
    setStatus(error.message, "error")
  } finally {
    lock(registerForm, false)
  }
})

loginForm.addEventListener("submit", async event => {
  event.preventDefault()
  lock(loginForm, true)
  try {
    setStatus("Verify on this device", "busy")
    const started = await call("/api/login/begin", {
      username: document.querySelector("#login-username").value
    })
    const credential = await navigator.credentials.get(requestOptions(started.options))
    const finished = await call(`/api/login/finish?session=${encodeURIComponent(started.sessionId)}`, loginCredential(credential))
    setStatus(`${finished.username} authenticated`, "ready")
  } catch (error) {
    setStatus(error.message, "error")
  } finally {
    lock(loginForm, false)
  }
})

if (!window.PublicKeyCredential) {
  setStatus("This browser does not support WebAuthn", "error")
  lock(registerForm, true)
  lock(loginForm, true)
} else {
  refresh().catch(error => setStatus(error.message, "error"))
}
