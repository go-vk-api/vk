package longpoll

// Stream represents a stream of events the VK User Longpoll API.
type Stream struct {
	lp *Longpoll

	TS int64

	Updates <-chan *Update
	Errors  <-chan error

	stop chan struct{}
}

// Start starts a channel for getting updates.
func (s *Stream) Start() error {
	updatesCh := make(chan *Update)
	errorsCh := make(chan error, 1)
	s.stop = make(chan struct{}, 1)

	s.Updates = updatesCh
	s.Errors = errorsCh

	go func() {
		defer func() {
			close(updatesCh)
			close(errorsCh)
			close(s.stop)
		}()

		for {
			select {
			case <-s.stop:
				return
			default:
			}

			updates, newTS, err := s.lp.Poll(s.TS)
			if err != nil {
				switch err {
				case ErrEventHistoryOutdated:
					s.TS = newTS
					continue
				case ErrKeyExpired, ErrUserInformationLost:
					newTS, err = s.lp.UpdateServer()
					if err != nil {
						errorsCh <- err
						return
					}
					s.TS = newTS
					continue
				default:
					errorsCh <- err
					return
				}
			}
			s.TS = newTS

			for _, update := range updates {
				select {
				case <-s.stop:
					return
				default:
				}

				updatesCh <- update
			}
		}
	}()

	return nil
}

// Stop stops the stream.
func (s *Stream) Stop() {
	s.stop <- struct{}{}
}
