package longpoll

import (
	"errors"
)

var (
	ErrEventHistoryOutdated = errors.New("event history outdated")
	ErrKeyExpired           = errors.New("key expired")
	ErrUserInformationLost  = errors.New("user information lost")
	ErrInvalidVersion       = errors.New("invalid version")
)
