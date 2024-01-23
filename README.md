# Simple rest timer

This is a simple rest timer for people who don't know how to sense time at work.

## How to use?

```shell
go run main.go -interval=25
```

_The flag sets the intervals for triggering the rest notification.
Intervals are set in minutes.
If you run the program without specifying intervals, the default value will be used - 25 minutes._

When you click the close rest notification button, the timer starts again.