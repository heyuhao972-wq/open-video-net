package protocol

import "testing"

func TestEncodeDecodeMessage(t *testing.T) {
	in := Message{
		Type: "chat",
		Data: []byte("hello"),
	}

	data, err := EncodeMessage(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	out, err := DecodeMessage(data)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Type != in.Type {
		t.Fatalf("expected type %s, got %s", in.Type, out.Type)
	}
	if string(out.Data) != string(in.Data) {
		t.Fatalf("expected data %s, got %s", string(in.Data), string(out.Data))
	}
}
