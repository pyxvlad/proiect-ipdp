package types

type AccountID int64
type BookID int64
type ProgressID int64
type PublisherID int64

const (
	InvalidBookID     BookID     = 0
	InvalidAccountID  AccountID  = 0
	InvalidProgressID ProgressID = 0
	InvalidPublisherID PublisherID = 0
)

type Status string

const (
	StatusToBeRead   Status = "to be read"
	StatusInProgress Status = "in progress"
	StatusRead       Status = "read"
	StatusDropped    Status = "dropped"
	StatusUncertain  Status = "uncertain"
)

func (s Status) String() string {
	return string(s)
}
