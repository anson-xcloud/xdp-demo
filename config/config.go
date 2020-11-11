package config

const (
	// XCloudAddr xcloud server addr
	// XCloudAddr = "http://127.0.0.1:4021" // debug addr
	XCloudAddr = "http://xcloud.singularityfuture.com.cn"
)

const (
	// APIAccessPoint xcloud api accesspoint url
	APIAccessPoint = "/app/accesspoint"

	// APIUserGetAccessPoint xcloud client get server tcp addr
	APIUserGetAccessPoint = "/user/accesspoint"
	APIUserLogin          = "/user/login"

	// APIClientXdpPrefix xcloud client transfer xdp url prefix
	APIClientXdpPrefix = "/xdp/"
)
