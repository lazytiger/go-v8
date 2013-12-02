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

func Test_GetVersion(t *testing.T) {
	t.Log(GetVersion())
}

func Test_HelloWorld(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		if cs.Eval([]byte("'Hello ' + 'World!'")).ToString() != "Hello World!" {
			t.Fatal("result not match")
		}
	})

	runtime.GC()
}

func Test_TryCatch(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		cs.TryCatch(true, func() {
			Default.Compile([]byte("a[=1"), nil, nil)
		})

		if cs.TryCatch(true, func() {
			cs.ThrowException("this is error")
		}) != "this is error" {
			t.Fatal("error message not match")
		}
	})

	runtime.GC()
}

func Test_PreCompile(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		// pre-compile
		code := []byte("'Hello ' + 'PreCompile!'")
		scriptData1 := Default.PreCompile(code)

		if scriptData1 == nil {
			t.Fatal("precompile failed")
		}

		// test save and load script data
		data := scriptData1.Data()
		scriptData2 := NewScriptData(data)

		if scriptData1 == nil {
			t.Fatal("load precompile data failed")
		}

		// test compile with script data
		script := Default.Compile(code, nil, scriptData2)
		value := script.Run()
		result := value.ToString()

		if result != "Hello PreCompile!" {
			t.Fatal("result not match")
		}
	})

	runtime.GC()
}

func Test_Values(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {

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

		if cs.NewBoolean(true).ToBoolean() != true {
			t.Fatal(`NewBoolean(true).ToBoolean() != true`)
		}

		if cs.NewNumber(12.34).ToNumber() != 12.34 {
			t.Fatal(`NewNumber(12.34).ToNumber() != 12.34`)
		}

		if cs.NewNumber(float64(maxNumber)).ToInteger() != maxNumber {
			t.Fatal(`NewNumber(float64(maxNumber)).ToInteger() != maxNumber`)
		}

		if cs.NewInteger(maxInt32).IsInt32() == false {
			t.Fatal(`NewInteger(maxInt32).IsInt32() == false`)
		}

		if cs.NewInteger(maxUint32).IsInt32() != false {
			t.Fatal(`NewInteger(maxUint32).IsInt32() != false`)
		}

		if cs.NewInteger(maxUint32).IsUint32() == false {
			t.Fatal(`NewInteger(maxUint32).IsUint32() == false`)
		}

		if cs.NewInteger(maxNumber).ToInteger() != maxNumber {
			t.Fatal(`NewInteger(maxNumber).ToInteger() != maxNumber`)
		}

		if cs.NewString("Hello World!").ToString() != "Hello World!" {
			t.Fatal(`NewString("Hello World!").ToString() != "Hello World!"`)
		}

		if cs.NewObject().IsObject() == false {
			t.Fatal(`NewObject().IsObject() == false`)
		}

		if cs.NewArray(5).IsArray() == false {
			t.Fatal(`NewArray(5).IsArray() == false`)
		}

		if cs.NewArray(5).ToArray().Length() != 5 {
			t.Fatal(`NewArray(5).Length() != 5`)
		}

		if cs.NewRegExp("foo", RF_None).IsRegExp() == false {
			t.Fatal(`NewRegExp("foo", RF_None).IsRegExp() == false`)
		}

		if cs.NewRegExp("foo", RF_Global).ToRegExp().Pattern() != "foo" {
			t.Fatal(`NewRegExp("foo", RF_Global).ToRegExp().Pattern() != "foo"`)
		}

		if cs.NewRegExp("foo", RF_Global).ToRegExp().Flags() != RF_Global {
			t.Fatal(`NewRegExp("foo", RF_Global).ToRegExp().Flags() != RF_Global`)
		}
	})

	runtime.GC()
}

