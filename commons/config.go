package commons

var Config = struct {
	JwtSecret string
}{
	JwtSecret: GetEnv("JWT_SECRET"),
}
