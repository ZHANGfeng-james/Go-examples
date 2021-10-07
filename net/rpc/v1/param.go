package v1

import (
	"errors"
	"time"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Que, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quotient *Quotient) error {
	if args.B == 0 {
		// string --> errors.New
		return errors.New("divided by zero")
	}
	quotient.Que = args.A / args.B
	quotient.Rem = args.A % args.B
	time.Sleep(2 * time.Second)
	return nil
}
