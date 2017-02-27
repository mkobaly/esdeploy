package elastic

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SchemaChange represents a schema change to apply to Elastic Search
type SchemaChange struct {
	Folder   string
	FileName string
	ID       string
	Action   Action
}

// NewSchemaChange will get keys (folder & filename)
func NewSchemaChange(file string) *SchemaChange {
	p := strings.Split(file, string(filepath.Separator))
	filename := filepath.Base(file)
	folder := p[len(p)-2]
	id := folder + "-" + filename
	s := new(SchemaChange)
	s.Folder = folder
	s.FileName = filename
	s.ID = id
	s.Action = getAction(file)
	return s
}

func getAction(esFile string) Action {
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

	var body bytes.Buffer
	for scanner.Scan() {
		body.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return Action{
		HTTPVerb: verb,
		URL:      url,
		JSON:     body.String(),
	}
}
