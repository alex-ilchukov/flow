// Package flow provides core concepts of data flow with support of error
// handling and possible cancellation of the flow via context mechanisms. The
// implementation is inspired by articles in [Go blog] and [Medium].
//
// [Go blog]: https://go.dev/blog/pipelines
// [Medium]: https://medium.com/amboss/applying-modern-go-concurrency-patterns-to-data-pipelines-b3b5327908d4
package flow
