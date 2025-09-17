package model

import "time"

type URL struct {
    Short       string     `json:"short"`
    Destination string     `json:"destination"`
    CreatedAt   time.Time  `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type ClickEvent struct {
    Short       string    `json:"short"`
    Destination string    `json:"destination"`
    Timestamp   time.Time `json:"timestamp"`
    ClientIP    string    `json:"client_ip"`
    UA          string    `json:"user_agent"`
    Referrer    string    `json:"referrer,omitempty"`
}
