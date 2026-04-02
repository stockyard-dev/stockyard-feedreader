package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Feed struct {
	ID string `json:"id"`
	Title string `json:"title"`
	URL string `json:"url"`
	SiteURL string `json:"site_url"`
	Category string `json:"category"`
	ItemCount int `json:"item_count"`
	UnreadCount int `json:"unread_count"`
	LastFetchedAt string `json:"last_fetched_at"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"feedreader.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS feeds(id TEXT PRIMARY KEY,title TEXT NOT NULL,url TEXT NOT NULL,site_url TEXT DEFAULT '',category TEXT DEFAULT '',item_count INTEGER DEFAULT 0,unread_count INTEGER DEFAULT 0,last_fetched_at TEXT DEFAULT '',status TEXT DEFAULT 'active',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Feed)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO feeds(id,title,url,site_url,category,item_count,unread_count,last_fetched_at,status,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.URL,e.SiteURL,e.Category,e.ItemCount,e.UnreadCount,e.LastFetchedAt,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Feed{var e Feed;if d.db.QueryRow(`SELECT id,title,url,site_url,category,item_count,unread_count,last_fetched_at,status,created_at FROM feeds WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.URL,&e.SiteURL,&e.Category,&e.ItemCount,&e.UnreadCount,&e.LastFetchedAt,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Feed{rows,_:=d.db.Query(`SELECT id,title,url,site_url,category,item_count,unread_count,last_fetched_at,status,created_at FROM feeds ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Feed;for rows.Next(){var e Feed;rows.Scan(&e.ID,&e.Title,&e.URL,&e.SiteURL,&e.Category,&e.ItemCount,&e.UnreadCount,&e.LastFetchedAt,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Feed)error{_,err:=d.db.Exec(`UPDATE feeds SET title=?,url=?,site_url=?,category=?,item_count=?,unread_count=?,last_fetched_at=?,status=? WHERE id=?`,e.Title,e.URL,e.SiteURL,e.Category,e.ItemCount,e.UnreadCount,e.LastFetchedAt,e.Status,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM feeds WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM feeds`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Feed{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (title LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["category"];ok&&v!=""{where+=" AND category=?";args=append(args,v)}
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,title,url,site_url,category,item_count,unread_count,last_fetched_at,status,created_at FROM feeds WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Feed;for rows.Next(){var e Feed;rows.Scan(&e.ID,&e.Title,&e.URL,&e.SiteURL,&e.Category,&e.ItemCount,&e.UnreadCount,&e.LastFetchedAt,&e.Status,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM feeds GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
