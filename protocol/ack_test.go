package protocol

import (
	"reflect"
	"testing"
)

func TestParseACKPayload(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want *ACKPayload
	}{
		{
			name: "ACK parse",
			args: args{
				buf: []byte{0x00, 0x00, 0x00, 0xab, 0x00, 0x00, 0x00},
			},
			want: &ACKPayload{
				SeqCounter: 0xab,
				Source: Contact{
					Callsign:       []byte{},
					CallsignLength: 0x00,
					IPs:            []ContactIP{},
					NumberIPs:      0,
					Type:           0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseACKPayload(tt.args.buf); !reflect.DeepEqual(got.Bytes(), tt.want.Bytes()) {
				t.Errorf("ParseACKPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestACKPayload_Bytes(t *testing.T) {
	type fields struct {
		Source     Contact
		SeqCounter uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Basic ack payload",
			fields: fields{
				SeqCounter: 0xab,
				Source: Contact{
					Callsign:       []byte{},
					CallsignLength: 0x00,
					IPs:            []ContactIP{},
					NumberIPs:      0,
					Type:           0,
				},
			},
			want: []byte{
				0x00, 0x00, 0x00, 0xab, 0x00, 0x00, 0x00,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ACKPayload{
				Source:     tt.fields.Source,
				SeqCounter: tt.fields.SeqCounter,
			}
			if got := a.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ACKPayload.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
