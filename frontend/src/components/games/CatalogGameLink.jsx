import { navigateToGameDetail } from '../../lib/catalogNavigation'

export default function CatalogGameLink({ href, slug, className, children, ...props }) {
  function handleClick(event) {
    if (!slug || event.defaultPrevented) {
      return
    }
    if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey || event.button !== 0) {
      return
    }
    event.preventDefault()
    navigateToGameDetail(slug)
  }

  return (
    <a href={href} className={className} onClick={handleClick} {...props}>
      {children}
    </a>
  )
}
