package ksync

import (
	"testing"

	"github.com/doctordesh/check"
)

func TestConfig(t *testing.T) {
	errs := Config{}.Validate()

	check.Assert(t, len(errs) == 0)

	errs = Config{Targets: Targets{Target{}}}.Validate()
	check.Assert(t, len(errs) == 1)
	check.EqualsWithMessage(t, "missing key 'name' for target 0", errs[0].Error(), "error message is not correct")

	errs = Config{Targets: Targets{Target{Name: "foobar"}}}.Validate()
	check.Assert(t, len(errs) == 4)
	check.Equals(t, "target 'foobar' is missing key 'source'", errs[0].Error())
	check.Equals(t, "target 'foobar' is missing key 'destination.user'", errs[1].Error())
	check.Equals(t, "target 'foobar' is missing key 'destination.host'", errs[2].Error())
	check.Equals(t, "target 'foobar' is missing key 'destination.path'", errs[3].Error())

	errs = Config{Targets: Targets{Target{}, Target{}}}.Validate()
	check.Assert(t, len(errs) == 2)
}
