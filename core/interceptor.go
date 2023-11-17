package core

type Interceptor func(in Message, out Message, opts ...PipeOption)
