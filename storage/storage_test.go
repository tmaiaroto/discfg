package storage

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tmaiaroto/discfg/config"
	"github.com/tmaiaroto/discfg/storage/mockdb"
	"testing"
)

func TestRegisterShipper(t *testing.T) {
	Convey("A new Shipper should be available for use once set", t, func() {
		RegisterShipper("mock", mockdb.MockShipper{})
		So(shippers["mock"], ShouldHaveSameTypeAs, mockdb.MockShipper{})
	})
}

func TestListShippers(t *testing.T) {
	Convey("Shippers should be returned", t, func() {
		shippers := ListShippers()
		So(shippers, ShouldNotBeEmpty)
	})
}

func TestCreateConfig(t *testing.T) {
	// Convey("A new Shipper should be available for use once set", t, func() {
	// RegisterShipper("mock", mockdb.MockShipper{})
	// So(shippers["mock"], ShouldHaveSameTypeAs, mockdb.MockShipper{})
	// })

	Convey("A valid Shipper must be used", t, func() {
		_, err := CreateConfig(config.Options{StorageInterfaceName: ""}, map[string]interface{}{})
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}

func TestDeleteConfig(t *testing.T) {
	Convey("A valid Shipper must be used", t, func() {
		_, err := DeleteConfig(config.Options{StorageInterfaceName: ""})
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}

func TestUpdateConfig(t *testing.T) {
	Convey("A valid Shipper must be used", t, func() {
		_, err := UpdateConfig(config.Options{StorageInterfaceName: ""}, map[string]interface{}{})
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}

func TestConfigState(t *testing.T) {
	Convey("A valid Shipper must be used", t, func() {
		_, err := ConfigState(config.Options{StorageInterfaceName: ""})
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}

func TestUpdate(t *testing.T) {
	Convey("A Shipper should set a key value and return the new node", t, func() {
		RegisterShipper("mock", mockdb.MockShipper{})
		opts := config.Options{
			StorageInterfaceName: "mock",
			Key:                  "testKey",
			Value:                []byte("testValue"),
		}
		resp, err := Update(opts)
		So(string(resp.Value.([]byte)), ShouldEqual, "testValue")
		So(err, ShouldHaveSameTypeAs, errors.New("error"))
	})

	Convey("A valid Shipper must be used", t, func() {
		opts := config.Options{
			StorageInterfaceName: "invalid",
			Key:                  "testKey",
			Value:                []byte("testValue"),
		}
		_, err := Update(opts)
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}

func TestGet(t *testing.T) {
	Convey("A Shipper should get a key value, returning the node", t, func() {
		RegisterShipper("mock", mockdb.MockShipper{})
		opts := config.Options{
			StorageInterfaceName: "mock",
			Key:                  "initial",
		}
		node, err := Get(opts)

		So(string(node.Value.([]byte)), ShouldEqual, "initial value for test")
		So(err, ShouldHaveSameTypeAs, errors.New("error"))
	})

	Convey("A valid Shipper must be used", t, func() {
		_, err := Get(config.Options{StorageInterfaceName: ""})
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}

func TestDelete(t *testing.T) {
	Convey("A Shipper should delete a key value and return the deleted node", t, func() {
		RegisterShipper("mock", mockdb.MockShipper{})
		opts := config.Options{
			StorageInterfaceName: "mock",
			Key:                  "initial_second",
		}
		node, err := Delete(opts)

		So(string(node.Value.([]byte)), ShouldEqual, "a second initial value for test")
		So(err, ShouldHaveSameTypeAs, errors.New("error"))

		So(mockdb.MockNodes["initial_second"], ShouldResemble, config.Node{})
	})

	Convey("A valid Shipper must be used", t, func() {
		_, err := Delete(config.Options{StorageInterfaceName: ""})
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}

func TestUpdateConfigVersion(t *testing.T) {
	Convey("A valid Shipper must be used", t, func() {
		err := UpdateConfigVersion(config.Options{StorageInterfaceName: ""})
		So(err.Error(), ShouldEqual, errMsgInvalidShipper)
	})
}
