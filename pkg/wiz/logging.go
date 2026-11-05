package wiz

// Logger defines the interface for logging messages at different levels.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// NopLogger is a logger implementation that does nothing.
type NopLogger struct {
}

func (n *NopLogger) Debugf(format string, args ...any) {}
func (n *NopLogger) Infof(format string, args ...any)  {}
func (n *NopLogger) Warnf(format string, args ...any)  {}
func (n *NopLogger) Errorf(format string, args ...any) {}
