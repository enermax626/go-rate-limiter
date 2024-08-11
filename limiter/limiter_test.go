package limiter_test

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/enermax626/go-ratelimiter/limiter"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {

	os.Setenv("LIMIT_IP", "5")            // Maximum 5 requests per IP
	os.Setenv("LIMIT_TOKEN", "10")        // Maximum 10 requests per token (not used in this test)
	os.Setenv("OVER_LIMIT_COOLDOWN", "4") // Cooldown period of 4 seconds

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Could not start miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	store := limiter.NewRedisStore(*rdb)
	rateLimiter := limiter.NewRateLimiter(store)

	router := chi.NewRouter()
	router.Use(rateLimiter.Middleware)
	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}

	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", server.URL+"/hello", nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	req, _ := http.NewRequest("GET", server.URL+"/hello", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	mr.FastForward(10 * time.Second)

	req, _ = http.NewRequest("GET", server.URL+"/hello", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
