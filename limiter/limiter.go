package limiter

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	store         Store
	limitIP       int
	limitApiKey   int
	blockDuration time.Duration
}

func NewRateLimiter(store Store) *RateLimiter {
	limitIp, _ := strconv.Atoi(os.Getenv("LIMIT_IP"))
	limitApiKey, _ := strconv.Atoi(os.Getenv("LIMIT_TOKEN"))
	blockDuration, _ := strconv.Atoi(os.Getenv("OVER_LIMIT_COOLDOWN"))

	return &RateLimiter{
		store:         store,
		limitIP:       limitIp,
		limitApiKey:   limitApiKey,
		blockDuration: time.Duration(blockDuration) * time.Second,
	}
}

func extractIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (rl *RateLimiter) IsValidRequest(ip, apiKey string) (bool, error) {
	var limit int

	key := ip
	if apiKey != "" && rl.limitApiKey > rl.limitIP {
		limit = rl.limitApiKey
		key = apiKey
	} else {
		limit = rl.limitIP
	}

	return rl.store.Allow(key, limit, rl.blockDuration)
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		apiKey := r.Header.Get("API_KEY")

		isValidRequest, err := rl.IsValidRequest(ip, apiKey)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !isValidRequest {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func InitializeRateLimiters() *RateLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-stack:6379",
	})

	redisStore := NewRedisStore(*rdb)
	return NewRateLimiter(redisStore)
}
