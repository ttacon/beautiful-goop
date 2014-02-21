package main

import (
	"bytes"
	"fmt"
	"github.com/ttacon/beautiful-goop/goop"
	"io/ioutil"
	"net/http"
)

func main() {

	resp, err := http.Get("http://www.google.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	goop, err := goop.BuildGoop(bytes.NewReader(data))
	fmt.Println(err)
	fmt.Printf("%#v\n", goop.Root)

	b := goop.FindAllElements("body")
	fmt.Println(b)

	c := goop.FindAllWithClass("gb1")
	for _, cEle := range c {
		fmt.Printf("%#v\n", cEle)
	}

	fmt.Printf("%#v\n", *goop.FindById("gbar").Node)
	fmt.Println(goop.FindById("gbar").Attributes())

	fmt.Println(goop.Find("a.cool.herro#yoyo.bazoo br#mamasita.homie .awesomest"))
}
