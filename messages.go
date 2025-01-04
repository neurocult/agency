package agency

type Message interface {
	Role() Role
	Content() []byte
	Kind() Kind // do we need this?
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

// --- Text Message ---

type TextMessage struct {
	role    Role
	content string
}

func (m TextMessage) Role() Role      { return m.role }
func (m TextMessage) Kind() Kind      { return TextKind }
func (m TextMessage) Content() []byte { return []byte(m.content) }

func NewTextMessage(role Role, content string) TextMessage {
	return TextMessage{
		content: content,
		role:    role,
	}
}

// --- Image Message ---

type ImageMessage struct {
	role        Role
	content     []byte
	description string
}

func (m ImageMessage) Role() Role          { return m.role }
func (m ImageMessage) Kind() Kind          { return ImageKind }
func (m ImageMessage) Content() []byte     { return m.content }
func (m ImageMessage) Description() string { return m.description }

// NewImageMessage creates new image message.
// Empty byte slice is NOT a valid content.
// Empty string IS valid description.
func NewImageMessage(role Role, content []byte, description string) ImageMessage {
	return ImageMessage{
		content:     content,
		description: description,
		role:        role,
	}
}

// --- Voice Message ---

type VoiceMessage struct {
	role    Role
	content []byte
}

func (m VoiceMessage) Role() Role      { return m.role }
func (m VoiceMessage) Kind() Kind      { return VoiceKind }
func (m VoiceMessage) Content() []byte { return m.content }

func NewVoiceMessage(role Role, content []byte) VoiceMessage {
	return VoiceMessage{
		content: content,
		role:    role,
	}
}

// --- Embedding Message ---

type EmbeddingMessage struct {
	role    Role
	content []byte
}

func (m EmbeddingMessage) Role() Role      { return m.role }
func (m EmbeddingMessage) Kind() Kind      { return EmbeddingKind }
func (m EmbeddingMessage) Content() []byte { return m.content }

func NewEmbeddingMessage(role Role, content []byte) EmbeddingMessage {
	return EmbeddingMessage{
		content: content,
		role:    role,
	}
}
