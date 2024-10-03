package tokenblacklist

import "sync"

var (
	blacklist = make(map[string]struct{})
	mu        sync.Mutex
)

// AddToken adds a token to the blacklist
func AddToken(token string) {
	mu.Lock()
	defer mu.Unlock()
	blacklist[token] = struct{}{}
}

// IsTokenBlacklisted checks if a token is in the blacklist
func IsTokenBlacklisted(token string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := blacklist[token]
	return exists
}
