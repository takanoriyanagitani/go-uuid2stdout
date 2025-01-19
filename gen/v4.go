package gen

import (
	"context"

	"github.com/google/uuid"
	us "github.com/takanoriyanagitani/go-uuid2stdout"
)

type Generate4 Generate

func (r RandomReader) ToGen4default() Generate4 {
	return func(_ context.Context) (us.Uuid, error) {
		u, e := uuid.NewRandom()

		return us.Uuid(u), e
	}
}

func (r RandomReader) ReaderToGen4() Generate4 {
	return func(_ context.Context) (us.Uuid, error) {
		u, e := uuid.NewRandomFromReader(r.Reader)

		return us.Uuid(u), e
	}
}

func (r RandomReader) ToGen4() Generate4 {
	switch r.Reader {
	case nil:
		return r.ToGen4default()
	default:
		return r.ReaderToGen4()
	}
}

var Generate4default Generate4 = RandomReaderEmpty.ToGen4default()
