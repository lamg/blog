# Error handling in [Go][0] (part 0)

[Go][0] has according part of its community a lack of adequate syntax for error handling, and that is the reason behind proposals like [adding check and handle statements][1]. I claim they are wrong, [Go][0] doesn't need more syntax for error handling. Errors in [Go][0] are just values, and the language can handle values and structures of them without problem. This is the way I do it:

Years ago I found the writings of [Edsger W. Dijkstra][2], and a very simple —but vital— lesson he taught is that programs can be made more simple to understand if arbitrary jumps are avoided, and instead rely on statements that make explicit and manageable the program's structure. Loops don't need the `break` command and procedures don't need the `return` command in the middle of their body ([Oberon-07][3] is an example), you can get the same effect using boolean values and the available semantic for conditional and repetitive statements.

Most of you when faced with a sequence of operations returning `error` values will use the return command after finding one of them different from `nil`. For example:

```go
func doSomething() (r *record, e error) {
	v, e := foo()
	if e != nil {
		return
	}
	u, e := bar()
	if e != nil {
		return
	}
	r = &record{V:v, U:u}
	return
}
```

And yes, there is a pattern that the language cannot express with brevity, but it appears because you keep sticking to the idea of jumping without structure. You will always find a way of jumping that will make the language look poor, because in general they are made for getting sequences of individual operations without having to write them all. The complex flow of operations you need can be achieved by combining simple statements that change the execution path by evaluating variables. Let's try that:

```go
func doSomething() (r *record, e error){
	v, e := foo()
	var u string
	if e == nil {
		u, e = bar()
	}
	if e == nil {
		r = &record{V:v, U:u}
	}
	return
}
```

Well, the pattern keeps appearing, and now if `e != nil` after calling `foo` or `bar`, at least a superfluous evaluation of `e == nil` will be made. But the way to avoid execution of statements after a condition it's false is by making a loop. You will argue that we cannot make a loop in this case since we have two different commands executed under the condition `e == nil` each time. But from the loop's perspective both look the same: set the value of `e` and something else. Using closures it's possible to hide or expose variables when they are irrelevant or relevant, respectively, for the context:

```go
func doSomething() (r *record, e error) {
	v, e := foo()
	var u string
	fs := []func(){
		func() { u, e = bar() },
		func() { r = &record{V: v, U: u} },
	}
	n := 0
	for e == nil && n != len(fs) {
		fs[n]()
		n = n + 1
	}
	return
}
```

Now there are no more superfluous evaluations of `e == nil` and the procedure keeps doing the same. Also the search for an `error` value not `nil` at the end of `doSomething` is an essential algorithm, the [Bounded Linear Search][4], and it doesn't depend on the `error` type but on a boolean function: 

```go
func bLnSearch(bf func(int) bool, n int) (ok bool, i int) {
	i = 0
	for i != n && !bf(i) {
		i = i + 1
	}
	ok = i != n
	return
}
```

And now `doSomething` is:

```go
func doSomething() (r *record, e error) {
	v, e := foo()
	var u string
	fs := []func(){
		func() { u, e = bar() },
		func() { r = &record{V: v, U: u} },
	}
	bLnSearch(
		func(i int) bool {
			b := e != nil
			if !b {
				fs[i]()
			}
			return b
		},
		len(fs),
	)
	return
}
```

Previously, the need for `b := e != nil` before `fs[i]()` came from not loosing the value of `e != nil` dependenig on the value of `e` returned after `v, e := foo()`, but if the latter is included in the loop's execution then the irregularity will be swept away:

```go
func doSomething() (r *record, e error) {
	var v int
	var u string
	fs := []func(){
		func() { v, e = foo() },
		func() { u, e = bar() },
		func() { r = &record{V: v, U: u} },
	}
	bLnSearch(
		func(i int) bool { fs[i](); return e != nil },
		len(fs),
	)
	return
}
``` 

## Conclusion

Once the source of the complexity is detected —execution paths are specified with irrelevant detail and not looking for structure— there is a nice way of not getting entangled with it, and it relies on basic resources like the [Bounded Linear Search][4] and closures, known since the beginning of computer programming. With that in mind, I suspect features like [generics][5] and the [check and handle statements][1] come from problems that haven't been exposed or analyzed with commitment to find the simplest solution.

[0]: https://golang.org
[1]: https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling-overview.md 
[2]: https://www.cs.utexas.edu/users/EWD
[3]: http://www.inf.ethz.ch/personal/wirth/Oberon/Oberon07.Report.pdf 
[4]: https://www.cs.utexas.edu/users/EWD/transcriptions/EWD09xx/EWD930.html
[5]: https://go.googlesource.com/proposal/+/master/design/go2draft-generics-overview.md 