func Test_Object(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		script := Default.Compile([]byte("a={};"), nil, nil)
		value := script.Run()
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

		// Test GetPropertyNames
		script = Default.Compile([]byte("a={x:10,y:20,z:30};"), nil, nil)
		value = script.Run()
		object = value.ToObject()

		names := object.GetPropertyNames()

		if names.Length() != 3 {
			t.Fatal(`names.Length() != 3`)
		}

		if names.GetElement(0).ToString() != "x" {
			t.Fatal(`names.GetElement(0).ToString() != "x"`)
		}

		if names.GetElement(1).ToString() != "y" {
			t.Fatal(`names.GetElement(1).ToString() != "y"`)
		}

		if names.GetElement(2).ToString() != "z" {
			t.Fatal(`names.GetElement(2).ToString() != "z"`)
		}
	})

	runtime.GC()
}

func Test_Array(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		script := Default.Compile([]byte("[1,2,3]"), nil, nil)
		value := script.Run()
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
	})

	runtime.GC()
}

func Test_Function(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		script := Default.Compile([]byte(`
			a = function(x,y,z){ 
				return x+y+z; 
			}
		`), nil, nil)

		value := script.Run()

		if value.IsFunction() == false {
			t.Fatal("value not a function")
		}

		result := value.ToFunction().Call(
			cs.NewInteger(1),
			cs.NewInteger(2),
			cs.NewInteger(3),
		)

		if result.IsNumber() == false {
			t.Fatal("result not a number")
		}

		if result.ToInteger() != 6 {
			t.Fatal("result != 6")
		}

		function := Default.NewFunctionTemplate(func(info FunctionCallbackInfo) {
			if info.Get(0).ToString() != "Hello World!" {
				t.Fatal(`info.Get(0).ToString() != "Hello World!"`)
			}
			info.ReturnValue().SetBoolean(true)
		}).NewFunction()

		if function == nil {
			t.Fatal("function == nil")
		}

		if function.ToFunction().Call(
			cs.NewString("Hello World!"),
		).IsTrue() == false {
			t.Fatal("callback return not match")
		}
	})

	runtime.GC()
}

func Test_Accessor(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		template := Default.NewObjectTemplate()
		var propertyValue int32

		template.SetAccessor(
			"abc",
			func(name string, info GetterCallbackInfo) {
				data := info.Data().(*int32)
				info.ReturnValue().SetInt32(*data)
			},
			func(name string, value *Value, info SetterCallbackInfo) {
				data := info.Data().(*int32)
				*data = value.ToInt32()
			},
			&propertyValue,
			PA_None,
		)

		template.SetProperty("def", cs.NewInteger(8888), PA_None)

		values := []*Value{
			template.NewObject(), // New
			cs.NewObject(),       // Wrap
		}
		template.WrapObject(values[1])

		for i := 0; i < 2; i++ {
			value := values[i]

			propertyValue = 1234

			object := value.ToObject()

			if object.GetProperty("abc").ToInt32() != 1234 {
				t.Fatal(`object.GetProperty("abc").ToInt32() != 1234`)
			}

			object.SetProperty("abc", cs.NewInteger(5678), PA_None)

			if propertyValue != 5678 {
				t.Fatal(`propertyValue != 5678`)
			}

			if object.GetProperty("abc").ToInt32() != 5678 {
				t.Fatal(`object.GetProperty("abc").ToInt32() != 5678`)
			}

			if object.GetProperty("def").ToInt32() != 8888 {
				t.Fatal(`object.GetProperty("def").ToInt32() != 8888`)
			}
		}
	})

	runtime.GC()
}

