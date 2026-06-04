import { describe, it, expect, vi } from 'vitest'
import { focusCodeInput } from './focusCodeInput'

describe('focusCodeInput', () => {
  it('focuses the input and moves the caret to the end', () => {
    const input = document.createElement('input')
    input.value = '12'
    const focus = vi.spyOn(input, 'focus')
    const setSelectionRange = vi.spyOn(input, 'setSelectionRange')

    focusCodeInput(input)

    expect(focus).toHaveBeenCalledWith({ preventScroll: true })
    expect(setSelectionRange).toHaveBeenCalledWith(2, 2)
  })
})
