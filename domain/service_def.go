package domain

type Follower interface {
	GetCurrentBlockHeight() (int64,error)
}