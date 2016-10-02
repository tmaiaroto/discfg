<img src="https://raw.githubusercontent.com/tmaiaroto/discfg/master/docs/logo.png?a" width="350" align="middle" alt="Distributed Config (discfg)" />

[![License Apache 2](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/tmaiaroto/discfg/blob/master/LICENSE) [![godoc discfg](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/tmaiaroto/discfg) [![Build Status](https://travis-ci.org/tmaiaroto/discfg.svg?branch=master)](https://travis-ci.org/tmaiaroto/discfg) [![goreport discfg](https://goreportcard.com/badge/github.com/tmaiaroto/discfg)](https://goreportcard.com/report/github.com/tmaiaroto/discfg)

_**NOTE** This project is under active development and is not considered production ready._ 
Feedback very much appreciated.

A serverless and distributed (key/value) configuration service built on top of Amazon Web Services. Specifically,
it aims to use Lambda, DyanmoDB, and API Gateway. Though it can be used with other services.

It can install to your system as a binary, so managing configuration from any machine is simple from 
a command line. You can also work with configurations via RESTful API.

### Command Line Interface

Assuming you built the binary to ```discfg``` and you have your AWS credentials under ```~/.aws``` 
because you've used the AWS CLI tool before and configured them...

```
./discfg create mycfg    
./discfg use mycfg    
./discfg set mykey '{"json": "works"}'    
./discfg set mykey -d file.json
```

That first command creates a configuration for you (a table in DynamoDB - US East region by default). 
After that you can set keys from values passed to the CLI or from a file on disk. All values ultimately 
get stored as binary data, so you could even store (small - DynamoDB size limits) files if you really 
wanted to; images, maybe icons, for example.

Note: If you did not want to call the ```use``` command or if you need to work with multiple configurations,
you can always get and set keys by passing the configuration name. So the following ```set``` command is
the same as the one above:

```
./discfg set mycfg mykey '{"json": "works"}'
```

Also note that the slash is optional. All keys without a forward slash will have one prepended automatically. 
That is to say they will be at the root level. Now to retrieve this value:

```
./discfg get mykey
```

To retrieve the value as a JSON response run (and jq is handy here; https://stedolan.github.io/jq):

```
./discfg get mykey -f json
```

You should see something like this:

```
{
  "action": "get",
  "item": {
    "version": 2,
    "key": "mykey",
    "value": {
      "foo": "test"
    }
  },
  "prevItem": {}
}
```

NOTE: You will only see ```prevItem``` populated upon an update. discfg does not store a history
of item values.

### Serverless API

The serverless API was built using the [Apex](http://apex.run/) framework along with [Terraform](https://www.terraform.io/).
This leverages AWS Lambda and API Gateway. Assuming you have AWS CLI setup and then setup Apex 
and Terraform, you could then easily deploy discfg (from the `apex` directory) with the following:

```
apex infra apply -var 'aws_account_id=XXXXX'
apex deploy
```

Where `XXXXX` has been replaced with your Amazon Web Services Account ID. Note that within the
`infrastructure` directory, you'll find all the `.tf` files. Feel free to adjust those in the
`variables.tf` to change simple things like the API name. You can also dig even deeper to change 
more complex things or course you can change things from the AWS web console once you've deployed
the default provided.

#### Example API Calls

You'll of course prepend these URL paths with your AWS API Gateway API's base URL.

**PUT /{name}/keys/{key}**

Body
```
any value
```

Would set the provided value from he PUT body for the config name and key name passed
in the API endpoint path parameters. There would then be a JSON response.

**GET /{name}/keys/{key}**

Would get the key value for the given key name and config name passed in the API endpoint
path parameters. The response would be a JSON message.


**PUT /{name}/cfg**

Body
```
{"WriteCapacityUnits": 2, "ReadCapacityUnits": 4}
```

Would create a table in DynamoDB with the provided name in the API endpoint path and would
also configure it with the given settings from the PUT body. In the case of DynamoDB these 
setings are the read and write capacity units (by default 1 write and 2 read).


### Running the API Server (on a server)

While discfg is meant to be a tool for a "serverless" architecture, it doesn't mean you can't
run it on your own server. Currently, there is no storage engine that would keep the data on 
the same server (and that defeats the whole purpose of being distributed), but the RESTful API 
can certainly be hosted on any server. So you can work with your configurations using JSON
instead of just on the CLI or having to bring discfg into your own Go package.

The API server can be on your local machine, or a remote server. Or both. The point is convenience.

Currently, discfg has no authentication built in. _Do not run it exposed to the world._ 
The point of relying on AWS is that Amazon provides you with the ability to control access.
From the API server exposed through API Gateway to the DynamoDB database.

You'll find the API server under the `server` directory. If you have the project cloned from
the repo, you could simply go to that directory and run `go main.go v1.go` to check it out.
You'll ultimatley want to build a binary and run it from where ever you need.

It runs on port `8899` by default, but you can change that with a `--port` flag. Also note
that discfg only uses AWS for storage engines right now so you should be sure to pay attention
to the AWS region. It's `us-east-1` by default, but you can change that too with a `region` flag.

## What prompted this tool?

The need for a serverless application configuration. When dealing with AWS Lambda, state and 
configuration become a regular issue with no consistent solution. This ends up becoming a bit
of boilerplate code in each of your functions.

Discfg solves this need by making it very easy to work with key/value configuration data.

Of course your application need not be serverless or run in AWS Lambda to benefit from discfg.

Managing a bunch of environment variables isn't scalable. It's annoying and when you go to deploy
or have a co-worker work on the project, it becomes a hassle. Even with tools like Docker. Things
change and keeping on top of configuration changes is simply annoying with environment variables.

Likewise, dropping an .ini or .json or .env file into a project is also not a terrific solution.
Configuration files also become quickly dated and it still doesn't help much when you need to
share configurations with others.

Essentially, discfg is an interface around DynamoDB (and other distributed storage solutions).
Opposed to some other configuration tools, it's not responsible for the distributed storage itself.

Etcd was a big inspiration for this project. However, etcd is not "serverless" and so it requires one
to set up and maintain a cluster of servers. This is a little less convenient, though the tool itself
is also much faster. There's a trade off for convenience. Discfg was meant for higher level application
use and so the performance factor wasn't a concern. The feature set of discfg also started to diverge
from etcd as well. Discfg is simply a tool with a different use case.
