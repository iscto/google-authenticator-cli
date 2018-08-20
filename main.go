package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"

	"github.com/chanyipiaomiao/hltool"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	secretDBName    = "twostep.db"
	secretTableName = "secret"
)

func SortMapByKey(m map[string][]byte) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, v := range keys {
		n, t, err := hltool.TwoStepAuthGenByKey(v)
		if err != nil {
			continue
		}
		fmt.Printf("%s %s %d\n", m[v], n, t)
	}
	return keys
}

type Secret struct {
	TwoStepDB *hltool.BoltDB
}

func NewSecret() (*Secret, error) {
	dbPath := path.Join(path.Dir(os.Args[0]), secretDBName)
	twostepDB, err := hltool.NewBoltDB(dbPath, secretTableName)
	if err != nil {
		return nil, err
	}
	return &Secret{TwoStepDB: twostepDB}, nil
}

func (s *Secret) Add(name, secret string) error {
	err := s.TwoStepDB.Set(map[string][]byte{
		secret: []byte(name),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Secret) Delete(secret string) error {
	err := s.TwoStepDB.Delete([]string{secret})
	if err != nil {
		return err
	}
	return nil
}

func (s *Secret) List(name string) error {
	if name != "all" {
		r, err := s.TwoStepDB.Get([]string{name})
		if err != nil {
			return err
		}
		SortMapByKey(r)
	} else {
		r, err := s.TwoStepDB.GetAll()
		if err != nil {
			return err
		}
		SortMapByKey(r)
	}

	return nil
}

var (
	newSecret *Secret
)

func init() {
	var err error
	newSecret, err = NewSecret()
	if err != nil {
		log.Fatalf(" NewSecret() error: %s\n", err)
	}
}

func cli() {
	app := kingpin.New("google-authenticator-cli", "模拟 Google Authenticator 验证器")

	add := app.Command("add", "添加secret")
	addName := add.Flag("name", "名称标识").Required().String()
	secret := add.Flag("secret", "二步验证里面生成的Secret,一般跟二维码一起展示").Required().String()

	del := app.Command("delete", "删除secret")
	deleteName := del.Flag("delete-secret", "名称标识").Required().String()

	show := app.Command("show", "显示所有的6位数字")
	showName := show.Flag("show-name", "显示指定的标识的6位数字").Default("all").String()

	c, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("parse cli args error: %s\n", err)
	}

	switch c {
	case "add":
		err := newSecret.Add(*addName, *secret)
		if err != nil {
			log.Fatalf("s.Add(*addName, *secret) error: %s\n", err)
		}
		fmt.Println("add ok.")
	case "delete":
		err := newSecret.Delete(*deleteName)
		if err != nil {
			log.Fatalf("s.Delete(*deleteName) error: %s\n", err)
		}
		fmt.Println("delete ok.")
	case "show":
		err := newSecret.List(*showName)
		if err != nil {
			log.Fatalf("s.List(*showName) error: %s\n", err)
		}
	}

}

func main() {
	cli()
}
