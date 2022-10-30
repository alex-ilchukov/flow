# Flow

Implementation of pipeline written in Go and inspired by 
[this article](https://go.dev/blog/pipelines).


## The design

### Basic pipeline

The core of any pipeline implementation in Go is flow of data through channel
or channels, if the data should be transformed, in non-blocking manner with
help of go-routines. So the very basic algorithm of pipeline construction can 
be described in the following way.

Emitting stage:
*   create a channel of data C1;
*   launch a go-routine, where the data is pushed into the channel;
*   return the channel C1.

Transforming stage 1:
*   take the channel C1 from the previous step;
*   create a channel of transformed data C2;
*   launch a go-routine, where the data is read from C1, transformed in
    some kind of way, and pushed to C2 in non-blocking way;
*   return the channel C2.

Transforming stage 2:
*   take the channel C2 from the previous step;
*   create a channel of transformed data C3;
*   launch a go-routine, where the data is read from C2, transformed in
    some kind of way, and pushed to C3 in non-blocking way;
*   return the channel C3.

…

Transforming stage n:
*   take the channel C(n-1) from the previous step;
*   create a channel of transformed data Cn;
*   launch a go-routine, where the data is read from C(n-1), transformed in
    some kind of way, and pushed to Cn in non-blocking way;
*   return the channel Cn.

Collecting stage:
*   take the channel Cn from the previous step;
*   launch a go-routine, where the data is read from Cn and collected in some
    manner in non-blocking-way.

There can be any number of transforming stages (even zero), but obviously there
are always emitting stage and collecting stage. All the launched go-routines 
happily go through the data in non-blocking way and actually know nothing of
other stages. This allows to describe the pipeline in very abstract, decoupled
manner.

### Error handling

All the stages described above can produce error values. That adds a bit to the 
basic algorithm. As all the processing goes in non-blocking way, the only 
proper enhancement would be to let the stages to return not only data channels, 
but also error channels. So the enhanced pipeline with error handling should 
gather all the error channels and listen them in the same non-blocking manner 
for possible error values.

### Context

Besides error values, there is also possible cancellation of the pipeline.
Standard library of Go provides implementation of 
[`context.Context`](https://pkg.go.dev/context@go1.19.2#Context) interface, and
its `Done()` method can be used to check on pipeline cancel during any data
pushing.

### Running a pipeline

After a pipeline is constructed, its should user should wait til all the 
go-routines are properly finished or an error happened. This could be reduced 
to just a one error channel, where all the error channels mentioned above push 
their errors to in non-blocking manner, and its listening in _blocking_ manner. 
If any error appears in the channel, user of the pipeline can decide if the
error is worth of cancelling the run (via context). If the channel is closed,
that means the whole success.

### The core design

So the project is about to implement algorithm of running a pipeline in generic 
way with proper error handling and use of the context. The following choices 
have been made for the implementation details.

#### Used terminology

The project uses "emitters", "collectors", and "transformers" to name entities
which describe stages. Other implementations use other words like "sources",
"drains", "sinks", and so on. The terminology is important, as it allows to
name interfaces and their methods in consistent way. An emitter obviously 
emits, a collector collects, a transformer transforms, but what a sink does? 
A drain?

#### Interfaces over functions and closures

The original article linked above is all about functions and closures, as the
same are many implementations of the mechanism. Nevertheless, the way is less
clean, as the functions call other functions, these ones call another ones,
every function makes closures et cetera — all in messy sheets of dozens lines
of code. Interfaces allow to describe things in more clean way, and their
implementations can use all the power of encapsulation and details hiding.

#### Squashing of transformer stages

It is pretty obvious, that a bundle of emitter and transformer with channels
linked is nothing new, but another emitter. That allows to describe the
pipeline in terms of emitter and collector only, which greatly simplifies
things and enables late introduction to transformers. (Technically, there could
be another way to introduce them via collectors, but that would mean backward
pipeline construction.)

## Installation

To install Flow package, you need to install Go (of version 1.19 or higher), 
get the package, and import it.

### Getting the package

```sh
go get -u github.com/alex-ilchukov/flow
```

### Import of the package

```go
import "github.com/alex-ilchukov/flow"
```


## Usage

TODO: fill the section
