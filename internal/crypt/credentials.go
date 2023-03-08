package crypt

import (
	"strings"
	"time"

	"github.com/tagus/crypt/internal/utils/sets"
)

type Credential struct {
	Id            string   `json:"id"`
	Service       string   `json:"service"`
	Domains       []string `json:"domains"`
	Email         string   `json:"email"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	Description   string   `json:"description"`
	Details       []Detail `json:"details"`
	Tags          []string `json:"tags"`
	CreatedAt     int64    `json:"created_at"`
	UpdatedAt     int64    `json:"updated_at"`
	AccessedAt    *int64   `json:"accessed_at"`
	AccessedCount int      `json:"accessed_count"`
	Fingerprint   string   `json:"fingerprint"`
}

// Detail models an arbitrary key value pair
type Detail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c *Credential) String() string {
	return c.Service
}

func (c *Credential) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

func (c *Credential) GetUpdatedAt() time.Time {
	return time.Unix(c.UpdatedAt, 0)
}

func (c *Credential) GetAccessedAt() *time.Time {
	if c.AccessedAt == nil {
		return nil
	}
	ts := time.Unix(*c.AccessedAt, 0)
	return &ts
}

func (c *Credential) GetShortFingerprint() string {
	if len(c.Fingerprint) < 8 {
		return c.Fingerprint
	}
	return c.Fingerprint[:8]
}

func (c *Credential) Merge(cred *Credential) *Credential {
	dirty := false
	if cred.Service != "" {
		val := strings.TrimSpace(cred.Service)
		if val != "" {
			dirty = true
			c.Service = val
		}
	}
	if cred.Email != "" {
		val := strings.TrimSpace(cred.Email)
		if val != "" {
			dirty = true
			c.Email = val
		}
	}
	if cred.Username != "" {
		val := strings.TrimSpace(cred.Username)
		if val != "" {
			dirty = true
			c.Username = val
		}
	}
	if cred.Password != "" {
		val := strings.TrimSpace(cred.Password)
		if val != "" {
			dirty = true
			c.Password = val
		}
	}
	if cred.Description != "" {
		val := strings.TrimSpace(cred.Description)
		if val != "" {
			dirty = true
			c.Description = val
		}
	}
	if len(cred.Tags) > 0 {
		ss := sets.NewString(c.Tags...)
		ss.Add(cred.Tags...)
		if ss.Len() > len(c.Tags) {
			c.Tags = ss.AsSlice()
			dirty = true
		}
	}
	if dirty {
		c.UpdatedAt = time.Now().Unix()
	}
	return c
}
