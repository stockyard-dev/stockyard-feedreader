package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-feedreader/internal/store")
func(s *Server)handleListFeeds(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListFeeds();if list==nil{list=[]store.Feed{}};writeJSON(w,200,list)}
func(s *Server)handleCreateFeed(w http.ResponseWriter,r *http.Request){var f store.Feed;json.NewDecoder(r.Body).Decode(&f);if f.Name==""||f.URL==""{writeError(w,400,"name and url required");return};if f.Type==""{f.Type="rss"};s.db.CreateFeed(&f);writeJSON(w,201,f)}
func(s *Server)handleDeleteFeed(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteFeed(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListItems(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);unread:=r.URL.Query().Get("unread")=="true";list,_:=s.db.ListItems(id,unread);if list==nil{list=[]store.FeedItem{}};writeJSON(w,200,list)}
func(s *Server)handleAddItem(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var item store.FeedItem;json.NewDecoder(r.Body).Decode(&item);item.FeedID=id;if item.Title==""{writeError(w,400,"title required");return};s.db.AddItem(&item);writeJSON(w,201,item)}
func(s *Server)handleMarkRead(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.MarkRead(id);writeJSON(w,200,map[string]string{"status":"read"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
