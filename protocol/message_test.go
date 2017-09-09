package protocol

import (
	"reflect"
	"testing"
)

func TestMessage_Bytes(t *testing.T) {
	type fields struct {
		Version       uint16
		SeqCounter    uint16
		Source        Contact
		PayloadType   PayloadType
		PayloadLenght uint8
		Payload       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Basic message",
			fields: fields{
				Version:    0x0a | (0x12 << 8),
				SeqCounter: 0x91 | (0x23 << 8),
				Source: Contact{
					Type:           0x01,
					CallsignLength: 0x00,
					Callsign:       []byte{},
					NumberIPs:      0,
					IPs:            []ContactIP{},
				},
				PayloadType:   0x91,
				PayloadLenght: 0x02,
				Payload:       []byte{0xaa, 0xbb},
			},
			want: []byte{0x0a, 0x12, 0x91, 0x23, 0x01, 0x00, 0x00, 0x91, 0x02, 0xaa, 0xbb},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:       tt.fields.Version,
				SeqCounter:    tt.fields.SeqCounter,
				Source:        tt.fields.Source,
				PayloadType:   tt.fields.PayloadType,
				PayloadLenght: tt.fields.PayloadLenght,
				Payload:       tt.fields.Payload,
			}
			if got := m.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		{
			name: "Basic parse",
			args: args{
				buf: []byte{0x0a, 0x12, 0x91, 0x23, 0x01, 0x00, 0x00, 0x91, 0x02, 0xaa, 0xbb},
			},
			want: Message{
				Version:    0x0a | (0x12 << 8),
				SeqCounter: 0x91 | (0x23 << 8),
				Source: Contact{
					Type:           0x01,
					CallsignLength: 0x00,
					Callsign:       []byte{},
					NumberIPs:      0,
					IPs:            []ContactIP{},
				},
				PayloadType:   0x91,
				PayloadLenght: 0x02,
				Payload:       []byte{0xaa, 0xbb},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseMessage(tt.args.buf); !reflect.DeepEqual(got.Bytes(), tt.want.Bytes()) {
				t.Errorf("ParseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
