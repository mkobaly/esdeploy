package elastic

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Seeder handles seeding elastic search with data
type Seeder struct {
	Creds      Creds
	HTTPClient *http.Client
	Directory  string
	ServerURL  string
}

// NewSeeder will initialize a new Seeder
func NewSeeder(directory string, serverURL string, creds Creds) *Seeder {

	return &Seeder{
		HTTPClient: http.DefaultClient,
		Creds:      creds,
		ServerURL:  serverURL,
		Directory:  directory,
	}
}

// Seed will examine all of the json files in a directory
// and apply that document against elastic search
func (s *Seeder) Seed() ([]string, error) {
	var results []string
	now := time.Now()
	p := filepath.Join(s.Directory, "poison", now.Format("20060102150405"))

	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, 0777)
		if err != nil {
			log.Panic(err)
		}
	}

	files := s.getFiles(s.Directory)
	for _, file := range files {

		pd := getPoisonSubDir(p, file)
		a := s.getAction(file)
		err := s.execute(a)

		if err != nil {
			results = append(results, "Error: "+file)
			_, f := filepath.Split(file)

			if _, err := os.Stat(pd); os.IsNotExist(err) {
				err := os.Mkdir(pd, 0777)
				if err != nil {
					log.Panic(err)
				}
			}

			err := os.Rename(file, filepath.Join(pd, f))
			if err != nil {
				results = append(results, "Error moving to poison folder: "+file+"\n"+err.Error())
			}
		} else {
			results = append(results, "Success: "+file)
			err := os.Remove(file)
			if err != nil {
				results = append(results, "Error deleting file: "+file+"\n"+err.Error())
			}
		}
	}
	return results, nil
}

// Apply will apply the schema change to Elastic Search
func (s *Seeder) execute(a Action) error {

	u := a.URL
	if !strings.HasPrefix(u, "/") {
		u = "/" + u
	}
	url := fmt.Sprintf("%s%s", s.ServerURL, u)
	var body io.Reader
	if a.JSON != "" {
		body = bytes.NewBuffer([]byte(a.JSON))
	}
	req, _ := http.NewRequest(a.HTTPVerb, url, body)
	req.Header.Add("Accept", "application/json")

	if s.Creds.AuthorizationNeeded() {
		req.SetBasicAuth(s.Creds.Username, s.Creds.Password)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil || resp.StatusCode != 201 {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return err2
		}
		bodyString := string(bodyBytes)
		return errors.New(bodyString)
	}
	defer resp.Body.Close()
	return nil
}

func (s *Seeder) getFiles(dir string) []string {
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

func (s *Seeder) getAction(esFile string) Action {
	file, err := os.Open(esFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Grab the url (line 1) and the body...rest of document
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	url := scanner.Text()

	var body bytes.Buffer
	for scanner.Scan() {
		body.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return Action{
		HTTPVerb: "PUT",
		URL:      url,
		JSON:     body.String(),
	}
}

func getPoisonSubDir(poisonDir, file string) string {
	dir := filepath.Dir(file)
	_, dn := filepath.Split(dir) // get the parent dir name only
	posDir := filepath.Join(poisonDir, dn)
	return posDir
}
