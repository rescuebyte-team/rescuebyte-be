package env

import (
	"time"

	"github.com/spf13/viper"

	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

type Env struct {
	AppEnv                  string        `mapstructure:"APP_ENV"`
	AppLang                 string        `mapstructure:"APP_LANG"`
	AppPort                 string        `mapstructure:"APP_PORT"`
	AppURL                  string        `mapstructure:"APP_URL"`
	FEAppURL                string        `mapstructure:"FE_APP_URL"`
	DBHost                  string        `mapstructure:"DB_HOST"`
	DBPort                  string        `mapstructure:"DB_PORT"`
	DBUser                  string        `mapstructure:"DB_USER"`
	DBPass                  string        `mapstructure:"DB_PASS"`
	DBName                  string        `mapstructure:"DB_NAME"`
	DBMaxConnection         int           `mapstructure:"DB_MAX_CONNECTION"`
	DBMaxIdleConnection     int           `mapstructure:"DB_MAX_IDLE_CONNECTION"`
	DBMaxLifetimeConnection time.Duration `mapstructure:"DB_MAX_LIFETIME_CONNECTION"`
	DBMaxIdleTimeConnection time.Duration `mapstructure:"DB_MAX_IDLE_TIME_CONNECTION"`
	RedisAddress            string        `mapstructure:"REDIS_ADDR"`
	RedisPassword           string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB                 int           `mapstructure:"REDIS_DB"`
	JwtSecretKey            string        `mapstructure:"JWT_SECRET_KEY"`
	JwtExpTime              time.Duration `mapstructure:"JWT_EXP_TIME"`
	JwtRefreshExpTime       time.Duration `mapstructure:"JWT_REFRESH_EXP_TIME"`
	GoogleClientID          string        `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret      string        `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURI       string        `mapstructure:"GOOGLE_REDIRECT_URI"`
	SupabaseURL             string        `mapstructure:"SUPABASE_URL"`
	SupabaseKey             string        `mapstructure:"SUPABASE_KEY"`
	SupabaseStorageBucket   string        `mapstructure:"SUPABASE_STORAGE_BUCKET"`
	MidtransServerKey       string        `mapstructure:"MIDTRANS_SERVER_KEY"`
	MidtransClientKey       string        `mapstructure:"MIDTRANS_CLIENT_KEY"`
	MidtransIsProduction    bool          `mapstructure:"MIDTRANS_IS_PRODUCTION"`
}

var AppEnv = getEnv()

func getEnv() *Env {
	env := &Env{}

	viper.SetConfigFile("./config/.env")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn(log.Fields{
			"error": err.Error(),
		}, "[Env] Failed to read config file, reading from environment variables")
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(env); err != nil {
		log.Fatal(log.Fields{
			"error": err.Error(),
		}, "[Env] Failed to unmarshal environment variables")
	}

	env.setDefaults()

	return env
}

func (env *Env) setDefaults() {
	if env.AppEnv == "" {
		env.AppEnv = "development"
	}
	if env.AppLang == "" {
		env.AppLang = "en"
	}
	if env.AppPort == "" {
		env.AppPort = "8080"
	}
	if env.DBHost == "" {
		env.DBHost = "localhost"
	}
	if env.DBPort == "" {
		env.DBPort = "5432"
	}
	if env.DBUser == "" {
		env.DBUser = "postgres"
	}
	if env.DBPass == "" {
		env.DBPass = "postgres"
	}
	if env.DBName == "" {
		env.DBName = "rescuebyte_db"
	}
	if env.DBMaxConnection == 0 {
		env.DBMaxConnection = 20
	}
	if env.DBMaxIdleConnection == 0 {
		env.DBMaxIdleConnection = 10
	}
	if env.DBMaxLifetimeConnection == 0 {
		env.DBMaxLifetimeConnection = 1 * time.Hour
	}
	if env.DBMaxIdleTimeConnection == 0 {
		env.DBMaxIdleTimeConnection = 15 * time.Minute
	}
	if env.RedisAddress == "" {
		env.RedisAddress = "localhost:6379"
	}
	if env.JwtSecretKey == "" {
		env.JwtSecretKey = "default-rescuebyte-secret-key-change-me"
	}
	if env.JwtExpTime == 0 {
		env.JwtExpTime = 24 * time.Hour
	}
	if env.JwtRefreshExpTime == 0 {
		env.JwtRefreshExpTime = 7 * 24 * time.Hour
	}
}
