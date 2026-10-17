export type Album = {
  id: number
  title: string
  artist: string
  track_count: number
}

export type Track = {
  id: number
  title: string
  artist: string
  album: string
  album_id: number
  track_no: number
}

export type AlbumDetail = Album & {
  tracks: Track[]
}
