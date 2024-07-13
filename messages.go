package agency

import "encoding/json"

type Message interface {
	Role() Role
	Content() []byte
	Kind() Kind
}

type Kind string

const (
	TextKind      Kind = "text"
	ImageKind     Kind = "image"
	VoiceKind     Kind = "voice"
	EmbeddingKind Kind = "embedding"
)

type Role string

const (
	UserRole      Role = "user"
	SystemRole    Role = "system"
	AssistantRole Role = "assistant"
	ToolRole      Role = "tool"
)

type BaseMessage struct {
	content []byte
	role    Role
	kind    Kind
}

func (bm BaseMessage) Role() Role {
	return bm.role
}

func (bm BaseMessage) Kind() Kind {
	return bm.kind
}
func (bm BaseMessage) Content() []byte {
	return bm.content
}

// NewMessage creates new `Message` with the specified `Role` and `Kind`
func NewMessage(role Role, kind Kind, content []byte) BaseMessage {
	return BaseMessage{
		content: content,
		role:    role,
		kind:    kind,
	}
}

// NewTextMessage creates new `Message` with Text kind and the specified `Role`
func NewTextMessage(role Role, content string) BaseMessage {
	return BaseMessage{
		content: []byte(content),
		role:    role,
		kind:    TextKind,
	}
}

// NewJsonMessage marshals content and creates new `Message` with text kind and the specified `Role`
func NewJsonMessage(role Role, content any) (BaseMessage, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return BaseMessage{}, err
	}

	return BaseMessage{
		content: data,
		role:    role,
		kind:    TextKind,
	}, nil
}

func GetStringContent(msg Message) string {
	if msg.Kind() == TextKind {
		return string(msg.Content())
	}

	return ""
}
