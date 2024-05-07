package slogerror

import (
	"log/slog"
	"os"
	"testing"
)

func TestFormatting(t *testing.T) {
	tests := []struct {
		attrs  []slog.Attr
		expect string
	}{
		{
			attrs:  nil,
			expect: " ",
		},
		{
			attrs:  []slog.Attr{slog.String("a", "a")},
			expect: `["a" = "a"] `,
		},
		{
			attrs: []slog.Attr{
				slog.String("a", "b"),
				slog.String("c", "d"),
			},
			expect: `["a" = "b"]["c" = "d"] `,
		},
		{
			attrs: []slog.Attr{
				slog.Group(
					"a",
					slog.String("b", "c"),
					slog.String("d", "e"),
				),
				slog.String("z", "y"),
			},
			expect: `["a.b" = "c"]["a.d" = "e"]["z" = "y"] `,
		},
		{
			attrs: []slog.Attr{
				slog.Group(
					"a",
					slog.Group(
						"x",
						slog.String("b", "c"),
						slog.String("d", "e"),
					),
				),
				slog.String("z", "y"),
			},
			expect: `["a.x.b" = "c"]["a.x.d" = "e"]["z" = "y"] `,
		},
	}
	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			l := slog.New(NewHandler(slog.NewTextHandler(os.Stderr, nil)))
			for _, attr := range test.attrs {
				l = l.With(attr)
			}
			err := Error(l, "")
			if err.Error() != test.expect {
				t.Fatalf("Expect '%s' but '%s'", test.expect, err.Error())
			}
		})
	}
}

func TestWithGroup(t *testing.T) {
	l := slog.New(NewHandler(slog.NewTextHandler(os.Stderr, nil)))
	l = l.WithGroup("G0")
	l = l.With(slog.String("A0", "V0"))
	err := Error(l, "0")
	const expect0 = `["G0.A0" = "V0"] 0`
	if err.Error() != expect0 {
		t.Fatalf("Expect '%s' but '%s'", expect0, err.Error())
	}
	l = l.WithGroup("G1")
	l = l.With(slog.String("A1", "V1"))
	err = Error(l, "1")
	const expect1 = `["G0.A0" = "V0"]["G1.A1" = "V1"] 1`
	if err.Error() != expect1 {
		t.Fatalf("Expect '%s' but '%s'", expect1, err.Error())
	}
}
