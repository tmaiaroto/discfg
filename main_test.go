package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSetOptsFromArgs(t *testing.T) {
	Convey("When 1 argument is passed", t, func() {
		Convey("A Key option should be set, a Value should not", func() {
			o := setOptsFromArgs([]string{"some key"})
			So(o.Key, ShouldEqual, "some key")
			So(string(o.Value), ShouldEqual, "")
		})
	})

	Convey("When 2 arguments are passed", t, func() {
		Convey("A Key and Value option should be set", func() {
			o := setOptsFromArgs([]string{"some key", "some value"})
			So(o.Key, ShouldEqual, "some key")
			So(string(o.Value), ShouldEqual, "some value")
		})
	})

	Convey("When 3 arguments are passed", t, func() {
		Convey("CfgName, Key, and Value options should be set", func() {
			o := setOptsFromArgs([]string{"testCfg", "some key", "some value"})
			So(o.CfgName, ShouldEqual, "testCfg")
			So(o.Key, ShouldEqual, "some key")
			So(string(o.Value), ShouldEqual, "some value")
		})
	})
}
