/**
* @program: exception
*
* @description:
*
* @author: lemo
*
* @create: 2021-05-22 00:12
**/

package exception

import "fmt"

type ass int

const Assert ass = iota

func (a ass) LastNil(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	if IsNil(v[len(v)-1]) {
		return
	}
	panic(fmt.Errorf("#exception# %v is not nil", v[len(v)-1]))
}

func (a ass) Nil(v interface{}) {
	if IsNil(v) {
		return
	}
	panic(fmt.Errorf("#exception# %v is not nil", v))
}

func (a ass) Equal(a1, a2 interface{}) {
	if a1 == a2 {
		return
	}
	panic(fmt.Errorf("#exception# %v is not equal %v", a1, a2))
}

func (a ass) True(cond bool) {
	if cond {
		return
	}
	panic(fmt.Errorf("#exception# %v is not true", cond))
}
