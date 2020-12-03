package entity

//FormSecret is helper struct for form values
type FormSecret struct {
	Secret           string
	ExpireAfter      string
	ExpireAfterViews string
}

//FormInt is helper struct for validated values
type FormInt struct {
	ExpireAfter      int32
	ExpireAfterViews int32
}
