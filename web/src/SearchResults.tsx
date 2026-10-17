import { useState, useEffect } from 'react'
import type { Track } from './types'

type Props = {
  query: string
  onPlay: (tracks: Track[], index: number) => void
}

export default function SearchResults({ query, onPlay }: Props) {
  const [result, setResult] = useState<{ query: string; tracks: Track[] } | null>(null)

  useEffect(() => {
    const controller = new AbortController()

    // Wait 300ms after the last keystroke before searching.
    const timer = setTimeout(() => {
      fetch(`/api/tracks?q=${encodeURIComponent(query)}`, {
        signal: controller.signal,
      })
        .then((res) => res.json())
        .then((tracks: Track[]) => setResult({ query, tracks }))
        .catch((err) => {
          if (err.name !== 'AbortError') console.error(err)
        })
    }, 300)

    // Runs before the effect re-runs (and when the component disappears).
    return () => {
      clearTimeout(timer)
      controller.abort()
    }
  }, [query])

  // Derived: we're loading if the results we hold are for an older query.
  const loading = result === null || result.query !== query

  if (result === null) {
    return <p>Searching…</p>
  }

  return (
    <div className={loading ? 'stale' : ''}>
      {result.tracks.length === 0 && !loading && <p>No matches.</p>}
      <ul>
        {result.tracks.map((track, index) => (
          <li key={track.id}>
            <button onClick={() => onPlay(result.tracks, index)}>
              {track.title} – {track.artist} ({track.album})
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}
