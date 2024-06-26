package types

type AccountID int64
type BookID int64
type ProgressID int64
type PublisherID int64
type SeriesID int64
type CollectionID int64

const (
	InvalidBookID       BookID       = 0
	InvalidAccountID    AccountID    = 0
	InvalidProgressID   ProgressID   = 0
	InvalidPublisherID  PublisherID  = 0
	InvalidSeriesID     SeriesID     = 0
	InvalidCollectionID CollectionID = 0
)

type Status string

const (
	StatusToBeRead   Status = "To be read"
	StatusInProgress Status = "In progress"
	StatusRead       Status = "Read"
	StatusDropped    Status = "Dropped"
	StatusUncertain  Status = "Uncertain"
)

func GetAllStatuses() []Status {
	return []Status{
		StatusToBeRead,
		StatusInProgress,
		StatusRead,
		StatusUncertain,
		StatusDropped,
	}
}

func (s Status) String() string {
	return string(s)
}
