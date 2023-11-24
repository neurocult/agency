package agency

import "fmt"

type Message struct {
	Role    Role
	Content []byte
}

func (m Message) String() string {
	return string(m.Content)
}

type Role string

const (
	UserRole      Role = "user"
	SystemRole    Role = "system"
	AssistantRole Role = "assistant"
)

// UserMessage creates new `Message` with the `Role` equal to `user`
func UserMessage(content string, args ...any) Message {
	s := fmt.Sprintf(content, args...)
	return Message{Role: UserRole, Content: []byte(s)}
}

// SystemMessage creates new `Message` with the `Role` equal to `system`
func SystemMessage(content string) Message {
	return Message{Role: SystemRole, Content: []byte(content)}
}
