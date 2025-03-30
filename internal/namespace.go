package internal

import "sync"

type NamespaceContext struct {
	mu     sync.RWMutex
	Topics map[string]*TopicContext
}

func NewNamespaceContext() *NamespaceContext {
	return &NamespaceContext{
		Topics: make(map[string]*TopicContext),
	}
}

func (nc *NamespaceContext) GetOrNewTopic(topic string) *TopicContext {
	nc.mu.RLock()
	topicContext, ok := nc.Topics[topic]
	nc.mu.RUnlock()
	if !ok {
		nc.mu.Lock()
		topicContext = NewTopicContext()
		nc.Topics[topic] = topicContext
		nc.mu.Unlock()
	}
	return topicContext
}

func (nc *NamespaceContext) DeleteTopic(topic string) bool {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	if _, ok := nc.Topics[topic]; !ok {
		return false
	}
	delete(nc.Topics, topic)
	return true
}
func (nc *NamespaceContext) DeleteAllTopics() {
	nc.mu.Lock()
	defer nc.mu.Unlock()
	for topic := range nc.Topics {
		delete(nc.Topics, topic)
	}
}
