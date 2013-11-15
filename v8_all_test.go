package v8

import (
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"testing"
	"time"
)

var Default = NewEngine()

func init() {
	// traceDispose = true
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			input, err := ioutil.ReadFile("test.cmd")

			if err == nil && len(input) > 0 {
				ioutil.WriteFile("test.cmd", []byte(""), 0744)

				cmd := strings.Trim(string(input), " \n\r\t")

				var p *pprof.Profile

				switch cmd {
				case "lookup goroutine":
					p = pprof.Lookup("goroutine")
				case "lookup heap":
					p = pprof.Lookup("heap")
				case "lookup threadcreate":
					p = pprof.Lookup("threadcreate")
				default:
					println("unknow command: '" + cmd + "'")
				}

				if p != nil {
					file, err := os.Create("test.out")
					if err != nil {
						println("couldn't create test.out")
					} else {
						p.WriteTo(file, 2)
					}
				}
			}

			time.Sleep(2 * time.Second)
		}
	}()
}

func Test_HelloWorld(t *testing.T) {
	context := Default.NewContext()
	script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
	value := script.Run(context)
	result := value.ToString()

	if result != "Hello World!" {
		t.Fatal("result not match")
	}

	runtime.GC()
}

func Test_PreCompile(t *testing.T) {
	// pre-compile
	code := []byte("'Hello ' + 'PreCompile!'")
	scriptData1 := Default.PreCompile(code)

	// test save and load script data
	data := scriptData1.Data()
	scriptData2 := NewScriptData(data)

	// test compile with script data
	context := Default.NewContext()
	script := Default.Compile(code, nil, scriptData2)
	value := script.Run(context)
	result := value.ToString()

	if result != "Hello PreCompile!" {
		t.Fatal("result not match")
	}

	runtime.GC()
}

func Test_TypeCheck(t *testing.T) {
	// TODO
}

func Test_Values(t *testing.T) {
	if !Default.Undefined().IsUndefined() {
		t.Fatal("Undefined() not match")
	}

	if !Default.Null().IsNull() {
		t.Fatal("Null() not match")
	}

	if !Default.True().IsTrue() {
		t.Fatal("True() not match")
	}

	if !Default.False().IsFalse() {
		t.Fatal("False() not match")
	}

	if Default.Undefined() != Default.Undefined() {
		t.Fatal("Undefined() != Undefined()")
	}

	if Default.Null() != Default.Null() {
		t.Fatal("Null() != Null()")
	}

	if Default.True() != Default.True() {
		t.Fatal("True() != True()")
	}

	if Default.False() != Default.False() {
		t.Fatal("False() != False()")
	}

	var (
		maxInt32  = int64(0x7FFFFFFF)
		maxUint32 = int64(0xFFFFFFFF)
		maxUint64 = uint64(0xFFFFFFFFFFFFFFFF)
		maxNumber = int64(maxUint64)
	)

	if Default.NewBoolean(true).ToBoolean() != true {
		t.Fatal(`NewBoolean(true).ToBoolean() != true`)
	}

	if Default.NewNumber(12.34).ToNumber() != 12.34 {
		t.Fatal(`NewNumber(12.34).ToNumber() != 12.34`)
	}

	if Default.NewNumber(float64(maxNumber)).ToInteger() != maxNumber {
		t.Fatal(`NewNumber(float64(maxNumber)).ToInteger() != maxNumber`)
	}

	if Default.NewInteger(maxInt32).IsInt32() == false {
		t.Fatal(`NewInteger(maxInt32).IsInt32() == false`)
	}

	if Default.NewInteger(maxUint32).IsInt32() != false {
		t.Fatal(`NewInteger(maxUint32).IsInt32() != false`)
	}

	if Default.NewInteger(maxUint32).IsUint32() == false {
		t.Fatal(`NewInteger(maxUint32).IsUint32() == false`)
	}

	if Default.NewInteger(maxNumber).ToInteger() != maxNumber {
		t.Fatal(`NewInteger(maxNumber).ToInteger() != maxNumber`)
	}

	if Default.NewString("Hello World!").ToString() != "Hello World!" {
		t.Fatal(`NewString("Hello World!").ToString() != "Hello World!"`)
	}

	if Default.NewObject().IsObject() == false {
		t.Fatal(`NewObject().IsObject() == false`)
	}

	if Default.NewArray(5).IsArray() == false {
		t.Fatal(`NewArray(5).IsArray() == false`)
	}

	if Default.NewArray(5).Length() != 5 {
		t.Fatal(`NewArray(5).Length() != 5`)
	}

	if Default.NewRegExp("foo", RF_None).IsRegExp() == false {
		t.Fatal(`NewRegExp("foo", RF_None).IsRegExp() == false`)
	}

	if Default.NewRegExp("foo", RF_Global).ToRegExp().Pattern() != "foo" {
		t.Fatal(`NewRegExp("foo", RF_Global).ToRegExp().Pattern() != "foo"`)
	}

	if Default.NewRegExp("foo", RF_Global).ToRegExp().Flags() != RF_Global {
		t.Fatal(`NewRegExp("foo", RF_Global).ToRegExp().Flags() != RF_Global`)
	}

	runtime.GC()
}

