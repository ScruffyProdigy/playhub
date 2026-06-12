/** Tab visibility helper — re-sync state when the user returns to the tab. */
export function onTabVisible(callback) {
  if (typeof document === 'undefined') {
    return () => {}
  }
  const handler = () => {
    if (document.visibilityState === 'visible') {
      callback()
    }
  }
  document.addEventListener('visibilitychange', handler)
  return () => document.removeEventListener('visibilitychange', handler)
}
