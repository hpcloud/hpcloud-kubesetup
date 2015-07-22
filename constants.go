package main

//OpenStack well-defined command line names
const (
	AuthURL    = "os-auth-url"
	TenantID   = "os-tenant-id"
	TenantName = "os-tenant-name"
	Username   = "os-username"
	Password   = "os-password"
	RegionName = "os-region-name"
	AuthToken  = "os-auth-token"
	CACert     = "os-cacert"
)

// OpenStack well-defined environment variable names
const (
	AuthURLEnv    = "OS_AUTH_URL"
	TenantIDEnv   = "OS_TENANT_ID"
	TenantNameEnv = "OS_TENANT_NAME"
	UsernameEnv   = "OS_USERNAME"
	PasswordEnv   = "OS_PASSWORD"
	RegionNameEnv = "OS_REGION_NAME"
	AuthTokenEnv  = "OS_AUTH_TOKEN"
	CACertEnv     = "OS_CACERT"
)

// OpenStack service types
const (
	Identity = "identity"
	Compute  = "compute"
	Network  = "network"
	Image    = "image"
)

// Commandline constants
const (
	Config            = "config"
	DefaultConfig     = "kubesetup.yml"
	Debug             = "debug"
	SkipSSLValidation = "skip-ssl-validation"
	Install           = "install"
	Status            = "status"
	Uninstall         = "uninstall"
)

// LogStringFormat1 and others are formats used to write to the log.
const (
	LogStringFormat1 string = "%-30s\n"
	LogStringFormat2 string = "%-30s - %s\n"
	LogStringFormat3 string = "%-30s - %s %s\n"
	LogStringFormat4 string = "%-30s - %s %s %s\n"
	LogValueFormat1  string = "%-30s - %s %v\n"
)
