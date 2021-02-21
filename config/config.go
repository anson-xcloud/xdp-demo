package config

var Env = &envRelease

type EnvConfig struct {
	Discription string

	XcloudAddr string
}

const (
	EnvDevDiscription     = "dev"
	EnvDebugDiscription   = "debug"
	EnvReleaseDiscription = "release"
)

var (
	envDev = EnvConfig{
		Discription: EnvDevDiscription,
		XcloudAddr:  "http://localhost:31181",
	}

	envDebug = EnvConfig{
		Discription: EnvDebugDiscription,
		XcloudAddr:  "http://localhost:31181",
	}

	envRelease = EnvConfig{
		Discription: EnvReleaseDiscription,
		XcloudAddr:  "http://xcloud.singularityfuture.com.cn",
	}
)

func SetEnv(env string) {
	switch env {
	case EnvDevDiscription:
		Env = &envDev
	case EnvDebugDiscription:
		Env = &envDebug
	case EnvReleaseDiscription:
		Env = &envRelease
	default:
		panic("invalid env")
	}
}

const (
	// APIAccessPoint xcloud api accesspoint url
	APIAccessPoint = "/app/ap"

	// APIUserGetAccessPoint xcloud client get server tcp addr
	APIUserGetAccessPoint = "/user/ap"
	APIUserLogin          = "/user/login"

	// APIClientXdpPrefix xcloud client transfer xdp url prefix
	APIClientXdpPrefix = "/xdp/"
)
