package gorse

import "time"

type Feedback struct {
	FeedbackType string `json:"FeedbackType"`
	Username     string `json:"Username"`
	ItemUuid     string `json:"ItemUuid"`
	Timestamp    string `json:"Timestamp"`
}

type Feedbacks struct {
	Cursor   string     `json:"Cursor"`
	Feedback []Feedback `json:"Feedback"`
}

type ErrorMessage string

func (e ErrorMessage) Error() string {
	return string(e)
}

type RowAffected struct {
	RowAffected int `json:"RowAffected"`
}

type Score struct {
	Id    string  `json:"Id"`
	Score float64 `json:"Score"`
}

type User struct {
	UserId    string   `json:"UserId"`
	Labels    []string `json:"Labels"`
	Subscribe []string `json:"Subscribe"`
	Comment   string   `json:"Comment"`
}

type Users struct {
	Cursor string `json:"Cursor"`
	Users  []User `json:"Users"`
}

type UserPatch struct {
	Labels    []string
	Subscribe []string
	Comment   *string
}

type Item struct {
	ItemUuid   string   `json:"ItemUuid"`
	IsHidden   bool     `json:"IsHidden"`
	Labels     []string `json:"Labels"`
	Categories []string `json:"Categories"`
	Timestamp  string   `json:"Timestamp"`
	Comment    string   `json:"Comment"`
}

type Items struct {
	Cursor string `json:"Cursor"`
	Items  []Item `json:"Items"`
}

type ItemPatch struct {
	IsHidden   *bool
	Categories []string
	Timestamp  *time.Time
	Labels     []string
	Comment    *string
}

type Health struct {
	CacheStoreConnected bool   `json:"CacheStoreConnected"`
	CacheStoreError     string `json:"CacheStoreError"`
	DataStoreConnected  bool   `json:"DataStoreConnected"`
	DataStoreError      string `json:"DataStoreError"`
	Ready               bool   `json:"Ready"`
}
