func v1GetKey debate...
Return base64 string inside JSON? For everything that was not passed as JSON?
Or return string for situations where a simple string was passed, but then base64
if it was an image in the value field.

// log.Println(err)
			//
			// After much debate, here's the solution:
			// Allow this API to return a JSON response (the default) OR other content types.
			// This would allow direct value responses (no meta data from the Node struct).
			// Plain text, images, etc.
			// Hit URL, get image. Pretty convenient and awesome. Not the default behavior.
			// JSON is the default and assumed transport method...Because discfg is for applications.
			//
			// Look at all the responses: http://labstack.com/echo/guide/response/
			//
			// TODO: Look into support for msgpack, protobuf, etc.
			// It would be nice to return msgpack. Which may be a gzip file.
			// So maybe Echo's File() response.
			// Eh, something is afoot: https://github.com/labstack/echo/blob/master/echo.go#L112
			//
			// Websocket sounds interesting too... Definitely not a Lambda thing (time limits),
			// but a very interesting solution for polling when not using Lambda.
			//
			// Granted when hooked up through Lambda + API Gateway... we might be a bit limited.
			// TODO: Look into API Gateway's response content-type. I imagine text and even images are ok.
			// But that may simply be restricted to JSON. Run discfg server elsewhere for more content-type
			// support.
			//
			// This means that we can always do string([]byte) in the response.
			// UNLESS another param to the API was passed that instructs otherwise.
			// So when reading images, for example, it's going to look ugly. It won't be base64.
			//
			// If the user wants data URIs in their JSON then they need to base64 convert the data
			// first and send JSON. Sending JSON (or a string - this is the convenince part) is
			// the preferred way to use discfg. However, since we store byte arrays...It's possible
			// to store images, etc.
			//

Heck, Chrome will render an image in the browser from the URL if it's returning the image as a string.
Not base64, but the []byte data in string representation. So ?type=text to the API endpoint may not
exactly result in the expected output in the browser. curl, JavaScript http request is a different story.
