package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/pipeline"
)

func main() {
	increment := core.NewPipe(incrementFunc)

	msg, err := pipeline.New(
		increment, increment, increment,
	).Execute(context.Background(), core.NewUserMessage("0"))

	if err != nil {
		panic(err)
	}

	fmt.Println(msg)
}

func incrementFunc(ctx context.Context, msg core.Message, _ ...core.PipeOption) (core.Message, error) {
	i, err := strconv.ParseInt(string(msg.Bytes()), 10, 10)
	if err != nil {
		return nil, err
	}
	inc := strconv.Itoa(int(i) + 1)
	return core.NewSystemMessage(inc), nil
}