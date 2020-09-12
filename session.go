package xdp

import "sync"

type Session struct {
	Addr      string
	OpenID    string
	SessionID string

	sv *xdpServer
}

func (s *Session) Send(data []byte) error {
	return s.sv.Send(s, data)
}

type sessionManager struct {
	sync.RWMutex

	sesses map[string]*Session
}

func (sm *sessionManager) Get(sid string) *Session {
	sm.RLock()
	defer sm.RUnlock()

	sess := sm.sesses[sid]
	return sess
}
