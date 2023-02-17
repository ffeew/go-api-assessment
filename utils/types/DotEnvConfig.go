package types

type Config struct {
	Debug bool `env:"DEBUG"`
	App   struct {
		Port               int16  `env:"PORT"`
		RefreshTokenSecret string `env:"REFRESH_TOKEN_SECRET"`
		AccessTokenSecret  string `env:"ACCESS_TOKEN_SECRET"`
	}
	Database struct {
		Name     string `env:"DB_NAME"`
		Port     int16  `env:"DB_PORT"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Host     string `env:"DB_HOST"`
	}
}
