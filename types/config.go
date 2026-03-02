// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"time"

	"github.com/EolaFam1828/SoloDev/blob"
	"github.com/EolaFam1828/SoloDev/events"
	gitenum "github.com/EolaFam1828/SoloDev/git/enum"
	"github.com/EolaFam1828/SoloDev/lock"
	"github.com/EolaFam1828/SoloDev/pubsub"

	gossh "golang.org/x/crypto/ssh"
)

// Config stores the system configuration.
type Config struct {
	// InstanceID specifis the ID of the Harness instance.
	// NOTE: If the value is not provided the hostname of the machine is used.
	InstanceID string `envconfig:"SOLODEV_INSTANCE_ID"`

	Debug bool `envconfig:"SOLODEV_DEBUG"`
	Trace bool `envconfig:"SOLODEV_TRACE"`

	// GracefulShutdownTime defines the max time we wait when shutting down a server.
	// 5min should be enough for most git clones to complete.
	GracefulShutdownTime time.Duration `envconfig:"SOLODEV_GRACEFUL_SHUTDOWN_TIME" default:"300s"`

	UserSignupEnabled   bool `envconfig:"SOLODEV_USER_SIGNUP_ENABLED" default:"true"`
	NestedSpacesEnabled bool `envconfig:"SOLODEV_NESTED_SPACES_ENABLED" default:"false"`

	// PublicResourceCreationEnabled specifies whether a user can create publicly accessible resources.
	PublicResourceCreationEnabled bool `envconfig:"SOLODEV_PUBLIC_RESOURCE_CREATION_ENABLED" default:"true"`

	Profiler struct {
		Type        string `envconfig:"SOLODEV_PROFILER_TYPE"`
		ServiceName string `envconfig:"SOLODEV_PROFILER_SERVICE_NAME" default:"gitness"`
	}

	// URL defines the URLs via which the different parts of the service are reachable by.
	URL struct {
		// Base is used to generate external facing URLs in case they aren't provided explicitly.
		// Value is derived from Server.HTTP Config unless explicitly specified (e.g. http://localhost:3000).
		Base string `envconfig:"SOLODEV_URL_BASE"`

		// Git defines the external URL via which the GIT API is reachable.
		// NOTE: for routing to work properly, the request path & hostname reaching gitness
		// have to satisfy at least one of the following two conditions:
		// - Path ends with `/git`
		// - Hostname is different to API hostname
		// (this could be after proxy path / header rewrite).
		// Value is derived from Base unless explicitly specified (e.g. http://localhost:3000/git).
		Git string `envconfig:"SOLODEV_URL_GIT"`

		// GitSSH defines the external URL via which the GIT SSH server is reachable.
		// Value is derived from Base or SSH Config unless explicitly specified (e.g. ssh://localhost).
		GitSSH string `envconfig:"SOLODEV_URL_GIT_SSH"`

		// API defines the external URL via which the rest API is reachable.
		// NOTE: for routing to work properly, the request path reaching Harness has to end with `/api`
		// (this could be after proxy path rewrite).
		// Value is derived from Base unless explicitly specified (e.g. http://localhost:3000/api).
		API string `envconfig:"SOLODEV_URL_API"`

		// UI defines the external URL via which the UI is reachable.
		// Value is derived from Base unless explicitly specified (e.g. http://localhost:3000).
		UI string `envconfig:"SOLODEV_URL_UI"`

		// Internal defines the internal URL via which the service is reachable.
		// Value is derived from HTTP.Server unless explicitly specified (e.g. http://localhost:3000).
		Internal string `envconfig:"SOLODEV_URL_INTERNAL"`

		// Container is the endpoint that can be used by running container builds to communicate
		// with Harness (for example while performing a clone on a local repo).
		// host.docker.internal allows a running container to talk to services exposed on the host
		// (either running directly or via a port exposed in a docker container).
		// Value is derived from HTTP.Server unless explicitly specified (e.g. http://host.docker.internal:3000).
		Container string `envconfig:"SOLODEV_URL_CONTAINER"`

		// Registry is used as a base to generate external facing URLs.
		// Value is derived from HTTP.Server unless explicitly specified (e.g. http://host.docker.internal:3000).
		Registry string `envconfig:"SOLODEV_URL_REGISTRY"`
	}

	// Git defines the git configuration parameters
	Git struct {
		// Trace specifies whether git operations should be traces.
		// NOTE: Currently limited to 'push' operation until we move to internal command package.
		Trace bool `envconfig:"SOLODEV_GIT_TRACE"`
		// DefaultBranch specifies the default branch for new repositories.
		DefaultBranch string `envconfig:"SOLODEV_GIT_DEFAULTBRANCH" default:"main"`
		// Root specifies the directory containing git related data (e.g. repos, ...)
		Root string `envconfig:"SOLODEV_GIT_ROOT"`
		// TmpDir (optional) specifies the directory for temporary data (e.g. repo clones, ...)
		TmpDir string `envconfig:"SOLODEV_GIT_TMP_DIR"`
		// HookPath points to the binary used as git server hook.
		HookPath string `envconfig:"SOLODEV_GIT_HOOK_PATH"`

		// LastCommitCache holds configuration options for the last commit cache.
		LastCommitCache struct {
			// Mode determines where the cache will be. Valid values are "inmemory" (default), "redis" or "none".
			Mode gitenum.LastCommitCacheMode `envconfig:"SOLODEV_GIT_LAST_COMMIT_CACHE_MODE" default:"inmemory"`

			// Duration defines cache duration of last commit.
			Duration time.Duration `envconfig:"SOLODEV_GIT_LAST_COMMIT_CACHE_DURATION" default:"12h"`
		}
	}

	// Encrypter defines the parameters for the encrypter
	Encrypter struct {
		Secret       string `envconfig:"SOLODEV_ENCRYPTER_SECRET"` // key used for encryption
		MixedContent bool   `envconfig:"SOLODEV_ENCRYPTER_MIXED_CONTENT"`
	}

	// HTTP defines the http server configuration parameters
	HTTP struct {
		Port  int    `envconfig:"SOLODEV_HTTP_PORT" default:"3000"`
		Host  string `envconfig:"SOLODEV_HTTP_HOST"`
		Proto string `envconfig:"SOLODEV_HTTP_PROTO" default:"http"`
	}

	// Acme defines Acme configuration parameters.
	Acme struct {
		Enabled bool   `envconfig:"SOLODEV_ACME_ENABLED"`
		Endpont string `envconfig:"SOLODEV_ACME_ENDPOINT"`
		Email   bool   `envconfig:"SOLODEV_ACME_EMAIL"`
		Host    string `envconfig:"SOLODEV_ACME_HOST"`
	}

	SSH struct {
		Enable bool   `envconfig:"SOLODEV_SSH_ENABLE" default:"false"`
		Host   string `envconfig:"SOLODEV_SSH_HOST"`
		Port   int    `envconfig:"SOLODEV_SSH_PORT" default:"3022"`
		// DefaultUser holds value for generating urls {user}@host:path and force check
		// no other user can authenticate unless it is empty then any username is allowed
		DefaultUser             string   `envconfig:"SOLODEV_SSH_DEFAULT_USER" default:"git"`
		Ciphers                 []string `envconfig:"SOLODEV_SSH_CIPHERS"`
		KeyExchanges            []string `envconfig:"SOLODEV_SSH_KEY_EXCHANGES"`
		MACs                    []string `envconfig:"SOLODEV_SSH_MACS"`
		ServerHostKeys          []string `envconfig:"SOLODEV_SSH_HOST_KEYS"`
		TrustedUserCAKeys       []string `envconfig:"SOLODEV_SSH_TRUSTED_USER_CA_KEYS"`
		TrustedUserCAKeysFile   string   `envconfig:"SOLODEV_SSH_TRUSTED_USER_CA_KEYS_FILENAME"`
		TrustedUserCAKeysParsed []gossh.PublicKey
		KeepAliveInterval       time.Duration `envconfig:"SOLODEV_SSH_KEEP_ALIVE_INTERVAL" default:"5s"`
		ServerKeyPath           string        `envconfig:"SOLODEV_SSH_SERVER_KEY_PATH" default:"ssh/gitness.rsa"`
	}

	// CI defines configuration related to build executions.
	CI struct {
		ParallelWorkers int `envconfig:"SOLODEV_CI_PARALLEL_WORKERS" default:"2"`
		// PluginsZipURL is a pointer to a zip containing all the plugins schemas.
		// This could be a local path or an external location.
		//nolint:lll
		PluginsZipURL string `envconfig:"SOLODEV_CI_PLUGINS_ZIP_URL" default:"https://github.com/bradrydzewski/plugins/archive/refs/heads/master.zip"`

		// ContainerNetworks is a list of networks that all containers created as part of CI
		// should be attached to.
		// This can be needed when we don't want to use host.docker.internal (eg when a service mesh
		// or proxy is being used) and instead want all the containers to run on the same network as
		// the Harness container so that they can interact via the container name.
		// In that case, SOLODEV_URL_CONTAINER should also be changed
		// (eg to http://<gitness_container_name>:<port>).
		ContainerNetworks []string `envconfig:"SOLODEV_CI_CONTAINER_NETWORKS"`
	}

	// Database defines the database configuration parameters.
	Database struct {
		Driver     string `envconfig:"SOLODEV_DATABASE_DRIVER" default:"sqlite3"`
		Datasource string `envconfig:"SOLODEV_DATABASE_DATASOURCE" default:"database.sqlite3"`
	}

	// BlobStore defines the blob storage configuration parameters.
	BlobStore struct {
		// MaxFileSize defines the maximum size of files that can be uploaded (in bytes)
		MaxFileSize int64 `envconfig:"SOLODEV_BLOBSTORE_MAX_FILE_SIZE" default:"10485760"` // 10MB default
		// Provider is a name of blob storage service like filesystem or gcs or cloudflare
		Provider blob.Provider `envconfig:"SOLODEV_BLOBSTORE_PROVIDER" default:"filesystem"`
		// Bucket is a path to the directory where the files will be stored when using filesystem blob storage,
		// in case of gcs provider this will be the actual bucket where the images are stored.
		Bucket string `envconfig:"SOLODEV_BLOBSTORE_BUCKET"`

		// In case of GCS provider, this is expected to be the path to the service account key file.
		KeyPath string `envconfig:"SOLODEV_BLOBSTORE_KEY_PATH" default:""`

		// Email ID of the google service account that needs to be impersonated
		TargetPrincipal string `envconfig:"SOLODEV_BLOBSTORE_TARGET_PRINCIPAL" default:""`

		ImpersonationLifetime time.Duration `envconfig:"SOLODEV_BLOBSTORE_IMPERSONATION_LIFETIME" default:"12h"`
	}

	// Token defines token configuration parameters.
	Token struct {
		CookieName string        `envconfig:"SOLODEV_TOKEN_COOKIE_NAME" default:"token"`
		Expire     time.Duration `envconfig:"SOLODEV_TOKEN_EXPIRE" default:"720h"`
	}

	Logs struct {
		// S3 provides optional storage option for logs.
		S3 struct {
			Bucket    string `envconfig:"SOLODEV_LOGS_S3_BUCKET"`
			Prefix    string `envconfig:"SOLODEV_LOGS_S3_PREFIX"`
			Endpoint  string `envconfig:"SOLODEV_LOGS_S3_ENDPOINT"`
			PathStyle bool   `envconfig:"SOLODEV_LOGS_S3_PATH_STYLE"`
		}
	}

	// Cors defines http cors parameters
	Cors struct {
		AllowedOrigins   []string `envconfig:"SOLODEV_CORS_ALLOWED_ORIGINS"   default:"*"`
		AllowedMethods   []string `envconfig:"SOLODEV_CORS_ALLOWED_METHODS"   default:"GET,POST,PATCH,PUT,DELETE,OPTIONS"`
		AllowedHeaders   []string `envconfig:"SOLODEV_CORS_ALLOWED_HEADERS"   default:"Origin,Accept,Accept-Language,Authorization,Content-Type,Content-Language,X-Requested-With,X-Request-Id"` //nolint:lll // struct tags can't be multiline
		ExposedHeaders   []string `envconfig:"SOLODEV_CORS_EXPOSED_HEADERS"   default:"Link"`
		AllowCredentials bool     `envconfig:"SOLODEV_CORS_ALLOW_CREDENTIALS" default:"true"`
		MaxAge           int      `envconfig:"SOLODEV_CORS_MAX_AGE"           default:"300"`
	}

	// Secure defines http security parameters.
	Secure struct {
		AllowedHosts          []string          `envconfig:"SOLODEV_HTTP_ALLOWED_HOSTS"`
		HostsProxyHeaders     []string          `envconfig:"SOLODEV_HTTP_PROXY_HEADERS"`
		SSLRedirect           bool              `envconfig:"SOLODEV_HTTP_SSL_REDIRECT"`
		SSLTemporaryRedirect  bool              `envconfig:"SOLODEV_HTTP_SSL_TEMPORARY_REDIRECT"`
		SSLHost               string            `envconfig:"SOLODEV_HTTP_SSL_HOST"`
		SSLProxyHeaders       map[string]string `envconfig:"SOLODEV_HTTP_SSL_PROXY_HEADERS"`
		STSSeconds            int64             `envconfig:"SOLODEV_HTTP_STS_SECONDS"`
		STSIncludeSubdomains  bool              `envconfig:"SOLODEV_HTTP_STS_INCLUDE_SUBDOMAINS"`
		STSPreload            bool              `envconfig:"SOLODEV_HTTP_STS_PRELOAD"`
		ForceSTSHeader        bool              `envconfig:"SOLODEV_HTTP_STS_FORCE_HEADER"`
		BrowserXSSFilter      bool              `envconfig:"SOLODEV_HTTP_BROWSER_XSS_FILTER"    default:"true"`
		FrameDeny             bool              `envconfig:"SOLODEV_HTTP_FRAME_DENY"            default:"true"`
		ContentTypeNosniff    bool              `envconfig:"SOLODEV_HTTP_CONTENT_TYPE_NO_SNIFF"`
		ContentSecurityPolicy string            `envconfig:"SOLODEV_HTTP_CONTENT_SECURITY_POLICY"`
		ReferrerPolicy        string            `envconfig:"SOLODEV_HTTP_REFERRER_POLICY"`
	}

	Principal struct {
		// System defines the principal information used to create the system service.
		System struct {
			UID         string `envconfig:"SOLODEV_PRINCIPAL_SYSTEM_UID"          default:"gitness"`
			DisplayName string `envconfig:"SOLODEV_PRINCIPAL_SYSTEM_DISPLAY_NAME" default:"Gitness"`
			Email       string `envconfig:"SOLODEV_PRINCIPAL_SYSTEM_EMAIL"        default:"system@gitness.io"`
		}
		// Pipeline defines the principal information used to create the pipeline service.
		Pipeline struct {
			UID         string `envconfig:"SOLODEV_PRINCIPAL_PIPELINE_UID"          default:"pipeline"`
			DisplayName string `envconfig:"SOLODEV_PRINCIPAL_PIPELINE_DISPLAY_NAME" default:"Gitness Pipeline"`
			Email       string `envconfig:"SOLODEV_PRINCIPAL_PIPELINE_EMAIL"        default:"pipeline@gitness.io"`
		}

		// Gitspace defines the principal information used to create the gitspace service.
		Gitspace struct {
			UID         string `envconfig:"SOLODEV_PRINCIPAL_GITSPACE_UID"          default:"gitspace"`
			DisplayName string `envconfig:"SOLODEV_PRINCIPAL_GITSPACE_DISPLAY_NAME" default:"Gitness Gitspace"`
			Email       string `envconfig:"SOLODEV_PRINCIPAL_GITSPACE_EMAIL"        default:"gitspace@gitness.io"`
		}

		// Admin defines the principal information used to create the admin user.
		// NOTE: The admin user is only auto-created in case a password and an email is provided.
		Admin struct {
			UID         string `envconfig:"SOLODEV_PRINCIPAL_ADMIN_UID"           default:"admin"`
			DisplayName string `envconfig:"SOLODEV_PRINCIPAL_ADMIN_DISPLAY_NAME"  default:"Administrator"`
			Email       string `envconfig:"SOLODEV_PRINCIPAL_ADMIN_EMAIL"`    // No default email
			Password    string `envconfig:"SOLODEV_PRINCIPAL_ADMIN_PASSWORD"` // No default password
		}
	}

	Redis struct {
		Endpoint           string `envconfig:"SOLODEV_REDIS_ENDPOINT"              default:"localhost:6379"`
		MaxRetries         int    `envconfig:"SOLODEV_REDIS_MAX_RETRIES"           default:"3"`
		MinIdleConnections int    `envconfig:"SOLODEV_REDIS_MIN_IDLE_CONNECTIONS"  default:"0"`
		Password           string `envconfig:"SOLODEV_REDIS_PASSWORD"`
		SentinelMode       bool   `envconfig:"SOLODEV_REDIS_USE_SENTINEL"          default:"false"`
		SentinelMaster     string `envconfig:"SOLODEV_REDIS_SENTINEL_MASTER"`
		SentinelEndpoint   string `envconfig:"SOLODEV_REDIS_SENTINEL_ENDPOINT"`
	}

	Events struct {
		Mode                  events.Mode `envconfig:"SOLODEV_EVENTS_MODE"                     default:"inmemory"`
		Namespace             string      `envconfig:"SOLODEV_EVENTS_NAMESPACE"                default:"gitness"`
		MaxStreamLength       int64       `envconfig:"SOLODEV_EVENTS_MAX_STREAM_LENGTH"        default:"10000"`
		ApproxMaxStreamLength bool        `envconfig:"SOLODEV_EVENTS_APPROX_MAX_STREAM_LENGTH" default:"true"`
	}

	Lock struct {
		// Provider is a name of distributed lock service like redis, memory, file etc...
		Provider      lock.Provider `envconfig:"SOLODEV_LOCK_PROVIDER"          default:"inmemory"`
		Expiry        time.Duration `envconfig:"SOLODEV_LOCK_EXPIRE"            default:"8s"`
		Tries         int           `envconfig:"SOLODEV_LOCK_TRIES"             default:"8"`
		RetryDelay    time.Duration `envconfig:"SOLODEV_LOCK_RETRY_DELAY"       default:"250ms"`
		DriftFactor   float64       `envconfig:"SOLODEV_LOCK_DRIFT_FACTOR"      default:"0.01"`
		TimeoutFactor float64       `envconfig:"SOLODEV_LOCK_TIMEOUT_FACTOR"    default:"0.25"`
		// AppNamespace is just service app prefix to avoid conflicts on key definition
		AppNamespace string `envconfig:"SOLODEV_LOCK_APP_NAMESPACE"     default:"gitness"`
		// DefaultNamespace is when mutex doesn't specify custom namespace for their keys
		DefaultNamespace string `envconfig:"SOLODEV_LOCK_DEFAULT_NAMESPACE" default:"default"`
	}

	PubSub struct {
		// Provider is a name of distributed lock service like redis, memory, file etc...
		Provider pubsub.Provider `envconfig:"SOLODEV_PUBSUB_PROVIDER"                default:"inmemory"`
		// AppNamespace is just service app prefix to avoid conflicts on channel definition
		AppNamespace string `envconfig:"SOLODEV_PUBSUB_APP_NAMESPACE"                default:"gitness"`
		// DefaultNamespace is custom namespace for their channels
		DefaultNamespace string        `envconfig:"SOLODEV_PUBSUB_DEFAULT_NAMESPACE" default:"default"`
		HealthInterval   time.Duration `envconfig:"SOLODEV_PUBSUB_HEALTH_INTERVAL"   default:"3s"`
		SendTimeout      time.Duration `envconfig:"SOLODEV_PUBSUB_SEND_TIMEOUT"      default:"60s"`
		ChannelSize      int           `envconfig:"SOLODEV_PUBSUB_CHANNEL_SIZE"      default:"100"`
	}

	BackgroundJobs struct {
		// MaxRunning is maximum number of jobs that can be running at once.
		MaxRunning int `envconfig:"SOLODEV_JOBS_MAX_RUNNING" default:"10"`

		// RetentionTime is the duration after which non-recurring,
		// finished and failed jobs will be purged from the DB.
		RetentionTime time.Duration `envconfig:"SOLODEV_JOBS_RETENTION_TIME" default:"120h"` // 5 days
	}

	Webhook struct {
		// UserAgentIdentity specifies the identity used for the user agent header
		// IMPORTANT: do not include version.
		UserAgentIdentity string `envconfig:"SOLODEV_WEBHOOK_USER_AGENT_IDENTITY" default:"Gitness"`
		// HeaderIdentity specifies the identity used for headers in webhook calls (e.g. X-Gitness-Trigger, ...).
		// NOTE: If no value is provided, the UserAgentIdentity will be used.
		HeaderIdentity      string `envconfig:"SOLODEV_WEBHOOK_HEADER_IDENTITY"`
		Concurrency         int    `envconfig:"SOLODEV_WEBHOOK_CONCURRENCY" default:"4"`
		MaxRetries          int    `envconfig:"SOLODEV_WEBHOOK_MAX_RETRIES" default:"3"`
		AllowPrivateNetwork bool   `envconfig:"SOLODEV_WEBHOOK_ALLOW_PRIVATE_NETWORK" default:"false"`
		AllowLoopback       bool   `envconfig:"SOLODEV_WEBHOOK_ALLOW_LOOPBACK" default:"false"`
		// RetentionTime is the duration after which webhook executions will be purged from the DB.
		RetentionTime  time.Duration `envconfig:"SOLODEV_WEBHOOK_RETENTION_TIME" default:"168h"` // 7 days
		InternalSecret string        `envconfig:"SOLODEV_WEBHOOK_INTERNAL_SECRET"`
	}

	Trigger struct {
		Concurrency int `envconfig:"SOLODEV_TRIGGER_CONCURRENCY" default:"4"`
		MaxRetries  int `envconfig:"SOLODEV_TRIGGER_MAX_RETRIES" default:"3"`
	}

	Branch struct {
		Concurrency int `envconfig:"SOLODEV_BRANCH_CONCURRENCY" default:"4"`
		MaxRetries  int `envconfig:"SOLODEV_BRANCH_MAX_RETRIES" default:"3"`
	}

	Metric struct {
		Enabled  bool   `envconfig:"SOLODEV_METRIC_ENABLED" default:"true"`
		Endpoint string `envconfig:"SOLODEV_METRIC_ENDPOINT" default:"https://stats.drone.ci/api/v1/gitness"`
		Token    string `envconfig:"SOLODEV_METRIC_TOKEN"`

		// PostHogEndpoint is URL to the PostHog service
		PostHogEndpoint string `envconfig:"SOLODEV_METRIC_POSTHOG_ENDPOINT" default:"https://us.i.posthog.com"`
		// PostHogProjectAPIKey (starts with "phc_") is public (can be exposed in frontend) token used to submit events.
		PostHogProjectAPIKey string `envconfig:"SOLODEV_METRIC_POSTHOG_PROJECT_APIKEY"`
		// PostHogPersonalAPIKey (starts with "phx_") is sensitive. It's used to access private access points.
		// It's not required for submitting events.
		PostHogPersonalAPIKey string `envconfig:"SOLODEV_METRIC_POSTHOG_PERSONAL_APIKEY"`
	}

	RepoSize struct {
		Enabled     bool          `envconfig:"SOLODEV_REPO_SIZE_ENABLED" default:"true"`
		CRON        string        `envconfig:"SOLODEV_REPO_SIZE_CRON" default:"0 0 * * *"`
		MaxDuration time.Duration `envconfig:"SOLODEV_REPO_SIZE_MAX_DURATION" default:"15m"`
		NumWorkers  int           `envconfig:"SOLODEV_REPO_SIZE_NUM_WORKERS" default:"5"`
	}

	Githook struct {
		DisableAuth bool `envconfig:"SOLODEV_GITHOOK_DISABLE_AUTH" default:"false"`
	}

	CodeOwners struct {
		FilePaths []string `envconfig:"SOLODEV_CODEOWNERS_FILEPATH" default:"CODEOWNERS,.harness/CODEOWNERS"`
	}

	SMTP struct {
		Host     string `envconfig:"SOLODEV_SMTP_HOST"`
		Port     int    `envconfig:"SOLODEV_SMTP_PORT"`
		Username string `envconfig:"SOLODEV_SMTP_USERNAME"`
		Password string `envconfig:"SOLODEV_SMTP_PASSWORD"`
		FromMail string `envconfig:"SOLODEV_SMTP_FROM_MAIL"`
		Insecure bool   `envconfig:"SOLODEV_SMTP_INSECURE"`
	}

	Notification struct {
		MaxRetries  int `envconfig:"SOLODEV_NOTIFICATION_MAX_RETRIES" default:"3"`
		Concurrency int `envconfig:"SOLODEV_NOTIFICATION_CONCURRENCY" default:"4"`
	}

	KeywordSearch struct {
		Concurrency int `envconfig:"SOLODEV_KEYWORD_SEARCH_CONCURRENCY" default:"4"`
		MaxRetries  int `envconfig:"SOLODEV_KEYWORD_SEARCH_MAX_RETRIES" default:"3"`
	}

	Repos struct {
		// DeletedRetentionTime is the duration after which deleted repositories will be purged.
		DeletedRetentionTime time.Duration `envconfig:"SOLODEV_REPOS_DELETED_RETENTION_TIME" default:"2160h"` // 90 days
	}

	Docker struct {
		// Host sets the url to the docker server.
		Host string `envconfig:"SOLODEV_DOCKER_HOST"`
		// APIVersion sets the version of the API to reach, leave empty for latest.
		APIVersion string `envconfig:"SOLODEV_DOCKER_API_VERSION"`
		// CertPath sets the path to load the TLS certificates from.
		CertPath string `envconfig:"SOLODEV_DOCKER_CERT_PATH"`
		// TLSVerify enables or disables TLS verification, off by default.
		TLSVerify string `envconfig:"SOLODEV_DOCKER_TLS_VERIFY"`
		// MachineHostName is the public host name of the machine on which the Docker.Host is running.
		// If not set, it parses the host from the URL.Base (e.g. localhost from http://localhost:3000).
		MachineHostName string `envconfig:"SOLODEV_DOCKER_MACHINE_HOST_NAME"`
	}

	IDE struct {
		VSCodeWeb struct {
			// Port is the port on which the VSCode Web will be accessible.
			Port int `envconfig:"SOLODEV_IDE_VSCODEWEB_PORT" default:"8089"`
		}

		VSCode struct {
			// Port is the port on which the SSH server for VSCode will be accessible.
			Port       int    `envconfig:"SOLODEV_IDE_VSCODE_PORT" default:"8088"`
			PluginName string `envconfig:"SOLODEV_IDE_VSCODE_Plugin_Name" default:"harness-inc.oss-gitspaces"`
		}

		Cursor struct {
			// Port is the port on which the SSH server for Cursor will be accessible.
			Port int `envconfig:"SOLODEV_IDE_CURSOR_PORT" default:"8098"`
		}

		Windsurf struct {
			// Port is the port on which the SSH server for Windsurf will be accessible.
			Port int `envconfig:"SOLODEV_IDE_WINDSURF_PORT" default:"8099"`
		}

		Intellij struct {
			// Port is the port on which the SSH server for IntelliJ will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_INTELLIJ_PORT" default:"8090"`
		}

		Goland struct {
			// Port is the port on which the SSH server for Goland will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_GOLAND_PORT" default:"8091"`
		}

		PyCharm struct {
			// Port is the port on which the SSH server for PyCharm will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_PYCHARM_PORT" default:"8092"`
		}

		WebStorm struct {
			// Port is the port on which the SSH server for WebStorm will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_WEBSTORM_PORT" default:"8093"`
		}

		CLion struct {
			// Port is the port on which the SSH server for CLion will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_CLION_PORT" default:"8094"`
		}

		PHPStorm struct {
			// Port is the port on which the SSH server for PHPStorm will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_PHPSTORM_PORT" default:"8095"`
		}

		RubyMine struct {
			// Port is the port on which the SSH server for RubyMine will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_RUBYMINE_PORT" default:"8096"`
		}

		Rider struct {
			// Port is the port on which the SSH server for Rider will be accessible
			Port int `envconfig:"CDE_MANAGER_GITSPACE_IDE_RIDER_PORT" default:"8097"`
		}
	}

	Gitspace struct {
		// DefaultBaseImage is used to create the Gitspace when no devcontainer.json is absent or doesn't have image.
		DefaultBaseImage string `envconfig:"SOLODEV_GITSPACE_DEFAULT_BASE_IMAGE" default:"mcr.microsoft.com/devcontainers/base:dev-ubuntu-24.04"` //nolint:lll

		Enable bool `envconfig:"SOLODEV_GITSPACE_ENABLE" default:"false"`

		AgentPort int `envconfig:"SOLODEV_GITSPACE_AGENT_PORT" default:"8083"`

		InfraTimeoutInMins int `envconfig:"SOLODEV_INFRA_TIMEOUT_IN_MINS" default:"60"`

		BusyActionInMins int `envconfig:"SOLODEV_BUSY_ACTION_IN_MINS" default:"15"`

		Events struct {
			Concurrency   int `envconfig:"SOLODEV_GITSPACE_EVENTS_CONCURRENCY" default:"4"`
			MaxRetries    int `envconfig:"SOLODEV_GITSPACE_EVENTS_MAX_RETRIES" default:"3"`
			TimeoutInMins int `envconfig:"SOLODEV_GITSPACE_EVENTS_TIMEOUT_IN_MINS" default:"45"`
		}
	}

	UI struct {
		ShowPlugin bool `envconfig:"SOLODEV_UI_SHOW_PLUGIN" default:"false"`
	}

	Registry struct {
		Enable  bool `envconfig:"SOLODEV_REGISTRY_ENABLED" default:"true"`
		Storage struct {
			// StorageType defines the type of storage to use for the registry. Options are: `filesystem`, `s3aws`, `gcs`
			StorageType string `envconfig:"SOLODEV_REGISTRY_STORAGE_TYPE" default:"filesystem"`

			// FileSystemStorage defines the configuration for the filesystem storage if StorageType is `filesystem`.
			FileSystemStorage struct {
				MaxThreads    int    `envconfig:"SOLODEV_REGISTRY_FILESYSTEM_MAX_THREADS" default:"100"`
				RootDirectory string `envconfig:"SOLODEV_REGISTRY_FILESYSTEM_ROOT_DIRECTORY"`
			}

			// S3Storage defines the configuration for the S3 storage if StorageType is `s3aws`.
			S3Storage struct {
				AccessKey                   string `envconfig:"SOLODEV_REGISTRY_S3_ACCESS_KEY"`
				SecretKey                   string `envconfig:"SOLODEV_REGISTRY_S3_SECRET_KEY"`
				Region                      string `envconfig:"SOLODEV_REGISTRY_S3_REGION"`
				RegionEndpoint              string `envconfig:"SOLODEV_REGISTRY_S3_REGION_ENDPOINT"`
				ForcePathStyle              bool   `envconfig:"SOLODEV_REGISTRY_S3_FORCE_PATH_STYLE" default:"true"`
				Accelerate                  bool   `envconfig:"SOLODEV_REGISTRY_S3_ACCELERATED" default:"false"`
				Bucket                      string `envconfig:"SOLODEV_REGISTRY_S3_BUCKET"`
				Encrypt                     bool   `envconfig:"SOLODEV_REGISTRY_S3_ENCRYPT" default:"false"`
				KeyID                       string `envconfig:"SOLODEV_REGISTRY_S3_KEY_ID"`
				Secure                      bool   `envconfig:"SOLODEV_REGISTRY_S3_SECURE" default:"true"`
				V4Auth                      bool   `envconfig:"SOLODEV_REGISTRY_S3_V4_AUTH" default:"true"`
				ChunkSize                   int    `envconfig:"SOLODEV_REGISTRY_S3_CHUNK_SIZE" default:"10485760"`
				MultipartCopyChunkSize      int    `envconfig:"SOLODEV_REGISTRY_S3_MULTIPART_COPY_CHUNK_SIZE" default:"33554432"`
				MultipartCopyMaxConcurrency int    `envconfig:"SOLODEV_REGISTRY_S3_MULTIPART_COPY_MAX_CONCURRENCY" default:"100"`
				MultipartCopyThresholdSize  int    `envconfig:"SOLODEV_REGISTRY_S3_MULTIPART_COPY_THRESHOLD_SIZE" default:"33554432"` //nolint:lll
				RootDirectory               string `envconfig:"SOLODEV_REGISTRY_S3_ROOT_DIRECTORY"`
				UseDualStack                bool   `envconfig:"SOLODEV_REGISTRY_S3_USE_DUAL_STACK" default:"false"`
				LogLevel                    string `envconfig:"SOLODEV_REGISTRY_S3_LOG_LEVEL" default:"info"`
				Delete                      bool   `envconfig:"SOLODEV_REGISTRY_S3_DELETE_ENABLED" default:"true"`
				Redirect                    bool   `envconfig:"SOLODEV_REGISTRY_S3_STORAGE_REDIRECT" default:"false"`
				Provider                    string `envconfig:"SOLODEV_REGISTRY_S3_PROVIDER" default:"cloudflare"`
			}

			// GCSStorage defines the configuration for the GCS storage if StorageType is `gcs`.
			// Authentication is handled via workload identity (google.DefaultTokenSource).
			GCSStorage struct {
				Bucket string `envconfig:"SOLODEV_REGISTRY_GCS_BUCKET"`
			}
		}

		HTTP struct {
			// SOLODEV_REGISTRY_HTTP_SECRET is used to encrypt the upload session details during docker push.
			// If not provided, a random secret will be generated. This may cause problems with uploads if multiple
			// registries are behind a load-balancer
			Secret string `envconfig:"SOLODEV_REGISTRY_HTTP_SECRET"`

			RelativeURL bool `envconfig:"SOLODEV_OCI_RELATIVE_URL" default:"false"`
		}

		//nolint:lll
		GarbageCollection struct {
			Enabled                     bool          `envconfig:"SOLODEV_REGISTRY_GARBAGE_COLLECTION_ENABLED" default:"false"`
			NoIdleBackoff               bool          `envconfig:"SOLODEV_REGISTRY_GARBAGE_COLLECTION_NO_IDLE_BACKOFF" default:"false"`
			MaxBackoffDuration          time.Duration `envconfig:"SOLODEV_REGISTRY_GARBAGE_COLLECTION_MAX_BACKOFF_DURATION" default:"10m"`
			InitialIntervalDuration     time.Duration `envconfig:"SOLODEV_REGISTRY_GARBAGE_COLLECTION_INITIAL_INTERVAL_DURATION" default:"5s"`     //nolint:lll
			TransactionTimeoutDuration  time.Duration `envconfig:"SOLODEV_REGISTRY_GARBAGE_COLLECTION_TRANSACTION_TIMEOUT_DURATION" default:"10s"` //nolint:lll
			BlobsStorageTimeoutDuration time.Duration `envconfig:"SOLODEV_REGISTRY_GARBAGE_COLLECTION_BLOB_STORAGE_TIMEOUT_DURATION" default:"5s"` //nolint:lll
		}
		SetupDetailsAuthHeaderPrefix string `envconfig:"SETUP_DETAILS_AUTH_PREFIX" default:"Authorization: Bearer"`

		PostProcessing struct {
			Concurrency   int  `envconfig:"SOLODEV_REGISTRY_POST_PROCESSING_CONCURRENCY" default:"4"`
			MaxRetries    int  `envconfig:"SOLODEV_REGISTRY_POST_PROCESSING_MAX_RETRIES" default:"3"`
			AllowLoopback bool `envconfig:"SOLODEV_REGISTRY_POST_PROCESSING_ALLOW_LOOPBACK" default:"false"`
		}
	}

	Auth struct {
		AnonymousUserSecret string `envconfig:"SOLODEV_ANONYMOUS_USER_SECRET"`
	}

	Instrumentation struct {
		Enable bool   `envconfig:"SOLODEV_INSTRUMENTATION_ENABLE" default:"false"`
		Cron   string `envconfig:"SOLODEV_INSTRUMENTATION_CRON" default:"0 0 * * *"`
	}

	UsageMetrics struct {
		Enabled       bool          `envconfig:"SOLODEV_USAGE_METRICS_ENABLED" default:"false"`
		FlushInterval time.Duration `envconfig:"SOLODEV_USAGE_METRICS_FLUSH_INTERVAL" default:"1m"`
	}

	Development struct {
		UISourceOverride string `envconfig:"SOLODEV_DEVELOPMENT_UI_SOURCE_OVERRIDE"`
	}

	SecurityScan struct {
		Enabled      bool          `envconfig:"SOLODEV_SECURITY_SCAN_ENABLED" default:"true"`
		SemgrepPath  string        `envconfig:"SOLODEV_SECURITY_SCAN_SEMGREP_PATH"`
		GitleaksPath string        `envconfig:"SOLODEV_SECURITY_SCAN_GITLEAKS_PATH"`
		TrivyPath    string        `envconfig:"SOLODEV_SECURITY_SCAN_TRIVY_PATH"`
		SemgrepRules string        `envconfig:"SOLODEV_SECURITY_SCAN_SEMGREP_RULES" default:"auto"`
		MaxDuration  time.Duration `envconfig:"SOLODEV_SECURITY_SCAN_MAX_DURATION" default:"5m"`
		WorkDir      string        `envconfig:"SOLODEV_SECURITY_SCAN_WORK_DIR"`
	}

	AIRemediation struct {
		Enabled         bool    `envconfig:"SOLODEV_AI_REMEDIATION_ENABLED" default:"false"`
		Provider        string  `envconfig:"SOLODEV_AI_REMEDIATION_PROVIDER"` // anthropic, openai, gemini, ollama
		APIKey          string  `envconfig:"SOLODEV_AI_REMEDIATION_API_KEY"`
		Model           string  `envconfig:"SOLODEV_AI_REMEDIATION_MODEL"`
		MaxTokens       int     `envconfig:"SOLODEV_AI_REMEDIATION_MAX_TOKENS" default:"4096"`
		Temperature     float64 `envconfig:"SOLODEV_AI_REMEDIATION_TEMPERATURE" default:"0.2"`
		CreateFixBranch bool    `envconfig:"SOLODEV_AI_REMEDIATION_CREATE_FIX_BRANCH" default:"false"`
	}
}
