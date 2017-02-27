package elastic

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Runner handles the coordination of applying elastic search schema changes
type Runner struct {
	SchemaChanger SchemaChanger
	Directory     string
}

// NewRunner will initialize a new Runner
func NewRunner(directory string, schemaChanger SchemaChanger) *Runner {

	return &Runner{
		//HTTPClient:    http.DefaultClient,
		SchemaChanger: schemaChanger,
		//ServerUrl:     serverUrl,
		Directory: directory,
	}
}

// Deploy will examine all of the files, verify
// they are valid and apply the changes to elastic search
func (r *Runner) Deploy() ([]string, error) {
	var results []string
	files := getFiles(r.Directory)
	for _, file := range files {
		s := NewSchemaChange(file)

		applied, err := r.SchemaChanger.WasApplied(s.ID)
		if err != nil {
			return results, err
		}
		path := fmt.Sprintf("%s\\%s", s.Folder, s.FileName)
		if applied {
			results = append(results, "Skipped: "+path)
		} else {
			err := r.SchemaChanger.Apply(s)
			if err != nil {
				results = append(results, "Error: "+path)
				return results, err
			}
			results = append(results, "Applied: "+path)
		}
	}
	return results, nil
}

// DryRun will examine all of the files, verify
// they are valid and ONLY list out the changes that
// would be applied to elastic search
func (r *Runner) DryRun() ([]string, error) {
	var results []string
	files := getFiles(r.Directory)
	for _, file := range files {
		s := NewSchemaChange(file)
		err := s.Action.Validate()
		if err != nil {
			return nil, err
		}
		applied, err := r.SchemaChanger.WasApplied(s.ID)
		if err != nil {
			return results, err
		}
		path := fmt.Sprintf("%s\\%s", s.Folder, s.FileName)
		if applied {
			results = append(results, "Skip: "+path)
		} else {
			results = append(results, "Apply: "+path)
		}
	}
	return results, nil
}

//esObjectExists will send a HEAD request to ES to verify if an object (index, document) exists
// func (r *Runner) esObjectExists(path string) (bool, error) {
// 	url := fmt.Sprintf("%s%s", r.ServerUrl, path)

// 	req, _ := http.NewRequest("HEAD", url, nil)

// 	resp, err := r.HTTPClient.Do(req)
// 	if err != nil {
// 		if resp.StatusCode == 200 {
// 			return true, nil
// 		}
// 	}
// 	return false, err
// }

// func (r *Runner) schemaChangeAlreadyRun(s *SchemaChange) (bool, error) {
// 	url := fmt.Sprintf("%s%s%s", r.ServerUrl, "version_info/", s.ID)
// 	req, _ := http.NewRequest("HEAD", url, nil)
// 	resp, err := r.HTTPClient.Do(req)
// 	if err != nil {
// 		if resp.StatusCode == 200 {
// 			return true, nil
// 		}
// 	}
// 	return false, err
// }

// func (r *Runner) applySchemaChange(a Action) (string, error) {

// 	url := fmt.Sprintf("%s%s", r.ServerUrl, a.URL)
// 	var body io.Reader
// 	if a.JSON != "" {
// 		body = bytes.NewBuffer([]byte(a.JSON))
// 	}
// 	req, _ := http.NewRequest(a.HTTPVerb, url, body)
// 	req.Header.Add("Accept", "application/json")

// 	if body != nil {
// 		req.Header.Add("Content-Type", "application/json")
// 	}

// 	resp, err := r.HTTPClient.Do(req)
// 	defer resp.Body.Close()

// 	if err != nil {
// 		if resp.StatusCode == 200 {
// 			return "", nil
// 		}
// 		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
// 		if err2 != nil {
// 			return "", err2
// 		}
// 		bodyString := string(bodyBytes)
// 		return bodyString, ErrSchemaChange
// 	}

// 	return "", err
// }

// func (r *Runner) SendRequest(method string, path string, data interface{}) ([]byte, error) {

// 	url := fmt.Sprintf("%s%s", r.ServerUrl, path)
// 	var body io.Reader
// 	if data != nil {
// 		jsonReq, err := json.Marshal(data)
// 		if err != nil {
// 			return nil, fmt.Errorf("marshaling data: %s", err)
// 		}
// 		body = bytes.NewBuffer(jsonReq)
// 	}

// 	req, _ := http.NewRequest(method, url, body)
// 	//req.SetBasicAuth(c.username, c.password)
// 	req.Header.Add("Accept", "application/json")

// 	if body != nil {
// 		req.Header.Add("Content-Type", "application/json")
// 	}

// 	resp, err := r.HTTPClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	return ioutil.ReadAll(resp.Body)
// }

func getFiles(dir string) []string {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".js" {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return fileList
}
