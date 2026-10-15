package main

type Album struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	TrackCount int    `json:"track_count"`
}
type AlbumDetail struct {
	Album
	Tracks []Track `json:"tracks"`
}

type Track struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
	AlbumID int64  `json:"album_id"`
	TrackNo int    `json:"track_no"`
}
