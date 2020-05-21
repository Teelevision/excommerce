package config

import (
	"log"
	"os"
	"time"
)

// app-wide configuration
var (
	CouponDefaultLifetime = 10 * time.Second
)

// parse COUPON_DEFAULT_LIFETIME
func init() {
	value := os.Getenv("COUPON_DEFAULT_LIFETIME")
	if value == "" {
		return
	}

	dur, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf(`Could not parse value %q of COUPON_DEFAULT_LIFETIME env: %s
Use values like "10s", "2.5m" or "1h30m" to express a duration.`, value, err)
	}

	if dur == time.Duration(0) {
		return
	}

	if dur < time.Second {
		log.Println("Notice: The value of COUPON_DEFAULT_LIFETIME is less than a second.")
	}

	CouponDefaultLifetime = dur
}