func Test_NamedPropertyHandler(t *testing.T) {
	obj_template := Default.NewObjectTemplate()

	var (
		get_called    = false
		set_called    = false
		query_called  = false
		delete_called = false
		enum_called   = false
	)

	obj_template.SetNamedPropertyHandler(
		func(name string, info PropertyCallbackInfo) {
			//t.Logf("get %s", name)
			get_called = get_called || name == "abc"
		},
		func(name string, value *Value, info PropertyCallbackInfo) {
			//t.Logf("set %s", name)
			set_called = set_called || name == "abc"
		},
		func(name string, info PropertyCallbackInfo) {
			//t.Logf("query %s", name)
			query_called = query_called || name == "abc"
		},
		func(name string, info PropertyCallbackInfo) {
			//t.Logf("delete %s", name)
			delete_called = delete_called || name == "abc"
		},
		func(info PropertyCallbackInfo) {
			//t.Log("enumerate")
			enum_called = true
		},
		nil,
	)

	func_template := Default.NewFunctionTemplate(func(info FunctionCallbackInfo) {
		info.ReturnValue().Set(obj_template.NewObject())
	})

	global_template := Default.NewObjectTemplate()

	global_template.SetAccessor("GetData", func(name string, info GetterCallbackInfo) {
		info.ReturnValue().Set(func_template.NewFunction())
	}, nil, nil, PA_None)

	Default.NewContext(global_template).Scope(func(cs ContextScope) {
		object := obj_template.NewObject().ToObject()

		object.GetProperty("abc")
		object.SetProperty("abc", cs.NewInteger(123), PA_None)
		object.GetPropertyAttributes("abc")

		cs.Eval([]byte(`
			var data = GetData();

			delete data.abc;

			for (var p in data) {
			}
		`))
	})

	if !(get_called && set_called && query_called && delete_called && enum_called) {
		t.Fatal(get_called, set_called, query_called, delete_called, enum_called)
	}

	runtime.GC()
}

func Test_IndexedPropertyHandler(t *testing.T) {
	obj_template := Default.NewObjectTemplate()

	var (
		get_called    = false
		set_called    = false
		query_called  = true // TODO
		delete_called = true // TODO
		enum_called   = true // TODO
	)

	obj_template.SetIndexedPropertyHandler(
		func(index uint32, info PropertyCallbackInfo) {
			//t.Logf("get %d", index)
			get_called = get_called || index == 1
		},
		func(index uint32, value *Value, info PropertyCallbackInfo) {
			//t.Logf("set %d", index)
			set_called = set_called || index == 1
		},
		func(index uint32, info PropertyCallbackInfo) {
			//t.Logf("query %d", index)
			query_called = query_called || index == 1
		},
		func(index uint32, info PropertyCallbackInfo) {
			//t.Logf("delete %d", index)
			delete_called = delete_called || index == 1
		},
		func(info PropertyCallbackInfo) {
			//t.Log("enumerate")
			enum_called = true
		},
		nil,
	)

	func_template := Default.NewFunctionTemplate(func(info FunctionCallbackInfo) {
		info.ReturnValue().Set(obj_template.NewObject())
	})

	global_template := Default.NewObjectTemplate()

	global_template.SetAccessor("GetData", func(name string, info GetterCallbackInfo) {
		info.ReturnValue().Set(func_template.NewFunction())
	}, nil, nil, PA_None)

	Default.NewContext(global_template).Scope(func(cs ContextScope) {
		object := obj_template.NewObject().ToObject()

		object.GetElement(1)
		object.SetElement(1, cs.NewInteger(123))

		cs.Eval([]byte(`
			var data = GetData();

			delete data[1];

			for (var p in data) {
			}
		`))
	})

	if !(get_called && set_called && query_called && delete_called && enum_called) {
		t.Fatal(get_called, set_called, query_called, delete_called, enum_called)
	}

	runtime.GC()
}

