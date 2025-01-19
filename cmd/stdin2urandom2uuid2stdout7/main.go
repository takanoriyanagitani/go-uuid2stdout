package main

import (
	"context"
	"iter"
	"log"
	"os"

	us "github.com/takanoriyanagitani/go-uuid2stdout"
	ug "github.com/takanoriyanagitani/go-uuid2stdout/gen"
	. "github.com/takanoriyanagitani/go-uuid2stdout/util"
	uw "github.com/takanoriyanagitani/go-uuid2stdout/writer"
)

var gen7 ug.Generate = ug.
	RandomReader{Reader: os.Stdin}.
	ReaderToGen7().
	AsGenerate()

var uuids IO[iter.Seq2[us.Uuid, error]] = gen7.ToIterDefault()

var wtr uw.Writer = uw.WriterTxtStdoutDefault

var uuids2stdout IO[Void] = Bind(
	uuids,
	wtr.WriteAll,
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return uuids2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
