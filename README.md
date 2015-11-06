# Distributed Config (discfg)

**NOTE** This is very early on and experimental. I would not use this for anything in production!
I'm just trying to get a proof of concept and organize my thoughts. I will start versioning this
and will note when I think it's something safe to play with. I will note when it is suitable
for production. PLEASE leave any feedback or issues in issues. Just keep in mind I probably
know of the obvious issues =) That doesn't mean they shoudln't be in there though. So flag away.

------

### Quick Example Usage

Assuming you built the binary to ```discfg``` and you have your AWS credentials under ```~/.aws``` 
because you've used the AWS CLI tool before and configured them...

```
./discfg create mycfg    
./discfg set /mykey '{"json": "works"}'    
```

That should create a configuration for you (a table in DynamoDB - US East region by default). 
The second command there should have set a key called "/mykey" at the root level.

Now to retrieve this value:

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
important differences and the target use case is a bit different.

The goal of discfg is to provide a configuration and service discovery solution for applications. 
The focus is not on creating a storage solution, but rather on solutions for configuration and
convention. Projects like etcd are **much** more complex because they handle the storage and 
quorum whereas discfg lets existing storage engines (DynamoDB for now) handle that.


When building _applications_ or services, configuration and state become a challenge. Especially 
in a distributed environment or when working with others.

You don't want to store the config inside the codebase or each microservice because that makes for
a tedious update process and provides many places where a compromise could occur. You also need
to handle updates without redeploying services. You also may want to share the configuration with 
your team members so it's easy for everyone to get the same page.

This tool will use Amazon DynamoDB to hold the configuration (for now). Though it was designed to 
be able to use other storage engines in the future.

Lambda with API Gateway can be used to work with the configuration (GET, PUT, DELETE) so it feels 
more like etcd. Or, you're of course free to put your own API in front of it.

## Motivations Behind discfg

Three main motivating factors behind making yet another distributed configuration service:

1. Cost. Using DynamoDB (and Lambda to access discfg through a REST API) is incredibly cost effective. 
There is no server in this situation so it is always cheaper than using something like etcd. Usually 
significantly so.

2. It's focused on AWS for a reason. When you build out applications using Lambda, a distributed 
configuration helps. So many decisions around discfg and its designed are based on using AWS services.

3. Developers. We need to share certain configuration data between each other and do a better job of 
keeping a lot of hard coded configuration out of our code. Not only do we see tons of sensitive credentials
all over public git repos, but it's also pain to manage configurations.

Wouldn't it be nice to set AWS credentials in a server environment variable and just be done with config? 
You could revoke the AWS IAM at any time. You could make it read only. Then all your application has to do 
is retrieve the confg, cache it locally (for efficiency) and you're done.

Need to update your application config? You can literally do it from your own machine's command line and
then wait for your application's cache to expire to get the new config (or no wait if you didn't cache it).

## Differences Between etcd and discfg

I think the most important difference to note is the distributed lock system. DynamoDB is an eventually
consistent storage solution. However. It does have operations for stronger consistency. It'll cost a little
bit more per operation, but it exists. Even still, it's not the same thing as a distributed lock system.
So it is technically possible to get stale information due to race conditions.

Related to that then is the fact that discfg has a loose sense of state. While etcd has an index, discfg
does not. It does not care about keeping a history (for reasons below with regard to AWS Lambda).

Fortunately, DynamoDB has conditional operations so it is not possible to make updates or delete out of order.
For example, if the value is not what it is expected to be when the operation gets to DynamoDB, it will fail.

So this may count discfg out for certain tasks and that's ok. Keep in mind what discfg was created for. 
It was not created to compete with or replace etcd. It's merely inspired by etcd.

Since Lambdas can't run forever, there is no reasonable way to listen for changes. Sure, other AWS services
could be used for this, but for now discfg is intended to have a simple scope. So while etcd has long polling 
for changed keys, discfg does not. For now...But SNS or SQS may be something to look into.

Related to listening for changes, another difference is that discfg does not store a history like etcd. 
This is useful in etcd because if the long polling got interrupted, it could continue where it left off.
Keeping a history of changes is not, currently, a goal of discfg nor is long polling for key changes.

discfg is using DynamoDB (for now) so keep in mind there's a document size limit of 400kb while etcd 
(currently) has no size limit.

There are many differences by design. Both good and bad depending on the use case. If the "good" differences
seem to make sense to you, then discfg might be the tool for you.

## "Good" Differences

Well, it's not that the above differences are bad...But those are, more or less, the major things lacking
in discfg when compared to etcd. Just to get them out of the way. Yes, there is a decisive trade off with 
eventual consistency and not having a distributed lock. So let's look at what you get in return.

I think the biggest benefits you gain is cost and ease of use. Without a doubt.

To run etcd you need (ok, you don't need, but should have) multiple servers to form a quorum and those servers 
run 24/7. This comes with a cost. The thing about discfg is that it runs using AWS Lambda, API Gateway, and 
DynamoDB. All services that carry a low cost pay as you go model. This makes discfg far cheaper to run.

Both are easy to use with a RESTful interface, but AWS adds some extra features for free. It makes security
fairly easy and straight forward. API Gateway adds rate limiting as well -- yea, maybe there's not much 
need there, but then again...Maybe there is.

While both are highly available, Lambda and DynamoDB really provide you with scalability. Lambda can
handle 100 invocations per second (and more upon request). Due to the way Lambda works, there is no 
server side daemon that might crash.

Due to how Lambda works, there's no concern for firewalls or other security risks you'd find with 
a traditional servers.

When we're talking about DynamoDB...An AWS service...We're talking about AWS IAM to access it. This is 
a wonderful security feature. It removes many concerns from the tool itself. It also makes it easy to manage
for the user because there are so many ways you can work with AWS IAM.