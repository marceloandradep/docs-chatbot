package registry

import (
	"chatbot/internal/env"
	"chatbot/internal/store"
	"strconv"
)

func (l *Locator) NewStore() *store.Store {
	dims, _ := strconv.Atoi(env.GetOrDefault(env.EmbeddingDimension, "3072"))
	return store.NewStore(l.db, dims)
}
