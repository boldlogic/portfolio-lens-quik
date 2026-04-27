package serviceregistry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultKeyPrefix = "pl:registry:v1:"

	KeyRelGRPCQuik     = "endpoint:grpc:quik"
	KeyRelGRPCCurrency = "endpoint:grpc:currency"
)

func NormalizePrefix(prefix string) string {
	p := strings.TrimSpace(prefix)
	if p == "" {
		p = DefaultKeyPrefix
	}
	if !strings.HasSuffix(p, ":") {
		p += ":"
	}
	return p
}

func FullKey(prefix, relative string) string {
	return NormalizePrefix(prefix) + relative
}

// ReadGRPCEndpoints returns host:port for quik and currency gRPC from Redis.
// Empty string means key missing or empty value (caller should fall back to bootstrap config).
func ReadGRPCEndpoints(ctx context.Context, rdb *redis.Client, prefix string) (quik, currency string, err error) {
	qkey := FullKey(prefix, KeyRelGRPCQuik)
	ckey := FullKey(prefix, KeyRelGRPCCurrency)

	s, e := rdb.Get(ctx, qkey).Result()
	if e != nil && e != redis.Nil {
		return "", "", fmt.Errorf("get %s: %w", qkey, e)
	}
	if e == redis.Nil {
		s = ""
	}
	quik = strings.TrimSpace(s)

	s, e = rdb.Get(ctx, ckey).Result()
	if e != nil && e != redis.Nil {
		return "", "", fmt.Errorf("get %s: %w", ckey, e)
	}
	if e == redis.Nil {
		s = ""
	}
	currency = strings.TrimSpace(s)
	return quik, currency, nil
}

// PublishGRPCEndpoint writes host:port for logicalName "quik" or "currency".
// ttl0 means no expiration (not recommended for liveness).
func PublishGRPCEndpoint(ctx context.Context, rdb *redis.Client, prefix, logicalName, hostPort string, ttl time.Duration) error {
	name := strings.ToLower(strings.TrimSpace(logicalName))
	var rel string
	switch name {
	case "quik":
		rel = KeyRelGRPCQuik
	case "currency":
		rel = KeyRelGRPCCurrency
	default:
		return fmt.Errorf("serviceregistry: unknown logical name %q", logicalName)
	}
	hostPort = strings.TrimSpace(hostPort)
	if hostPort == "" {
		return fmt.Errorf("serviceregistry: empty host:port")
	}
	key := FullKey(prefix, rel)
	if ttl > 0 {
		return rdb.Set(ctx, key, hostPort, ttl).Err()
	}
	return rdb.Set(ctx, key, hostPort, 0).Err()
}
