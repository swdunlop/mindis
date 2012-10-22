package mindis

import (
	"encoding/json"
	"testing"
)

func Test_UnmarshalItem_Json(t *testing.T) {
	x := MockJson(true)
	var v MockJson
	testUnmarshalItem(t, &v, "0n")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}

	x = MockJson(false)
	testUnmarshalItem(t, &v, "0ff")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalItem_String(t *testing.T) {
	x := "hello"
	var v string
	testUnmarshalItem(t, &v, "hello")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalItem_Bool(t *testing.T) {
	x := false
	var v bool
	testUnmarshalItem(t, &v, "false")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalItem_Int(t *testing.T) {
	x := 123
	var v int
	testUnmarshalItem(t, &v, "123")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalItem_Int32(t *testing.T) {
	x := 123
	var v int32
	testUnmarshalItem(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalItem_Int64(t *testing.T) {
	x := 123
	var v int64
	testUnmarshalItem(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}
func Test_UnmarshalItem_Uint(t *testing.T) {
	x := 123
	var v uint
	testUnmarshalItem(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalItem_Uint32(t *testing.T) {
	x := 123
	var v uint32
	testUnmarshalItem(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalItem_Uint64(t *testing.T) {
	x := 123
	var v uint64
	testUnmarshalItem(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func testUnmarshalItem(t *testing.T, v interface{}, data string) {
	err := unmarshalItem([]byte(data), v)
	if err != nil {
		t.Fatal("unexpected error:", err.Error())
	}
}

func Test_UnmarshalData_Json(t *testing.T) {
	x := MockJson(true)
	var v MockJson
	testUnmarshalData(t, &v, "0n")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}

	x = MockJson(false)
	testUnmarshalData(t, &v, "0ff")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_String(t *testing.T) {
	x := "hello"
	var v string
	testUnmarshalData(t, &v, "hello")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_Bool(t *testing.T) {
	x := false
	var v bool
	testUnmarshalData(t, &v, "false")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_Int(t *testing.T) {
	x := 123
	var v int
	testUnmarshalData(t, &v, "123")
	if v != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_Int32(t *testing.T) {
	x := 123
	var v int32
	testUnmarshalData(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_Int64(t *testing.T) {
	x := 123
	var v int64
	testUnmarshalData(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}
func Test_UnmarshalData_Unt(t *testing.T) {
	x := 123
	var v uint
	testUnmarshalData(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_Uint32(t *testing.T) {
	x := 123
	var v uint32
	testUnmarshalData(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_Uint64(t *testing.T) {
	x := 123
	var v uint64
	testUnmarshalData(t, &v, "123")
	if int(v) != x {
		t.Fatal("expected:", x, "got:", js(v))
	}
}

func Test_UnmarshalData_MapString(t *testing.T) {
	v := make(map[string]string)
	testUnmarshalData(t, v, "a", "alpha", "b", "beta")
	vv, ok := v["a"]
	if !ok {
		t.Fatal("missing 'a' in map", js(v))
	}
	if vv != "alpha" {
		t.Fatal("expected 'alpha'", "got:", js(vv))
	}
	vv, ok = v["b"]
	if !ok {
		t.Fatal("missing 'b' in map", v)
	}
	if vv != "beta" {
		t.Fatal("expected 'beta'", "got:", js(vv))
	}
}

func Test_UnmarshalData_Redis(t *testing.T) {
	var v MockRedis
	testUnmarshalData(t, &v, "alpha", "beta")
	if string(v) != "[alpha][beta]" {
		t.Fatal("expected '[alpha][beta]', got:", js(v))
	}
}

func testUnmarshalData(t *testing.T, v interface{}, data ...string) {
	b := make([][]byte, len(data))
	for i, s := range data {
		b[i] = []byte(s)
	}

	err := unmarshalData(b, v)
	if err != nil {
		t.Fatal("unexpected error:", err.Error())
	}
}

type MockJson bool

func (mj *MockJson) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "0n":
		*mj = MockJson(true)
	case "0ff":
		*mj = MockJson(false)
	}
	return nil
}

type MockRedis string

func (mr *MockRedis) UnmarshalRedis(bb [][]byte) error {
	s := ""
	for _, b := range bb {
		s += "[" + string(b) + "]"
	}
	*mr = MockRedis(s)
	return nil
}

func js(v interface{}) string {
	b, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return string(b)
}
