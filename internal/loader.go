package internal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

type Loader struct {
	pathFiles string
	template  string
}

func NewLoader(template string) *Loader {
	ex, err := os.Executable()
	if err != nil {
		log.Panicf("ERROR: %s\n", err.Error())
	}
	pathFiles := filepath.Dir(ex)

	return &Loader{pathFiles, template}
}

func (l *Loader) CreateTemplate(name string) error {
	pro := Project{
		Name: name,
		Stages: []Stage{
			{
				Type: ENVIRONMENTS,
				Environment: map[string]string{
					"MY_VAR": "Value my var",
				},
			},
			{
				Type: ENTRIES,
				Entries: []Entry{
					{
						IsDir: true,
						Path:  "test_folder",
					},
					{
						Path:    "test_folder/test_file.txt",
						Content: "Example content {{.Name}}",
					},
					{
						Encoded: true,
						Path:    "test_folder/test_file_encoded.txt",
						Content: "RXhhbXBsZSBjb250ZW50IHt7Lk5hbWV9fQ==",
					},
				},
			},
			{
				Type: COMMANDS,
				Commands: []string{
					"ls",
					"pwd",
				},
			},
		},
	}

	bf, err := json.MarshalIndent(pro, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(l.pathFiles, fmt.Sprintf("%s.json", l.template)), bf, fs.ModePerm)
}

func (l *Loader) Build(projectName string) error {
	current, err := os.Getwd()
	if err != nil {
		return err
	}

	pro, err := l.loadTemplate(projectName)
	if err != nil {
		return err
	}

	for _, stage := range pro.Stages {
		switch stage.Type {
		case ENVIRONMENTS:
			for k, v := range stage.Environment {
				err = os.Setenv(k, v)
				if err != nil {
					return err
				}
			}
		case ENTRIES:
			for _, entry := range stage.Entries {
				if entry.IsDir {
					err := os.Mkdir(filepath.Join(current, entry.Path), fs.ModePerm)
					if err != nil {
						return err
					}
				} else {
					f, err := os.Create(filepath.Join(current, entry.Path))
					if err != nil {
						return err
					}

					_, err = f.WriteString(entry.Content)
					if err != nil {
						return err
					}
					f.Close()
				}
			}
		case COMMANDS:
			for _, command := range stage.Commands {
				fmt.Printf("%s\n", command)
				if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
					c := exec.Command("bash", "-c", command)
					c.Stdout = os.Stdout
					err := c.Run()
					if err != nil {
						return err
					}
				} else {
					c := exec.Command(command)
					c.Stdout = os.Stdout
					err := c.Run()
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (l *Loader) loadTemplate(projectName string) (*Project, error) {
	pathFile := ""

	files, err := ioutil.ReadDir(l.pathFiles)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() && f.Name() == fmt.Sprintf("%s.json", l.template) {
			pathFile = filepath.Join(l.pathFiles, f.Name())
			break
		}
	}

	f, err := os.Open(pathFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var pro Project
	err = json.Unmarshal(b, &pro)
	if err != nil {
		return nil, err
	}

	pro.Name = projectName
	for i := 0; i < len(pro.Stages); i++ {
		if pro.Stages[i].Type == ENTRIES {
			for j := 0; j < len(pro.Stages[i].Entries); j++ {
				if !pro.Stages[i].Entries[j].IsDir && pro.Stages[i].Entries[j].Encoded {
					bcontent, err := base64.StdEncoding.DecodeString(pro.Stages[i].Entries[j].Content)
					if err != nil {
						return nil, err
					}
					pro.Stages[i].Entries[j].Content = string(bcontent)
				}
			}
		}
	}

	b, err = json.Marshal(&pro)
	if err != nil {
		return nil, err
	}

	temp, err := template.New("template.json").Parse(string(b))
	if err != nil {
		return nil, err
	}

	bfr := new(bytes.Buffer)
	err = temp.Execute(bfr, &pro)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bfr.Bytes(), &pro)
	if err != nil {
		return nil, err
	}

	return &pro, nil
}
