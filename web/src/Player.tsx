import type { Track } from './types'

type Props = {
  track: Track | null
  hasPrev: boolean
  hasNext: boolean
  onPrev: () => void
  onNext: () => void
}

export default function Player({ track, hasPrev, hasNext, onPrev, onNext }: Props) {
  return (
    <div className="player">
      <div>
        {track ? (
          <>
            <strong>{track.title}</strong>
            <div>{track.artist}</div>
          </>
        ) : (
          <span>Nothing playing</span>
        )}
      </div>
      <button onClick={onPrev} disabled={!hasPrev}>⏮</button>
      <button onClick={onNext} disabled={!hasNext}>⏭</button>
      <audio
        src={track ? `/stream/${track.id}` : undefined}
        controls
        autoPlay
        onEnded={onNext}
      />
    </div>
  )
}
