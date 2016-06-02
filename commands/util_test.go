package commands

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestOut(t *testing.T) {
}
func TestGetDiscfgNameFromFile(t *testing.T) {
	Convey("When a .discfg file is present", t, func() {
		Convey("The current working config name should be returned", func() {
			_ = ioutil.WriteFile(".discfg", []byte("testcfg"), 0644)

			c := GetDiscfgNameFromFile()
			So(c, ShouldEqual, "testcfg")

			_ = os.Remove(".discfg")
		})
	})

	Convey("When a .discfg file is not present", t, func() {
		Convey("An empty string should be returned", func() {
			ce := GetDiscfgNameFromFile()
			So(ce, ShouldEqual, "")
		})
	})
}
