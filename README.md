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

### Quick Example Usage

Assuming you built the binary to ```discfg``` and you have your AWS credentials under ```~/.aws``` 
because you've used the AWS CLI tool before and configured them...

```
./discfg create mycfg    
./discfg use mycfg    
./discfg set /mykey '{"json": "works"}'    
```

That should create a configuration for you (a table in DynamoDB - US East region by default). 
The second command there should have set a key called "/mykey" at the root level.

Note: If you did not want to call the ```use``` command or if you need to work with multiple configurations,
you can always get and set keys by passing the configuration name. So the following ```set``` command is
the same as the one above:

```
./discfg set mycfg /mykey '{"json": "works"}'
```

Also note that the slash is optional. All keys without a forward slash will have one prepended automatically. 
That is to say they will be at the root level. Now to retrieve this value:

```
./discfg get /mykey
```

To retrieve the value as a JSON response run (and jq is handy here; https://stedolan.github.io/jq):

```
./discfg get /mykey -f json
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

## What Is It?

A distributed configuration service built on top of Amazon Web Services. Specifically, it uses 
Lambda, DyanmoDB, and API Gateway to access it all. Though it can also be used via command line.

It's heavily inspired by [etcd](https://github.com/coreos/etcd). However, there are some very
important differences and the intended use case is a bit different.

The goal of discfg is to provide a configuration and service discovery solution for applications. 
The focus is not on creating a storage solution, but rather on solutions for configuration and
convention. Projects like etcd are **much** more complex because they handle the storage and 
quorum whereas discfg reaches for other services (DynamoDB for example) to handle these concerns. 
This of course means that some control is given up...But there are good things in return.

When building _applications_ or services, configuration and state become a challenge. Especially 
in a distributed environment or when working with others.

This tool will use Amazon DynamoDB to hold the configuration, but it was designed to be able to use 
other storage engines in the future.

Lambda with API Gateway can be used to work with the configuration (GET, PUT, DELETE) so it feels 
more like etcd. However, discfg can be used from the command line or you can run your own REST API
server(s) as well. Or any combination of those interfaces.

## Motivations Behind discfg

Three main motivating factors behind making yet another distributed configuration service:

1. Cost. Using DynamoDB (and Lambda to access discfg through a REST API) is incredibly cost effective. 
There is no server in this situation so it is almost always cheaper than using something like etcd. 
Usually significantly so.

2. Convenience and simplicity. Leveragnig AWS Services makes things easy. The original intent was 
for working with Lambda and API Gateway, microservice, based applications. If you're not an AWS fan,
sorry...

3. Good habits. As developers, we need to share certain configuration data between each other and our
applications all the time. All too often we keep a lot of hard coded configuration in our code. Not good.
Especially not good when it accidentally gets committed to a public repository. Aside from the oops, it's
also a pain to manage configuration for applications.

Wouldn't it be nice to set AWS credentials in a server environment variable and just be done with config? 
You could revoke the AWS IAM at any time. You could make it read only. Then all your application has to do 
is retrieve the confg, cache it locally (for efficiency...or not) and you're done.

Need to update your application config? You can literally do it from your own machine's command line and 
have your applications reflect that change in real-time. Configuration management is a breeze.

## Differences Between etcd and discfg (where discfg is weaker)

I think the most important difference to note is the distributed lock system. DynamoDB is an eventually
consistent storage solution. However, it does have operations for stronger consistency. It'll cost a little
bit more per operation, but it exists. Even still, it's not the same thing as a distributed lock system. 
It's not quite as dependable for certain concerns.

Related to that then is the fact that discfg has a loose sense of state. While etcd has an index, discfg
does not. It does not care about keeping a history (for reasons below with regard to AWS Lambda). It does
keep a simple version counter on both the entire configuration and each key.

Fortunately, DynamoDB has conditional operations so it is not possible to make updates or delete out of 
order when using this feature. It helps with those race condition scenarios.

Currently there is no support for a tree structure. It's completely possible to use keys with slashes 
and imply a hiearchy, but there is no recursive functionality right now. DynamoDB doesn't lend itself
well to this task given a few limits. However, there are some ways to go about it and it should be 
available in a future version.

Since Lambdas can't run forever, there is no reasonable way to listen for changes. So there is no feature
for long polling to listen for change and there is no history kept (there's an alternative for this though).

discfg is using DynamoDB (for now) so keep in mind there's a document size limit of 400kb while etcd 
(currently) has no size limit.

There are many differences by design. Both good and bad depending on the use case. You're going to need
to make a decision based on the needs of your project. It's ok if discfg isn't the solution for your 
project, but I do suspect it is a good solution for many. So for the above trade offs, let's see what
we get in return by leveraging AWS.

## "Good" Differences (where discfg is perhaps stronger)

I think the biggest benefits you gain is cost and ease of use. Without a doubt.

To run etcd you need (ok, you don't need, but should have) multiple servers to form a quorum and those servers 
run 24/7. This comes with a cost. The thing about discfg is that it can run using AWS Lambda, API Gateway, and 
DynamoDB. All services that carry a low cost pay as you go model. This makes discfg far cheaper to run and faster
to setup.

Both are easy to use with a RESTful interface, but AWS adds some extra features for free. There's a few 
convenient features magically taken care of jsut by using AWS. Security, rate limiting, auto-scaling, 
caching, access control, and more. Think about AWS IAM users. Think about how you can now easily set 
discfg up to run in a read-only capacity. Where ```get``` works, but ```set``` and ```delete``` don't. 
Also think about how easy it is to revoke credentials and create new ones at any time.

Lambda and DynamoDB really provide you with great scalability without much effort. Lambda can handle 100 
invocations per second (and more upon request). Due to the way Lambda works, there is no server side daemon 
that might crash. No server to setup and no server to hack (unless you prefer to run it on a server of course).

Speaking of AWS services...How about Kinesis? While discfg won't have a feature to listen for changes, it will 
optionally send changes to a Kinesis stream. This stream can invoke another Lambda at any interval or otherwise 
be read from by an application. This comes at an additional cost, but could still weigh in at a lower total cost
than running some other configuration service.