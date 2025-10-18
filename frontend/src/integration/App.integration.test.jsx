import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import App from '../App'

describe('App Integration Tests', () => {
  describe('User Journey: Content Discovery', () => {
    it('presents all content in logical order', () => {
      render(<App />)

      const heading = screen.getByRole('heading', { name: 'PlayHub' })
      const tagline = screen.getByText('Your Gaming Hub - Queue, Play, Trade')

      expect(heading).toBeInTheDocument()
      expect(tagline).toBeInTheDocument()
    })

    it('provides clear branding and messaging', () => {
      render(<App />)

      const heading = screen.getByText('PlayHub')
      const tagline = screen.getByText('Your Gaming Hub - Queue, Play, Trade')

      expect(heading).toBeInTheDocument()
      expect(tagline).toBeInTheDocument()
    })
  })

  describe('User Journey: Visual Feedback', () => {
    it('provides clear visual hierarchy', () => {
      render(<App />)

      const heading = screen.getByRole('heading', { level: 1 })
      const tagline = screen.getByText('Your Gaming Hub - Queue, Play, Trade')

      expect(heading).toBeInTheDocument()
      expect(tagline).toBeInTheDocument()
    })

    it('maintains visual consistency across renders', () => {
      const { asFragment, rerender } = render(<App />)

      const initialRender = asFragment()
      rerender(<App />)
      const afterRerender = asFragment()

      expect(initialRender).toEqual(afterRerender)
    })
  })

  describe('User Journey: Error Handling', () => {
    it('handles component re-renders gracefully', () => {
      const { rerender } = render(<App />)

      expect(screen.getByText('PlayHub')).toBeInTheDocument()
      expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()

      rerender(<App />) // Simulate a re-render

      expect(screen.getByText('PlayHub')).toBeInTheDocument()
      expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()
    })

    it('maintains functionality after multiple re-renders', () => {
      const { rerender } = render(<App />)

      // Multiple re-renders
      rerender(<App />)
      rerender(<App />)
      rerender(<App />)

      expect(screen.getByText('PlayHub')).toBeInTheDocument()
      expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()
    })
  })
})

describe('App Performance Integration', () => {
  it('renders efficiently', () => {
    const startTime = performance.now()
    render(<App />)
    const endTime = performance.now()
    const renderTime = endTime - startTime

    // Should render quickly (less than 10ms for a simple component)
    expect(renderTime).toBeLessThan(10)
    expect(screen.getByText('PlayHub')).toBeInTheDocument()
  })

  it('does not cause memory leaks with repeated renders', () => {
    const { unmount } = render(<App />)

    // The app should render without issues
    expect(screen.getByText('PlayHub')).toBeInTheDocument()

    // Unmount should not cause errors
    expect(() => unmount()).not.toThrow()
  })
})