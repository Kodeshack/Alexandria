package models

import (
	"net/http"
	"time"

	"alexandria.app/crypto"
)

const (
	// SessionCookieName is the name of the cookie which is used for authenticating a user.
	SessionCookieName = "AlexandriaUserSession"
	// SessionContextKey is the key used in the request's context to save and retrieve a session.
	SessionContextKey = "UserSession"
	// Cookie duration in days.
	cookieDuration = 30
)

// A Session represents a logged-in user.
type Session struct {
	User      *User
	sessionID string
	createdAt time.Time
}

// Cookie creates a new cookie that can be passed along with the HTTP response.
// It is highly recommended to always set isHTTPS to true.
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

// NewSession creates a new session for a specific user, most likeley for the user of the current request.
// Generates a cryptographically secure random id which is not guaranteed to be unique.
func NewSession(u *User) *Session {
	id, _ := crypto.GetRandomString(32)
	return &Session{
		User:      u,
		sessionID: id,
		createdAt: time.Now(),
	}
}

// SessionStorage holds all currently active sessions.
// Note: The storage is not persistent across restarts.
type SessionStorage struct {
	sessions  []*Session
	SpawnedAt time.Time
}

// AddSession adds a session to the session storage.
func (sstrg *SessionStorage) AddSession(sess *Session) {
	sstrg.sessions = append(sstrg.sessions, sess)
}

// RemoveSession removes a session to the session storage.
func (sstrg *SessionStorage) RemoveSession(sess *Session) {
	for i, s := range sstrg.sessions {
		if s.sessionID == sess.sessionID {
			sstrg.sessions = append(sstrg.sessions[:i], sstrg.sessions[i+1:]...)
			break
		}
	}

}

// RemoveSessionsForUser removes all sessions associated with a user.
// As a user can have several sessions across multiple devices when deleting a user
// it is best to invalidate and remove all sessions associated with that user.
func (sstrg *SessionStorage) RemoveSessionsForUser(id uint32) {
	for i, s := range sstrg.sessions {
		if s.User.ID == id {
			sstrg.sessions = append(sstrg.sessions[:i], sstrg.sessions[i+1:]...)
		}
	}
}

// GetSession retrieves the session associated with the token from the storage.
// To prevent data leaking only the token is sent by cookie, not the entire session.
func (sstrg *SessionStorage) GetSession(token string) *Session {
	for _, s := range sstrg.sessions {
		if s.sessionID == token {
			return s
		}
	}

	return nil
}

// GetRequestSession extracts the session from a request object.
func GetRequestSession(r *http.Request) *Session {
	val := r.Context().Value(SessionContextKey)

	if val != nil {
		return val.(*Session)
	}

	return nil
}

// GetRequestUser extracts the (logged-in) user from a request object.
func GetRequestUser(r *http.Request) *User {
	if session := GetRequestSession(r); session != nil {
		return session.User
	}

	return nil
}

// NewSessionStorage creates a new SessionStorage struct.
// There should really only be one of these at any given point.
func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		sessions:  []*Session{},
		SpawnedAt: time.Now(),
	}
}

// RemoveSessionCookie creates a new cookie with the same name as the session cookie
// with an expire date so far in the past that devices will delete the cookie thus invalidating the session.
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
