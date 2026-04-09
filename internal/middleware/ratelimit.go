package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var clients = make(map[string]int)
var mu sync.Mutex//mutex to prevent race conditions
//allowa max 100 requests per IP per min. If exceeded returns 429 too many requests
func RateLimitMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		ip,_,_:= net.SplitHostPort(r.RemoteAddr)//extract client ip
		mu.Lock()//lock map
		clients[ip]++//increase request count
		count:=clients[ip]//store count
		mu.Unlock()//unlock map
		if count >100{
			http.Error(w,"Too many requests", http.StatusTooManyRequests)
			return 
		}
		next.ServeHTTP(w,r)
	})
}

func init(){//runs automatically when app starts
	go func(){
		for{
			time.Sleep(time.Minute)
			mu.Lock()
			clients=make(map[string]int)
			mu.Unlock()
		}
	}()//clears map every 1 min and resets the counters
}
//This middleware tracks request count per IP using a shared map and blocks requests exceeding 100 per minute, 
// resetting counts every minute using a background goroutine.