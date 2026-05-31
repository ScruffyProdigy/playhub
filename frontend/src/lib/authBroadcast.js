const AUTH_BROADCAST_CHANNEL = 'joinquest-auth'

export function notifyAuthComplete() {
  try {
    const channel = new BroadcastChannel(AUTH_BROADCAST_CHANNEL)
    channel.postMessage({ type: 'signed-in' })
    channel.close()
  } catch {
    // BroadcastChannel is unavailable in some browsers/contexts.
  }
}

export function onAuthComplete(callback) {
  try {
    const channel = new BroadcastChannel(AUTH_BROADCAST_CHANNEL)
    channel.onmessage = (event) => {
      if (event.data?.type === 'signed-in') {
        callback()
      }
    }
    return () => channel.close()
  } catch {
    return () => {}
  }
}
