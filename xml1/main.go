package main

import (
	"fmt"
	"encoding/xml"
	"io/ioutil"
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
}

func main() {
	data, _ := ioutil.ReadFile("licenses.xml")
	s := &LicenseSummary{}

	_ = xml.Unmarshal([]byte(data), &s)

	for _, dep := range s.Dependencies {
		var l License
		if len(dep.Licenses) > 0 {
			// TODO: select preferred license
			l = dep.Licenses[0]
		} else {
			// TODO: read license from another list
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
}
