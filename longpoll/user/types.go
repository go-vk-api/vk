package longpoll

// Event codes.
const (
	EventNewMessage = 4
)

// NewMessage represents the information an event EventNewMessage contains.
type NewMessage struct {
	ID          int64
	Flags       int64
	PeerID      int64
	Timestamp   int64
	Text        string
	Attachments map[string]string
	RandomID    int64
}
