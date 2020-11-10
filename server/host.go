package server

import apipb "github.com/anson-xcloud/xdp-demo/api"

type HostSetting struct {
	// MaxRequestPerSecond int
	// MaxRequestPerDay    int

	AllowUserApis   map[string]struct{}
	AllowServerApis map[string]struct{}
}

func (h *HostSetting) isAllow(remote *apipb.Remote, api string) bool {
	var allowed map[string]struct{}
	if remote.Sid != "" {
		allowed = h.AllowUserApis
	} else {
		allowed = h.AllowServerApis
	}

	_, ok := allowed[api]
	return ok
}

func (h *HostSetting) isServerAllow(api string) bool {
	_, ok := h.AllowServerApis[api]
	return ok
}
