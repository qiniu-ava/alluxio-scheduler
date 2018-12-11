package util

import "os"

func ExitWithMessage(no int, msg string) {
	if no != 0 {
		switch no / 100 {
		case 1:
			Logger.Fatalf(msg)
		default:
			Logger.Errorf(msg)
		}
	}
	os.Exit(no)
}
