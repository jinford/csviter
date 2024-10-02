package csviter

import (
	"encoding/csv"
	"errors"
	"io"
	"iter"
)

type ReaderRes struct {
	RecordPos    int
	Record       []string
	InputOffset  int64
	FieldPosList []FieldPos
}

type FieldPos struct {
	Line   int
	Column int
}

type reader struct {
	*csv.Reader
	withFiledPos bool
}

func NewReader(r io.Reader, opts ...ReaderOption) iter.Seq2[*ReaderRes, error] {
	csvr := &reader{
		Reader:       csv.NewReader(r),
		withFiledPos: false,
	}
	for _, opt := range opts {
		opt.apply(csvr)
	}

	return func(yield func(*ReaderRes, error) bool) {
		var recordPos int
		for {
			record, err := csvr.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				if !yield(nil, err) {
					break
				}
			}

			recordPos++

			res := &ReaderRes{
				RecordPos:    recordPos,
				Record:       record,
				InputOffset:  csvr.InputOffset(),
				FieldPosList: []FieldPos{},
			}

			if csvr.withFiledPos {
				fieldPosList := make([]FieldPos, len(record))
				for i := range record {
					line, columns := csvr.FieldPos(i)
					fieldPosList[i] = FieldPos{
						Line:   line,
						Column: columns,
					}
				}

				res.FieldPosList = fieldPosList
			}

			if !yield(res, nil) {
				break
			}
		}
	}
}

type ReaderOption interface {
	apply(*reader)
}

type readerOptionFunc func(*reader)

func (f readerOptionFunc) apply(r *reader) {
	f(r)
}

func Comma(comma rune) ReaderOption {
	return readerOptionFunc(func(r *reader) {
		r.Comma = comma
	})
}

func Comment(comment rune) ReaderOption {
	return readerOptionFunc(func(r *reader) {
		r.Comment = comment
	})
}

func FieldsPerRecord(fieldsPerRecord int) ReaderOption {
	return readerOptionFunc(func(r *reader) {
		r.FieldsPerRecord = fieldsPerRecord
	})
}

func LazyQuotes(lazyQuotes bool) ReaderOption {
	return readerOptionFunc(func(r *reader) {
		r.LazyQuotes = lazyQuotes
	})
}

func TrimLeadingSpace(trimLeadingSpace bool) ReaderOption {
	return readerOptionFunc(func(r *reader) {
		r.TrimLeadingSpace = trimLeadingSpace
	})
}

func ReuseRecord(reuseRecord bool) ReaderOption {
	return readerOptionFunc(func(r *reader) {
		r.ReuseRecord = reuseRecord
	})
}

func WithFieldPos(withFiledPos bool) ReaderOption {
	return readerOptionFunc(func(r *reader) {
		r.withFiledPos = withFiledPos
	})
}
