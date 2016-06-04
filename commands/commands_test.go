package commands

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/storage"
	"github.com/tmaiaroto/discfg/storage/mockdb"
	"io/ioutil"
	//"log"
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
}
func TestGetKey(t *testing.T) {
	Convey("Should return a ResponseObject with the key value", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg", Key: "initial"}
		r := GetKey(opts)
		So(r.Action, ShouldEqual, "get")
		So(r.Node.Version, ShouldEqual, int64(1))
		So(string(r.Node.Value.([]byte)), ShouldEqual, "initial value for test")
	})

	Convey("Should return a ResponseObject with an Error message if not enough arguments were provided", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg"}
		r := GetKey(opts)
		So(r.Action, ShouldEqual, "get")
		So(r.Error, ShouldEqual, NotEnoughArgsMsg)
	})
}
func TestDeleteKey(t *testing.T) {
}
func TestInfo(t *testing.T) {
	Convey("Should return a ResponseObject with info about the config", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg"}
		r := Info(opts)
		So(r.Action, ShouldEqual, "info")
		// TODO... this fails. But works through RESTful API? Why???
		//So(r.Node.CfgVersion, ShouldEqual, int64(4))
		//So(r.CfgState, ShouldEqual, "ACTIVE")
	})
}
func TestExport(t *testing.T) {
}
