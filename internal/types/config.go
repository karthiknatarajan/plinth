package types

import (
	"fmt"
	"time"

	"github.com/karthiknatarajan/plinth/internal/platform/pubsub"
)

// Config stores the system configuration.
type Config struct {
	// InstanceID specifis the ID of the cloudness instance.
	// NOTE: If the value is not provided the hostname of the machine is used.
	InstanceID  string `envconfig:"PLINTH_INSTANCE_ID"`
	Environment string `envconfig:"PLINTH_ENVIRONMENT"`

	Log struct {
		Level string `envconfig:"PLINTH_LOG_LEVEL"`
	}

	// GracefulShutdownTime defines the max time we wait when shutting down a server.
	// 5min should be enough for most git clones to complete.
	GracefulShutdownTime time.Duration `envconfig:"PLINTH_GRACEFUL_SHUTDOWN_TIME" default:"300s"`

	Profiler struct {
		Type        string `envconfig:"PLINTH_PROFILER_TYPE"`
		ServiceName string `envconfig:"PLINTH_PROFILER_SERVICE_NAME" default:"cloudness"`
	}

	Instance struct {
		AllowNewTenantCreation bool `envconfig:"PLINTH_INSTANCE_ALLOW_NEW_TENANT" default:"true"`
	}

	ACME struct {
		ACMEUrl    string `envconfig:"PLINTH_ACME_URL" default:"https://acme-v02.api.letsencrypt.org/directory"`
		Email      string `envconfig:"PLINTH_ACME_EMAIL" default:"cloudness@localhost.com"`
		UseStaging bool   `envconfig:"PLINTH_ACME_USE_STAGING" default:"false"`
	}

	// Token defines token configuration parameters.
	Token struct {
		CookieName string        `envconfig:"PLINTH_TOKEN_COOKIE_NAME" default:"token"`
		Expire     time.Duration `envconfig:"PLINTH_TOKEN_EXPIRE" default:"720h"`
	}

	// Database defines the database configuration parameters.
	Database struct {
		Driver     string `envconfig:"PLINTH_DATABASE_DRIVER"     default:"sqlite3"`
		Datasource string `envconfig:"PLINTH_DATABASE_DATASOURCE" default:"database.sqlite3"`
		Host       string `envconfig:"PLINTH_DATABASE_HOST"`
		Port       string `envconfig:"PLINTH_DATABASE_PORT"`
		Name       string `envconfig:"PLINTH_DATABASE_NAME"`
		User       string `envconfig:"PLINTH_DATABASE_USER"`
		Password   string `envconfig:"PLINTH_DATABASE_PASSWORD"`
		SSLMode    string `envconfig:"PLINTH_DATABASE_SSL_MODE"`
	}

	PubSub struct {
		Provider         pubsub.Provider `envconfig:"PLINTH_PUBSUB_PROVIDER"          default:"inmemory"`
		AppNamespace     string          `envconfig:"PLINTH_PUBSUB_APP_NAMESPACE"     default:"cloudness"`
		DefaultNamespace string          `envconfig:"PLINTH_PUBSUB_DEFAULT_NAMESPACE" default:"default"`
		HealthInterval   time.Duration   `envconfig:"PLINTH_PUBSUB_HEALTH_INTERVAL"   default:"3s"`
		SendTimeout      time.Duration   `envconfig:"PLINTH_PUBSUB_SEND_TIMEOUT"      default:"60s"`
		ChannelSize      int             `envconfig:"PLINTH_PUBSUB_CHANNEL_SIZE"      default:"100"`
	}

	Server struct {
		// HTTP defines the http configuration parameters
		HTTP struct {
			Port              int           `envconfig:"PLINTH_HTTP_PORT" default:"7722"`
			Proto             string        `envconfig:"PLINTH_HTTP_PROTO" default:"http"`
			ReadHeaderTimeout time.Duration `envconfig:"PLINTH_HTTP_READ_HEADER_TIMEOUT"   default:"2s"`
			ReadTimeout       time.Duration `envconfig:"PLINTH_HTTP_READ_TIMEOUT"          default:"5s"`
			WriteTimeout      time.Duration `envconfig:"PLINTH_HTTP_WRITE_TIMEOUT"         default:"15s"`
			IdleTimeout       time.Duration `envconfig:"PLINTH_HTTP_IDLE_TIMEOUT"          default:"0"`
		}

		// Acme defines Acme configuration parameters.
		Acme struct {
			Enabled bool   `envconfig:"PLINTH_ACME_ENABLED"`
			Endpont string `envconfig:"PLINTH_ACME_ENDPOINT"`
			Email   bool   `envconfig:"PLINTH_ACME_EMAIL"`
			Host    string `envconfig:"PLINTH_ACME_HOST"`
		}
	}

	// Cors defines http cors parameters
	Cors struct {
		AllowedOrigins   []string `envconfig:"PLINTH_CORS_ALLOWED_ORIGINS"   default:"*"`
		AllowedMethods   []string `envconfig:"PLINTH_CORS_ALLOWED_METHODS"   default:"GET,POST,PATCH,PUT,DELETE,OPTIONS"`
		AllowedHeaders   []string `envconfig:"PLINTH_CORS_ALLOWED_HEADERS"   default:"Origin,Accept,Accept-Language,Authorization,Content-Type,Content-Language,X-Requested-With,X-Request-Id"` //nolint:lll // struct tags can't be multiline
		ExposedHeaders   []string `envconfig:"PLINTH_CORS_EXPOSED_HEADERS"   default:"Link"`
		AllowCredentials bool     `envconfig:"PLINTH_CORS_ALLOW_CREDENTIALS" default:"true"`
		MaxAge           int      `envconfig:"PLINTH_CORS_MAX_AGE"           default:"300"`
	}

	Redis struct {
		Endpoint           string `envconfig:"PLINTH_REDIS_ENDPOINT"              default:"localhost:6379"`
		MaxRetries         int    `envconfig:"PLINTH_REDIS_MAX_RETRIES"           default:"3"`
		MinIdleConnections int    `envconfig:"PLINTH_REDIS_MIN_IDLE_CONNECTIONS"  default:"0"`
		Username           string `envconfig:"PLINTH_REDIS_USERNAME"`
		Password           string `envconfig:"PLINTH_REDIS_PASSWORD"`
	}

	// Secure defines http security parameters.
	Secure struct {
		AllowedHosts          []string          `envconfig:"PLINTH_HTTP_ALLOWED_HOSTS"`
		HostsProxyHeaders     []string          `envconfig:"PLINTH_HTTP_PROXY_HEADERS"`
		SSLRedirect           bool              `envconfig:"PLINTH_HTTP_SSL_REDIRECT"`
		SSLTemporaryRedirect  bool              `envconfig:"PLINTH_HTTP_SSL_TEMPORARY_REDIRECT"`
		SSLHost               string            `envconfig:"PLINTH_HTTP_SSL_HOST"`
		SSLProxyHeaders       map[string]string `envconfig:"PLINTH_HTTP_SSL_PROXY_HEADERS"`
		STSSeconds            int64             `envconfig:"PLINTH_HTTP_STS_SECONDS"`
		STSIncludeSubdomains  bool              `envconfig:"PLINTH_HTTP_STS_INCLUDE_SUBDOMAINS"`
		STSPreload            bool              `envconfig:"PLINTH_HTTP_STS_PRELOAD"`
		ForceSTSHeader        bool              `envconfig:"PLINTH_HTTP_STS_FORCE_HEADER"`
		BrowserXSSFilter      bool              `envconfig:"PLINTH_HTTP_BROWSER_XSS_FILTER"    default:"true"`
		FrameDeny             bool              `envconfig:"PLINTH_HTTP_FRAME_DENY"            default:"true"`
		ContentTypeNosniff    bool              `envconfig:"PLINTH_HTTP_CONTENT_TYPE_NO_SNIFF"`
		ContentSecurityPolicy string            `envconfig:"PLINTH_HTTP_CONTENT_SECURITY_POLICY"`
		ReferrerPolicy        string            `envconfig:"PLINTH_HTTP_REFERRER_POLICY"`
	}
}

// Process performs post-processing on the configuration after it has been loaded.
func (c *Config) Process() {
	// If the database driver is postgres and a full datasource is not already provided,
	// construct it from the individual parts.
	if c.Database.Driver == "postgres" && c.Database.Datasource == "" {
		c.Database.Datasource = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
			c.Database.SSLMode,
		)
	}
}
