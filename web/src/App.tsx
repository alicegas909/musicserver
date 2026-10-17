import { useState, useEffect } from 'react'
import type { Album, Track } from './types'
import AlbumView from './AlbumView'
import Player from './Player'

export default function App() {
  const [albums, setAlbums] = useState<Album[]>([])
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [queue, setQueue] = useState<Track[]>([])
  const [queueIndex, setQueueIndex] = useState(-1)

  useEffect(() => {
    fetch('/api/albums')
      .then((res) => res.json())
      .then((data: Album[]) => setAlbums(data))
  }, [])

  const currentTrack = queueIndex >= 0 ? queue[queueIndex] : null
  const hasPrev = queueIndex > 0
  const hasNext = queueIndex < queue.length - 1

  function playFrom(tracks: Track[], index: number) {
    setQueue(tracks)
    setQueueIndex(index)
  }

  function prev() {
    if (hasPrev) setQueueIndex(queueIndex - 1)
  }

  function next() {
    if (hasNext) setQueueIndex(queueIndex + 1)
  }

  return (
    <div>
      <h1>My Music</h1>
      {selectedId === null ? (
        <ul>
          {albums.map((album) => (
            <li key={album.id}>
              <button onClick={() => setSelectedId(album.id)}>
                {album.artist} – {album.title} ({album.track_count} tracks)
              </button>
            </li>
          ))}
        </ul>
      ) : (
        <AlbumView
          id={selectedId}
          onBack={() => setSelectedId(null)}
          onPlay={playFrom}
        />
      )}
      <Player
        track={currentTrack}
        hasPrev={hasPrev}
        hasNext={hasNext}
        onPrev={prev}
        onNext={next}
      />
    </div>
  )
}
