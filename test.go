package main

import (
	"fmt"
	"regexp"
)

func main() {
	domain := "myshoplaza.com"
	shopRegexp := regexp.MustCompile("^[a-zA-Z0-9-]+." + domain + "$")
	fmt.Println(shopRegexp, shopRegexp.MatchString("puping.myshoplaza.com"))
}
