# Design Decisions

I decided to use an event-based architecture for this project. The domain logic is divided into two main concepts; Sensors and Sinks. 
A sensor runs in a goroutine and runs a function periodically generating an event containing metrics. This event is published onto an [internal bus](https://github.com/adampointer/eventbus).
Listening for events on the bus are Sinks. In this example the only Sink is a filesystem output which writes all events, serialised to JSON to a file.

Using this approach I can easily add new sensor types to the application and as they run concurrently and isolated, performance will not be affected. Equally, new sinks can be added which also run in isolated goroutines. 
For example, as well as logging metrics to the local filesystem, you could broadcast them simultaneously to an external message broker like RabbitMQ or Kafka.

For the sensors I chose to use off-the-shelf libraries where possible. For the `algod_metrics` sensor I used parts of the official golang Prometheus client.
For the `proc` sensor I used [procfs](https://github.com/prometheus/procfs). I chose this particular library as it has a set of test fixtures included. As I only have a macbook which does not have a proc filesystem I could use this to develop against
and to use for automated tests.

I usually try to structure an application by separating out the pure domain logic away from the interface. This lets me use simple unit tests to test the important domain logic
and any convoluted edge cases. Interface code which talks to databases, remote APIs etc. can then be kept as simple as possible meaning that tests covering this code can be kept 
to the happy path. A simple set of integration tests can then test that all code works together against some mocks or stubs.
For this application there really is no pure domain logic and almost everything interacts with some external resource. Because of this, I have kept the structure quite simple
and leaned on test fixtures quite heavily. Due to time constraints, testing is mostly limited to the happy path.

# Productionising

In order to do something useful, the file sink would need to be replaced. The most performant approach would be to shovel all metrics as they are emitted from the different sources
onto an external message broker. If any further processing or aggregation is required, a second service could be listening to the broker and consume the raw events before processing and re-publishing
onto another topic.

To ensure the shoveling is performant and fault-tolerant, one would make sure a connection pool is available with a pool of connections which matched the number of sink goroutines.
In the event of a send failure, an exponential backoff retry should be in place to retry failed sends with a sensible maximum.
If a total loss of connection to the external broker occurs, a circuit-breaker needs to trip and events would need be cached locally to be broadcasted once the broker is available again. A background process would need to monitor the connection status.
The developer would need to make sure that suitable metrics and alerting are in place to warn on-call support if this were to happen.
When the broker connectivity is restored then the cached metrics will need to be replayed to the broker from the disk before new metrics can be sent. A configurable drain rate limit would need to be set so the broker is not overwhelmed by messages being replayed.
However, this rate will need to be set higher than the normal peak message rate or the local cache will never clear due to new metrics being cached during the drain period. 
Once the local cache is drained, new metrics can be published directly to the broker as usual. If the client is not sensitive to order it would be simpler just to backfill the cached metrics once connectivity is restored.

If the number of sensors and the sample frequency were dramatically increased, one could change the internal bus
topology slightly to support multiple concurrent sink goroutines. Currently, they would all broadcast the same sets of metrics at the same time so would be a bit of a useless change.
However, it would be trivial to shard sink goroutines either by source label or, for example, by taking the modulus of an additional message number field.

In a high-throughput production version, I would not use JSON to encode the metrics. Assuming you chose a message broker which supported it, and your clients could support it, you could use a more efficient binary message format.
You could use something proprietary, an established standard such as ISO8583 or even something as simple as gob encoded structs. Personally I would use gob encoding as it is lightning-fast and built into the language. However, that would mean
that any clients listening to the metrics stream would have to also be written in Go which may or may not be possible due to business requirements. If this was not possible, I would choose ISO8583 which is a standard common to 
financial applications such as high volume payments systems. ISO8583 codecs are available in most languages and can efficiently encode structured data.