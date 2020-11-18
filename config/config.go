package config

var Env = &envRelease

type EnvConfig struct {
	Discription string

	XcloudAddr string
}

const (
	envDevDiscription     = "dev"
	envDebugDiscription   = "debug"
	envReleaseDiscription = "release"
)

var (
	envDev = EnvConfig{
		Discription: envDevDiscription,
		XcloudAddr:  "http://127.0.0.1:4021",
	}

	envDebug = EnvConfig{
		Discription: envDebugDiscription,
		XcloudAddr:  "http://127.0.0.1:4021",
	}

	envRelease = EnvConfig{
		Discription: envReleaseDiscription,
		XcloudAddr:  "http://xcloud.singularityfuture.com.cn",
	}
)

func SetEnv(env string) {
	switch env {
	case envDevDiscription:
		Env = &envDev
	case envDebugDiscription:
		Env = &envDebug
	case envReleaseDiscription:
		Env = &envRelease
	default:
		panic("invalid env")
	}
}

const (
	// APIAccessPoint xcloud api accesspoint url
	APIAccessPoint = "/app/accesspoint"

	// APIUserGetAccessPoint xcloud client get server tcp addr
	APIUserGetAccessPoint = "/user/accesspoint"
	APIUserLogin          = "/user/login"

	// APIClientXdpPrefix xcloud client transfer xdp url prefix
	APIClientXdpPrefix = "/xdp/"
)
