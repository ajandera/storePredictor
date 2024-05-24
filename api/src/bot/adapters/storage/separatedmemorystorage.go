package storage

import (
	"encoding/gob"
	"main/bot/nlp"
	"os"
)

type SeparatedMemoryStorage struct {
	filepath           string
	declarativeStorage GobStorage
	questionStorage    GobStorage
}

func NewSeparatedMemoryStorage(filepath string) (*SeparatedMemoryStorage, error) {
	var declarativeStorage, questionStorage GobStorage

	if _, err := os.Stat(filepath); err == nil {
		f, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		decoder := gob.NewDecoder(f)
		if declarativeStorage, err = RestoreMemoryStorage(decoder); err != nil {
			return nil, err
		}

		if questionStorage, err = RestoreMemoryStorage(decoder); err != nil {
			return nil, err
		}
	} else {
		declarativeStorage = NewMemoryStorage()
		questionStorage = NewMemoryStorage()
	}

	return &SeparatedMemoryStorage{
		filepath:           filepath,
		declarativeStorage: declarativeStorage,
		questionStorage:    questionStorage,
	}, nil
}

func (storage *SeparatedMemoryStorage) BuildIndex() {
	storage.declarativeStorage.BuildIndex()
	storage.questionStorage.BuildIndex()
}

func (storage *SeparatedMemoryStorage) Count() int {
	return storage.declarativeStorage.Count() + storage.questionStorage.Count()
}

func (storage *SeparatedMemoryStorage) Find(sentence string) (map[string]int, bool) {
	if nlp.IsQuestion(sentence) {
		return storage.questionStorage.Find(sentence)
	} else {
		return storage.declarativeStorage.Find(sentence)
	}
}

func (storage *SeparatedMemoryStorage) Search(sentence string) []string {
	if nlp.IsQuestion(sentence) {
		return storage.questionStorage.Search(sentence)
	} else {
		return storage.declarativeStorage.Search(sentence)
	}
}

func (storage *SeparatedMemoryStorage) Remove(sentence string) {
	if nlp.IsQuestion(sentence) {
		storage.questionStorage.Remove(sentence)
	} else {
		storage.declarativeStorage.Remove(sentence)
	}
}

func (storage *SeparatedMemoryStorage) Sync() error {
	f, err := os.Create(storage.filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)

	storage.declarativeStorage.SetOutput(encoder)
	if err := storage.declarativeStorage.Sync(); err != nil {
		return err
	}

	storage.questionStorage.SetOutput(encoder)
	return storage.questionStorage.Sync()
}

func (storage *SeparatedMemoryStorage) Update(sentence string, responses map[string]int) {
	if nlp.IsQuestion(sentence) {
		storage.questionStorage.Update(sentence, responses)
	} else {
		storage.declarativeStorage.Update(sentence, responses)
	}
}
