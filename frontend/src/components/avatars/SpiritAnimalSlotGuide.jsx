import { JOURNEY_SLOTS } from '../../lib/spiritAnimalSlots'

export default function SpiritAnimalSlotGuide({ highlightKey = null, compact = false }) {
  if (compact) {
    return (
      <ul className="spirit-animal__slot-strip" role="list" aria-label="Journey slots">
        {JOURNEY_SLOTS.map((slot, index) => (
          <li
            key={slot.key}
            className={
              highlightKey === slot.key
                ? 'spirit-animal__slot-strip-item spirit-animal__slot-strip-item--active'
                : 'spirit-animal__slot-strip-item'
            }
            title={`${index + 1}. ${slot.name} — ${slot.prompt}`}
          >
            <img src={slot.imageUrl} alt="" className="spirit-animal__slot-strip-icon" />
            <span className="spirit-animal__slot-strip-name">{slot.name}</span>
          </li>
        ))}
      </ul>
    )
  }

  return (
    <ol className="spirit-animal__slot-guide" aria-label="Five chapters of the journey">
      {JOURNEY_SLOTS.map((slot, index) => {
        const active = highlightKey === slot.key
        return (
          <li
            key={slot.key}
            className={active ? 'spirit-animal__slot-guide-item spirit-animal__slot-guide-item--active' : 'spirit-animal__slot-guide-item'}
          >
            <span className="spirit-animal__slot-guide-step" aria-hidden="true">
              {index + 1}
            </span>
            <img src={slot.imageUrl} alt="" className="spirit-animal__slot-guide-icon" />
            <div className="spirit-animal__slot-guide-copy">
              <p className="spirit-animal__slot-guide-name">{slot.name}</p>
              <p className="spirit-animal__slot-guide-prompt">{slot.prompt}</p>
            </div>
          </li>
        )
      })}
    </ol>
  )
}
