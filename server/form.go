package server

import (
	"fmt"
	"strings"
)

// Message represents data coming from form and validation errors.
type Message struct {
	Secret           string
	ExpireAfterViews string
	ExpireAfter      string
	Errors           map[string]string
}

// Validate validates data submitted to form.
func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	if strings.TrimSpace(msg.Secret) == "" {
		msg.Errors["Secret"] = "Secret cannot be empty"
	}

	return len(msg.Errors) == 0
}

// PrintErrors pretty prints Errors map.
func (msg *Message) PrintErrors() string {
	var sb strings.Builder

	for k, v := range msg.Errors {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	return sb.String()
}
