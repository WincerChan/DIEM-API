package tools

func CheckError(err error, shouldCrash bool) {
	if err == nil {
		return
	}
	if shouldCrash {
		panic(err)
	}
	println(err)
}
