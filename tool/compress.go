package main

import (
	"fmt"
	"io/ioutil"

	"github.com/qiwenilli/auv.kit/utils"
)

func main() {

	str, _ := ioutil.ReadFile("index.html")

	//
	s := string(str)
	s1, err := utils.FlateEncode(s)
	if err != nil {
		panic(err)
	}

	//
	fmt.Println("package internal")
	fmt.Println("")
	fmt.Printf("var SwggerHtml = []byte{")
	for i, x := range s1 {
		if i%16 == 0 {
			fmt.Println("")
			fmt.Print("\t")
		}
		fmt.Printf("0x%02x, ", x)
	}
	fmt.Println("\n}")
	return

	// fmt.Printf("%#v \n ", s1)
	// flate

	//
	enflated, err := utils.FlateDecode(s1)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(enflated))

}
