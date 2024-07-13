package openai

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type TextToStreamParams struct {
	TextToTextParams
	StreamHandler func(delta, total string, isFirst, isLast bool) error
}

func (p Provider) TextToStream(params TextToStreamParams) *agency.Operation {
	openAITools := castFuncDefsToOpenAITools(params.FuncDefs)

	return agency.NewOperation(
		func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
			openAIMessages, err := agencyToOpenaiMessages(cfg, msg)
			if err != nil {
				return nil, fmt.Errorf("text to stream: %w", err)
			}

			for { // streaming loop
				openAIResponse, err := p.client.CreateChatCompletionStream(
					ctx,
					openai.ChatCompletionRequest{
						Model:       params.Model,
						Temperature: nullableToFloat32(params.Temperature),
						MaxTokens:   params.MaxTokens,
						Messages:    openAIMessages,
						Tools:       openAITools,
						Stream:      params.StreamHandler != nil,
						ToolChoice:  params.ToolCallRequired(),
						Seed:        params.Seed,
					},
				)
				if err != nil {
					return nil, fmt.Errorf("create chat completion stream: %w", err)
				}

				var content string
				var accumulatedStreamedFunctions = make([]openai.ToolCall, 0, len(openAITools))
				var isFirstDelta = true
				var isLastDelta = false
				var lastDelta string

				for {
					recv, err := openAIResponse.Recv()
					isLastDelta = errors.Is(err, io.EOF)

					if len(lastDelta) > 0 || (isLastDelta && len(content) > 0) {
						if err = params.StreamHandler(lastDelta, content, isFirstDelta, isLastDelta); err != nil {
							return nil, fmt.Errorf("handing stream: %w", err)
						}

						isFirstDelta = false
					}

					if isLastDelta {
						if len(accumulatedStreamedFunctions) == 0 {
							return agency.NewTextMessage(
								agency.AssistantRole,
								content,
							), nil
						}

						break
					}

					if err != nil {
						return nil, err
					}

					if len(recv.Choices) < 1 {
						return nil, errors.New("no choice")
					}

					firstChoice := recv.Choices[0]

					if len(firstChoice.Delta.Content) > 0 {
						lastDelta = firstChoice.Delta.Content
						content += lastDelta
					} else {
						lastDelta = ""
					}

					for index, toolCall := range firstChoice.Delta.ToolCalls {
						if len(accumulatedStreamedFunctions) < index+1 {
							accumulatedStreamedFunctions = append(accumulatedStreamedFunctions, openai.ToolCall{
								Index: toolCall.Index,
								ID:    toolCall.ID,
								Type:  toolCall.Type,
								Function: openai.FunctionCall{
									Name:      toolCall.Function.Name,
									Arguments: toolCall.Function.Arguments,
								},
							})
						}
						accumulatedStreamedFunctions[index].Function.Arguments += toolCall.Function.Arguments
					}

					if firstChoice.FinishReason != openai.FinishReasonToolCalls {
						continue
					}

					// Saving tool call to history
					openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
						Role:      openai.ChatMessageRoleAssistant,
						ToolCalls: accumulatedStreamedFunctions,
					})

					for _, call := range accumulatedStreamedFunctions {
						toolResponse, err := callTool(ctx, call, params.FuncDefs)
						if err != nil {
							return nil, fmt.Errorf("text to text call tool: %w", err)
						}

						if toolResponse.Role() != agency.ToolRole {
							return toolResponse, nil
						}

						openAIMessages = append(openAIMessages, toolMessageToOpenAI(toolResponse, call.ID))
					}
				}

				openAIResponse.Close()
			}
		},
	)
}

func messageToOpenAI(message agency.Message) openai.ChatCompletionMessage {
	wrappedMessage := openai.ChatCompletionMessage{
		Role: string(message.Role()),
	}

	switch message.Kind() {

	case agency.ImageKind:
		wrappedMessage.MultiContent = append(
			wrappedMessage.MultiContent,
			openAIBase64ImageMessage(message.Content()),
		)
	default:
		wrappedMessage.Content = string(message.Content())
	}

	return wrappedMessage
}

func toolMessageToOpenAI(message agency.Message, toolID string) openai.ChatCompletionMessage {
	wrappedMessage := messageToOpenAI(message)
	wrappedMessage.ToolCallID = toolID

	return wrappedMessage
}
