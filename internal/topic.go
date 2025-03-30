package internal

import (
	"fmt"
	"sync"
)

type TopicContext struct {
	mu          sync.RWMutex
	Messages    []*Message
	Index       int64
	Subscribers []chan *Message
}

func NewTopicContext() *TopicContext {
	return &TopicContext{
		Messages:    make([]*Message, 0, 100),
		Subscribers: make([]chan *Message, 0, 100),
	}
}

func (tc *TopicContext) addMessage(msg *Message) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.Messages = append(tc.Messages, msg)
	tc.Index++
}

func (tc *TopicContext) AddSubscriber(sub chan *Message) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.Subscribers = append(tc.Subscribers, sub)
}

func (tc *TopicContext) GetMessages() []*Message {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.Messages
}

func (tc *TopicContext) GetSubscribers() []chan *Message {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.Subscribers
}

func (tc *TopicContext) GetMessagesFromID(msgID int64) (result []*Message) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	
	for _, v := range tc.Messages {
		if v.ID >= msgID {
			result = append(result, v)
		}
	}
	return result
}

func (tc *TopicContext) Publish(m *Message) {
	tc.addMessage(m)
	fmt.Println("do publish", m, len(tc.Subscribers))
	for _, sub := range tc.Subscribers {
		sub <- m
	}
}
