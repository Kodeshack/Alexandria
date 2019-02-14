package models

import (
	"net/http"
	"time"

	"alexandria.app/crypto"
)

const (
	SessionCookieName = "AlexandriaUserSession"
	SessionContextKey = "UserSession"
	// Cookie duration in days.
	cookieDuration = 30
)

type Session struct {
	User      *User
	sessionID string
	createdAt time.Time
}

func (s *Session) Cookie(isHTTPS bool) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    s.sessionID,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   isHTTPS,
		HttpOnly: true,
		Expires:  time.Now().Add(cookieDuration * 24 * time.Hour),
	}
}

func NewSession(u *User) *Session {
	id, _ := crypto.GetRandomString(32)
	return &Session{
		User:      u,
		sessionID: id,
		createdAt: time.Now(),
	}
}

type SessionStorage struct {
	sessions  []*Session
	SpawnedAt time.Time
}

func (sstrg *SessionStorage) AddSession(sess *Session) {
	sstrg.sessions = append(sstrg.sessions, sess)
}

func (sstrg *SessionStorage) RemoveSession(sess *Session) {
	for i, s := range sstrg.sessions {
		if s.sessionID == sess.sessionID {
			sstrg.sessions = append(sstrg.sessions[:i], sstrg.sessions[i+1:]...)
			break
		}
	}

}

func (sstrg *SessionStorage) RemoveSessionsForUser(id uint32) {
	for i, s := range sstrg.sessions {
		if s.User.ID == id {
			sstrg.sessions = append(sstrg.sessions[:i], sstrg.sessions[i+1:]...)
		}
	}
}

func (sstrg *SessionStorage) GetSession(token string) *Session {
	for _, s := range sstrg.sessions {
		if s.sessionID == token {
			return s
		}
	}

	return nil
}

func GetRequestSession(r *http.Request) *Session {
	val := r.Context().Value(SessionContextKey)

	if val != nil {
		return val.(*Session)
	}

	return nil
}

func GetRequestUser(r *http.Request) *User {
	if session := GetRequestSession(r); session != nil {
		return session.User
	}

	return nil
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		sessions:  []*Session{},
		SpawnedAt: time.Now(),
	}
}

func RemoveSessionCookie(isHTTPS bool) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		SameSite: http.SameSiteStrictMode,
		Secure:   isHTTPS,
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
