import { ref } from 'vue'

/**
 * useGoogleAuth — wraps Google Identity Services (GIS) for real Google Sign-In.
 *
 * Loads the GIS script on demand, then uses the OAuth token-client + ID token
 * flow so the backend receives a genuine, verifiable `id_token`. No mock tokens.
 *
 * Requires VITE_GOOGLE_CLIENT_ID to be set. If it's missing, `available` stays
 * false and callers should hide the Google button rather than send a fake token.
 */
const GIS_SRC = 'https://accounts.google.com/gsi/client'

const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID || ''
const available = ref(!!clientId)

let scriptPromise = null

function loadScript() {
  if (window.google?.accounts?.id) return Promise.resolve()
  if (scriptPromise) return scriptPromise
  scriptPromise = new Promise((resolve, reject) => {
    const existing = document.querySelector(`script[src="${GIS_SRC}"]`)
    if (existing) {
      existing.addEventListener('load', () => resolve())
      existing.addEventListener('error', () => reject(new Error('Không tải được Google SDK')))
      return
    }
    const s = document.createElement('script')
    s.src = GIS_SRC
    s.async = true
    s.defer = true
    s.onload = () => resolve()
    s.onerror = () => reject(new Error('Không tải được Google SDK'))
    document.head.appendChild(s)
  })
  return scriptPromise
}

export function useGoogleAuth() {
  /**
   * Trigger the Google Sign-In popup and resolve with a verified id_token.
   * Uses the ID-token callback flow (initialize + prompt) which returns a JWT
   * credential the Go backend verifies against Google's tokeninfo endpoint.
   */
  const signIn = () => {
    return new Promise(async (resolve, reject) => {
      if (!clientId) {
        reject(new Error('Google Sign-In chưa được cấu hình (thiếu VITE_GOOGLE_CLIENT_ID).'))
        return
      }
      try {
        await loadScript()
      } catch (e) {
        reject(e)
        return
      }

      const google = window.google
      if (!google?.accounts?.id) {
        reject(new Error('Google SDK chưa sẵn sàng.'))
        return
      }

      google.accounts.id.initialize({
        client_id: clientId,
        callback: (response) => {
          if (response?.credential) {
            resolve(response.credential) // this is the id_token (JWT)
          } else {
            reject(new Error('Không nhận được thông tin xác thực từ Google.'))
          }
        },
      })

      // Render a one-tap / popup prompt.
      google.accounts.id.prompt((notification) => {
        if (notification.isNotDisplayed?.() || notification.isSkippedMoment?.()) {
          reject(new Error('Cửa sổ đăng nhập Google đã bị đóng hoặc bị chặn.'))
        }
      })
    })
  }

  return { available, signIn }
}
