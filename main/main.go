package main

import (
	"fmt"
	"log"
	"reflect"
)

func plusOne(n int) int {
	return n + 1
}
func double(n int) int {
	return 2 * n
}
func logF(in interface{}) string {
	return fmt.Sprintf("log(): %v", in)
}
func print(in string) {
	fmt.Println(in)
}

func NewPipe(fns ...interface{}) (*Pipe, error) {
	p := Pipe{}
	p.fns = make([]reflect.Value, len(fns))
	for i, f := range fns {
		p.fns[i] = reflect.ValueOf(f)
	}
	for i := 1; i < len(fns); i++ {
		err := checkTypes(fns[i-1], fns[i])
		if err != nil {
			return nil, err
		}
	}
	return &p, nil
}

type Pipe struct {
	fns           []reflect.Value
	subscriptions []reflect.Value
}

func (p *Pipe) Subscribe(f interface{}) error {
	err := checkTypes(p.fns[len(p.fns)-1].Interface(), f)
	if err != nil {
		return err
	}
	p.subscriptions = append(p.subscriptions, reflect.ValueOf(f))
	return nil
}
func (p *Pipe) Send(args ...interface{}) {
	send := make([]reflect.Value, len(args))
	for i, arg := range args {
		send[i] = reflect.ValueOf(arg)
	}
	out := p.send(send)
	for _, sub := range p.subscriptions {
		sub.Call(out)
	}
}
func (p *Pipe) send(args []reflect.Value) []reflect.Value {
	for _, f := range p.fns {
		args = f.Call(args)
	}

	return args
}

func checkTypes(f1 interface{}, f2 interface{}) error {
	t1 := reflect.TypeOf(f1)
	if t1.Kind() != reflect.Func {
		return fmt.Errorf("f1 is not a func")
	}
	t2 := reflect.TypeOf(f2)
	if t2.Kind() != reflect.Func {
		return fmt.Errorf("f2 is not a func")
	}
	num := t1.NumOut()
	if num != t2.NumIn() {
		return fmt.Errorf("argument count mismatch")
	}
	for i := 0; i < num; i++ {
		out := t1.Out(i)
		in := t2.In(i)
		if out == in {
			continue
		}
		if out.ConvertibleTo(in) {
			continue
		}
		return fmt.Errorf("arg mismatch between %v and %v", out.Name(), in.Name())
	}
	return nil
}

func main() {
	fmt.Println("Checking f1 to f2")
	p, err := NewPipe(double, plusOne, logF, print)
	if err != nil {
		log.Fatalf("Unable to create pipe: %v\n", err)
	}
	p.Subscribe(func() { fmt.Println("Subscription 1 called") })
	p.Subscribe(func() { fmt.Println("Subscription 2 called") })
	p.Send(1)

}
