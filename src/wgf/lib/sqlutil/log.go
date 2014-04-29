package sqlutil

var Logger func(log string)

func log(body string) {
	if nil!=Logger {
		Logger(body)
	}
}
