package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type TextToStreamParams struct {
	Model       string
	Temperature NullableFloat32
	MaxTokens   int
	FuncDefs    []FuncDef
	Stream      chan<- string
}

var ToolAnswerShouldBeFinal = errors.New("tool answer should be final")

func (p Provider) TextToStream(params TextToStreamParams) *agency.Operation {
	openAITools := castFuncDefsToOpenAITools(params.FuncDefs)

	return agency.NewOperation(
		func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
			openAIMessages := make([]openai.ChatCompletionMessage, 0, len(cfg.Messages)+2)

			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: cfg.Prompt,
			})

			for _, textMsg := range cfg.Messages {
				openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
					Role:    string(textMsg.Role),
					Content: string(textMsg.Content),
				})
			}

			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: msg.String(),
			})

			for {
				openAIResponse, err := p.client.CreateChatCompletionStream(
					ctx,
					openai.ChatCompletionRequest{
						Model:       params.Model,
						Temperature: nullableToFloat32(params.Temperature),
						MaxTokens:   params.MaxTokens,
						Messages:    openAIMessages,
						Tools:       openAITools,
						Stream:      params.Stream != nil,
					},
				)
				if err != nil {
					return agency.Message{}, fmt.Errorf("create chat completion stream: %w", err)
				}

				var content string
				var accumulatedStreamedFunctions = make([]openai.ToolCall, 0, len(openAITools))
				for {
					recv, err := openAIResponse.Recv()
					if errors.Is(err, io.EOF) {
						if len(accumulatedStreamedFunctions) == 0 {
							return agency.Message{
								Role:    agency.AssistantRole,
								Content: []byte(content),
							}, nil
						}

						break
					}

					if err != nil {
						return agency.Message{}, err
					}

					if len(recv.Choices) < 1 {
						return agency.Message{}, errors.New("no choice")
					}

					firstChoice := recv.Choices[0]

					if len(firstChoice.Delta.Content) > 0 {
						params.Stream <- firstChoice.Delta.Content
						content += firstChoice.Delta.Content
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

					openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
						Role:      openai.ChatMessageRoleAssistant,
						ToolCalls: accumulatedStreamedFunctions,
					})

					for _, toolCall := range accumulatedStreamedFunctions {

						funcToCall := getFuncDefByName(params.FuncDefs, toolCall.Function.Name)
						if funcToCall == nil {
							return agency.Message{}, errors.New("function not found")
						}

						funcResult, err := funcToCall.Body(ctx, []byte(toolCall.Function.Arguments))
						var isFinal = errors.Is(err, ToolAnswerShouldBeFinal)
						if err != nil && !isFinal {
							return agency.Message{}, fmt.Errorf("call function %s: %w", funcToCall.Name, err)
						}

						escapedFuncResult, err := json.Marshal(funcResult)
						if err != nil {
							return agency.Message{}, fmt.Errorf("marshal function result: %w", err)
						}

						if isFinal {
							params.Stream <- string(escapedFuncResult)

							return agency.Message{
								Role:    openai.ChatMessageRoleTool,
								Content: escapedFuncResult,
							}, nil
						}

						openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
							Role:       openai.ChatMessageRoleTool,
							Content:    string(escapedFuncResult),
							Name:       toolCall.Function.Name,
							ToolCallID: toolCall.ID,
						})
					}
				}

				openAIResponse.Close()
			}
		},
	)
}