func Test_Object(t *testing.T) {
	context := Default.NewContext()
	script := Default.Compile([]byte("a={};"), nil, nil)
	value := script.Run(context)
	object := value.ToObject()

	// Test get/set property
	if prop := object.GetProperty("a"); prop != nil {
		if !prop.IsUndefined() {
			t.Fatal("property 'a' value not match")
		}
	} else {
		t.Fatal("could't get property 'a'")
	}

	if !object.SetProperty("b", Default.True(), PA_None) {
		t.Fatal("could't set property 'b'")
	}

	if prop := object.GetProperty("b"); prop != nil {
		if !prop.IsBoolean() || !prop.IsTrue() {
			t.Fatal("property 'b' value not match")
		}
	} else {
		t.Fatal("could't get property 'b'")
	}

	// Test get/set non-ascii property
	if !object.SetProperty("中文字段", Default.False(), PA_None) {
		t.Fatal("could't set non-ascii property")
	}

	if prop := object.GetProperty("中文字段"); prop != nil {
		if !prop.IsBoolean() || !prop.IsFalse() {
			t.Fatal("non-ascii property value not match")
		}
	} else {
		t.Fatal("could't get non-ascii property")
	}

	// Test get/set element
	if elem := object.GetElement(0); elem != nil {
		if !elem.IsUndefined() {
			t.Fatal("element 0 value not match")
		}
	} else {
		t.Fatal("could't get element 0")
	}

	if !object.SetElement(0, Default.True()) {
		t.Fatal("could't set element 0")
	}

	if elem := object.GetElement(0); elem != nil {
		if !elem.IsTrue() {
			t.Fatal("element 0 value not match")
		}
	} else {
		t.Fatal("could't get element 0")
	}

	// Test GetPropertyAttributes
	if !object.SetProperty("x", Default.True(), PA_DontDelete|PA_ReadOnly) {
		t.Fatal("could't set property with attributes")
	}

	attris := object.GetPropertyAttributes("x")

	if attris&(PA_DontDelete|PA_ReadOnly) != PA_DontDelete|PA_ReadOnly {
		t.Fatal("property attributes not match")
	}

	// Test ForceSetProperty
	if !object.ForceSetProperty("x", Default.False(), PA_None) {
		t.Fatal("could't force set property 'x'")
	}

	if prop := object.GetProperty("x"); prop != nil {
		if !prop.IsBoolean() || !prop.IsFalse() {
			t.Fatal("property 'x' value not match")
		}
	} else {
		t.Fatal("could't get property 'x'")
	}

	// Test HasProperty and DeleteProperty
	if object.HasProperty("a") {
		t.Fatal("property 'a' exists")
	}

	if !object.HasProperty("b") {
		t.Fatal("property 'b' not exists")
	}

	if !object.DeleteProperty("b") {
		t.Fatal("could't delete property 'b'")
	}

	if object.HasProperty("b") {
		t.Fatal("delete property 'b' failed")
	}

	// Test ForceDeleteProperty
	if !object.ForceDeleteProperty("x") {
		t.Fatal("could't delete property 'x'")
	}

	if object.HasProperty("x") {
		t.Fatal("delete property 'x' failed")
	}

	// Test HasElement and DeleteElement
	if object.HasElement(1) {
		t.Fatal("element 1 exists")
	}

	if !object.HasElement(0) {
		t.Fatal("element 0 not exists")
	}

	if !object.DeleteElement(0) {
		t.Fatal("could't delete element 0")
	}

	if object.HasElement(0) {
		t.Fatal("delete element 0 failed")
	}

	runtime.GC()
}

func Test_Array(t *testing.T) {
	context := Default.NewContext()
	script := Default.Compile([]byte("[1,2,3]"), nil, nil)
	value := script.Run(context)
	result := value.ToArray()

	if result.Length() != 3 {
		t.Fatal("array length not match")
	}

	if elem := result.GetElement(0); elem != nil {
		if !elem.IsNumber() || elem.ToNumber() != 1 {
			t.Fatal("element 0 value not match")
		}
	} else {
		t.Fatal("could't get element 0")
	}

	if elem := result.GetElement(1); elem != nil {
		if !elem.IsNumber() || elem.ToNumber() != 2 {
			t.Fatal("element 1 value not match")
		}
	} else {
		t.Fatal("could't get element 1")
	}

	if elem := result.GetElement(2); elem != nil {
		if !elem.IsNumber() || elem.ToNumber() != 3 {
			t.Fatal("element 2 value not match")
		}
	} else {
		t.Fatal("could't get element 2")
	}

	if !result.SetElement(0, Default.True()) {
		t.Fatal("could't set element")
	}

	if elem := result.GetElement(0); elem != nil {
		if !elem.IsTrue() {
			t.Fatal("element 0 value not match")
		}
	} else {
		t.Fatal("could't get element 0")
	}

	runtime.GC()
}

