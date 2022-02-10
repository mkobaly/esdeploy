package elastic

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ValidationResult is the result of validating a schema file
type ValidationResult struct {
	File    string
	IsValid bool
}

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
func (r *Runner) Deploy(shards, replicas int) ([]string, error) {
	var results []string
	files := getFiles(r.Directory)
	for _, file := range files {
		s := NewSchemaChange(file, shards, replicas)

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
		s := NewSchemaChange(file, -1, -1)
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

//Validate will ensure all schema files are following
//the required format and are valid
func (r *Runner) Validate() []ValidationResult {
	var results []ValidationResult
	files := getFiles(r.Directory)
	for _, file := range files {
		s := NewSchemaChange(file, -1, -1)
		err := s.Action.Validate()
		if err != nil {
			results = append(results, ValidationResult{File: file, IsValid: false})
			continue
		}
		results = append(results, ValidationResult{File: file, IsValid: true})
	}
	return results
}

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
