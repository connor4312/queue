/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/
package queue

import (
	"errors"
)

const minQueueLen = 16

// Queue represents a single instance of the queue data structure.
type Queue struct {
	buf               []interface{}
	head, tail, count int
}

// New constructs and returns a new Queue.
func New() *Queue {
	return &Queue{
		buf: make([]interface{}, minQueueLen),
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() {
	newBuf := make([]interface{}, q.count*2)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an element on the end of the queue.
func (q *Queue) Add(elem interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *Queue) Peek() (interface{}, error) {
	if q.count <= 0 {
		return nil, errors.New("queue: Peek() called on empty queue")
	}
	return q.buf[q.head], nil
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic.
func (q *Queue) Get(i int) (interface{}, error) {
	if i < 0 || i >= q.count {
		return nil, errors.New("queue: Get() called with index out of range")
	}
	return q.buf[(q.head+i)%len(q.buf)], nil
}

// Gets and returns the first item from the queue.
func (q *Queue) Pop() (interface{}, error) {
	item, err := q.Peek()
	if err != nil {
		return nil, err
	}

	return item, q.Remove()
}

// Remove removes the element from the front of the queue. If you actually
// want the element, call Peek first. This call panics if the queue is empty.
func (q *Queue) Remove() error {
	if q.count <= 0 {
		return errors.New("queue: Remove() called on empty queue")
	}
	q.buf[q.head] = nil
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	if len(q.buf) > minQueueLen && q.count*4 == len(q.buf) {
		q.resize()
	}

	return nil
}
