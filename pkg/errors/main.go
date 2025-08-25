package errors

import(
	"github.com/sirupsen/logrus"
)
var logger = logrus.New()

var (
	Debug  = logger.Debug
	Debugf = logger.Debugf
	Info   = logger.Info
	Infof  = logger.Infof
	Warn   = logger.Warn
	Warnf  = logger.Warnf
	Error  = logger.Error
	Errorf = logger.Errorf
	Fatal  = logger.Fatal
	Fatalf = logger.Fatalf
)