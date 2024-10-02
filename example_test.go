package csviter_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/jinford/csviter"
)

func ExampleNewReader() {
	in := `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`
	for res, err := range csviter.NewReader(strings.NewReader(in)) {
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.RecordPos, res.Record)
	}

	// Output:
	// 1 [first_name last_name username]
	// 2 [Rob Pike rob]
	// 3 [Ken Thompson ken]
	// 4 [Robert Griesemer gri]
}

func ExampleNewReader_options() {
	in := `first_name;last_name;username
"Rob";"Pike";rob
# lines beginning with a # character are ignored
Ken;Thompson;ken
"Robert";"Griesemer";"gri"
`

	for res, err := range csviter.NewReader(strings.NewReader(in),
		csviter.Comma(';'),
		csviter.Comment('#'),
	) {
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.RecordPos, res.Record)
	}

	// Output:
	// 1 [first_name last_name username]
	// 2 [Rob Pike rob]
	// 3 [Ken Thompson ken]
	// 4 [Robert Griesemer gri]
}
