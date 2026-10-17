import { useState, useEffect } from 'react'
import type { AlbumDetail, Track } from './types'

type Props = {
  id: number
  onBack: () => void
  onPlay: (tracks: Track[], index: number) => void
}

export default function AlbumView({ id, onBack, onPlay }: Props) {
  const [album, setAlbum] = useState<AlbumDetail | null>(null)

  useEffect(() => {
    fetch(`/api/albums/${id}`)
      .then((res) => res.json())
      .then((data: AlbumDetail) => setAlbum(data))
  }, [id])

  if (album === null) {
    return <p>Loading…</p>
  }

  return (
    <div>
      <button onClick={onBack}>← Back</button>
      <h2>{album.title}</h2>
      <p>{album.artist}</p>
      <ol>
        {album.tracks.map((track, index) => (
          <li key={track.id}>
            <button onClick={() => onPlay(album.tracks, index)}>
              {track.title}
            </button>
          </li>
        ))}
      </ol>
    </div>
  )
}
