/** Focus the sign-in code field and place the caret at the end (for keyboard / OTP autofill). */
export function focusCodeInput(input) {
  if (!input) return
  input.focus({ preventScroll: true })
  try {
    const end = input.value.length
    input.setSelectionRange(end, end)
  } catch {
    // setSelectionRange is not supported on all input types
  }
}
