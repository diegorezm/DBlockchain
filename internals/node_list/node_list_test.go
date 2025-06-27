package nodelist

import "testing"

func TestNodeList_Add(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "Test if size is correct",
			data: "testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n NodeList[*string]
			n.Add(&tt.data)
			if n.Size() != 1 {
				t.Errorf("Size() = %v, want %v\n", n.Size(), 1)
			}

			n.Add(&tt.data)
			if n.Size() != 2 {
				t.Errorf("Size() = %v, want %v\n", n.Size(), 1)
			}

		})
	}
}

func TestNodeList_Peek(t *testing.T) {
	n1 := NewNodeList[*string]()
	data := "testing"
	n1.Add(&data)
	if n1.Peek() == nil {
		t.Errorf("Peek() = %v, want %v\n", n1.Peek(), &data)
	}

	n2 := NewNodeList[*string]()

	if n2.Peek() != nil {
		t.Error("N2 Peek() should be null.")
	}
}

func TestNodeList_Pop(t *testing.T) {
	n := NewNodeList[*string]()
	data := "testing"

	n.Add(&data)
	n.Add(&data)
	_ = n.Pop()
	if n.Size() != 1 {
		t.Errorf("Size() = %v, want %v\n", n.Size(), 1)
	}

	n.Pop()
	if n.Size() != 0 {
		t.Errorf("Size() = %v, want %v\n", n.Size(), 0)
	}

	if n.Pop() != nil {
		t.Errorf("Pop() = %v, want %v\n", n.Size(), 0)
	}

}
