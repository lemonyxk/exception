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

	"github.com/lemoyxk/exception"
)

func main() {

	fmt.Println(exception.New("hello error").String())

	exception.Try(func() {
		panic(1)
	}).Catch(func(err exception.Error) {
		fmt.Printf("%s\n", err.String())
	}).Finally(func(err exception.Error) {
		fmt.Printf("%s\n", err.String())
	}).String()

}
