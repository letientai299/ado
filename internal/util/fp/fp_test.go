package fp

import (
	"reflect"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("int transformations", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []int
			fn       func(int) int
			expected []int
		}{
			{
				name:     "double values",
				input:    []int{1, 2, 3, 4, 5},
				fn:       func(n int) int { return n * 2 },
				expected: []int{2, 4, 6, 8, 10},
			},
			{
				name:     "square values",
				input:    []int{1, 2, 3, 4, 5},
				fn:       func(n int) int { return n * n },
				expected: []int{1, 4, 9, 16, 25},
			},
			{
				name:     "add constant",
				input:    []int{10, 20, 30},
				fn:       func(n int) int { return n + 5 },
				expected: []int{15, 25, 35},
			},
			{
				name:     "empty slice",
				input:    []int{},
				fn:       func(n int) int { return n * 2 },
				expected: []int{},
			},
			{
				name:     "single element",
				input:    []int{42},
				fn:       func(n int) int { return n - 10 },
				expected: []int{32},
			},
			{
				name:     "negative numbers",
				input:    []int{-1, -2, -3},
				fn:       func(n int) int { return -n },
				expected: []int{1, 2, 3},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Map(tt.input, tt.fn)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Map() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("type conversions", func(t *testing.T) {
		tests := []struct {
			name string
			test func(t *testing.T)
		}{
			{
				name: "int to string",
				test: func(t *testing.T) {
					input := []int{1, 2, 3}
					fn := func(n int) string { return strings.Repeat("x", n) }
					expected := []string{"x", "xx", "xxx"}

					result := Map(input, fn)
					if !reflect.DeepEqual(result, expected) {
						t.Errorf("Map() = %v, want %v", result, expected)
					}
				},
			},
			{
				name: "string to int",
				test: func(t *testing.T) {
					input := []string{"hello", "hi", "golang"}
					fn := func(s string) int { return len(s) }
					expected := []int{5, 2, 6}

					result := Map(input, fn)
					if !reflect.DeepEqual(result, expected) {
						t.Errorf("Map() = %v, want %v", result, expected)
					}
				},
			},
			{
				name: "string to bool",
				test: func(t *testing.T) {
					input := []string{"", "hello", "", "world"}
					fn := func(s string) bool { return s != "" }
					expected := []bool{false, true, false, true}

					result := Map(input, fn)
					if !reflect.DeepEqual(result, expected) {
						t.Errorf("Map() = %v, want %v", result, expected)
					}
				},
			},
			{
				name: "struct to field",
				test: func(t *testing.T) {
					type person struct {
						name string
						age  int
					}
					input := []person{
						{"Alice", 30},
						{"Bob", 25},
						{"Charlie", 35},
					}
					fn := func(p person) string { return p.name }
					expected := []string{"Alice", "Bob", "Charlie"}

					result := Map(input, fn)
					if !reflect.DeepEqual(result, expected) {
						t.Errorf("Map() = %v, want %v", result, expected)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, tt.test)
		}
	})

	t.Run("string transformations", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []string
			fn       func(string) string
			expected []string
		}{
			{
				name:     "uppercase",
				input:    []string{"hello", "world", "go"},
				fn:       strings.ToUpper,
				expected: []string{"HELLO", "WORLD", "GO"},
			},
			{
				name:     "lowercase",
				input:    []string{"HELLO", "World", "GO"},
				fn:       strings.ToLower,
				expected: []string{"hello", "world", "go"},
			},
			{
				name:     "add prefix",
				input:    []string{"one", "two", "three"},
				fn:       func(s string) string { return "prefix-" + s },
				expected: []string{"prefix-one", "prefix-two", "prefix-three"},
			},
			{
				name:     "empty strings",
				input:    []string{"", "", ""},
				fn:       func(s string) string { return s + "x" },
				expected: []string{"x", "x", "x"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Map(tt.input, tt.fn)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Map() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("preserves capacity", func(t *testing.T) {
		input := make([]int, 3, 10)
		input[0], input[1], input[2] = 1, 2, 3

		result := Map(input, func(n int) int { return n * 2 })

		if cap(result) != len(input) {
			t.Errorf("Map() capacity = %d, want %d", cap(result), len(input))
		}
	})

	t.Run("nil input returns empty slice", func(t *testing.T) {
		var input []int
		result := Map(input, func(n int) int { return n * 2 })

		if result == nil {
			t.Error("Map(nil) returned nil, want empty slice")
		}
		if len(result) != 0 {
			t.Errorf("Map(nil) length = %d, want 0", len(result))
		}
	})
}

func TestNot(t *testing.T) {
	t.Run("integer predicates", func(t *testing.T) {
		isEven := func(n int) bool { return n%2 == 0 }
		isPositive := func(n int) bool { return n > 0 }

		tests := []struct {
			name      string
			predicate func(int) bool
			input     int
			expected  bool
		}{
			{"not even (odd number)", isEven, 3, true},
			{"not even (even number)", isEven, 4, false},
			{"not positive (negative)", isPositive, -5, true},
			{"not positive (positive)", isPositive, 5, false},
			{"not positive (zero)", isPositive, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				notPredicate := Not(tt.predicate)
				result := notPredicate(tt.input)
				if result != tt.expected {
					t.Errorf("Not(predicate)(%d) = %v, want %v", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("string predicates", func(t *testing.T) {
		tests := []struct {
			name      string
			predicate func(string) bool
			input     string
			expected  bool
		}{
			{
				name:      "not has prefix",
				predicate: func(s string) bool { return strings.HasPrefix(s, "test") },
				input:     "test123",
				expected:  false,
			},
			{
				name:      "not has prefix (no match)",
				predicate: func(s string) bool { return strings.HasPrefix(s, "test") },
				input:     "hello",
				expected:  true,
			},
			{
				name:      "not empty",
				predicate: func(s string) bool { return s == "" },
				input:     "hello",
				expected:  true,
			},
			{
				name:      "not empty (empty string)",
				predicate: func(s string) bool { return s == "" },
				input:     "",
				expected:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				notPredicate := Not(tt.predicate)
				result := notPredicate(tt.input)
				if result != tt.expected {
					t.Errorf("Not(predicate)(%q) = %v, want %v", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("double negation", func(t *testing.T) {
		isPositive := func(n int) bool { return n > 0 }
		doubleNegation := Not(Not(isPositive))

		tests := []struct {
			input    int
			expected bool
		}{
			{-5, false},
			{0, false},
			{5, true},
			{10, true},
		}

		for _, tt := range tests {
			result := doubleNegation(tt.input)
			original := isPositive(tt.input)
			if result != original {
				t.Errorf("Not(Not(isPositive))(%d) = %v, want %v", tt.input, result, original)
			}
		}
	})

	t.Run("constant predicates", func(t *testing.T) {
		tests := []struct {
			name      string
			predicate func(any) bool
			expected  bool
		}{
			{"not always true", func(any) bool { return true }, false},
			{"not always false", func(any) bool { return false }, true},
		}

		testValues := []any{1, "test", true, nil, []int{1, 2}}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				notPredicate := Not(tt.predicate)
				for _, val := range testValues {
					result := notPredicate(val)
					if result != tt.expected {
						t.Errorf("Not(predicate)(%v) = %v, want %v", val, result, tt.expected)
					}
				}
			})
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("integer filtering", func(t *testing.T) {
		tests := []struct {
			name      string
			input     []int
			predicate func(int) bool
			expected  []int
		}{
			{
				name:      "filter even numbers",
				input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				predicate: func(n int) bool { return n%2 == 0 },
				expected:  []int{2, 4, 6, 8, 10},
			},
			{
				name:      "filter positive numbers",
				input:     []int{-3, -1, 0, 1, 2, 3},
				predicate: func(n int) bool { return n > 0 },
				expected:  []int{1, 2, 3},
			},
			{
				name:      "filter greater than 5",
				input:     []int{1, 3, 5, 7, 9, 11},
				predicate: func(n int) bool { return n > 5 },
				expected:  []int{7, 9, 11},
			},
			{
				name:      "no matches",
				input:     []int{1, 3, 5, 7, 9},
				predicate: func(n int) bool { return n%2 == 0 },
				expected:  []int{},
			},
			{
				name:      "all match",
				input:     []int{2, 4, 6, 8, 10},
				predicate: func(n int) bool { return n%2 == 0 },
				expected:  []int{2, 4, 6, 8, 10},
			},
			{
				name:      "empty input",
				input:     []int{},
				predicate: func(n int) bool { return n > 0 },
				expected:  []int{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Filter(tt.input, tt.predicate)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Filter() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("string filtering", func(t *testing.T) {
		tests := []struct {
			name      string
			input     []string
			predicate func(string) bool
			expected  []string
		}{
			{
				name:      "filter by prefix",
				input:     []string{"apple", "banana", "apricot", "cherry", "avocado"},
				predicate: func(s string) bool { return strings.HasPrefix(s, "a") },
				expected:  []string{"apple", "apricot", "avocado"},
			},
			{
				name:      "filter by length",
				input:     []string{"go", "java", "rust", "c", "python"},
				predicate: func(s string) bool { return len(s) <= 3 },
				expected:  []string{"go", "c"},
			},
			{
				name:      "filter non-empty",
				input:     []string{"", "hello", "", "world", ""},
				predicate: func(s string) bool { return s != "" },
				expected:  []string{"hello", "world"},
			},
			{
				name:      "filter contains substring",
				input:     []string{"hello", "world", "hello world", "goodbye"},
				predicate: func(s string) bool { return strings.Contains(s, "o") },
				expected:  []string{"hello", "world", "hello world", "goodbye"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Filter(tt.input, tt.predicate)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Filter() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("struct filtering", func(t *testing.T) {
		type person struct {
			name string
			age  int
		}

		people := []person{
			{"Alice", 30},
			{"Bob", 17},
			{"Charlie", 25},
			{"Diana", 16},
			{"Eve", 21},
		}

		tests := []struct {
			name      string
			predicate func(person) bool
			expected  []person
		}{
			{
				name:      "filter adults",
				predicate: func(p person) bool { return p.age >= 18 },
				expected: []person{
					{"Alice", 30},
					{"Charlie", 25},
					{"Eve", 21},
				},
			},
			{
				name:      "filter by name length",
				predicate: func(p person) bool { return len(p.name) <= 4 },
				expected: []person{
					{"Bob", 17},
					{"Eve", 21},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Filter(people, tt.predicate)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Filter() = %v, want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("preserves order", func(t *testing.T) {
		input := []int{5, 2, 8, 1, 9, 4, 7, 3, 6}
		predicate := func(n int) bool { return n > 4 }
		expected := []int{5, 8, 9, 7, 6}

		result := Filter(input, predicate)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Filter() order = %v, want %v", result, expected)
		}
	})

	t.Run("integration with Not", func(t *testing.T) {
		tests := []struct {
			name      string
			input     []int
			predicate func(int) bool
			expected  []int
		}{
			{
				name:      "filter odd numbers using Not",
				input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				predicate: Not(func(n int) bool { return n%2 == 0 }),
				expected:  []int{1, 3, 5, 7, 9},
			},
			{
				name:      "filter non-positive using Not",
				input:     []int{-3, -1, 0, 1, 2, 3},
				predicate: Not(func(n int) bool { return n > 0 }),
				expected:  []int{-3, -1, 0},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Filter(tt.input, tt.predicate)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Filter() with Not = %v, want %v", result, tt.expected)
				}
			})
		}
	})
}
