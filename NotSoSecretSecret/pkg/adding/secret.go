package adding

// Secret adefines properties needed for secret creation
type Secret struct {
	SecretText       string `schema:"secret,required"`
	ExpireAfterViews int32  `schema:"expireAfterViews,required"`
	ExpireAfter      int32  `schema:"expireAfter,required"`
}
