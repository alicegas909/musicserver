package main

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

func scanLibrary(db *sql.DB, root string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	count := 0
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("skipping %s: %v", path, err)
			return nil
		}
		if d.IsDir() || !strings.EqualFold(filepath.Ext(path), ".mp3") {
			return nil
		}
		if err := indexFile(tx, path); err != nil {
			log.Printf("could not index %s: %v", path, err)
			return nil
		}
		count++
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("scan finished: %d tracks indexed", count)
	return tx.Commit()
}

func indexFile(tx *sql.Tx, path string) error {

	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	artist, album, trackNo := "Unknown Artist", "Unknown Album", 0

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if m, err := tag.ReadFrom(f); err == nil {
		if s := strings.TrimSpace(m.Title()); s != "" {
			title = s
		}
		// Prefer album artist so compilations stay grouped together.
		if s := strings.TrimSpace(m.AlbumArtist()); s != "" {
			artist = s
		} else if s := strings.TrimSpace(m.Artist()); s != "" {
			artist = s
		}
		if s := strings.TrimSpace(m.Album()); s != "" {
			album = s
		}
		trackNo, _ = m.Track()

	}
	var artistID, albumID int64
	err = tx.QueryRow(`
		INSERT INTO artists (name) VALUES (?)
		ON CONFLICT(name) DO UPDATE SET name = excluded.name
		RETURNING id`, artist).Scan(&artistID)
	if err != nil {
		return err
	}

	err = tx.QueryRow(`
		INSERT INTO albums (title, artist_id) VALUES (?, ?)
		ON CONFLICT(title, artist_id) DO UPDATE SET title = excluded.title
		RETURNING id`, album, artistID).Scan(&albumID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO tracks (path, title, album_id, track_no) VALUES (?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			title    = excluded.title,
			album_id = excluded.album_id,
			track_no = excluded.track_no`,
		path, title, albumID, trackNo)

	return err
}
