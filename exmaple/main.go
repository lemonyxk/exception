/**
* @program: exception
*
* @description:
*
* @author: lemo
*
* @create: 2020-07-11 13:35
**/

package main

import (
	"fmt"
	"log"

	"github.com/lemonyxk/exception"
)

func main() {

	fmt.Println(exception.New("hello error").String())

	var err = exception.Try(func() {
		panic(1111111111111)
	}).Catch(func(err exception.Error) {
		fmt.Printf("%s\n", err.String())
	}).Error()

	log.Println(err)
}
