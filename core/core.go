package core

import (
	"context"
	"fmt"
)

type Pipe func(context.Context, Message) (Message, error)

func (p Pipe) Chain(next Pipe) Pipe {
	return func(ctx context.Context, bb Message) (Message, error) {
		bb, err := p(ctx, bb)
		if err != nil {
			return nil, err
		}
		return next(ctx, bb)
	}
}

type Message interface {
	Bytes() []byte
}

type TextMessage struct {
	Role    Role
	Content string
}

func (t TextMessage) Bytes() []byte {
	return []byte(t.Content)
}

func (t TextMessage) Bind(args ...any) TextMessage {
	return TextMessage{
		Role:    t.Role,
		Content: fmt.Sprintf(t.Content, args...),
	}
}

type ImageMessage struct {
	bb []byte
}

func (i ImageMessage) Bytes() []byte {
	return i.bb
}

func NewImageMessage(bb []byte) ImageMessage {
	return ImageMessage{bb}
}

type Role string

const (
	UserRole   Role = "user"
	SystemRole Role = "system"
)

func UserMessage(content string) TextMessage {
	return TextMessage{Role: UserRole, Content: content}
}

func SystemMessage(content string) TextMessage {
	return TextMessage{Role: SystemRole, Content: content}
}

type PipeFactory interface {
	TextToText(prefix ...TextMessage) Pipe
	TextToImage() Pipe
}
