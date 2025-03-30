package internal

import "sync"

type PubSubContext struct {
	mu         sync.RWMutex
	Namespaces map[string]*NamespaceContext
	Blacklist  map[string]struct{}
}

func NewPubSubContext() *PubSubContext {
	return &PubSubContext{
		Namespaces: make(map[string]*NamespaceContext),
		Blacklist:  make(map[string]struct{}),
	}
}
func (ps *PubSubContext) GetOrNewNamespace(namespace string) *NamespaceContext {
	ps.mu.RLock()
	namespaceContext, ok := ps.Namespaces[namespace]
	ps.mu.RUnlock()
	if ok {
		return namespaceContext
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.Blacklist[namespace]; ok {
		return nil
	}
	namespaceContext = NewNamespaceContext()
	ps.Namespaces[namespace] = namespaceContext
	return namespaceContext
}

func (ps *PubSubContext) DeleteNamespace(namespace string) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ns, ok := ps.Namespaces[namespace]
	if !ok {
		return false
	}
	ns.DeleteAllTopics()
	delete(ps.Namespaces, namespace)
	return true
}

func (ps *PubSubContext) GetAllNamespaces() map[string]*NamespaceContext {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.Namespaces
}

func (ps *PubSubContext) DeleteAllNamespaces() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for namespace := range ps.Namespaces {
		ns := ps.Namespaces[namespace]
		ns.DeleteAllTopics()
		delete(ps.Namespaces, namespace)
	}
}