func Test_ObjectConstructor(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		ftConstructor := Default.NewFunctionTemplate(nil)
		ftConstructor.SetClassName("MyClass")

		obj_template := ftConstructor.InstanceTemplate()

		var (
			get_called    = false
			set_called    = false
			query_called  = false
			delete_called = false
			enum_called   = false
		)

		obj_template.SetNamedPropertyHandler(
			func(name string, info PropertyCallbackInfo) {
				//t.Logf("get %s", name)
				get_called = get_called || name == "abc"
			},
			func(name string, value *Value, info PropertyCallbackInfo) {
				//t.Logf("set %s", name)
				set_called = set_called || name == "abc"
			},
			func(name string, info PropertyCallbackInfo) {
				//t.Logf("query %s", name)
				query_called = query_called || name == "abc"
			},
			func(name string, info PropertyCallbackInfo) {
				//t.Logf("delete %s", name)
				delete_called = delete_called || name == "abc"
			},
			func(info PropertyCallbackInfo) {
				//t.Log("enumerate")
				enum_called = true
			},
			nil,
		)

		cs.Global().SetProperty("MyClass", ftConstructor.NewFunction(), PA_None)

		if !cs.Eval([]byte("(new MyClass) instanceof MyClass")).IsTrue() {
			t.Fatal("(new MyClass) instanceof MyClass == false")
		}

		object := cs.Eval([]byte(`
			var data = new MyClass;
			var temp = data.abc;
			data.abc = 1;
			delete data.abc;
			for (var p in data) {
			}
			data;
		`)).ToObject()

		object.GetPropertyAttributes("abc")

		if !(get_called && set_called && query_called && delete_called && enum_called) {
			t.Fatal(get_called, set_called, query_called, delete_called, enum_called)
		}
	})

	runtime.GC()
}

func Test_Context(t *testing.T) {
	script1 := Default.Compile([]byte("typeof(Test_Context) == 'undefined';"), nil, nil)
	script2 := Default.Compile([]byte("Test_Context = 1;"), nil, nil)
	script3 := Default.Compile([]byte("Test_Context = Test_Context + 7;"), nil, nil)

	test_func := func(cs ContextScope) {
		if script1.Run().IsFalse() {
			t.Fatal(`script1.Run(c).IsFalse()`)
		}

		if script2.Run().ToInteger() != 1 {
			t.Fatal(`script2.Run(c).ToInteger() != 1`)
		}

		if script3.Run().ToInteger() != 8 {
			t.Fatal(`script3.Run(c).ToInteger() != 8`)
		}
	}

	Default.NewContext(nil).Scope(func(cs ContextScope) {
		Default.NewContext(nil).Scope(test_func)
		Default.NewContext(nil).Scope(test_func)
		test_func(cs)
	})

	functionTemplate := Default.NewFunctionTemplate(func(info FunctionCallbackInfo) {
		for i := 0; i < info.Length(); i++ {
			println(info.Get(i).ToString())
		}
	})

	// Test Global Template
	globalTemplate := Default.NewObjectTemplate()

	globalTemplate.SetAccessor("log", func(name string, info GetterCallbackInfo) {
		info.ReturnValue().Set(functionTemplate.NewFunction())
	}, nil, nil, PA_None)

	Default.NewContext(globalTemplate).Scope(func(cs ContextScope) {
		cs.Eval([]byte(`log("Hello World!")`))
	})

	// Test Global Object
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		global := cs.Global()

		if !global.SetProperty("println", functionTemplate.NewFunction(), PA_None) {
		}

		global = cs.Global()

		if !global.HasProperty("println") {
			t.Fatal(`!global.HasProperty("println")`)
			return
		}

		cs.Eval([]byte(`println("Hello World!")`))
	})

	runtime.GC()
}

func Test_UnderscoreJS(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		code, err := ioutil.ReadFile("labs/underscore.js")

		if err != nil {
			return
		}

		script := Default.Compile(code, nil, nil)
		script.Run()

		test := []byte(`
			_.find([1, 2, 3, 4, 5, 6], function(num) { 
				return num % 2 == 0; 
			});
		`)
		testScript := Default.Compile(test, nil, nil)
		value := testScript.Run()

		if value == nil || value.IsNumber() == false {
			t.FailNow()
		}

		result := value.ToNumber()

		if result != 2 {
			t.FailNow()
		}
	})

	runtime.GC()
}

