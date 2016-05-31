<img src="https://raw.githubusercontent.com/tmaiaroto/discfg/master/docs/logo.png?a" width="350" align="middle" alt="Distributed Config (discfg)" />

[![License Apache 2](https://img.shields.io/badge/license-Apache%202-blue.svg)]() [![godoc discfg](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/tmaiaroto/discfg) [![goreport discfg](https://goreportcard.com/badge/github.com/tmaiaroto/discfg)](https://goreportcard.com/report/github.com/tmaiaroto/discfg) 

_**NOTE** This project is under active development and is not considered production ready._

A serverless and distributed configuration service built on top of Amazon Web Services. Specifically,
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
  "node": {
    "version": 2,
    "value": {
      "foo": "test"
    }
  },
  "prevNode": {}
}
```

NOTE: You will only see ```prevNode``` populated upon an update. discfg does not store a history
of values.

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

## Why Yet Another One?

The goal is not to re-invent the wheel. There are many other solutions out there that work well. 
However, they are mostly "self-host" solutions. As a result, there is a bit of maintenance involved
and additional cost to get the redundancy. Plus, many of these solutions don't really take access 
control into account (discfg included). By leveraging AWS, there's a lot of convenience and options
for then securing access to these tools. In this case, something like API Gateway is a no-brainer. 
_A big goal for discfg is to provide a serverless option for application configuration._

The idea is that discfg provides a quicker, cheaper, and perhaps more convenient, option in the mix. 
Yes, you'll have to make some concessions for that...But that doesn't mean you still can't get 
a highly available solution that costs less and is easier to maintain.

Originally, this project was heavily inspired by the wonderful [etcd](https://github.com/coreos/etcd). 
The goal was to create an alternative that would be cheaper to host, leverage AWS infrastructure, 
and have a flexible storage engine which would allow you to choose how to storage the data. 

However, the project has since deviated far away from etcd. There are some very important 
differences and the intended use case is a bit different. People are using etcd (and raft) 
for some fantastic things. The target use case for discfg isn't quite the same, though there
is some cross-over. There's some ideas that discfg borrows from etcd.

In fact, due to discfg's flexiblity, there may be other possible uses beyond the original intent
of (micro)services and application configuration. 

When building _applications_ or (micro)services, configuration and state become a challenge. 
However an eventually consistent database like DynamoDB may work just fine. We may or may not have 
1,000's of writes per second. Either way, since we are leveraging Amazon services we should be able
to scale in a cost effective manner. There's just some trade offs for that cost and convenience.

This tool will use Amazon DynamoDB to store data for each configuration, but it was designed with the
ability to use other storage engines in the future (such as S3). Each storage solution is going to come
with different pros and cons.

To be completely serverless, Lambda with API Gateway can be used to work with the configuration. 
Additionally, discfg can be used from the command line or you can run your own REST API server(s). 
You could even import the package into your own Go application. _Very flexible and convenient._

Essentially, discfg is basically an interface around DynamoDB (and other distributed storage solutions).
It's not responsible for the distributed storage itself. Theoretically, this means it could even use
etcd as a storage engine.