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
	context := DefaultEngine.NewContext()
	script := context.CompileScript("'Hello ' + 'World!'")
	value := script.Run()
	result := value.ToString()
	if result != "Hello World!" {
		t.FailNow()
	}
	println(result)
}

func Test_ThreadSafe1(t *testing.T) {
	fail := false

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			context := DefaultEngine.NewContext()
			script := context.CompileScript("'Hello ' + 'World!'")
			value := script.Run()
			result := value.ToString()
			fail = fail || result != "Hello World!"
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
	context := DefaultEngine.NewContext()

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			script := context.CompileScript("'Hello ' + 'World!'")
			value := script.Run()
			result := value.ToString()
			fail = fail || result != "Hello World!"
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
	context := DefaultEngine.NewContext()
	script := context.CompileScript("'Hello ' + 'World!'")

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			value := script.Run()
			result := value.ToString()
			fail = fail || result != "Hello World!"
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
	context := DefaultEngine.NewContext()
	script := context.CompileScript("'Hello ' + 'World!'")
	value := script.Run()

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			result := value.ToString()
			fail = fail || result != "Hello World!"
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
	myrand := rand.New(rand.NewSource(time.Now().UnixNano()))

	fail := false
	gonum := 100
	contextChan := make(chan *Context, gonum)
	scriptChan := make(chan *Script, gonum)
	valueChan := make(chan *Value, gonum)

	wg := new(sync.WaitGroup)
	for i := 0; i < gonum; i++ {
		wg.Add(1)

		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			contextChan <- DefaultEngine.NewContext()
		}()

		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			context := <-contextChan
			scriptChan <- context.CompileScript("'Hello ' + 'World!'")
		}()

		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			script := <-scriptChan
			valueChan <- script.Run()
		}()

		go func() {
			time.Sleep(time.Duration(myrand.Intn(500)) * time.Millisecond)
			value := <-valueChan
			result := value.ToString()
			fail = fail || result != "Hello World!"
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
		DefaultEngine.NewContext()
	}
}

func Benchmark_CompileScript(b *testing.B) {
	b.StartTimer()
	context := DefaultEngine.NewContext()
	scripts := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		scripts[i] = "'Hello ' + '" + strconv.Itoa(i) + "'"
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		context.CompileScript(scripts[i])
	}
}

func Benchmark_RunScript(b *testing.B) {
	b.StartTimer()
	context := DefaultEngine.NewContext()
	script := context.CompileScript("'Hello ' + 'World!'")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		script.Run()
	}
}
