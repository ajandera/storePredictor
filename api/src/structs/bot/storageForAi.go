package structs

import "main/bot/adapters/storage"

type StorageForAI struct {
	Name    string
	Storage *storage.SeparatedMemoryStorage
}
