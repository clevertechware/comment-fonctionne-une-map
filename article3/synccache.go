package article3

import "sync"

type User struct {
	ID   string
	Name string
}

func syncMapStoreAndLoadExample(user User) (User, bool) {
	var cache sync.Map

	cache.Store("user:42", user)
	v, ok := cache.Load("user:42")
	if !ok {
		return User{}, false
	}
	return v.(User), true
}
