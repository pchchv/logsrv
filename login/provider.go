package login

// Factory method for creation of login backends
type Provider func(config map[string]string) (Backend, error)

var provider = map[string]Provider{}
var providerDescription = map[string]*ProviderDescription{}

// GetProviderDescription returns the metainfo for a provider
func GetProviderDescription(providerName string) (*ProviderDescription, bool) {
	p, exist := providerDescription[providerName]
	return p, exist
}

// Returns the names of all registered provider
func ProviderList() []string {
	list := make([]string, 0, len(provider))
	for k := range provider {
		list = append(list, k)
	}
	return list
}
