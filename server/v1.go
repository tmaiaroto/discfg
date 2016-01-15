// API Version 1
package main

import (
	"encoding/json"
	//"encoding/base64"
	"github.com/labstack/echo"
	"github.com/tmaiaroto/discfg/commands"
	"github.com/tmaiaroto/discfg/config"
	//"github.com/ugorji/go/codec" // <-- may eventually be used as a different output format. have to see what echo supports.
	//"log"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Set the routes for V1 API
func v1Routes(e *echo.Echo) {
	e.Put("/v1/:name/key/:key", v1SetKey)
	e.Get("/v1/:name/key/:key", v1GetKey)
	e.Delete("/v1/:name/key/:key", v1DeleteKey)

	e.Put("/v1/cfg/:name", v1CreateCfg)
	e.Delete("/v1/cfg/:name", v1DeleteCfg)
	e.Options("/v1/cfg/:name", v1OptionsCfg)
}

// Gets a key from discfg
func v1GetKey(c *echo.Context) error {
	Options.CfgName = c.Param("name")
	Options.Key = c.Param("key")
	resp := commands.GetKey(Options)
	// Since this option is not needed for anything else, it's not held on the Options struct.
	contentType := c.Query("type")

	// This is very awesome. Very interesting possibilties now.
	switch contentType {
	case "text", "text/plain", "string":
		return c.String(http.StatusOK, string(resp.Node.Value.([]byte)))
		break
	// This one is going to be interesting. Weird? Bad practice? I don't know, but I dig it and it starts giving me wild ideas.
	case "html", "text/html":
		return c.HTML(http.StatusOK, string(resp.Node.Value.([]byte)))
		break
	// TODO:
	//case "jsonp":
	//break
	case "json", "application/json":
		resp = formatJsonValue(resp)
		break
	default:
		resp = formatJsonValue(resp)
		break
	}
	// default response
	return c.JSON(http.StatusOK, resp)
}

// Sets the Node Value (an interface{}) as a map[string]interface{} (from []byte which is how it's stored - at least for now)
// so that it can be converted to JSON in the Echo response. If it can't be represented in a map, then it'll be set as a string.
// For example, a string was set as the value, we can still represent that in JSON. However, if an image was stored...Then it's
// going to look ugly. It won't be a base64 string, it'll be the string representation of the binary data. Which apparently Chrome
// will render if given...But still. Not so hot. The user should know what they are setting and getting though and this should
// still technically return JSON with a usable value. Valid JSON at that. Just with some funny looking characters =)
func formatJsonValue(resp config.ResponseObject) config.ResponseObject {
	// Don't attempt to Unmarshal or anything if the Value is empty. We wouldn't want to create a panic now.
	if resp.Node.Value == nil {
		return resp
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(resp.Node.Value.([]byte), &dat); err != nil {
		resp.Node.Value = string(resp.Node.Value.([]byte))
	} else {
		resp.Node.Value = dat
	}
	return resp
}

// Sets a key in discfg
func v1SetKey(c *echo.Context) error {
	Options.CfgName = c.Param("name")
	Options.Key = c.Param("key")
	resp := config.ResponseObject{
		Action: "set",
	}

	// Allow the value to be passed via querystring param.
	Options.Value = []byte(c.Query("value"))

	// Overwrite that if the request body passes a value that can be read, preferring that.
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		// If reading the body failed and we don't have a value from the querystring parameter
		// then we have a (potential) problem.
		//
		// Some data stores may be ok with an empty key value. DynamoDB is not. It will only
		// return a ValidationException error. Plus, even if it was allowed, it would really
		// confuse the user. Some random error reading the body of a request and poof, the data
		// vanishes? That'd be terrible UX.
		// log.Println(err)
		resp.Error = err.Error()
		resp.Message = "Something went wrong reading the body of the request."
		// resp.ErrorCode = 500 <-- TODO: I need to come up with discfg specific error codes. Keep as const somewhere.
		//	return c.JSON(http.StatusOK, resp)
		//	Or maybe return an HTTP status message... This is outside discfg's concern. It's not an error message/code
		//	that would ever be seen from the CLI, right? Or maybe it would. Maybe a more generic, "error parsing key value" ...
	} else if len(b) > 0 {
		Options.Value = b
	}

	resp = commands.SetKey(Options)

	return c.JSON(http.StatusOK, resp)
}

// Deletes a key in discfg
func v1DeleteKey(c *echo.Context) error {
	Options.CfgName = c.Param("name")
	Options.Key = c.Param("key")
	return c.JSON(http.StatusOK, commands.DeleteKey(Options))
}

// Creates a new configuration
func v1CreateCfg(c *echo.Context) error {
	Options.CfgName = c.Param("name")

	// Optionally set write capacity units (for DynamoDB, defaults to 1)
	wu, err := strconv.ParseInt(c.Query("writeUnits"), 10, 64)
	if err == nil {
		Options.Storage.WriteCapacityUnits = wu
	}
	// Optionally set read capacity units (for DynamoDB, defaults to 2)
	ru, err := strconv.ParseInt(c.Query("readUnits"), 10, 64)
	if err == nil {
		Options.Storage.ReadCapacityUnits = ru
	}

	return c.JSON(http.StatusOK, commands.CreateCfg(Options))
}

// Deletes a configuration
func v1DeleteCfg(c *echo.Context) error {
	Options.CfgName = c.Param("name")
	return c.JSON(http.StatusOK, commands.DeleteCfg(Options))
}

// Gets/sets options for a configuration
func v1OptionsCfg(c *echo.Context) error {
	Options.CfgName = c.Param("name")
	// If nothing that would change a configuration's options/settings are passed,
	// then return commands.Info() ... Otherwise, call a new commands.UpdateInfo() function...
	// That then adjusts settings and returns Info() anyway.
	// Or maybe commands.Info() can be a getter and setter...Eh... maybe.
	//
	// OR. Make another route and handler for adjusting configuration settings using PUT.
	// https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html ...says OPTIONS can contain
	// some data in the body request, though defines no use for it. So it's up in the air.
	//
	// It might then also be neat to hav an OPTIONS for other routes.
	// These responses would all be handled here and be just for self-documenting the API.
	// OPTIONS request on a key path for example would return information about how to
	// update or delete that specific key. On /v1/:name/key it might then just be
	// instructions for creating a key.
	//
	// Need to think this over. Sleep on it.
	// YES. All of it. OPTIONS all the things!
	//
	// There will only be so many options that can be changed after creation.
	// For example, if we wanted to gzip the data ... That couldn't be turned on/off after creation.
	// Doing so would mean all existing values would need to change.
	// Though for gzip, I think it might be a nice additional value on the record that specifies the data is gzipped.
	// So it can be done on a key by key basis. So if it was sent gzipped, then when it is read there'll be a simple
	// if/then that unzips. So users can decide how the data is stored. ... Think about that too. I'm not sure they
	// should care or have to say on each key. It likely should just be done always. Or on a per configuration basis.
	// So again, that kinda stuff must be set up front upon creation and can not change.
	//
	// So the options being set here are going to be quite limited.
	// Though I don't want to handle them all here... EAch adapter should.
	//
	// So take a map of options. map[string]interface{}
	// Then pass that to each storage adapter. Then each adapter can figure out and choose which things to use/set.
	// Which means the gzip setting, and any other setting that can't be changed after creation, is simply ignored
	// at that point even if passed.
	//
	//
	return c.JSON(http.StatusOK, commands.Info(Options, map[string]interface{}{}))
}
