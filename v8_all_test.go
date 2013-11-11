package v8

import "sync"
import "testing"
import "runtime"
import "time"
import "math/rand"
import "strconv"

func init() {
	//traceDispose = true
}

func Test_HelloWorld(t *testing.T) {
	context := Default.NewContext()
	script := context.Compile("'Hello ' + 'World!'", nil, nil)
	value := script.Run(context)
	result := value.ToString()
	if result != "Hello World!" {
		t.FailNow()
	}
	//println(result)
}

func Test_ThreadSafe1(t *testing.T) {
	fail := false

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			context := Default.NewContext()
			script := context.Compile("'Hello ' + 'World!'", nil, nil)
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
	myrand := rand.New(rand.NewSource(time.Now().UnixNano()))

	fail := false
	context := Default.NewContext()

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			script := context.Compile("'Hello ' + 'World!'", nil, nil)
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
	myrand := rand.New(rand.NewSource(time.Now().UnixNano()))

	fail := false
	context := Default.NewContext()
	script := context.Compile("'Hello ' + 'World!'", nil, nil)

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
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
	myrand := rand.New(rand.NewSource(time.Now().UnixNano()))

	fail := false
	context := Default.NewContext()
	script := context.Compile("'Hello ' + 'World!'", nil, nil)
	value := script.Run(context)

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
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
	contextChan := make(chan *Context, gonum*2)
	scriptChan := make(chan *Script, gonum)
	valueChan := make(chan *Value, gonum)

	wg := new(sync.WaitGroup)
	for i := 0; i < gonum; i++ {
		wg.Add(1)

		go func() {
			contextChan <- Default.NewContext()
		}()

		go func() {
			context := Default.NewContext()
			scriptChan <- context.Compile("'Hello ' + 'World!'", nil, nil)
		}()

		go func() {
			context := <-contextChan
			script := <-scriptChan
			valueChan <- script.Run(context)
		}()

		go func() {
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

func Test_PreCompile(t *testing.T) {
	code := "'Hello ' + 'PreCompile!'"
	scriptData1 := Default.PreCompile(code)

	data := scriptData1.Data()
	scriptData2 := NewScriptData(data)
	context := Default.NewContext()
	script := context.Compile(code, nil, scriptData2)
	value := script.Run(context)
	result := value.ToString()
	if result != "Hello PreCompile!" {
		t.FailNow()
	}
	//println(result)
}

func Benchmark_NewContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Default.NewContext()
	}
}

func Benchmark_Compile(b *testing.B) {
	b.StartTimer()
	context := Default.NewContext()
	scripts := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		scripts[i] = `function myfunc(a, b) { 
			return 'Hello ' + '` + strconv.Itoa(i) + `' + a + b
		}`
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.Compile(scripts[i], nil, nil)
	}
}

func Benchmark_PreCompile(b *testing.B) {
	b.StartTimer()
	context := Default.NewContext()
	scripts := make([]string, b.N)
	scriptDatas := make([]*ScriptData, b.N)
	for i := 0; i < b.N; i++ {
		scripts[i] = `function myfunc(a, b) { 
			return 'Hello ' + '` + strconv.Itoa(i) + `' + a + b
		}`
		scriptDatas[i] = Default.PreCompile(scripts[i])
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.Compile(scripts[i], nil, scriptDatas[i])
	}
}

func Benchmark_RunScript(b *testing.B) {
	b.StartTimer()
	context := Default.NewContext()
	script := context.Compile("'Hello ' + 'World!'", nil, nil)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		script.Run(context)
	}
}
