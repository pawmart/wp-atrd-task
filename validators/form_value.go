package validators

import (
	"net/http"
	"strconv"
)

func FormValidator(r *http.Request) []string {
	required := []string{"secret", "expireAfterViews", "expireAfter"}

	errorMessage := make([]string, 0)
	for _, i := range required {
		if r.FormValue(i) == "" {
			errorMessage = append(errorMessage, i)
		}
	}

	return errorMessage
}

func ExpireViewsValidator(r *http.Request) bool {
	v, err := strconv.Atoi(r.FormValue("expireAfterViews"))
	if err != nil {
		return false
	}

	if v < 1 {
		return false
	}

	return true
}
