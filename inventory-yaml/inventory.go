package inventory_yaml

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/goccy/go-yaml"
)

const (
	AVHEADER = `$ANSIBLE_VAULT;1.1;AES256`
)

var (
	FLAG_RM_EMPTY_DICT_VAL = true
)

type Inventory struct {
	mu      sync.Mutex
	Path    string      `json:"-" yaml:"-"`
	Mode    fs.FileMode `json:"-" yaml:"-"`
	Yml     string      `json:"-" yaml:"-"`
	Bin     string      `json:"-" yaml:"-"`
	All     *Group      `json:"all" yaml:"all"`
	Disable bool        `json:"-" yaml:"-"`
}

func NewInventory() *Inventory {
	return &Inventory{
		All:  &Group{},
		Path: os.Getenv("PWD"),
		Yml:  "inventory.yml",
		Bin:  "inventory.bin",
		Mode: 0755,
	}
}

func (o *Inventory) IsDisabled() bool {
	isDisabled, ok := os.LookupEnv("INVENTORY_DISABLE")
	if ok {
		isDisabledBool, err := strconv.ParseBool(isDisabled)
		if err != nil {
			log.Println(err)
			o.Disable = false
		}
		o.Disable = isDisabledBool
	} else {
		o.Disable = false
	}
	return o.Disable
}

func (o *Inventory) ToYML() ([]byte, error) {
	data, err := yaml.Marshal(o)
	if err != nil {
		return nil, err
	}
	replace := regexp.MustCompile(`\|\-`)
	data = replace.ReplaceAll(data, []byte("|"))

	if FLAG_RM_EMPTY_DICT_VAL {
		replace := regexp.MustCompile(`(.*): {}`)
		data = replace.ReplaceAll(data, []byte("$1:"))
	}

	pattern := fmt.Sprintf("(.*): \\|\\n(.*)\\%s", AVHEADER)
	newval := fmt.Sprintf("$1: !vault |\n$2$%s", AVHEADER)
	avPattern := regexp.MustCompile(pattern)
	data = avPattern.ReplaceAll(data, []byte(newval))

	data = append([]byte("---\n"), data...)
	return data, nil
}

func (o *Inventory) ToBIN() error {
	if o.Disable {
		return nil
	}
	buff := bytes.Buffer{}
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(o)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(o.Path, o.Bin), buff.Bytes(), o.Mode)
	if err != nil {
		return err
	}
	return nil
}

func (o *Inventory) FromBIN() error {
	if o.Disable {
		return nil
	}
	filePointer, err := os.Open(filepath.Join(o.Path, o.Bin))
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(filePointer)
	err = dec.Decode(&o)
	if err != nil {
		return err
	}
	return nil
}

func (o *Inventory) Save() error {
	if o.Disable {
		return nil
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	data, err := o.ToYML()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(o.Path, o.Yml), data, o.Mode)
	if err != nil {
		return err
	}
	return nil
}