func Test_JSON(t *testing.T) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		json := `{"a":1,"b":2,"c":"xyz","e":true,"f":false,"g":null,"h":[4,5,6]}`

		value := cs.ParseJSON(json)

		if value == nil {
			t.Fatal(`value == nil`)
		}

		if value.IsObject() == false {
			t.Fatal(`value == false`)
		}

		if string(ToJSON(value)) != json {
			t.Fatal(`string(ToJSON(value)) != json`)
		}

		object := value.ToObject()

		if object.GetProperty("a").ToInt32() != 1 {
			t.Fatal(`object.GetProperty("a").ToInt32() != 1`)
		}

		if object.GetProperty("b").ToInt32() != 2 {
			t.Fatal(`object.GetProperty("b").ToInt32() != 2`)
		}

		if object.GetProperty("c").ToString() != "xyz" {
			t.Fatal(`object.GetProperty("c").ToString() != "xyz"`)
		}

		if object.GetProperty("e").IsTrue() == false {
			t.Fatal(`object.GetProperty("e").IsTrue() == false`)
		}

		if object.GetProperty("f").IsFalse() == false {
			t.Fatal(`object.GetProperty("f").IsFalse() == false`)
		}

		if object.GetProperty("g").IsNull() == false {
			t.Fatal(`object.GetProperty("g").IsNull() == false`)
		}

		array := object.GetProperty("h").ToArray()

		if array.Length() != 3 {
			t.Fatal(`array.Length() != 3`)
		}

		if array.GetElement(0).ToInt32() != 4 {
			t.Fatal(`array.GetElement(0).ToInt32() != 4`)
		}

		if array.GetElement(1).ToInt32() != 5 {
			t.Fatal(`array.GetElement(1).ToInt32() != 5`)
		}

		if array.GetElement(2).ToInt32() != 6 {
			t.Fatal(`array.GetElement(2).ToInt32() != 6`)
		}

		json = `"\"\/\r\n\t\b\\"`

		if string(ToJSON(cs.ParseJSON(json))) != json {
			t.Fatal(`ToJSON(cs.ParseJSON(json)) != json`)
		}
	})

	runtime.GC()
}

func rand_sched(max int) {
	for j := rand.Intn(max); j > 0; j-- {
		runtime.Gosched()
	}
}

