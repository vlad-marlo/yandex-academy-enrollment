package logger

import "go.uber.org/zap"

// New prepares new zap logger and replacing global logger with newly initialized.
func New(opts ...zap.Option) (*zap.Logger, error) {
	l, _ := zap.NewProduction(opts...)

	l.Info("replacing global logger")
	zap.ReplaceGlobals(l)
	return l, nil
}
