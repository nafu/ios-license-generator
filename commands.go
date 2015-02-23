package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
    "github.com/mono0926/ios-license-generator/io"
    "github.com/mono0926/ios-license-generator/generator"
    "github.com/mono0926/ios-license-generator/parser"
    "io/ioutil"
    "fmt"
)

var Commands = []cli.Command{
	commandGenerate,
}

var commandGenerate = cli.Command{
	Name:  "generate",
	Usage: "",
	Description: `
`,
	Action: doGenerate,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doGenerate(c *cli.Context) {
    args := c.Args()
    if len(args) != 2 {
        log.Println("Specify license path and output path!")
        return
    }
    path := args[0]
    outDir := args[1]

    if err := os.MkdirAll(outDir, 0777); err != nil {
        log.Println("cant create output dir")
        return
    }

    lines, e := reader.ReadLines(path)
    if e != nil {
        log.Println("invalid license file")
        return
    }

    names := make([]string, 1)
    for _, l := range lines {

        match, name, body := parser.Fetch(l)

        if !match {
            log.Println("invalid license line")
            continue
        }
        s := generator.GenerateLicense("template/License.plist", name, body)
        log.Println(s)
        ioutil.WriteFile(fmt.Sprintf("%s/%s.plist", outDir, name), []byte(s), 0644)
        names = append(names, name)
    }
    s := generator.GenerateLicenseList("template/LicenseList.plist", "template/LicenseListItem.plist", names)
    log.Println(s)
    ioutil.WriteFile(fmt.Sprintf("%s/LicenseList.plist", outDir), []byte(s), 0644)
}
