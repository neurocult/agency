package core

// Interceptor is a function that is called after one pipe and before another.
type Interceptor func(in Message, out Message, opts ...PipeOption)
