package d1

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/tinyredglasses/workers/jsutil"
	"syscall/js"
	"time"
)

type stmt struct {
	stmtObj js.Value
}

var (
	_ driver.Stmt             = (*stmt)(nil)
	_ driver.StmtExecContext  = (*stmt)(nil)
	_ driver.StmtQueryContext = (*stmt)(nil)
)

func (s *stmt) Close() error {
	// do nothing
	return nil
}

// NumInput is not supported and always returns -1.
func (s *stmt) NumInput() int {
	return -1
}

func (s *stmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, errors.New("d1: Exec is deprecated and not implemented")
}

// ExecContext executes prepared statement.
// Given []drier.NamedValue's `Name` field will be ignored because Cloudflare D1 client doesn't support it.
func (s *stmt) ExecContext(_ context.Context, args []driver.NamedValue) (driver.Result, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from", r)
		}
	}()

	argValues := make([]any, len(args))
	for i, arg := range args {
		argValues[i] = arg.Value
	}

	fmt.Println(s.stmtObj.Get("bind").String())
	fmt.Println(s.stmtObj.Get("run").String())

	resultPromise := s.stmtObj.Call("bind", argValues...).Call("run")
	resultObj, err := jsutil.AwaitPromise(resultPromise)
	if err != nil {
		return nil, err
	}
	return &result{
		resultObj: resultObj,
	}, nil
}

func (s *stmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, errors.New("d1: Query is deprecated and not implemented")
}

func (s *stmt) QueryContext(_ context.Context, args []driver.NamedValue) (driver.Rows, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from", r)
		}
	}()

	fmt.Println("stmt1")
	argValues := make([]any, len(args))
	fmt.Println("stmt2")
	fmt.Println(s.stmtObj)

	for i, arg := range args {
		argValues[i] = arg.Value
		fmt.Println(arg.Value)
	}

	fmt.Println(argValues)
	//fmt.Println(argValues)

	//v1 := s.stmtObj.Get("all")
	fmt.Println(s.stmtObj.Get("bind").String())
	fmt.Println(s.stmtObj.Get("all").String())
	resultPromise := s.stmtObj.Call("bind").Call("all")
	fmt.Println("stmt3")

	fmt.Println(resultPromise)

	//rowsObj, err := jsutil.AwaitPromise(resultPromise)
	time.Sleep(time.Second)
	fmt.Println("stmt4")

	//if err != nil {
	//	return nil, err
	//}
	fmt.Println("stmt5")
	//fmt.Println(rowsObj)

	//if !rowsObj.Get("success").Bool() {
	//	return nil, errors.New("d1: failed to query")
	//}
	fmt.Println("stmt6")

	return &rows{
		//rowsObj: rowsObj.Get("results"),
		rowsObj: js.Value{},
	}, nil
}
