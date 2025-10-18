import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import App from './App'

describe('App Component', () => {
  it('renders without crashing', () => {
    render(<App />)
    expect(screen.getByText('PlayHub')).toBeInTheDocument()
  })

  it('displays the main heading', () => {
    render(<App />)
    expect(screen.getByRole('heading', { name: 'PlayHub' })).toBeInTheDocument()
  })

  it('displays the tagline', () => {
    render(<App />)
    expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()
  })

  it('has proper heading structure', () => {
    render(<App />)
    const heading = screen.getByRole('heading', { level: 1 })
    expect(heading).toHaveTextContent('PlayHub')
  })

  it('renders all expected content', () => {
    render(<App />)
    
    expect(screen.getByText('PlayHub')).toBeInTheDocument()
    expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()
  })

  it('has accessible content structure', () => {
    render(<App />)
    
    const heading = screen.getByRole('heading')
    const paragraph = screen.getByText('Your Gaming Hub - Queue, Play, Trade')
    
    expect(heading).toBeInTheDocument()
    expect(paragraph).toBeInTheDocument()
  })

  it('maintains content across re-renders', () => {
    const { rerender } = render(<App />)
    
    expect(screen.getByText('PlayHub')).toBeInTheDocument()
    expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()
    
    // Re-render the component
    rerender(<App />)
    
    expect(screen.getByText('PlayHub')).toBeInTheDocument()
    expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()
  })
})

describe('App Component Integration', () => {
  it('renders all expected elements in correct order', () => {
    render(<App />)

    const heading = screen.getByRole('heading', { name: 'PlayHub' })
    const tagline = screen.getByText('Your Gaming Hub - Queue, Play, Trade')

    expect(heading).toBeInTheDocument()
    expect(tagline).toBeInTheDocument()

    // Basic order check
    expect(heading.compareDocumentPosition(tagline)).toBe(Node.DOCUMENT_POSITION_FOLLOWING)
  })

  it('has proper semantic structure', () => {
    render(<App />)
    
    const heading = screen.getByRole('heading', { level: 1 })
    const paragraph = screen.getByText('Your Gaming Hub - Queue, Play, Trade')
    
    expect(heading.tagName).toBe('H1')
    expect(paragraph.tagName).toBe('P')
  })
})