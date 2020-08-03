package commons

import "time"

// Configuration for some parts of the application
var Config = struct {
	JwtSecret         string
	JwtTokenExpiresIn time.Duration
	EncryptionCost    int
}{
	JwtSecret:         GetEnv("JWT_SECRET"),
	JwtTokenExpiresIn: 2 * time.Hour,
	EncryptionCost:    14,
}
