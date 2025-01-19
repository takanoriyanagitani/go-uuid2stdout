package gen

import (
	"context"
	"errors"
	"io"
	"iter"

	us "github.com/takanoriyanagitani/go-uuid2stdout"
	. "github.com/takanoriyanagitani/go-uuid2stdout/util"
)

type Generate IO[us.Uuid]

type StopType int

const (
	StopTypeUnspecified StopType = iota
	StopTypeStop        StopType = iota
	StopTypeCont        StopType = iota
)

type Stopper func(error) (StopType, error)

func (g Generate) ToIter(stp Stopper) IO[iter.Seq2[us.Uuid, error]] {
	return func(ctx context.Context) (iter.Seq2[us.Uuid, error], error) {
		return func(yield func(us.Uuid, error) bool) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				uid, err := g(ctx)
				typ, err := stp(err)

				switch typ {
				case StopTypeStop:
					return
				case StopTypeCont:
					yield(uid, err)
				case StopTypeUnspecified:
					yield(uid, err)

					return
				}
			}
		}, nil
	}
}

func (g Generate) ToIterDefault() IO[iter.Seq2[us.Uuid, error]] {
	return g.ToIter(StopperDefault)
}

func StopperEof(err error) (StopType, error) {
	switch {
	case nil == err:
		return StopTypeCont, nil
	case errors.Is(err, io.EOF):
		return StopTypeStop, nil
	default:
		return StopTypeUnspecified, err
	}
}

var StopperDefault Stopper = StopperEof
