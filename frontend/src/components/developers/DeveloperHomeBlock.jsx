import YourGamesStrip from './YourGamesStrip'

export default function DeveloperHomeBlock() {
  return (
    <section className="panel-card developer-home-block" aria-labelledby="developer-home-block-heading">
      <div className="developer-home-block__cta">
        <p id="developer-home-block-heading" className="developer-home-block__copy">
          <strong>Making a game?</strong> We&apos;d love to host it.
        </p>
        <a className="developer-home-block__link" href="/developers">
          Get started for developers →
        </a>
      </div>
      <YourGamesStrip embedded />
    </section>
  )
}
