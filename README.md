# Distributed Config (discfg)

**NOTE** This project is under active development and is not considered production ready.
In fact, under no circumstances should you use this in production. However, I always appreciate
feedback. The goals of this project are constantly changing and while I first compared and 
contrasted it with etcd, it's quickly departing from the same feature set.

discfg is being created for a very simple reason. Shared application configuration and state 
within AWS services. The expectation is that you're developing some distributed application. 
Perhaps a series of micro/services.

In the future, maybe it will have the ability to run independent of AWS. However, the first
cut of this will heavily rely upon AWS. This is for convenience and cost. The drawback is 
speed and consistency. Solutions like etcd will ultimately be faster and better geared for
systems level needs. discfg is looking more toward applications.


------

A serverless and distributed configuration service built on top of Amazon Web Services. Specifically,
it uses Lambda, DyanmoDB, and API Gateway to access it all. Though it can also be used via command line.


### Quick Example Usage

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

### Running the API Server

While discfg is meant to be a tool for a "serverless" architecture, it doesn't mean you can't
run it on your own server. Currently, there is no storage engine that would keep the data on 
the same server, but the RESTful API can certainly be hosted there. So you can work with your 
configurations using JSON messages instead of just on the CLI or having to bring discfg into 
your own Go package.

In other words, discfg is language agnostic. You can use it with any application you build
so long as that application can make HTTP requests.

The API server can be on your local machine, or a remote server. The point is convenience. 
In fact, it's exactly what will run in a Lambda and be exposed through API Gateway to create
a serverless version of it.

Currently, discfg has no authentication built in. _Do not run it exposed to the world._ 
The point of using API Gateway is that Amazon provides you with the ability to protect your API.

You'll find the API server under the `server` directory. If you have the project cloned from
the repo, you could simply go to that directory and run `go main.go v1.go` or you could build
a binary and run it from where ever you like.

It runs on port `8899` by default, but you can change that with a `--port` flag. Also note
that discfg only uses AWS for storage engines right now so you should be sure to pay attention
to the AWS region. It's `us-east-1` by default, but you can change that too with a `region` flag.

Note: Releases of discfg will include both binaries, so you'll basically just be able to take
that and run with it.

## Why Yet Another One?

The goal is not to re-invent the wheel. There are many other solutions out there that work well. 
However, they are mostly self-host solutions. As a result, there is a bit of maintenance involved
and additional cost to get the redundancy. Plus, many of these solutions don't really take access 
control into account. That's up to you to manage.

There has to be a cheaper, more convenient, option in the mix. Yes, you'll have to make some 
concessions for that...But that doesn't mean you still can't get a highly available solution
on a tight budget. discfg is so cheap and so convenient that you'll want to use it for basically
every application you build. Well, maybe not every project...But that's the rationale and that's
a big part of why it's not wheel re-invention. You will use discfg in a different manner to 
accomplish _many_ of the same goals, but not all.

This project was heavily inspired by [etcd](https://github.com/coreos/etcd). The goal was to create 
an alternative that would be cheaper to host, leverage AWS infrastructure, and have a flexible 
storage engine which would allow you to choose how to storage the data. However, the project 
has since deviated far away from etcd. There are some very important differences and the intended 
use case is a bit different.

The goal of discfg is to provide a configuration and service discovery solution for **applications.** 
The focus is not on creating a storage solution for systems level use, but rather a conventional 
solution for application configuration and state.

In fact, due to it's flexiblity, there may be other possible uses beyond the original intent. 

When building _applications_ or services, configuration and state become a challenge. Especially 
in a distributed environment or when working with others. However it's not always "mission critial"
and an eventually consistent database like DynamoDB may work just fine. We may or may not have 
1,000's of writes per second. Either way, since we are leveraging Amazon services we should be able
to scale in a cost effective manner. There's just some trade offs for the cost and convenience.

This tool will use Amazon DynamoDB to hold the configuration, but it was designed to be able to use 
other storage engines in the future (such as S3).

Lambda with API Gateway can be used to work with the configuration (GET, PUT, DELETE) so it feels 
more like etcd. However, discfg can be used from the command line or you can run your own REST API
server(s) as well. Or any combination of those interfaces. You could even import the package into 
your own Go application.