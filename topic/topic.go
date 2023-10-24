package topic

import (
	"errors"
)

var (
	ErrUnsupported    = errors.New("unsupported topic")
	ErrNoTopicEnabled = errors.New("no topic enabled")
)

// Topic is a callback message that can be named.
type Topic interface {
	Name() string
}
