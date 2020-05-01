package mock

type ILog interface {
	Error(err error, v ...interface{})
	Debugf(format string, v ...interface{})
	Printf(format string, v ...interface{})
}
