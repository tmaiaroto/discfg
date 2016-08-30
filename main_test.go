package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSetOptsFromArgs(t *testing.T) {
	Convey("When 1 argument is passed", t, func() {
		Convey("A Key option should be set, a Value should not", func() {
			setOptsFromArgs([]string{"some key"})
			So(Options.Key, ShouldEqual, "some key")
			So(string(Options.Value), ShouldEqual, "")
		})
	})

	Convey("When 2 arguments are passed", t, func() {
		Convey("A Key and Value option should be set", func() {
			setOptsFromArgs([]string{"some key", "some value"})
			So(Options.Key, ShouldEqual, "some key")
			So(string(Options.Value), ShouldEqual, "some value")
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
