package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

// Shared SELECT for anything that returns tracks; callers add WHERE/ORDER BY.
const trackSelect = `
	SELECT t.id, t.title, ar.name, al.title, al.id, t.track_no
	FROM tracks t
	JOIN albums al ON al.id = t.album_id
	JOIN artists ar ON ar.id = al.artist_id`

// registerRoutes needs Go 1.22+ for the "GET /path/{id}" pattern syntax.
func registerRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("GET /api/albums", listAlbums(db))
	mux.HandleFunc("GET /api/albums/{id}", getAlbum(db))
	mux.HandleFunc("GET /api/tracks", searchTracks(db))
	mux.HandleFunc("GET /stream/{id}", streamTrack(db))
}

func listAlbums(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT al.id, al.title, ar.name, COUNT(t.id)
			FROM albums al
			JOIN artists ar ON ar.id = al.artist_id
			LEFT JOIN tracks t ON t.album_id = al.id
			GROUP BY al.id
			ORDER BY ar.name, al.title`)
		if err != nil {
			serverError(w, err)
			return
		}
		defer rows.Close()

		albums := []Album{} // empty slice so JSON is [] rather than null
		for rows.Next() {
			var a Album
			if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.TrackCount); err != nil {
				serverError(w, err)
				return
			}
			albums = append(albums, a)
		}
		if err := rows.Err(); err != nil {
			serverError(w, err)
			return
		}
		writeJSON(w, albums)
	}
}

func getAlbum(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var a AlbumDetail
		err = db.QueryRow(`
			SELECT al.id, al.title, ar.name
			FROM albums al
			JOIN artists ar ON ar.id = al.artist_id
			WHERE al.id = ?`, id).Scan(&a.ID, &a.Title, &a.Artist)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			serverError(w, err)
			return
		}

		rows, err := db.Query(trackSelect+` WHERE al.id = ? ORDER BY t.track_no, t.title`, id)
		if err != nil {
			serverError(w, err)
			return
		}
		a.Tracks, err = scanTracks(rows)
		if err != nil {
			serverError(w, err)
			return
		}
		a.TrackCount = len(a.Tracks)
		writeJSON(w, a)
	}
}

// searchTracks handles /api/tracks?q=something, matching title, artist or album.
func searchTracks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		like := "%" + r.URL.Query().Get("q") + "%"

		rows, err := db.Query(trackSelect+`
			WHERE t.title LIKE ? OR ar.name LIKE ? OR al.title LIKE ?
			ORDER BY ar.name, al.title, t.track_no
			LIMIT 200`, like, like, like)
		if err != nil {
			serverError(w, err)
			return
		}
		tracks, err := scanTracks(rows)
		if err != nil {
			serverError(w, err)
			return
		}
		writeJSON(w, tracks)
	}
}

// streamTrack looks the file up by ID, so clients never supply a file path.
// http.ServeContent handles Range requests, which is what makes seeking work.
func streamTrack(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var path string
		err = db.QueryRow(`SELECT path FROM tracks WHERE id = ?`, id).Scan(&path)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			serverError(w, err)
			return
		}

		f, err := os.Open(path)
		if err != nil {
			http.NotFound(w, r) // file was probably moved or deleted
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			serverError(w, err)
			return
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		http.ServeContent(w, r, info.Name(), info.ModTime(), f)
	}
}

// scanTracks reads every row from a trackSelect query and closes the rows.
func scanTracks(rows *sql.Rows) ([]Track, error) {
	defer rows.Close()
	tracks := []Track{}
	for rows.Next() {
		var t Track
		if err := rows.Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.AlbumID, &t.TrackNo); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, rows.Err()
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("json encode: %v", err)
	}
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("error: %v", err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}
