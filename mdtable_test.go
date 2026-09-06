package mdtable

import "testing"

func TestFormat(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "messy spacing gets aligned",
			input: "| Name | Age |\n" +
				"|---|---|\n" +
				"| Bob |30|\n" +
				"|Alice   | 4 |\n",
			want: "| Name  | Age |\n" +
				"| ----- | --- |\n" +
				"| Bob   | 30  |\n" +
				"| Alice | 4   |\n",
		},
		{
			name: "ragged rows are padded to the widest row",
			input: "| A | B |\n" +
				"|---|---|\n" +
				"| 1 |\n" +
				"| 2 | 3 | 4 |\n",
			want: "| A   | B   |     |\n" +
				"| --- | --- | --- |\n" +
				"| 1   |     |     |\n" +
				"| 2   | 3   | 4   |\n",
		},
		{
			name: "alignment markers control padding side",
			input: "|Left|Center|Right|\n" +
				"|:---|:---:|---:|\n" +
				"|a|b|c|\n" +
				"|aaaa|bbbb|cccc|\n",
			want: "| Left | Center | Right |\n" +
				"| :--- | :----: | ----: |\n" +
				"| a    |   b    |     c |\n" +
				"| aaaa |  bbbb  |  cccc |\n",
		},
		{
			name: "rows without outer pipes are normalized",
			input: "Name | Score\n" +
				"--- | ---\n" +
				"Ann | 9\n",
			want: "| Name | Score |\n" +
				"| ---- | ----- |\n" +
				"| Ann  | 9     |\n",
		},
		{
			name: "escaped pipe inside a cell survives round trip",
			input: "| Expr | Result |\n" +
				"|---|---|\n" +
				"| a \\| b | true |\n",
			want: "| Expr   | Result |\n" +
				"| ------ | ------ |\n" +
				"| a \\| b | true   |\n",
		},
		{
			name: "column width is measured in runes, not bytes",
			input: "| Word | Note |\n" +
				"|---|---|\n" +
				"| café | ok |\n" +
				"| x | ok |\n",
			// "café" is 4 runes but 5 bytes (é is two UTF-8 bytes);
			// byte-based padding would misalign the column by one space.
			want: "| Word | Note |\n" +
				"| ---- | ---- |\n" +
				"| café | ok   |\n" +
				"| x    | ok   |\n",
		},
		{
			name: "wide CJK runes are measured as two columns wide",
			input: "| Word | Note |\n" +
				"|---|---|\n" +
				"| 日本語 | ok |\n" +
				"| x | ok |\n",
			// "日本語" is 3 runes but occupies 6 display columns, so the
			// column needs to be as wide as "Note" plus two more spaces.
			want: "| Word   | Note |\n" +
				"| ------ | ---- |\n" +
				"| 日本語 | ok   |\n" +
				"| x      | ok   |\n",
		},
		{
			name: "surrounding prose and blank lines are left alone",
			input: "# Heading\n" +
				"\n" +
				"Some text before.\n" +
				"\n" +
				"|A|B|\n" +
				"|---|---|\n" +
				"|1|2|\n" +
				"\n" +
				"Some text after.\n",
			want: "# Heading\n" +
				"\n" +
				"Some text before.\n" +
				"\n" +
				"| A   | B   |\n" +
				"| --- | --- |\n" +
				"| 1   | 2   |\n" +
				"\n" +
				"Some text after.\n",
		},
		{
			name: "empty cells stay empty but padded to column width",
			input: "| A | B |\n" +
				"|---|---|\n" +
				"|   | x |\n",
			want: "| A   | B   |\n" +
				"| --- | --- |\n" +
				"|     | x   |\n",
		},
		{
			name: "short content still gets a minimum three dash separator",
			input: "|A|\n" +
				"|-|\n" +
				"|1|\n",
			want: "| A   |\n" +
				"| --- |\n" +
				"| 1   |\n",
		},
		{
			name: "two independent tables in one document are each formatted",
			input: "|A|B|\n|---|---|\n|1|2|\n\n|X|Y|\n|---|---|\n|9|8|\n",
			want: "| A   | B   |\n" +
				"| --- | --- |\n" +
				"| 1   | 2   |\n" +
				"\n" +
				"| X   | Y   |\n" +
				"| --- | --- |\n" +
				"| 9   | 8   |\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Format(tc.input)
			if got != tc.want {
				t.Errorf("Format() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tc.want)
			}
		})
	}
}

func TestFormatIsIdempotent(t *testing.T) {
	input := "| Name | Age |\n|---|---|\n| Bob |30|\n|Alice   | 4 |\n"
	once := Format(input)
	twice := Format(once)
	if once != twice {
		t.Errorf("Format is not idempotent\n--- once ---\n%s\n--- twice ---\n%s", once, twice)
	}
}
