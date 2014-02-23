package main

import (
	"fmt"
	"strings"

	"github.com/ttacon/beautiful-goop/goop"
)

const webpage = `
<html>
  <head>
    <title>
      Super basic example
    </title>
  </head>
  <body>
    <div>
      Div 1
    </div>
    <div>
      Div 2
    </div>
    <div class="fun-element">
      Div 3
    </div>
    <div class="fun-element">
      Div 4
    </div>
    <div id="onlyOne">
      Div 5
    </div>
  </body>
</html>
`

func main() {
	g, err := goop.BuildGoop(strings.NewReader(webpage))
	if err != nil {
		fmt.Printf("ruh, roh! something bad happened! err: %v\n", err)
		return
	}

	// Get all divs
	divs := g.FindAllElements("div")
	fmt.Println("divs:")
	for _, div := range divs {
		// Do something with div...
		fmt.Println(div)
	}

	funElements := g.FindAllWithClass("fun-element")
	fmt.Println("\nfunElements:")
	for _, funElement := range funElements {
		// Do something with funElement...
		fmt.Println(funElement)
	}

	theOnlyOne := g.FindById("onlyOne")
	fmt.Println("\ntheOnlyOne:")
	// Do something with the only one...
	fmt.Println(theOnlyOne)
}
