package mindis

import (
	"testing"
)

func Test_MarshalItem_Redis(t *testing.T) {
	//TODO
}

func Test_MarshalItem_JSON(t *testing.T) {
	//TODO
}

func Test_MarshalItem_Bytes(t *testing.T) {
	testMarshalItem(t, []byte{1, 2, 3}, "\001\002\003")
}

func Test_MarshalItem_String(t *testing.T) {
	testMarshalItem(t, "hello", "hello")
}

func Test_MarshalItem_Int(t *testing.T) {
	testMarshalItem(t, 1, "1")
}

func Test_MarshalItem_Int32(t *testing.T) {
	testMarshalItem(t, int32(1), "1")
}

func Test_MarshalItem_Int64(t *testing.T) {
	testMarshalItem(t, int64(1), "1")
}

func Test_MarshalItem_Uint(t *testing.T) {
	testMarshalItem(t, uint(1), "1")
}

func Test_MarshalItem_Uint32(t *testing.T) {
	testMarshalItem(t, uint32(1), "1")
}

func Test_MarshalItem_Uint64(t *testing.T) {
	testMarshalItem(t, uint(1), "1")
}

func Test_MarshalItem_Float32(t *testing.T) {
	testMarshalItem(t, float32(1.1), "1.1")
}

func Test_MarshalItem_Float64(t *testing.T) {
	testMarshalItem(t, float64(1.1), "1.1")
}

func Test_MarshalItem_Bool(t *testing.T) {
	testMarshalItem(t, true, "true")
}

func testMarshalItem(t *testing.T, value interface{}, data string) {
	b, e := marshalItem(value)
	if e != nil {
		t.Fatal("unexpected error:", e.Error())
	}

	if string(b) != data {
		t.Fatalf("expected encoding: %#v got: %#v", data, string(b))
	}
}
