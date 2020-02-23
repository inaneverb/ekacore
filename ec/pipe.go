package ec

//
func Pipe(pipedFunctions ...func() EC) EC {

	for _, pipedFunction := range pipedFunctions {
		if pipedFunction != nil {
			if e := pipedFunction(); e != EOK {
				return e
			}
		}
	}
	return EOK
}

//
func Pipecxt(pipedFunctions ...func() ECXT) ECXT {

	for _, pipedFunction := range pipedFunctions {
		if pipedFunction != nil {
			if e := pipedFunction(); !e.EOK() {
				return e
			}
		}
	}
	return EOK.ECXT()
}
