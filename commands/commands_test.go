package commands

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmaiaroto/discfg/config"
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
	Convey("A ResponseObject with the current working config should be returned", t, func() {
		var opts = config.Options{StorageInterfaceName: "dynamodb", Version: "0.0.0", CfgName: "thisCfg"}
		w := Which(opts)
		So(w.Action, ShouldEqual, "which")
	})
}

func TestSetKey(t *testing.T) {
}
func TestGetKey(t *testing.T) {
}
func TestDeleteKey(t *testing.T) {
}
func TestInfo(t *testing.T) {
}
func TestExport(t *testing.T) {
}