// use one engine in different threads
//
func Test_ThreadSafe1(t *testing.T) {
	fail := false

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Default.NewContext(nil).Scope(func(cs ContextScope) {
				script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
				value := script.Run()
				result := value.ToString()
				fail = fail || result != "Hello World!"
				runtime.GC()
				wg.Done()
			})
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

// use one context in different threads
//
func Test_ThreadSafe2(t *testing.T) {
	fail := false
	context := Default.NewContext(nil)

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			context.Scope(func(cs ContextScope) {
				rand_sched(200)

				script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
				value := script.Run()
				result := value.ToString()
				fail = fail || result != "Hello World!"
				runtime.GC()
				wg.Done()
			})
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

// use one script in different threads
//
func Test_ThreadSafe3(t *testing.T) {
	fail := false
	script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Default.NewContext(nil).Scope(func(cs ContextScope) {
				rand_sched(200)

				value := script.Run()
				result := value.ToString()
				fail = fail || result != "Hello World!"
				runtime.GC()
				wg.Done()
			})
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

// use one context and one script in different threads
//
func Test_ThreadSafe4(t *testing.T) {
	fail := false
	script := Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
	context := Default.NewContext(nil)

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			context.Scope(func(cs ContextScope) {
				rand_sched(200)

				value := script.Run()
				result := value.ToString()
				fail = fail || result != "Hello World!"
				runtime.GC()
				wg.Done()
			})
		}()
	}
	wg.Wait()
	runtime.GC()

	if fail {
		t.FailNow()
	}
}

// ....
//
func Test_ThreadSafe6(t *testing.T) {
	var (
		fail        = false
		gonum       = 100
		scriptChan  = make(chan *Script, gonum)
		contextChan = make(chan *Context, gonum)
	)

	for i := 0; i < gonum; i++ {
		go func() {
			rand_sched(200)

			scriptChan <- Default.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
		}()
	}

	for i := 0; i < gonum; i++ {
		go func() {
			rand_sched(200)

			contextChan <- Default.NewContext(nil)
		}()
	}

	for i := 0; i < gonum; i++ {
		go func() {
			rand_sched(200)

			context := <-contextChan
			script := <-scriptChan

			context.Scope(func(cs ContextScope) {
				result := script.Run().ToString()
				fail = fail || result != "Hello World!"
			})
		}()
	}

	runtime.GC()

	if fail {
		t.FailNow()
	}
}

func Benchmark_NewContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Default.NewContext(nil)
	}

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_NewInteger(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.NewInteger(int64(i))
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_NewString(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.NewString("Hello World!")
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_NewObject(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.NewObject()
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_NewArray0(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.NewArray(0)
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_NewArray5(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.NewArray(5)
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_NewArray20(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.NewArray(20)
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_NewArray100(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.NewArray(100)
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_Compile(b *testing.B) {
	b.StopTimer()
	code, err := ioutil.ReadFile("labs/underscore.js")
	if err != nil {
		return
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Default.Compile(code, nil, nil)
	}

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_PreCompile(b *testing.B) {
	b.StopTimer()
	code, err := ioutil.ReadFile("labs/underscore.js")
	if err != nil {
		return
	}
	scriptData := Default.PreCompile(code)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Default.Compile(code, nil, scriptData)
	}

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_RunScript(b *testing.B) {
	b.StopTimer()
	context := Default.NewContext(nil)
	script := Default.Compile([]byte("1+1"), nil, nil)
	b.StartTimer()

	context.Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			script.Run()
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_JsFunction(b *testing.B) {
	b.StopTimer()

	script := Default.Compile([]byte(`
		a = function(){ 
			return 1; 
		}
	`), nil, nil)

	Default.NewContext(nil).Scope(func(cs ContextScope) {
		value := script.Run()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			value.ToFunction().Call()
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_GoFunction(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		b.StopTimer()
		value := Default.NewFunctionTemplate(func(info FunctionCallbackInfo) {
			info.ReturnValue().SetInt32(123)
		}).NewFunction()
		function := value.ToFunction()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			function.Call()
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_Getter(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		b.StopTimer()
		var propertyValue int32 = 1234

		template := Default.NewObjectTemplate()

		template.SetAccessor(
			"abc",
			func(name string, info GetterCallbackInfo) {
				data := info.Data().(*int32)
				info.ReturnValue().SetInt32(*data)
			},
			func(name string, value *Value, info SetterCallbackInfo) {
				data := info.Data().(*int32)
				*data = value.ToInt32()
			},
			&propertyValue,
			PA_None,
		)

		object := template.NewObject().ToObject()

		b.StartTimer()

		for i := 0; i < b.N; i++ {
			object.GetProperty("abc")
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_Setter(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		b.StopTimer()

		var propertyValue int32 = 1234

		template := Default.NewObjectTemplate()

		template.SetAccessor(
			"abc",
			func(name string, info GetterCallbackInfo) {
				data := info.Data().(*int32)
				info.ReturnValue().SetInt32(*data)
			},
			func(name string, value *Value, info SetterCallbackInfo) {
				data := info.Data().(*int32)
				*data = value.ToInt32()
			},
			&propertyValue,
			PA_None,
		)

		object := template.NewObject().ToObject()

		b.StartTimer()

		for i := 0; i < b.N; i++ {
			object.SetProperty("abc", cs.NewInteger(5678), PA_None)
		}
	})

	b.StopTimer()
	runtime.GC()
	b.StartTimer()
}

func Benchmark_TryCatch(b *testing.B) {
	Default.NewContext(nil).Scope(func(cs ContextScope) {
		for i := 0; i < b.N; i++ {
			cs.TryCatch(false, func() {
				cs.Eval([]byte("a[=1;"))
			})
		}
	})
}
