package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-feedreader/internal/server";"github.com/stockyard-dev/stockyard-feedreader/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./feedreader-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("feedreader: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Feedreader — Self-hosted RSS and news aggregator\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("feedreader: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
