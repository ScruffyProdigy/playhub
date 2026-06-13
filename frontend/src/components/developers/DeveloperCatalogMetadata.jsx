import { useEffect, useState } from 'react'
import {
  fetchCatalogTagTaxonomy,
  updateMyGameMetadata,
} from '../../lib/developers'

export default function DeveloperCatalogMetadata({ game, onSaved }) {
  const [shortDescription, setShortDescription] = useState(game?.shortDescription ?? '')
  const [longDescription, setLongDescription] = useState(game?.longDescription ?? '')
  const [howToPlay, setHowToPlay] = useState(game?.howToPlay ?? '')
  const [selectedTags, setSelectedTags] = useState(game?.tags ?? [])
  const [taxonomy, setTaxonomy] = useState([])
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState('')

  useEffect(() => {
    setShortDescription(game?.shortDescription ?? '')
    setLongDescription(game?.longDescription ?? '')
    setHowToPlay(game?.howToPlay ?? '')
    setSelectedTags(game?.tags ?? [])
  }, [game])

  useEffect(() => {
    void fetchCatalogTagTaxonomy()
      .then(setTaxonomy)
      .catch(() => setTaxonomy([]))
  }, [])

  function toggleTag(tagId) {
    setSelectedTags((prev) => {
      if (prev.includes(tagId)) {
        return prev.filter((t) => t !== tagId)
      }
      if (prev.length >= 3) {
        return prev
      }
      return [...prev, tagId]
    })
  }

  async function handleSubmit(event) {
    event.preventDefault()
    if (!game?.id) {
      return
    }
    setStatus('saving')
    setError('')
    try {
      const updated = await updateMyGameMetadata({
        gameId: game.id,
        shortDescription: shortDescription.trim(),
        longDescription: longDescription.trim(),
        howToPlay: howToPlay.trim(),
        tags: selectedTags,
      })
      onSaved?.(updated)
      setStatus('idle')
    } catch (err) {
      setError(err.message || 'Could not save catalog listing.')
      setStatus('idle')
    }
  }

  return (
    <section className="panel-card" aria-labelledby="catalog-metadata-heading">
      <h2 id="catalog-metadata-heading">Catalog listing</h2>
      <p className="panel-copy">
        Player-facing copy for your catalog card and detail page. Agents can draft these fields
        after the discovery interview — edit and save here.
      </p>
      <form className="developer-form" onSubmit={(event) => void handleSubmit(event)}>
        <div className="developer-form__field">
          <label htmlFor="dev-short-description">Short description</label>
          <textarea
            id="dev-short-description"
            rows={2}
            maxLength={200}
            value={shortDescription}
            onChange={(event) => setShortDescription(event.target.value)}
          />
        </div>
        <div className="developer-form__field">
          <label htmlFor="dev-long-description">Long description</label>
          <textarea
            id="dev-long-description"
            rows={5}
            value={longDescription}
            onChange={(event) => setLongDescription(event.target.value)}
          />
        </div>
        <div className="developer-form__field">
          <label htmlFor="dev-how-to-play">How to play</label>
          <textarea
            id="dev-how-to-play"
            rows={4}
            value={howToPlay}
            onChange={(event) => setHowToPlay(event.target.value)}
            placeholder="Step-by-step for first-time players"
          />
        </div>
        {taxonomy.length > 0 ? (
          <fieldset className="developer-form__field developer-tag-picker">
            <legend>Tags (up to 3)</legend>
            <ul className="developer-tag-picker__list">
              {taxonomy.map((tag) => (
                <li key={tag.id}>
                  <label className="developer-tag-picker__option">
                    <input
                      type="checkbox"
                      checked={selectedTags.includes(tag.id)}
                      onChange={() => toggleTag(tag.id)}
                    />
                    <span>{tag.label}</span>
                  </label>
                </li>
              ))}
            </ul>
          </fieldset>
        ) : null}
        <button type="submit" className="button-secondary" disabled={status === 'saving'}>
          {status === 'saving' ? 'Saving…' : 'Save listing'}
        </button>
        {error ? (
          <p className="status-message status-message-error" role="alert">
            {error}
          </p>
        ) : null}
      </form>
    </section>
  )
}
