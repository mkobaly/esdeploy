package elastic

import "testing"

func TestValidateBadVerb(t *testing.T) {
	a := Action{
		HTTPVerb: "foo",
		URL:      "foo/",
	}

	err := a.Validate()
	if err != ErrBadHTTPVerb {
		t.Error("Was expecting ErrBadHTTPVerb but got a different error")
	}
}

func TestValidateUrl(t *testing.T) {
	a := Action{
		HTTPVerb: "POST",
	}

	err := a.Validate()
	if err != ErrEmptyURL {
		t.Error("Was expecting ErrEmptyURL but got a different error")
	}
}

func TestValidateJson(t *testing.T) {
	a := Action{
		HTTPVerb: "POST",
		URL:      "foo/",
		JSON:     "bad json",
	}

	err := a.Validate()
	if err != ErrBadJSON {
		t.Error("Was expecting ErrBadJSON but got a different error")
	}
}

func TestValidateValidAction(t *testing.T) {
	a := Action{
		HTTPVerb: "POST",
		URL:      "foo/",
		JSON: `{
                    "settings" : {
                        "index" : {
                            "number_of_shards" : 5,
                            "number_of_replicas" : 1,
                            "mapper": {
                                "dynamic":false
                            }       
                        }
                    }
                }`,
	}

	err := a.Validate()
	if err != nil {
		t.Error("Was expecting no eror parsing valid Action")
	}
}
