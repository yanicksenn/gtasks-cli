# Date proto

A simple proto representing gregorian dates and time.

## Usage

```build
# In BUILD files:
//protos/datetime:datetime-proto
//protos/datetime:datetime-proto-java
```

```java
// In java files:
Date date = Date.newBuilder()
    .setYear(2024)
    .setMonth(7)
    .setDay(14)
    .build();
```

```protobuf
// In proto files.
import "protos/datetime/datetime.proto";
```