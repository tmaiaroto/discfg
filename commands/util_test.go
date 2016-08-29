package commands

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/storage"
	"github.com/tmaiaroto/discfg/storage/mockdb"
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

func TestFormatJSONValue(t *testing.T) {
	Convey("Should handle basic string values", t, func() {
		storage.RegisterShipper("mock", mockdb.MockShipper{})
		var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg", Key: "initial"}
		r := GetKey(opts)
		rFormatted := FormatJSONValue(r)
		So(rFormatted.Item.Value.(string), ShouldEqual, "initial value for test")
	})

	// Not yet
	// Convey("Should handle base64 encoded string values", t, func() {
	// 	storage.RegisterShipper("mock", mockdb.MockShipper{})
	// 	var opts = config.Options{StorageInterfaceName: "mock", Version: "0.0.0", CfgName: "mockcfg", Key: "encoded"}
	// 	r := GetKey(opts)
	// 	rFormatted := FormatJSONValue(r)
	// 	mapValue := map[string]interface{}{"updated": "friday"}
	// 	So(rFormatted.Item.Value.(map[string]interface{}), ShouldResemble, mapValue)
	// })
}
