package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// SchemaChanger is the interface that handles applying schema changes
// to backend storage systems
type SchemaChanger interface {
	WasApplied(id string) (bool, error)
	Apply(s *SchemaChange) error
}

// EsSchemaChanger handles applying schema changes for Elastic Search
type EsSchemaChanger struct {
	ServerURL  string
	HTTPClient *http.Client
}

// NewEsSchemaChanger creates Elastic Search Schema changer
func NewEsSchemaChanger(serverURL string) *EsSchemaChanger {
	if !strings.HasSuffix(serverURL, "/") {
		serverURL += "/"
	}
	sc := &EsSchemaChanger{
		HTTPClient: http.DefaultClient,
		ServerURL:  serverURL,
	}
	sc.initialize()
	return sc
}

// WasApplied determins if the schema change has already been applied or not
func (s *EsSchemaChanger) WasApplied(id string) (bool, error) {
	url := fmt.Sprintf("%s%s/%s/%s", s.ServerURL, index, esType, id)
	req, _ := http.NewRequest("HEAD", url, nil)
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

// Apply will apply the schema change to Elastic Search
func (s *EsSchemaChanger) Apply(sc *SchemaChange) error {

	url := fmt.Sprintf("%s%s", s.ServerURL, sc.Action.URL)
	var body io.Reader
	if sc.Action.JSON != "" {
		body = bytes.NewBuffer([]byte(sc.Action.JSON))
	}
	req, _ := http.NewRequest(sc.Action.HTTPVerb, url, body)
	req.Header.Add("Accept", "application/json")

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := s.HTTPClient.Do(req)
	defer resp.Body.Close()

	if err != nil || resp.StatusCode != 200 {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return err2
		}
		bodyString := string(bodyBytes)
		return ErrSchemaChange{
			Message: bodyString,
		}
	}

	// successfully applied schema so now track its completed
	return s.markScheamaChangeComplete(sc)
}

func (s *EsSchemaChanger) markScheamaChangeComplete(sc *SchemaChange) error {
	url := fmt.Sprintf("%s%s/%s/%s", s.ServerURL, index, esType, sc.ID)
	h, _ := os.Hostname()
	v := VersionInfo{
		ID:         sc.ID,
		Folder:     sc.Folder,
		File:       sc.FileName,
		Machine:    h,
		DateRunUtc: time.Now().UTC(),
	}
	json, _ := json.Marshal(v)
	body := bytes.NewBuffer(json)
	req, _ := http.NewRequest("POST", url, body)
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		b, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			return err1
		}
		bs := string(b)
		return ErrSchemaChange{
			Message: bs,
		}
	}
	return nil
}

func (s *EsSchemaChanger) initialize() {
	url := fmt.Sprintf("%s%s/%s", s.ServerURL, index, esType)
	req, _ := http.NewRequest("HEAD", url, nil)
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 404 {
		body := bytes.NewBuffer([]byte(typeDefinition))
		req, _ = http.NewRequest("PUT", url, body)
		resp, err = s.HTTPClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	}

}

const index = "esdeploy_v1"
const esType = "version_info"
const typeDefinition = `
{
	"metadata": {
		"_id" : { "path" : "id" },
		"properties": {
			"id": { "type": "string", "index" : "not_analyzed" },
			"folder": { "type": "string", "index" : "not_analyzed" },			
			"file": { "type": "string", "index" : "not_analyzed" },			
			"machine": { "type": "string", "index" : "not_analyzed" },		
			"dateRunUtc": {"type": "date", "format": "dateOptionalTime", "index" : "not_analyzed" }
		}
	}
}
`
const alias = `
{
    "actions" : [
        { "add" : { "index" : "esdeploy_v1", "alias" : "esdeploy" } }
    ]
}
`
