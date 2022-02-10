package elastic

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// SchemaChange represents a schema change to apply to Elastic Search
type SchemaChange struct {
	Folder   string
	FileName string
	ID       string
	Action   Action
	Retrys   int
	Shards   int //Number of shards to use per index. Only if user used tokenized value {{shards}}
	Replicas int //Number of replicas for shards. Only if user used tokenized value {{replicas}}
}

// NewSchemaChange will get keys (folder & filename)
func NewSchemaChange(file string, shards, replicas int) *SchemaChange {
	p := strings.Split(file, string(filepath.Separator))
	filename := filepath.Base(file)
	folder := p[len(p)-2]
	id := folder + "-" + filename

	if shards <= 0 {
		shards = 5 //default what ES 6 was doing
	}

	if replicas < 0 {
		replicas = 1 //default to match what ES does
	}

	s := new(SchemaChange)
	s.Folder = folder
	s.FileName = filename
	s.ID = id
	s.Shards = shards
	s.Replicas = replicas
	s.Action, s.Retrys = s.parseFile(file)

	return s
}

func (s *SchemaChange) parseFile(esFile string) (Action, int) {
	file, err := os.Open(esFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Grab the verb (line 1) and Url (line 2) from document
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	verb := scanner.Text()
	scanner.Scan()
	url := scanner.Text()

	url, retry, err := parseURL(url)
	if err != nil {
		log.Fatal(err)
	}

	var body bytes.Buffer
	for scanner.Scan() {
		//Apply both supported token replacements if present in the file for shards and replicas
		tmp := scanner.Text()
		tmp = strings.ReplaceAll(tmp, "{{shards}}", strconv.Itoa(s.Shards))
		tmp = strings.ReplaceAll(tmp, "{{replicas}}", strconv.Itoa(s.Replicas))
		body.WriteString(tmp)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return Action{
		HTTPVerb: verb,
		URL:      url,
		JSON:     body.String(),
	}, retry
}

// func parseFile(esFile string) (Action, int) {
// 	file, err := os.Open(esFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	//Grab the verb (line 1) and Url (line 2) from document
// 	scanner := bufio.NewScanner(file)
// 	scanner.Scan()
// 	verb := scanner.Text()
// 	scanner.Scan()
// 	url := scanner.Text()

// 	url, retry, err := parseURL(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var body bytes.Buffer
// 	for scanner.Scan() {
// 		body.WriteString(scanner.Text())
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	return Action{
// 		HTTPVerb: verb,
// 		URL:      url,
// 		JSON:     body.String(),
// 	}, retry
// }

//parseURL will take a URL and look for the retry option
//Returns URL, retry count, error
func parseURL(url string) (string, int, error) {
	parts := strings.Split(url, "retry=")
	if len(parts) == 2 {
		//see if ?retry=x or &retry=x present
		u := strings.TrimSuffix(parts[0], "&")
		u = strings.TrimSuffix(u, "?")
		retry, err := strconv.Atoi(parts[1])
		if err != nil {
			return u, 0, err
		}
		return u, retry, nil
	}
	return parts[0], 0, nil
}

func parseToken(body, token string, replacement string) string {
	return strings.ReplaceAll(body, token, replacement)
}
