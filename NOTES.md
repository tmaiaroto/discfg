Just some notes with thoughts to be figured out.

## Think about using DynamoDB streams.
It says "near real time" ... So I'm a bit skeptical. Though it has a lot of convenience. It will show the old version of the item when modified. 
However, it only keeps history for 24hrs (or longer, but no guarantees). So it's not exatly going to be great for getting a whole history of modifications.

DynamoDB streams conveniently integrate with Lambda though. So we could get the notifications pushed out in a very easy way by using them with the triggers feature.

discfg does not have a listener like etcd, but it could have a notification service. Of course a notification service can still be built into discfg even
without the use of streams and triggers.

Might be an interesting configurable option.

UPDATE on this: Definitely going to use Kinesis streams. It's really going to make for an interesting feature.

## JSON Support

DynamoDB is supposed to support JSON and querying into objects.    
http://www.allthingsdistributed.com/2014/10/document-model-dynamodb.html    

How? Can it be used for to setup a tree hiearchy for keys? Or would that mean one config document period? I'd rather have multiple items for the config
because of size restrictions in DynamoDB. Though that was bumped up to 400kb apparently. Which is pretty darn big but still.

But then I go and read: http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DataFormat.html and it says JSON is a used as a transport
protocol only.

UPDATE on this: Now storing binary data for so many reasons. JSON can be stored and Go handles with json.RawMessage() when possible.
However, there was no querying within the JSON.

## Snowflakes

One approach I thought about taking was to use snowflake ids. These are sequential (mostly) and that kept the use of DynamoDB as append only.
Only new values would be added and the snowflake could be used as a RANGE key which would help distribute data well.

The challenge here is some of the benefits of DynamoDB would have been lost. I would have needed to re-implement certain features that I would
otherwise get for free.

There would be more queries as a result and this would mean more DynamoDB utilization which would mean more cost. ...Which would go against
the goals of the project.

I love the idea of append only. I dislike updates. The cool thing is the query would sort by this RANGE key part and by limiting the results
to just one, it would always return the latest. But at this point I don't want to think about things like expiring old items and making additional
queries for conditional updates.

UPDATE on this: I don't think this is worth the tradeoffs.