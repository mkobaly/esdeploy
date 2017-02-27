package elastic

import "encoding/json"

var verbs = [4]string{"POST", "PUT", "DELETE", "HEAD"}

// Action contains the actual changes to apply to elastic search
type Action struct {
	HTTPVerb string
	URL      string
	JSON     string
}

// Validate will ensure the Action is properly formated and syntactically correct
func (a Action) Validate() error {
	if a.URL == "" {
		return ErrEmptyURL
	}
	err := a.verbValid()
	if err != nil {
		return err
	}
	err = a.jsonValid()
	if err != nil {
		return err
	}
	return nil
}

func (a Action) verbValid() error {
	for _, v := range verbs {
		if v == a.HTTPVerb {
			return nil
		}
	}
	return ErrBadHTTPVerb
}

func (a Action) jsonValid() error {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(a.JSON), &js)
	if err == nil {
		return nil
	}
	return ErrBadJSON

}
