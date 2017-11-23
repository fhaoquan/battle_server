package room_1v1

import "github.com/sirupsen/logrus"

func (me *room) room_log_dbg(args ...interface{}){
	logrus.Debug(args...)
}
func (me *room) room_log_inf(args ...interface{}){
	logrus.Info(args...)
}
func (me *room) room_log_err(args ...interface{}){
	logrus.Error(args...)
}
