package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStatusText(t *testing.T) {
	Convey("Should return status text string for passed code int", t, func() {
		s := StatusText(100)
		So(s, ShouldEqual, "Continue")
	})
}
