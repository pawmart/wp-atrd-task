package adding

// Secret adefines properties needed for secret creation
type Secret struct {
	SecretText       string
	ExpireAfterViews int32
	ExpireAfter      int32
}
