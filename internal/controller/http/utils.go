package http

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

// respond writes data to response writer.
//
// Passing nil data will write text status code to w.
func (srv *Controller) respond(w http.ResponseWriter, code int, data interface{}, fields ...zap.Field) {
	var lvl zapcore.Level
	switch {
	case code >= http.StatusInternalServerError:
		lvl = zap.ErrorLevel
	case code >= http.StatusBadRequest:
		lvl = zap.WarnLevel
	default:
		lvl = zap.DebugLevel
	}
	w.Header().Set("content-type", "application/json")

	if data == nil {
		data = http.StatusText(code)
	}

	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fields = append(fields, zap.Error(err))
	}

	if len(fields) > 0 {
		srv.log.Log(lvl, "respond", fields...)
	}
}
