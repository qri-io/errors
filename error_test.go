package errors

import (
	"fmt"
	"testing"
)

func ExampleNewWithMessage() {
	// an example function that returns an error
	errFunc := func(id int) error {
		return fmt.Errorf("not found")
	}

	id := 1
	err := errFunc(id)
	if err != nil {
		err = WrapFriendly(CodeNotFound, err, "", "couldn't find data with the id", id)
	}

	// when it's time to deal with the error, we to check to see if the error is a
	// qri error to get access to special error methods
	// the else condition won't happen in this example, but we have to handle
	// normal errors in the real world
	if qriErr, ok := err.(*Error); ok {
		fmt.Println(qriErr.Friendly())
	} else {
		fmt.Println(err.Error())
	}

	// Output:
	// missing: couldn't find data with the id 1.
}

func TestNewWithMessage(t *testing.T) {
	e := NewFriendly(CodeGeneric, "machine error", "message")

	expectMsg := "error: message"
	if e.Friendly() != expectMsg {
		t.Errorf("error in Error struct function `Message()`: expected: %s, got: %s", expectMsg, e.Friendly())
	}

	expectErr := "error: machine error"
	if e.Error() != expectErr {
		t.Errorf("error in Error struct function `Error()`: expected: %s, got: %s", expectErr, e.Error())
	}
}

func TestError(t *testing.T) {
	val := "a"
	a := New(CodeGeneric, "so bad")
	b := Wrap(CodeInvalidArgs, a, val)
	c := WrapFriendlyFix(CodeForbidden, b, "you're not allowed access to the value", "please try a different value", val)

	switch Cause(c).(type) {
	case *Error:
		t.Errorf("wrong error type returned")
	}

	if a.Friendly() != "" {
		t.Errorf("expected no friendly message to be empty")
	}

	if c.Code() != CodeForbidden {
		t.Errorf("wrong code. expected: %d, got: %d", CodeForbidden, c.Code())
	}

	fix := "are you sure you're logged in?"
	d := NewFriendlyFix(CodeForbidden, "forbidden", "you don't have access to the following things:", fix, "apples", "oranges")
	if d.Fix() != fix {
		t.Errorf("fix mismatch. expected: %s, got: %s", fix, c.Fix())
	}

	expFriendly := "auth: you don't have access to the following things: apples, oranges. are you sure you're logged in?"
	if d.Friendly() != expFriendly {
		t.Errorf("friendly mismatch. expected: %s, got: %s", expFriendly, d.Friendly())
	}
}

func TestRegisterCode(t *testing.T) {
	if err := RegisterCode(CodeForbidden, 200, "forbidden"); err == nil {
		t.Error("expected registring an already-existing Code to error")
	}
	CodeNoDatabase := Code(100)
	RegisterCode(CodeNoDatabase, 504, "database")

	expStatus := 504
	status := CodeHTTPStatus(CodeNoDatabase)
	if expStatus != status {
		t.Errorf("status mismatch. expected: %d, got: %d", 504, status)
	}
}

func TestCodeVals(t *testing.T) {
	expect := "error"
	got := CodeString(Code(-1))
	if got != expect {
		t.Errorf("code string mismatch. expected %s, got: %s", expect, got)
	}

	expectHTTP := 500
	gotHTTP := CodeHTTPStatus(Code(-1))
	if gotHTTP != expectHTTP {
		t.Errorf("http code mismatch. expected: %d, got: %d", expectHTTP, gotHTTP)
	}
}
