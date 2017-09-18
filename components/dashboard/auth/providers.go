package auth

import (
	"sync"

	"github.com/markbates/goth"
)

var mutex sync.RWMutex

func UseProviders(viders ...goth.Provider) {
	mutex.Lock()
	defer mutex.Unlock()

	goth.UseProviders(viders...)
}
func GetProviders() goth.Providers {
	mutex.RLock()
	defer mutex.RUnlock()

	return goth.GetProviders()
}

func GetProvider(name string) (goth.Provider, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	return goth.GetProvider(name)
}

func ClearProviders() {
	mutex.Lock()
	defer mutex.Unlock()

	goth.ClearProviders()
}
