package logger

import "go.uber.org/zap"

func New() (*zap.Logger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	l.Info("replacing global logger")
	zap.ReplaceGlobals(l)
	return l, nil
}
