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
}

type Dependency struct {
//	XMLName xml.Name `xml:"dependency"`
	GroupId string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Licenses []License `xml:"licenses>license"`
}

func main() {
	data, _ := ioutil.ReadFile("licenses.xml")
	s := &LicenseSummary{}

	_ = xml.Unmarshal([]byte(data), &s)

	for _, dep := range s.Dependencies {
		fmt.Println(dep)
	}
}
