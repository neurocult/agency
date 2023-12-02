package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/neurocult/agency"
)

func main() {
	increment := agency.NewOperation(incrementFunc)

	msg, err := agency.NewProcess(
		increment, increment, increment,
	).Execute(context.Background(), agency.UserMessage("0"))

	if err != nil {
		panic(err)
	}

	fmt.Println(msg)
}

func incrementFunc(ctx context.Context, msg agency.Message, _ *agency.OperationConfig) (agency.Message, error) {
	i, err := strconv.ParseInt(string(msg.Content), 10, 10)
	if err != nil {
		return agency.Message{}, err
	}
	inc := strconv.Itoa(int(i) + 1)
	return agency.SystemMessage(inc), nil
}
