package longpoll

// Option is a configuration option to initialize a longpoll.
type Option func(*Longpoll) error

// WithMode overrides the longpoll mode with the specified one.
func WithMode(m Mode) Option {
	return func(lp *Longpoll) error {
		lp.Mode = m

		return nil
	}
}
