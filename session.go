package xdp

import "sync"

// Session use session
type Session struct {
	Addr      string
	OpenID    string
	SessionID string

	sv *Server
}

// Send do send session data
func (s *Session) Send(data []byte) error {
	return s.sv.Send(s, data)
}

type sessionManager struct {
	sync.RWMutex

	sesses map[string]*Session
}

func newSessionManager() *sessionManager {
	sm := new(sessionManager)
	sm.sesses = make(map[string]*Session)
	return sm
}

func (sm *sessionManager) Get(sid string) *Session {
	sm.RLock()
	defer sm.RUnlock()

	sess := sm.sesses[sid]
	return sess
}

func (sm *sessionManager) Add(sess *Session) {
	sm.Lock()
	defer sm.Unlock()

	sm.sesses[sess.SessionID] = sess
}

func (sm *sessionManager) Del(sess *Session) {
	sm.Lock()
	defer sm.Unlock()

	delete(sm.sesses, sess.SessionID)
}
