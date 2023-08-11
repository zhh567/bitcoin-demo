package main

import (
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	preHash := [32]byte{0}
	tests := []Block{
		*NewBlock("xxx", preHash[:]),
	}

	for _, item := range tests {
		t.Log(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n", item)
		b, err := item.Serialize()
		if err != nil {
			t.Fatal("serialize fail: ", err)
		}

		b2, err2 := Deserialize(b)
		if err2 != nil {
			t.Fatal("deserialize fail: ", err2)
		}
		t.Log(b2, "\n<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")

		if !reflect.DeepEqual(item, *b2) {
			t.Error("Deserialized data don't equal wanted")
		}
	}
}
