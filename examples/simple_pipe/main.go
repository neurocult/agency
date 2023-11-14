package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eqtlab/lib/core"
)

func main() {
	var incPipe core.Pipe = incrementPipe

	incPipe = incPipe.Chain(incPipe).Chain(incPipe)

	msg, err := incPipe(context.Background(), core.NewUserMessage("0"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(msg.Bytes()))
}

func incrementPipe(ctx context.Context, msg core.Message) (core.Message, error) {
	i, err := strconv.ParseInt(string(msg.Bytes()), 10, 10)
	if err != nil {
		return nil, err
	}
	inc := strconv.Itoa(int(i) + 1)
	return core.NewSystemMessage(inc), nil
}
