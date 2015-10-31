DynamoDB is the only intended storage solution for now. This directory/package exists mainly
for file/code organization. It is a personal preference over having a ton of go files all in 
one directory. However, since it's being built with an interface, it does leave the door open
for other storage options in the future.

Would that be so strange? DynamoDB was decided upon for a few reasons, but it may not work 
for all cases. A different solution may work better. This leaves the door open.

Also, configurations are meant to be shareable. If that's the case, then it's feasible that
we would move a configuration from one storage solution to another.

Maybe we'll use other data stores like Cockroach DB in the future. Who knows. DynamoDB was
just the obvious choice given the initial goals of the project (which included cost).

Or maybe it's just a local SQLite or file based storage. At that point it's not distributed...
But maybe the tool can do a little more.