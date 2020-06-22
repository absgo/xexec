package xexec

import "strconv"

func getSleepCmdLine(sleepSec int) []string {
	return []string{"./bin/sleep", strconv.Itoa(sleepSec)}
}
