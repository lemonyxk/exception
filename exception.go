/**
* @program: lemo
*
* @description:
*
* @author: lemo
*
* @create: 2019-09-25 20:37
**/

package exception

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type exception struct {
	time  time.Time
	file  string
	line  int
	error string
}

func (e exception) Time() time.Time {
	return e.time
}

func (e exception) File() string {
	return e.file
}

func (e exception) Line() int {
	return e.line
}

func (e exception) Error() string {
	return e.error
}

func (e exception) String() string {
	return fmt.Sprintf(
		`{"time":"%s","file":"%s","line":%d,"error":"%s"}`,
		e.time.Format("2006-01-02 15:04:05"), e.file, e.line, e.error,
	)
}

func NewException(time time.Time, file string, line int, err string) Error {
	return exception{time: time, file: file, line: line, error: err}
}

type Error interface {
	Time() time.Time
	File() string
	Line() int
	Error() string
	String() string
}

type catchFunc func(err Error)

type finallyFunc func(err Error)

type catch struct {
	fn func(catchFunc) *finally
}

func (c *catch) Catch(fn catchFunc) *finally {
	return c.fn(fn)
}

type finally struct {
	fn  func(finallyFunc)
	err Error
}

func (f *finally) Finally(fn finallyFunc) Error {
	f.fn(fn)
	return f.err
}

func (f *finally) Error() Error {
	return f.err
}

func Try(fn func()) (c *catch) {

	defer func() {
		if err := recover(); err != nil {
			var d = 1
			var e = fmt.Errorf("%v", err)
			if strings.HasPrefix(e.Error(), "#exception# ") {
				d = 2
			}
			var stacks = newStackErrorFromDeep(strings.Replace(e.Error(), "#exception# ", "", 1), d)
			c = &catch{fn: func(f catchFunc) *finally {
				f(stacks)
				return &finally{
					err: stacks,
					fn:  func(ff finallyFunc) { ff(stacks) },
				}
			}}
		}
	}()

	fn()

	return &catch{fn: func(f catchFunc) *finally {
		return &finally{
			fn: func(ff finallyFunc) { ff(nil) },
		}
	}}
}

func Throw(v interface{}) {
	panic(fmt.Errorf("#exception# %v", v))
}

func Eat(v ...interface{}) error {
	if len(v) == 0 {
		return nil
	}
	if IsNil(v[len(v)-1]) {
		return nil
	}
	if err, ok := v[len(v)-1].(error); ok {
		return err
	}
	return nil
}

func AssertError(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	if IsNil(v[len(v)-1]) {
		return
	}
	panic(fmt.Errorf("#exception# %v", v[len(v)-1]))
}

// func AssertEqual(condition bool, v ...interface{}) {
// 	if !condition {
// 		return
// 	}
// 	var str = fmt.Sprintln(v...)
// 	panic(fmt.Errorf("#exception# %s", str[:len(str)-1]))
// }

func New(v interface{}) Error {
	return newErrorFromDeep(v)
}

func NewMany(v ...interface{}) Error {
	var str = fmt.Sprintln(v...)
	return newErrorFromDeep(str[0 : len(str)-1])
}

func NewFormat(format string, v ...interface{}) Error {
	var str = fmt.Sprintf(format, v...)
	return newErrorFromDeep(str)
}

func newErrorFromDeep(v interface{}) Error {
	file, line := Caller()
	return newErrorWithFileAndLine(v, file, line)
}

func newStackErrorFromDeep(v interface{}, deep int) Error {
	deep = 10 + deep*2
	var file, line = Stack(deep)
	return newErrorWithFileAndLine(v, file, line)
}

func newErrorWithFileAndLine(v interface{}, file string, line int) Error {
	var err string

	switch v.(type) {
	case error:
		err = v.(error).Error()
	case string:
		err = v.(string)
	case Error:
		return v.(Error)
	default:
		err = fmt.Sprintf("%v", v)
	}

	return NewException(time.Now(), file, line, err)
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	return (*eFace)(unsafe.Pointer(&i)).data == nil
}

type eFace struct {
	_type unsafe.Pointer
	data  unsafe.Pointer
}

var rootPath, _ = os.Getwd()

func Caller() (string, int) {

	var file, line = "", 0

	// 0 for opt
	for skip := 0; true; skip++ {
		_, codePath, codeLine, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		if !strings.HasPrefix(codePath, rootPath) {
			break
		}

		file, line = codePath, codeLine
	}

	return clipFileAndLine(file, line)
}

func Stack(deep int) (string, int) {
	var list = strings.Split(string(debug.Stack()), "\n")
	var info = strings.TrimSpace(list[deep])
	var flInfo = strings.Split(strings.Split(info, " ")[0], ":")
	var file, l = flInfo[0], flInfo[1]
	var line, _ = strconv.Atoi(l)
	return clipFileAndLine(file, line)
}

func GetFuncName(fn interface{}) string {
	t := reflect.ValueOf(fn).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	}
	return t.String()
}

func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func clipFileAndLine(file string, line int) (string, int) {
	if file == "" || line == 0 {
		return "", 0
	}

	if runtime.GOOS == "windows" {
		rootPath = strings.Replace(rootPath, "\\", "/", -1)
	}

	if rootPath == "/" {
		return file, line
	}

	if strings.HasPrefix(file, rootPath) {
		file = file[len(rootPath)+1:]
	}

	return file, line
}
