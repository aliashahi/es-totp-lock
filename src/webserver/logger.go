package webserver

import (
	"fmt"
	"time"
)

func Logger(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	s = fmt.Sprintf("%s : %s", time.Now().Format(time.DateTime), s)
	fmt.Println(s)
	var lastId int64 = 0
	if len(logs) != 0 {
		lastId = logs[len(logs)-1].ID
	}
	logs = append(logs, &Log{
		ID:      lastId + 1,
		Message: s,
	})
}