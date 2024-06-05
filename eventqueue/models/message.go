package models

type IMessage interface {
	GetKey() any
	GetValue() any
}

type Message struct {
	key   any
	value any
}

func (m *Message) GetKey() any {
	return m.key
}

func (m *Message) GetValue() any {
	return m.value
}
