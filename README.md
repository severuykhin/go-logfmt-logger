# go-logfmt-logger



### Version
3.0.1

### Description
The package implements simple logger for golang projects. Logger struct generates lines of output to an io.Writer  in logfmt format - [logfmt standards](https://www.cloudbees.com/blog/logfmt-a-log-format-thats-easy-to-read-and-write)  

Supports configurable verbosity level


### Usage

Install  
```
go get github.com/severuykhin/logfmt/v3
```

Import  
```
import "github.com/severuykhin/logfmt/v3"
```

Initialize  
```
// Set any output source that implements io.Writer interface
// And set verbosity level using package constants
logger := logfmt.New(os.Stdout, logfmt.L_DEBUG)

// OR integer 
logger := logfmt.New(os.Stdout, 1)

// OR built-in type
logger := logfmt.New(os.Stdout, logfmt.Level(2))
```

Use for basic log  
```
// Base information
logger.Error(42, "error message")
``` 
Output will be  
```
datetime=2022-02-22T14:14:48+03:00 level=ERROR code=42 message="error message"
```

Use for extended log  
```
// Add some context  - where rest optional Log method params will be represented as a key=value pair
logger.Debug(42, "error message", "param1", "value1", "param2", 1703)
```
Output will be  
```
datetime=2022-02-22T14:14:48+03:00 level=DEBUG code=42 message="error message" param1="value1" param2=1703
```

Only strings and numbers are supported as additional context parameters  

### Optional params   

#### `WithAppName`

You can pass appName as additional field, that will be appended in all log lines
```
logger := logfmt.New(os.Stdout, logfmt.L_DEBUG, logfmt.WithAppName("node1"))  
logger.Debug(42, "message")
```  

Output will be 
```
datetime=2022-02-22T14:14:48+03:00 level=DEBUG code=42 message="message" appName=node1
```

#### `WithFatalHook`
You can pass function that will be called in log.Fatal to stop the application
```
signals := make(chan os.Signal, 1)
signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

logger := logfmt.New(os.Stdout, logfmt.L_DEBUG, logfmt.WithFatalHook(func() {
    // some logic
    signals <- syscall.SIGINT // or just os.Exit(1)
}))

logger.Fatal(123, "fatality")

<-signals

// gaceful shutdown
```
`os.Exit(1)` is will be called by default

### Logging Levels  
```
L_DEBUG Level = 1 // DEBUG
L_INFO  Level = 2 // INFO
L_WARN  Level = 3 // WARN
L_ERROR Level = 4 // ERROR
L_FATAL Level = 5 // FATAL
```



