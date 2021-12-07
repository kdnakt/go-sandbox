package main

import (
	"flag"
	"fmt"
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/jszwec/csvutil"
)

type LicenseSummary struct {
	XMLName xml.Name `xml:"licenseSummary"`
	Dependencies []Dependency `xml:"dependencies>dependency"`
}

type License struct {
//	XMLName xml.Name `xml:"license"`
	Name string `xml:"name"`
	Url string `xml:"url"`
}

type Dependency struct {
//	XMLName xml.Name `xml:"dependency"`
	GroupId string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version string `xml:"version"`
	Licenses []License `xml:"licenses>license"`

	Package string `csv:"module name"`
	License string `csv:"license"`
	Repository string `csv:"repository"`
}

func main() {
	flag.Parse()
	fileName := flag.Arg(0)
	data, _ := ioutil.ReadFile(fileName)

	if strings.HasSuffix(fileName, ".xml") {
		s := &LicenseSummary{}
		_ = xml.Unmarshal([]byte(data), &s)
		for _, dep := range s.Dependencies {
			var l License
			if len(dep.Licenses) > 0 {
				l = dep.Licenses[0]
			} else {
				l = License{
					Name: "NOT FOUND",
					Url: "NOT FOUND",
				}
			}
			fmt.Printf(`
[ライブラリ名] %s
[ バージョン ] %s
[ ライセンス ] %s
[    URL     ] %s
`,
				dep.ArtifactId,
				dep.Version,
				l.Name,
				l.Url)
		}
	} else if strings.HasSuffix(fileName, ".csv") {
		var deps []Dependency
		_ = csvutil.Unmarshal([]byte(data), &deps)
		for _, dep := range deps {
			p := strings.Split(dep.Package, "@")
			l := strings.ReplaceAll(dep.License, "*", "")
			fmt.Printf(`
[ライブラリ名] %s
[ バージョン ] %s
[ ライセンス ] %s
[    URL     ] %s
`,
				p[0], p[1], l, dep.Repository)
		}
	}
}
