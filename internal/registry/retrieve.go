package registry

import (
	"chatbot/internal/env"
	"chatbot/internal/retrieve"
)

const (
	defaultChatModel = "gpt-4.1-mini"
)

func (l *Locator) NewRetrievalHandler() *retrieve.Handler {
	return retrieve.NewHandler(l.NewRetriever())
}

func (l *Locator) NewRetriever() *retrieve.Retriever {
	chatModel := env.GetOrDefault(env.ChatModel, defaultChatModel)
	return retrieve.NewRetriever(l.db, l.oaiCli, l.NewIngestor(), chatModel)
}
