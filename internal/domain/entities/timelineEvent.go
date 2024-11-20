package entities

type EventType int

const (
	_ EventType = iota
	TimelineAccessed
	PostCreated
	PostDeleted
)

type TimelineEvent struct {
	EventType EventType
	Posts     []*Post
}
