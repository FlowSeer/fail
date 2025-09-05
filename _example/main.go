package main

import "github.com/FlowSeer/fail"

func main() {
	err := fail.New().
		Tag("tag1", "tag2").
		Attribute("key1", "value1").
		Attribute("key2", "value2").
		Cause(fail.Msg("cause")).
		Msg("test message")

	fail.PrintPretty(err)
}
