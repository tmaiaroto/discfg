package commands

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/storage"
	"github.com/tmaiaroto/discfg/storage/mockdb"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCreateCfg(t *testing.T) {
}

func TestDeleteCfg(t *testing.T) {
}
func TestUpdateCfg(t *testing.T) {
}
func TestUse(t *testing.T) {
}

func TestWhich(t *testing.T) {
	Convey("Should return an error when no .discfg exists", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0"}
		r := Which(opts)
		So(r.Action, ShouldEqual, "which")
		So(r.Error, ShouldEqual, NoCurrentWorkingCfgMsg)
	})

	Convey("Should return a ResponseObject with the current working config", t, func() {
		_ = ioutil.WriteFile(".discfg", []byte("testcfg"), 0644)

		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0"}
		r := Which(opts)
		So(r.Action, ShouldEqual, "which")
		So(r.CurrentDiscfg, ShouldEqual, "testcfg")

		_ = os.Remove(".discfg")
	})
}

func TestSetKey(t *testing.T) {
	Convey("Should return a ResponseObject with an Error message if no value was provided", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0"}
		r := SetKey(opts)
		So(r.Action, ShouldEqual, "set")
		So(r.Error, ShouldEqual, ValueRequiredMsg)
	})

	Convey("Should return a ResponseObject with an Error message if no key name was provided", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg", Value: []byte("test")}
		r := SetKey(opts)
		So(r.Action, ShouldEqual, "set")
		So(r.Error, ShouldEqual, MissingKeyNameMsg)
	})

	Convey("Should return a ResponseObject with an Error message if no config name was provided", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", Value: []byte("test"), Key: "test"}
		r := SetKey(opts)
		So(r.Action, ShouldEqual, "set")
		So(r.Error, ShouldEqual, MissingCfgNameMsg)
	})
}

func TestGetKey(t *testing.T) {
	Convey("Should return a ResponseObject with the key value", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg", Key: "initial"}
		r := GetKey(opts)
		So(r.Action, ShouldEqual, "get")
		So(r.Item.Version, ShouldEqual, int64(1))
		So(string(r.Item.Value.([]byte)), ShouldEqual, "initial value for test")
	})

	Convey("Should return a ResponseObject with the key value", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg", Key: "encoded"}
		r := GetKey(opts)
		So(r.Action, ShouldEqual, "get")
		So(r.Item.Version, ShouldEqual, int64(1))
		//So(string(r.Item.Value.([]byte)), ShouldEqual, "initial value for test")
		log.Println(FormatJSONValue(r).Item.Value)
	})

	Convey("Should return a ResponseObject with an Error message if no key name was provided", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg"}
		r := GetKey(opts)
		So(r.Action, ShouldEqual, "get")
		So(r.Error, ShouldEqual, MissingKeyNameMsg)
	})
}

func TestDeleteKey(t *testing.T) {
	Convey("Should return a ResponseObject with an Error message if not enough arguments were provided", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0"}
		r := DeleteKey(opts)
		So(r.Action, ShouldEqual, "delete")
		So(r.Error, ShouldEqual, NotEnoughArgsMsg)
	})
}

func TestInfo(t *testing.T) {
	Convey("Should return a ResponseObject with info about the config", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg"}
		r := Info(opts)
		So(r.Action, ShouldEqual, "info")
		So(r.CfgVersion, ShouldEqual, int64(4))
		So(r.CfgState, ShouldEqual, "ACTIVE")
		So(r.CfgModifiedNanoseconds, ShouldEqual, int64(1464675792991825937))
		So(r.CfgModified, ShouldEqual, int64(1464675792))
		So(r.CfgModifiedParsed, ShouldEqual, "2016-05-30T23:23:12-07:00")
	})

	Convey("Should return a ResponseObject with an Error message if not enough arguments were provided", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0"}
		r := Info(opts)
		So(r.Action, ShouldEqual, "info")
		So(r.Error, ShouldEqual, NotEnoughArgsMsg)
	})
}

func TestExport(t *testing.T) {
}
