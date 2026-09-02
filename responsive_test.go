package main

import "testing"

func TestReferenceColumnsWidth(t *testing.T) {
	columns := []refColumn{
		refTeeth,
		refFull,
		refHalf,
		refOffset,
		refCompValue,
		refCompressors,
	}

	tests := []struct {
		spacing int
		want    int
	}{
		{4, 73},
		{3, 68},
		{2, 63},
		{1, 58},
		{0, 53},
	}

	for _, tt := range tests {
		got := referenceColumnsWidth(columns, tt.spacing)

		if got != tt.want {
			t.Errorf(
				"spacing %d: got width %d, want %d",
				tt.spacing,
				got,
				tt.want,
			)
		}
	}
}

func TestReferenceSpacingForWidth(t *testing.T) {
	columns := []refColumn{
		refTeeth,
		refFull,
		refHalf,
		refOffset,
		refCompValue,
		refCompressors,
	}

	tests := []struct {
		panelWidth int
		want       int
	}{
		{100, 4},
		{77, 4},
		{76, 3},
		{72, 3},
		{71, 2},
		{66, 1},
		{65, 1},
		{60, 0},
		{59, 0},
		{57, 0},
	}

	for _, tt := range tests {
		spacing := 4

		for len(columns) > 1 {
			requiredWidth :=
				referenceColumnsWidth(
					columns,
					spacing,
				)

			if requiredWidth <= tt.panelWidth-4 ||
				spacing == 0 {
				break
			}

			spacing--
		}

		if spacing != tt.want {
			t.Errorf(
				"panel width %d: got spacing %d, want %d",
				tt.panelWidth,
				spacing,
				tt.want,
			)
		}
	}
}
