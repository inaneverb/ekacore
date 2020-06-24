![ekago_logo](https://user-images.githubusercontent.com/16417743/85555445-9a32e900-b62e-11ea-9a38-464199ff08e5.jpg)

# `[ekaerr]` An error management you deserve

```go
import "github.com/qioalice/ekago/ekaerr"
```

## What?

`ekaerr` is an error generating and managing package which is extremely linked with [`ekalog`: _intelligence logging package_](../ekalog/). 

Now you can:
-  Generate stacktrace based error objects with a possibility to adding fields (arguments), additional messages to **each** stackframe;
- Add _Public error message_ which you can use later for example as a part of API response;
- Have an unique `error_id` for each error object. UUID to be more precision;
- Divide errors by _error classes_, use prepared built-in error classes like `DataUnavailable`, `IllegalArgument`, and others or declare your own.


## Why? 

Golang has very primitive error management by default. It's just error interface with `Error() string` method. All you can do w/o extern libraries is create a simple error message with some text, nothing more: `errors.New(...)` or `fmt.Errorf(...)`. 

What if you want stacktrace? Message's fields like IP addr, user's ID? Or add auxiliary messages for each stacktrace's frame that could be describe what happened in details?

Imagine you define a func, that may be failed:
```go
func foo1() error { ... }
```
What do you do if func has been completed with `err != nil` ? Log it, i guess. But what if `foo1` has been called inside `foo2` that has another functions calls and also returns an `error` object? Like
```go
func foo2(arg1, arg2 interface{}) error {
	// You has code here
	if err := foo1(); err != nil {
		// Log? Like "foo2 has been failed because of foo1"?
		// But what's about arguments 'arg1', 'arg2'? 
		// Are they important?
	}
	// And here you also has a code
}
```
So, it's already complicated, isn't? What if you have `foo3`, `foo4`, `foo5`? Nested? Or instead of code placeholders in the example above?

Moreover. 
Ok, let's pretend you logged it somehow. What's about API? Will you just return HTTP 500 code to your user if it's web server? Or how do you make it clear for user what's happened when you need
- Log as more info as you can to perfectly understand what happens and fix it easiest way
- Tell to the user only common regular info, nothing more, especially your runtime private data

at the same time?

And now imagine how it would be great if:
- You accumulate as more error related data as you want and all that as one error object
- You add your custom arguments, messages
- You log it
- You send its unique ID to the user's (no more unclear HTTP 500 or _"an error has been occurred"_)
- User can write to your support about error using unique ID
- Support can find all related error's info by it's ID and send it to your tech specialists
- Tech specialists (of course it's you) do fixes

Imagine, 
Send it to a client and imagine how easy it will be to find an associated error object later, especially when you will have only one error object

## How?

< in development >

-----

Inspired by [joomcode/errorx](https://github.com/joomcode/errorx).