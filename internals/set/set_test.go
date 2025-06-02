package set

import (
	"testing"
)

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  int
	}{
		{
			name:  "Check if correct number of items is added",
			want:  2,
			items: []int{2, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet[int]()
			for _, n := range tt.items {
				s.Add(n)
			}
			if s.Size() > tt.want {
				t.Error("The incorrect number of items as added.\n")
			}
		})
	}
}

func TestSet_Delete(t *testing.T) {
	tests := []struct {
		name     string
		items    []int
		toDelete int
		want     int
	}{
		{
			name:     "Test if the item is deleted.",
			items:    []int{1, 2, 3, 4},
			toDelete: 2,
			want:     3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet[int]()
			for _, n := range tt.items {
				s.Add(n)
			}
			s.Delete(tt.toDelete)
			if s.Size() > tt.want {
				t.Errorf("The item %d was not deleted.\n", tt.toDelete)
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  int
	}{
		{
			name:  "Test if the set contains a certain item.",
			items: []int{1, 2, 3, 4},
			want:  3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet[int]()
			for _, n := range tt.items {
				s.Add(n)
			}
			if !s.Contains(tt.want) {
				t.Errorf("The item %d is not present in the set.\n", tt.want)
			}
		})
	}
}

func TestSet_ToSlice(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  int
	}{
		{
			name:  "Test if the list length is correct.",
			items: []int{1, 2, 3, 4},
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet[int]()
			for _, n := range tt.items {
				s.Add(n)
			}
			slice := s.ToSlice()
			if len(slice) != tt.want {
				t.Errorf("The list does not comply with the expected length of %d.\n", tt.want)
			}
		})
	}
}

func TestSet_Size(t *testing.T) {
	tests := []struct {
		name  string
		items []int
		want  int
	}{
		{
			name:  "Test if the length is correct.",
			items: []int{1, 2, 3, 4},
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet[int]()
			for _, n := range tt.items {
				s.Add(n)
			}
			if s.Size() != tt.want {
				t.Errorf("The set does not comply with the expected length of %d.\n", tt.want)
			}
		})
	}
}

func TestSet_AddStruct(t *testing.T) {
	type items struct {
		name string
	}

	tests := []struct {
		name  string
		items []items
		want  int
	}{
		{
			name:  "Test if the length is correct.",
			items: []items{{"a"}, {"b"}, {"c"}, {"d"}, {"d"}},
			want:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet[items]()
			for _, i := range tt.items {
				s.Add(i)
			}
			if s.Size() != tt.want {
				t.Errorf("The set does not comply with the expected length of %d. got %d\n", tt.want, s.Size())
			}
		})
	}
}
