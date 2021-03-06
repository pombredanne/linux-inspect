// generate-etc generates 'etc' struct based on the schema.
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gyuho/linux-inspect/etc"
	"github.com/gyuho/linux-inspect/pkg/fileutil"
	"github.com/gyuho/linux-inspect/pkg/timeutil"
	"github.com/gyuho/linux-inspect/schema"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	exp := filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect")
	if wd != exp {
		panic(fmt.Errorf("must be run in repo root %q, but run at %q", exp, wd))
	}

	buf := new(bytes.Buffer)
	buf.WriteString(`package etc

// updated at ` + timeutil.NowPST().String() + `

`)

	// '/etc/mtab'
	buf.WriteString(`// Mtab is '/etc/mtab' in Linux.
type Mtab struct {
`)
	buf.WriteString(schema.Generate(etc.MtabSchema))
	buf.WriteString("}\n\n")

	txt := buf.String()
	if err := fileutil.ToFile(txt, filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect/etc/generated.go")); err != nil {
		panic(err)
	}
	if err := os.Chdir(filepath.Join(os.Getenv("GOPATH"), "src/github.com/gyuho/linux-inspect/etc")); err != nil {
		panic(err)
	}
	if err := exec.Command("go", "fmt", "./...").Run(); err != nil {
		panic(err)
	}

	fmt.Println("DONE")
}
