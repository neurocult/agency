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
	Model               string
	Temperature         NullableFloat32
	MaxTokens           int
	FuncDefs            []FuncDef
	StreamHandler       func(delta, total string, isFirst, isLast bool) error
	IsToolsCallRequired bool
}

var ToolAnswerAsModelsAnswer = errors.New("tool answer should be final")

func (p Provider) TextToStream(params TextToStreamParams) *agency.Operation {
	openAITools := castFuncDefsToOpenAITools(params.FuncDefs)

	var toolChoice *string
	if params.IsToolsCallRequired {
		v := "required"
		toolChoice = &v
	}

	return agency.NewOperation(
		func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
			openAIMessages := make([]openai.ChatCompletionMessage, 0, len(cfg.Messages)+2)

			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: cfg.Prompt,
			})

			for _, cfgMsg := range cfg.Messages {
				openaiCfgMsg, err := messageToOpenAI(cfgMsg)
				if err != nil {
					return nil, fmt.Errorf("openAI msg mapping: %w", err)
				}

				openAIMessages = append(openAIMessages, openaiCfgMsg)
			}

			openaiMsg := openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleUser,
			}

			switch msg.Kind() {
			case agency.TextKind:
				openaiMsg.Content = string(msg.Content())
			case agency.ImageKind:
				openaiMsg.MultiContent = append(
					openaiMsg.MultiContent,
					openAIBase64ImageMessage(msg.Content()),
				)
			default:
				return nil, fmt.Errorf("text to stream doesn't support %s kind", msg.Kind())
			}

			openAIMessages = append(openAIMessages, openaiMsg)

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
						ToolChoice:  toolChoice,
						StreamOptions: &openai.StreamOptions{
							IncludeUsage: true,
						},
					},
				)
				if err != nil {
					return nil, fmt.Errorf("create chat completion stream: %w", err)
				}

				var content string
				var accumulatedStreamedFunctions = make([]openai.ToolCall, 0, len(openAITools))
				var usage openai.Usage
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
							// TODO update operation API and return usage along with message
							_ = usage

							return agency.NewMessage(
								agency.AssistantRole,
								agency.TextKind,
								[]byte(content),
							), nil
						}

						break
					}

					if err != nil {
						return nil, err
					}

					if recv.Usage != nil { // penultimate message
						usage = *recv.Usage
						continue
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

					for _, toolCall := range accumulatedStreamedFunctions {
						funcToCall := getFuncDefByName(params.FuncDefs, toolCall.Function.Name)
						if funcToCall == nil {
							return nil, errors.New("function not found")
						}

						var funcResult agency.Message
						funcResult, err = funcToCall.Body(ctx, []byte(toolCall.Function.Arguments))
						var isFunctionCallAsModelAnswer = errors.Is(err, ToolAnswerAsModelsAnswer)
						if err != nil && !isFunctionCallAsModelAnswer {
							return nil, fmt.Errorf("call function %s: %w", funcToCall.Name, err)
						}

						if isFunctionCallAsModelAnswer {
							return funcResult, nil
						}

						var openaiFuncResult openai.ChatCompletionMessage
						openaiFuncResult, err = messageToOpenAI(funcResult)
						if err != nil {
							return nil, fmt.Errorf("openAI msg mapping: %w", err)
						}

						openaiFuncResult.ToolCallID = toolCall.ID
						openaiFuncResult.Name = toolCall.Function.Name

						openAIMessages = append(openAIMessages, openaiFuncResult)
					}
				}

				openAIResponse.Close()
			}
		},
	)
}

func messageToOpenAI(message agency.Message) (openai.ChatCompletionMessage, error) {
	wrappedMessage := openai.ChatCompletionMessage{
		Role: string(message.Role()),
	}

	switch message.Kind() {
	case agency.TextKind:
		wrappedMessage.Content = string(message.Content())
	case agency.ImageKind:
		wrappedMessage.MultiContent = append(
			wrappedMessage.MultiContent,
			openAIBase64ImageMessage(message.Content()),
		)
	default:
		return openai.ChatCompletionMessage{}, fmt.Errorf("text to stream doesn't support %s kind", message.Kind())
	}

	return wrappedMessage, nil
}
