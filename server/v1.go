// API Version 1
package main

// import (
// 	"encoding/json"
// 	//"encoding/base64"
// 	"github.com/labstack/echo"
// 	"github.com/tmaiaroto/discfg/commands"
// 	"github.com/tmaiaroto/discfg/config"
// 	//"github.com/ugorji/go/codec" // <-- may eventually be used as a different output format. have to see what echo supports.
// 	"io/ioutil"
// 	//"log"
// 	"net/http"
// )

// // Set the routes for V1 API
// func v1Routes(e *echo.Echo) {
// 	e.Put("/v1/:name/keys/:key", v1SetKey)
// 	e.Get("/v1/:name/keys/:key", v1GetKey)
// 	e.Delete("/v1/:name/keys/:key", v1DeleteKey)

// 	e.Put("/v1/:name/cfg", v1CreateCfg)
// 	e.Delete("/v1/:name/cfg", v1DeleteCfg)
// 	e.Patch("/v1/:name/cfg", v1PatchCfg)
// 	e.Options("/v1/:name/cfg", v1OptionsCfg)
// }

// // Gets a key from discfg
// func v1GetKey(c *echo.Context) error {
// 	options.CfgName = c.Param("name")
// 	options.Key = c.Param("key")
// 	resp := commands.GetKey(options)
// 	// Since this option is not needed for anything else, it's not held on the Options struct.
// 	contentType := c.Query("type")

// 	// This is very awesome. Very interesting possibilties now.
// 	switch contentType {
// 	case "text", "text/plain", "string":
// 		return c.String(http.StatusOK, string(resp.Node.Value.([]byte)))
// 		break
// 	// This one is going to be interesting. Weird? Bad practice? I don't know, but I dig it and it starts giving me wild ideas.
// 	case "html", "text/html":
// 		return c.HTML(http.StatusOK, string(resp.Node.Value.([]byte)))
// 		break
// 	// TODO:
// 	//case "jsonp":
// 	//break
// 	case "json", "application/json":
// 		resp = commands.FormatJSONValue(resp)
// 		break
// 	default:
// 		resp = commands.FormatJSONValue(resp)
// 		break
// 	}
// 	// default response
// 	return c.JSON(http.StatusOK, resp)
// }

// // Sets a key in discfg
// func v1SetKey(c *echo.Context) error {
// 	options.CfgName = c.Param("name")
// 	options.Key = c.Param("key")
// 	resp := config.ResponseObject{
// 		Action: "set",
// 	}

// 	// Allow the value to be passed via querystring param.
// 	options.Value = []byte(c.Query("value"))

// 	// Overwrite that if the request body passes a value that can be read, preferring that.
// 	b, err := ioutil.ReadAll(c.Request().Body)
// 	if err != nil {
// 		// If reading the body failed and we don't have a value from the querystring parameter
// 		// then we have a (potential) problem.
// 		//
// 		// Some data stores may be ok with an empty key value. DynamoDB is not. It will only
// 		// return a ValidationException error. Plus, even if it was allowed, it would really
// 		// confuse the user. Some random error reading the body of a request and poof, the data
// 		// vanishes? That'd be terrible UX.
// 		// log.Println(err)
// 		resp.Error = err.Error()
// 		resp.Message = "Something went wrong reading the body of the request."
// 		// resp.ErrorCode = 500 <-- TODO: I need to come up with discfg specific error codes. Keep as const somewhere.
// 		//	return c.JSON(http.StatusOK, resp)
// 		//	Or maybe return an HTTP status message... This is outside discfg's concern. It's not an error message/code
// 		//	that would ever be seen from the CLI, right? Or maybe it would. Maybe a more generic, "error parsing key value" ...
// 	} else if len(b) > 0 {
// 		options.Value = b
// 	}

// 	resp = commands.SetKey(options)

// 	return c.JSON(http.StatusOK, resp)
// }

// // Deletes a key in discfg
// func v1DeleteKey(c *echo.Context) error {
// 	options.CfgName = c.Param("name")
// 	options.Key = c.Param("key")
// 	return c.JSON(http.StatusOK, commands.DeleteKey(options))
// }

// // Creates a new configuration
// func v1CreateCfg(c *echo.Context) error {
// 	options.CfgName = c.Param("name")
// 	resp := config.ResponseObject{
// 		Action: "create cfg",
// 	}

// 	// Any settings to pass along to the storage interface (for example, ReadCapacityUnits and WriteCapacityUnits for DynamoDB).
// 	var settings map[string]interface{}
// 	b, err := ioutil.ReadAll(c.Request().Body)
// 	if err != nil {
// 		resp.Error = err.Error()
// 		resp.Message = "Something went wrong reading the body of the request."
// 		// resp.ErrorCode = 500 <-- TODO
// 	} else if len(b) > 0 {
// 		resp := config.ResponseObject{
// 			Action: "create cfg",
// 		}
// 		//options.Value = b
// 		if err := json.Unmarshal(b, &settings); err != nil {
// 			resp.Error = err.Error()
// 			resp.Message = "Something went wrong reading the body of the request."
// 			return c.JSON(http.StatusOK, resp)
// 		}
// 	}

// 	return c.JSON(http.StatusOK, commands.CreateCfg(options, settings))
// }

// // Deletes a configuration
// func v1DeleteCfg(c *echo.Context) error {
// 	options.CfgName = c.Param("name")
// 	return c.JSON(http.StatusOK, commands.DeleteCfg(options))
// }

// // Sets options for a configuration
// func v1PatchCfg(c *echo.Context) error {
// 	options.CfgName = c.Param("name")
// 	resp := config.ResponseObject{
// 		Action: "info",
// 	}

// 	var settings map[string]interface{}
// 	b, err := ioutil.ReadAll(c.Request().Body)
// 	if err != nil {
// 		resp.Error = err.Error()
// 		resp.Message = "Something went wrong reading the body of the request."
// 		// resp.ErrorCode = 500 <-- TODO
// 		return c.JSON(http.StatusOK, resp)
// 	} else if len(b) > 0 {
// 		resp := config.ResponseObject{
// 			Action: "update cfg",
// 		}
// 		//options.Value = b
// 		if err := json.Unmarshal(b, &settings); err != nil {
// 			resp.Error = err.Error()
// 			resp.Message = "Something went wrong reading the body of the request."
// 			return c.JSON(http.StatusOK, resp)
// 		}
// 	}

// 	return c.JSON(http.StatusOK, commands.UpdateCfg(options, settings))
// }

// func v1OptionsCfg(c *echo.Context) error {
// 	options.CfgName = c.Param("name")
// 	// resp := config.ResponseObject{
// 	// 	Action: "info",
// 	// }

// 	return c.JSON(http.StatusOK, commands.Info(options))
// }
