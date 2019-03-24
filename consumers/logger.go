package consumers

type NoopLogger struct{}

// Output allows us to implement the nsq.Logger interface
func (l *NoopLogger) Output(calldepth int, s string) error {
	return nil
}
