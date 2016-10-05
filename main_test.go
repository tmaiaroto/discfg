package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSetOptsFromArgs(t *testing.T) {
	Convey("When 1 argument is passed", t, func() {
		Convey("A CfgName option should be set, a Key and Value should not", func() {
			setOptsFromArgs([]string{"testCfg"})
			So(Options.CfgName, ShouldEqual, "testCfg")
			So(string(Options.Key), ShouldEqual, "")
			So(string(Options.Value), ShouldEqual, "")
		})
	})

	Convey("When 2 arguments are passed", t, func() {
		Convey("A CfgName and Key option should be set", func() {
			setOptsFromArgs([]string{"testCfg", "some key"})
			So(Options.CfgName, ShouldEqual, "testCfg")
			So(string(Options.Value), ShouldEqual, "some key")
		})
	})

	Convey("When 3 arguments are passed", t, func() {
		Convey("CfgName, Key, and Value options should be set", func() {
			setOptsFromArgs([]string{"testCfg", "some key", "some value"})
			So(Options.CfgName, ShouldEqual, "testCfg")
			So(Options.Key, ShouldEqual, "some key")
			So(string(Options.Value), ShouldEqual, "some value")
		})
	})
}
