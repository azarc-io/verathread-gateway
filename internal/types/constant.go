package types

import "time"

const (
	ShellConfigurationUpdatedSubject = "gateway.shell.v1.configuration.rebuilt"
	TargetURLKey                     = "targetUrl"
	AppNameKey                       = "appName"
	KeySpaceExpiryChannel            = "__key*__:expired"
	KeepAliveKeySpacePrefix          = "app:keepalive"
)

var (
	CacheShards         = 3
	CacheCleanupFreq    = time.Second * 30
	TargetCacheDuration = time.Minute * 2
	KeepAliveTTL        = time.Second * 10
)
