package goautoneg

import (
	"reflect"
	"testing"
)

var empty = make(map[string]string)

func accept(mainType, subType string, q float64, params map[string]string) Accept {
	return Accept{
		Type:    mainType,
		SubType: subType,
		Q:       q,
		Params:  params,
	}
}

func TestParseAccept(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   []Accept
	}{{
		name:   "empty",
		header: "",
		want:   []Accept{},
	}, {
		name:   "broken",
		header: "text/html/and/other, text/plain",
		want:   []Accept{accept("text", "plain", 1.0, empty)},
	}, {
		name:   "single asterisk",
		header: "*",
		want:   []Accept{accept("*", "*", 1.0, empty)},
	}, {
		name:   "broken q",
		header: "text/plain;q=foo",
		want:   []Accept{accept("text", "plain", 0.0, empty)},
	}, {
		name:   "broken parameter",
		header: "text/plain;q=0.3;doggo",
		want:   []Accept{},
	}, {
		name:   "rfc2616 simple",
		header: "audio/*; q=0.2, audio/basic",
		want: []Accept{
			accept("audio", "basic", 1.0, empty),
			accept("audio", "*", 0.2, empty),
		},
	}, {
		name:   "rfc2616 elaborate",
		header: "text/plain; q=0.5, text/html, text/x-dvi; q=0.8, text/x-c",
		want: []Accept{
			accept("text", "html", 1.0, empty),
			accept("text", "x-c", 1.0, empty),
			accept("text", "x-dvi", 0.8, empty),
			accept("text", "plain", 0.5, empty),
		},
	}, {
		name:   "rfc2616 complex",
		header: "text/*;q=0.3, text/html;q=0.7, text/html;level=1, text/html;level=2;q=0.4, */*;q=0.5",
		want: []Accept{
			accept("text", "html", 1.0, map[string]string{"level": "1"}),
			accept("text", "html", 0.7, empty),
			accept("*", "*", 0.5, empty),
			accept("text", "html", 0.4, map[string]string{"level": "2"}),
			accept("text", "*", 0.3, empty),
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ParseAccept(test.header)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Unexpected value\ngot:  %+v\nwant: %+v", got, test.want)
			}
		})
	}
}
