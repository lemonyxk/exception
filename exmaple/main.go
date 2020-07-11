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

	"github.com/Lemo-yxk/exception"
)

func main() {

	exception.Try(func() {
		panic(1)
	}).Catch(func(err exception.Error) {
		fmt.Printf("%s", err.String())
	})

}
