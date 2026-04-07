package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Feed is a single tracked RSS/Atom feed subscription. ItemCount tracks
// the total number of items ever seen. UnreadCount tracks how many of
// those haven't been marked read yet. Status is one of: active, paused,
// error, archived.
type Feed struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	SiteURL       string `json:"site_url"`
	Category      string `json:"category"`
	ItemCount     int    `json:"item_count"`
	UnreadCount   int    `json:"unread_count"`
	LastFetchedAt string `json:"last_fetched_at"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "feedreader.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS feeds(
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		site_url TEXT DEFAULT '',
		category TEXT DEFAULT '',
		item_count INTEGER DEFAULT 0,
		unread_count INTEGER DEFAULT 0,
		last_fetched_at TEXT DEFAULT '',
		status TEXT DEFAULT 'active',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_status ON feeds(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_category ON feeds(category)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_last_fetched ON feeds(last_fetched_at)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *Feed) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Status == "" {
		e.Status = "active"
	}
	_, err := d.db.Exec(
		`INSERT INTO feeds(id, title, url, site_url, category, item_count, unread_count, last_fetched_at, status, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Title, e.URL, e.SiteURL, e.Category, e.ItemCount, e.UnreadCount, e.LastFetchedAt, e.Status, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Feed {
	var e Feed
	err := d.db.QueryRow(
		`SELECT id, title, url, site_url, category, item_count, unread_count, last_fetched_at, status, created_at
		 FROM feeds WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Title, &e.URL, &e.SiteURL, &e.Category, &e.ItemCount, &e.UnreadCount, &e.LastFetchedAt, &e.Status, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Feed {
	rows, _ := d.db.Query(
		`SELECT id, title, url, site_url, category, item_count, unread_count, last_fetched_at, status, created_at
		 FROM feeds ORDER BY title ASC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Feed
	for rows.Next() {
		var e Feed
		rows.Scan(&e.ID, &e.Title, &e.URL, &e.SiteURL, &e.Category, &e.ItemCount, &e.UnreadCount, &e.LastFetchedAt, &e.Status, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Feed) error {
	_, err := d.db.Exec(
		`UPDATE feeds SET title=?, url=?, site_url=?, category=?, item_count=?, unread_count=?, last_fetched_at=?, status=?
		 WHERE id=?`,
		e.Title, e.URL, e.SiteURL, e.Category, e.ItemCount, e.UnreadCount, e.LastFetchedAt, e.Status, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM feeds WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM feeds`).Scan(&n)
	return n
}

// MarkFetched is the dedicated endpoint for the most common automation:
// after a feed reader checks a URL, it POSTs the new item count, new
// unread count, and stamps last_fetched_at without touching title/url/
// category/site_url. Status is also updated (e.g. to "error" if the
// fetch failed).
func (d *DB) MarkFetched(id string, itemCount, unreadCount int, status string) error {
	if status == "" {
		status = "active"
	}
	_, err := d.db.Exec(
		`UPDATE feeds SET item_count=?, unread_count=?, status=?, last_fetched_at=? WHERE id=?`,
		itemCount, unreadCount, status, now(), id,
	)
	return err
}

// MarkAllRead atomically zeroes a feed's unread_count without touching
// item_count or status. The original implementation had no way to do
// this except by sending a full PUT, which would have triggered the
// 'if UnreadCount == 0 preserve' bug.
func (d *DB) MarkAllRead(id string) error {
	_, err := d.db.Exec(
		`UPDATE feeds SET unread_count=0 WHERE id=?`,
		id,
	)
	return err
}

// IncrementUnread atomically bumps unread_count by N (typically called
// when a fetch sees N new items).
func (d *DB) IncrementUnread(id string, by int) error {
	_, err := d.db.Exec(
		`UPDATE feeds SET unread_count = unread_count + ?, item_count = item_count + ? WHERE id=?`,
		by, by, id,
	)
	return err
}

func (d *DB) Search(q string, filters map[string]string) []Feed {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (title LIKE ? OR url LIKE ? OR site_url LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s)
	}
	if v, ok := filters["category"]; ok && v != "" {
		where += " AND category=?"
		args = append(args, v)
	}
	if v, ok := filters["status"]; ok && v != "" {
		where += " AND status=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, title, url, site_url, category, item_count, unread_count, last_fetched_at, status, created_at
		 FROM feeds WHERE `+where+`
		 ORDER BY unread_count DESC, title ASC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Feed
	for rows.Next() {
		var e Feed
		rows.Scan(&e.ID, &e.Title, &e.URL, &e.SiteURL, &e.Category, &e.ItemCount, &e.UnreadCount, &e.LastFetchedAt, &e.Status, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns total feeds, total unread items across all feeds,
// total items across all feeds, plus by_status and by_category
// breakdowns. The original returned just total + by_status.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":        d.Count(),
		"total_unread": 0,
		"total_items":  0,
		"with_unread":  0,
		"by_status":    map[string]int{},
		"by_category":  map[string]int{},
	}

	var totalUnread, totalItems, withUnread int
	d.db.QueryRow(`SELECT COALESCE(SUM(unread_count), 0) FROM feeds`).Scan(&totalUnread)
	d.db.QueryRow(`SELECT COALESCE(SUM(item_count), 0) FROM feeds`).Scan(&totalItems)
	d.db.QueryRow(`SELECT COUNT(*) FROM feeds WHERE unread_count > 0`).Scan(&withUnread)
	m["total_unread"] = totalUnread
	m["total_items"] = totalItems
	m["with_unread"] = withUnread

	if rows, _ := d.db.Query(`SELECT status, COUNT(*) FROM feeds GROUP BY status`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_status"] = by
	}

	if rows, _ := d.db.Query(`SELECT category, COUNT(*) FROM feeds WHERE category != '' GROUP BY category`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_category"] = by
	}

	return m
}

// ─── Extras ───────────────────────────────────────────────────────

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
