package gen

import (
	"context"
	"io"

	"github.com/google/uuid"
	us "github.com/takanoriyanagitani/go-uuid2stdout"
)

type Generate7 Generate

func (g Generate7) AsGenerate() Generate { return Generate(g) }

type RandomReader struct {
	io.Reader
}

var RandomReaderEmpty RandomReader

func (r RandomReader) ToGen7default() Generate7 {
	return func(_ context.Context) (us.Uuid, error) {
		u, e := uuid.NewV7()

		return us.Uuid(u), e
	}
}

func (r RandomReader) ReaderToGen7() Generate7 {
	return func(_ context.Context) (us.Uuid, error) {
		u, e := uuid.NewV7FromReader(r.Reader)

		return us.Uuid(u), e
	}
}

func (r RandomReader) ToGen7() Generate7 {
	switch r.Reader {
	case nil:
		return r.ToGen7default()
	default:
		return r.ReaderToGen7()
	}
}

var Generate7default Generate7 = RandomReaderEmpty.ToGen7default()
