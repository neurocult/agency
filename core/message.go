package core

import "fmt"

// Message represents abstract message
type Message interface {
	Bytes() []byte
	String() string
}

type TextMessage struct {
	Role    Role
	Content string
}

func (t TextMessage) Bytes() []byte {
	return []byte(t.Content)
}

func (t TextMessage) String() string {
	return t.Content
}

// Bind allows to use prompt as a template by replacing printf directives like `%s` with the given `args`
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

func (ImageMessage) String() string {
	return "<ImageMessage>"
}

func NewImageMessage(bb []byte) ImageMessage {
	return ImageMessage{bb}
}

type Role string

const (
	UserRole      Role = "user"
	SystemRole    Role = "system"
	AssistantRole Role = "assistant"
)

// NewUserMessage creates new `TextMessage` with the `Role` equal to `user`
func NewUserMessage(content string) TextMessage {
	return TextMessage{Role: UserRole, Content: content}
}

// NewSystemMessage creates new `TextMessage` with the `Role` equal to `system`
func NewSystemMessage(content string) TextMessage {
	return TextMessage{Role: SystemRole, Content: content}
}

type SpeechMessage struct {
	bb []byte
}

func (s SpeechMessage) Bytes() []byte {
	return s.bb
}

func (SpeechMessage) String() string {
	return "<SpeechMessage>"
}

func NewSpeechMessage(bb []byte) SpeechMessage {
	return SpeechMessage{
		bb: bb,
	}
}
