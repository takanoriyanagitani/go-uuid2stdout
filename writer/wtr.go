package writer

import (
	"context"
	"errors"
	"io"
	"iter"
	"os"

	"github.com/google/uuid"
	us "github.com/takanoriyanagitani/go-uuid2stdout"
	. "github.com/takanoriyanagitani/go-uuid2stdout/util"
)

type Writer func(us.Uuid) IO[Void]

func (w Writer) WriteAll(ids iter.Seq2[us.Uuid, error]) IO[Void] {
	return func(ctx context.Context) (Void, error) {
		for uid, err := range ids {
			select {
			case <-ctx.Done():
				return Empty, ctx.Err()
			default:
			}

			if nil != err {
				return Empty, err
			}

			_, err = w(uid)(ctx)
			if nil != err {
				return Empty, err
			}
		}

		return Empty, nil
	}
}

type UuidToString func(us.Uuid) IO[string]

func Uuid2str(u us.Uuid) IO[string] {
	return func(_ context.Context) (string, error) {
		var raw [16]byte = u

		var i uuid.UUID = raw

		return i.String(), nil
	}
}

var UuidToStringDefault UuidToString = Uuid2str

// Human readable UUID text.
type UuidText []byte

type UuidToText func(us.Uuid) IO[UuidText]

func Uuid2txt(u us.Uuid) IO[UuidText] {
	return func(_ context.Context) (UuidText, error) {
		var raw [16]byte = u

		var id uuid.UUID = raw

		return id.MarshalText()
	}
}

var UuidToTextDefault UuidToText = Uuid2txt

type TextWriter func(UuidText) IO[Void]

func (t TextWriter) ToWriter(u2t UuidToText) Writer {
	return func(uid us.Uuid) IO[Void] {
		return Bind(
			u2t(uid),
			t,
		)
	}
}

func (t TextWriter) ToWriterDefault() Writer {
	return t.ToWriter(UuidToTextDefault)
}

// A string to be appended after writing the [UuidText].
type TextWriterSuffix string

func (s TextWriterSuffix) ToWriter(wtr io.Writer) TextWriter {
	var converted []byte = []byte(s)

	return func(txt UuidText) IO[Void] {
		var raw []byte = txt

		return Bind(
			func(_ context.Context) (int, error) {
				ir, er := wtr.Write(raw)
				ic, ec := wtr.Write(converted)

				return ir + ic, errors.Join(er, ec)
			},
			Lift(func(_ int) (Void, error) { return Empty, nil }),
		)
	}
}

const TextWriterSuffixDefault TextWriterSuffix = "\n"

var TextWriterStdoutDefault TextWriter = TextWriterSuffixDefault.
	ToWriter(os.Stdout)

var WriterTxtStdoutDefault Writer = TextWriterStdoutDefault.ToWriterDefault()
