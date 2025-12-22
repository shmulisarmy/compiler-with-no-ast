package assert

func Assert(considtion bool, messages ...string) {
	if len(messages) == 0 {
		messages = append(messages, "Assertion failed")
	}
	if !considtion {
		panic(messages[0])
	}

}
