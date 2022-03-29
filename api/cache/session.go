package cache

import (
	"api/models"
	"fmt"
	"time"
)

func (c *Cache) AddSessionToCache(session models.Session) bool {
	err := c.Client.Add(fmt.Sprintf("session_%v", session.SessionID), session, 15 * time.Minute)
	if err != nil {
		return false
	}

	return true
}

func (c *Cache) SetSession(newSession models.Session) bool {
	session := c.GetSession(newSession.SessionID)
	if session != nil {
		c.Client.Set(fmt.Sprintf("session_%v", session.SessionID), newSession, 15 * time.Minute)
		return true
	}

	return false
}

func (c *Cache) GetSession(session string) *models.Session {
	storedSession, ok := c.Client.Get(fmt.Sprintf("session_%v", session))
	if !ok {
		return nil
	}

	s := storedSession.(models.Session)
	return &s
}

func (c *Cache) DeleteSession(session string) {
	c.Client.Delete(fmt.Sprintf("session_%v", session))
}