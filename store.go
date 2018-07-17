package main

import (
	"container/list"
	"sync"
	"time"
)

// TagedStore implementing captha Store interface
// that can be registered with SetCustomStore
// function to handle storage and retrieval of captcha ids and solutions for
// them, replacing the default memory store.
//
// It is the responsibility of an object to delete expired and used captchas
// when necessary (for example, the default memory store collects them in Set
// method after the certain amount of captchas has been stored.)
type TagedStore interface {
	// Set sets the digits for the captcha id.
	Set(id string, digits []byte)

	// Get returns stored digits for the captcha id.
	//Clear originally indicates whether the captcha must be deleted from the store.
	//now not in use, use remove instead
	Get(id string, clear bool) (digits []byte)

	// SetTag set tag, ok if id valid and tag was empty
	SetTag(id string, tag string) (ok bool)

	// GetTag gets tag, ok if id valid
	GetTag(id string) (tag string, ok bool)

	// Del remove catcha item
	Del(id string)
}

// expValue stores timestamp and id of captchas. It is used in the list inside
// memoryStore for indexing generated captchas by timestamp to enable garbage
// collection of expired captchas.
type idByTimeValue struct {
	timestamp time.Time
	id        string
}

// StoreItem is struct to hold captcha value and some string (tag).
type StoreItem struct {
	digits []byte
	tag    string
}

// memoryTagedStore is an  store for captcha ids and their values.
type tagedMemoryStore struct {
	sync.RWMutex
	itemsByID map[string]StoreItem
	idByTime  *list.List
	// Number of items stored since last collection.
	numStored int
	// Number of saved items that triggers collection.
	collectNum int
	// Expiration time of captchas.
	expiration time.Duration
}

// NewTagedStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func NewTagedStore(collectNum int, expiration time.Duration) TagedStore {
	s := new(tagedMemoryStore)
	s.itemsByID = make(map[string]StoreItem)
	s.idByTime = list.New()
	s.collectNum = collectNum
	s.expiration = expiration
	return s
}

func (s *tagedMemoryStore) Set(id string, digits []byte) {
	s.Lock()
	s.itemsByID[id] = StoreItem{digits, ""}
	s.idByTime.PushBack(idByTimeValue{time.Now(), id})
	s.numStored++
	if s.numStored <= s.collectNum {
		s.Unlock()
		return
	}
	s.Unlock()
	go s.collect()
}

func (s *tagedMemoryStore) SetTag(id string, tag string) (ok bool) {
	s.Lock()
	defer s.Unlock()

	item, ok := s.itemsByID[id]
	if !ok {
		return
	}
	if len(item.tag) > 0 {
		ok = false
	} else {
		item.tag = tag
	}
	return
}

func (s *tagedMemoryStore) Get(id string, clear bool) (digits []byte) {
	//if !clear {
	// When we don't need to clear captcha, acquire read lock.
	s.RLock()
	defer s.RUnlock()
	/*
		} else {
			s.Lock()
			defer s.Unlock()
		}
	*/
	item, ok := s.itemsByID[id]
	if !ok {
		return
	}
	digits = item.digits
	/*
		if clear {
			delete(s.itemsByID, id)
			// XXX(dchest) Index (s.idByTime) will be cleaned when
			// collecting expired captchas.  Can't clean it here, because
			// we don't store reference to expValue in the map.
			// Maybe store it?
		}
	*/
	return
}

func (s *tagedMemoryStore) Del(id string) {
	s.Lock()
	defer s.Unlock()
	_, ok := s.itemsByID[id]
	if !ok {
		return
	}
	delete(s.itemsByID, id)
	// XXX(dchest) Index (s.idByTime) will be cleaned when
	// collecting expired captchas.  Can't clean it here, because
	// we don't store reference to expValue in the map.
	// Maybe store it?
	return
}

func (s *tagedMemoryStore) GetTag(id string) (tag string, ok bool) {
	s.RLock()
	defer s.RUnlock()

	item, ok := s.itemsByID[id]
	if !ok {
		tag = ""
		return
	}
	tag = item.tag
	return
}

func (s *tagedMemoryStore) collect() {
	now := time.Now()
	s.Lock()
	defer s.Unlock()
	s.numStored = 0
	for e := s.idByTime.Front(); e != nil; {
		ev, ok := e.Value.(idByTimeValue)
		if !ok {
			return
		}
		if ev.timestamp.Add(s.expiration).Before(now) {
			delete(s.itemsByID, ev.id)
			next := e.Next()
			s.idByTime.Remove(e)
			e = next
		} else {
			return
		}
	}
}
