package commons

import "time"

// Config for some parts of the application
var Config = struct {
	JwtSecret         string
	JwtTokenExpiresIn time.Duration
	EncryptionCost    int
	ScrapingEnabled   bool
}{
	JwtSecret:         GetEnv("JWT_SECRET"),
	JwtTokenExpiresIn: 2 * time.Hour,
	EncryptionCost:    14,
	ScrapingEnabled:   GetEnvBool("SCRAPING_ENABLED"),
}