func Test_UnderscoreJS(t *testing.T) {
	// Need download underscore.js from:
	// https://raw.github.com/jashkenas/underscore/master/underscore.js
	code, err := ioutil.ReadFile("labs/underscore.js")

	if err != nil {
		return
	}

	context := Default.NewContext()
	script := Default.Compile(code, nil, nil)
	script.Run(context)

	test := []byte("_.find([1, 2, 3, 4, 5, 6], function(num){ return num % 2 == 0; });")
	testScript := Default.Compile(test, nil, nil)
	value := testScript.Run(context)

	if value == nil || value.IsNumber() == false {
		t.FailNow()
	}

	result := value.ToNumber()

	if result != 2 {
		t.FailNow()
	}
}

func rand_sched(max int) {
	for j := rand.Intn(max); j > 0; j-- {
		runtime.Gosched()
	}
}

func Test_ThreadSafe1(t *testing.T) {
	fail := false

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			context := Default.NewContext()
			script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
			value := script.Run(context)
			result := value.ToString()
			fail = fail || result != "Hello World!"
			runtime.GC()
			wg.Done()
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

func Test_ThreadSafe2(t *testing.T) {
	fail := false
	context := Default.NewContext()

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			rand_sched(200)

			script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
			value := script.Run(context)
			result := value.ToString()
			fail = fail || result != "Hello World!"
			runtime.GC()
			wg.Done()
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

func Test_ThreadSafe3(t *testing.T) {
	fail := false
	context := Default.NewContext()
	script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			rand_sched(200)

			value := script.Run(context)
			result := value.ToString()
			fail = fail || result != "Hello World!"
			runtime.GC()
			wg.Done()
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

func Test_ThreadSafe4(t *testing.T) {
	fail := false
	context := Default.NewContext()
	script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
	value := script.Run(context)

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			rand_sched(200)

			result := value.ToString()
			fail = fail || result != "Hello World!"
			runtime.GC()
			wg.Done()
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

func Test_ThreadSafe5(t *testing.T) {
	fail := false
	gonum := 100
	contextChan := make(chan *Context, gonum)
	scriptChan := make(chan *Script, gonum)
	valueChan := make(chan *Value, gonum)

	for i := 0; i < gonum; i++ {
		go func() {
			rand_sched(200)

			contextChan <- Default.NewContext()
		}()
	}

	for i := 0; i < gonum; i++ {
		go func() {
			rand_sched(200)

			scriptChan <- Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
		}()
	}

	for i := 0; i < gonum; i++ {
		go func() {
			rand_sched(200)

			context := <-contextChan
			script := <-scriptChan
			valueChan <- script.Run(context)
		}()
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < gonum; i++ {
		wg.Add(1)
		go func() {
			rand_sched(200)

			value := <-valueChan
			result := value.ToString()
			fail = fail || result != "Hello World!"
			runtime.GC()
			wg.Done()
		}()
	}
	wg.Wait()

	runtime.GC()

	if fail {
		t.FailNow()
	}
}

func Benchmark_NewContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Default.NewContext()
	}

	b.StopTimer()
	runtime.GC()
}

func Benchmark_Compile(b *testing.B) {
	// Need download underscore.js from:
	// https://raw.github.com/jashkenas/underscore/master/underscore.js
	code, err := ioutil.ReadFile("labs/underscore.js")

	if err != nil {
		return
	}

	for i := 0; i < b.N; i++ {
		Default.Compile(code, nil, nil)
	}

	b.StopTimer()
	runtime.GC()
}

func Benchmark_PreCompile(b *testing.B) {
	// Need download underscore.js from:
	// https://raw.github.com/jashkenas/underscore/master/underscore.js
	code, err := ioutil.ReadFile("labs/underscore.js")

	if err != nil {
		return
	}

	b.StopTimer()
	scriptData := Default.PreCompile(code)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Default.Compile(code, nil, scriptData)
	}

	b.StopTimer()
	runtime.GC()
}

func Benchmark_RunScript(b *testing.B) {
	b.StartTimer()
	context := Default.NewContext()
	script := Default.Compile([]byte("1+1"), nil, nil)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		script.Run(context)
	}

	b.StopTimer()
	runtime.GC()
}
