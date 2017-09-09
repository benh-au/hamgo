package protocol

import (
	"reflect"
	"testing"
)

func TestContactIP_Bytes(t *testing.T) {
	type fields struct {
		Type   ContactIPType
		Length uint8
		Data   []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Basic",
			fields: fields{
				Type:   0x01,
				Length: 0x02,
				Data:   []byte{0x01, 0x02},
			},
			want: []byte{0x01, 0x02, 0x01, 0x02},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ContactIP{
				Type:   tt.fields.Type,
				Length: tt.fields.Length,
				Data:   tt.fields.Data,
			}
			if got := c.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContactIP.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseContactIP(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name  string
		args  args
		want  ContactIP
		want1 int
	}{
		{
			name: "Basic parse",
			args: args{
				buf: []byte{0x01, 0x02, 0x01, 0x02},
			},
			want: ContactIP{
				Type:   0x01,
				Length: 0x02,
				Data:   []byte{0x01, 0x02},
			},
			want1: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseContactIP(tt.args.buf)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseContactIP() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseContactIP() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestContact_Bytes(t *testing.T) {
	type fields struct {
		Type           ContactType
		CallsignLength uint8
		Callsign       []byte
		NumberIPs      uint8
		IPs            []ContactIP
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Basic contact serialize",
			fields: fields{
				Type:           0x01,
				CallsignLength: 0x02,
				Callsign:       []byte{0x03, 0x04},
				NumberIPs:      0x02,
				IPs: []ContactIP{
					ContactIP{
						Type:   0x05,
						Length: 0x03,
						Data:   []byte{0x08, 0x09, 0x0a},
					},
					ContactIP{
						Type:   0x06,
						Length: 0x01,
						Data:   []byte{0x01},
					},
				},
			},
			want: []byte{0x01, 0x02, 0x03, 0x04, 0x02, 0x05, 0x03, 0x08, 0x09, 0x0a, 0x06, 0x01, 0x01},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Contact{
				Type:           tt.fields.Type,
				CallsignLength: tt.fields.CallsignLength,
				Callsign:       tt.fields.Callsign,
				NumberIPs:      tt.fields.NumberIPs,
				IPs:            tt.fields.IPs,
			}
			if got := c.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Contact.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseContact(t *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name  string
		args  args
		want  Contact
		want1 []byte
	}{
		{
			name: "Basic contact parse",
			args: args{
				msg: []byte{0x01, 0x02, 0x03, 0x04, 0x02, 0x05, 0x03, 0x08, 0x09, 0x0a, 0x06, 0x01, 0x01, 0xaa},
			},
			want: Contact{
				Type:           0x01,
				CallsignLength: 0x02,
				Callsign:       []byte{0x03, 0x04},
				NumberIPs:      0x02,
				IPs: []ContactIP{
					ContactIP{
						Type:   0x05,
						Length: 0x03,
						Data:   []byte{0x08, 0x09, 0x0a},
					},
					ContactIP{
						Type:   0x06,
						Length: 0x01,
						Data:   []byte{0x01},
					},
				},
			},
			want1: []byte{
				0xaa,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseContact(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseContact() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseContact() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestContact_Compare(t *testing.T) {
	type fields struct {
		Type           ContactType
		CallsignLength uint8
		Callsign       []byte
		NumberIPs      uint8
		IPs            []ContactIP
	}
	type args struct {
		other *Contact
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Basic equal",
			fields: fields{
				Type:           0x01,
				CallsignLength: 0x02,
				Callsign:       []byte{0x01, 0x02},
				NumberIPs:      0x01,
				IPs: []ContactIP{
					{
						Data:   []byte{0x01, 0x02},
						Length: 0x02,
						Type:   0x01,
					},
				},
			},
			args: args{
				other: &Contact{
					Type:           0x01,
					CallsignLength: 0x02,
					Callsign:       []byte{0x01, 0x02},
					NumberIPs:      0x01,
					IPs: []ContactIP{
						{
							Data:   []byte{0x01, 0x02},
							Length: 0x02,
							Type:   0x01,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Type differs",
			fields: fields{
				Type:           0x01,
				CallsignLength: 0x02,
				Callsign:       []byte{0x01, 0x02},
				NumberIPs:      0x01,
				IPs: []ContactIP{
					{
						Data:   []byte{0x01, 0x02},
						Length: 0x02,
						Type:   0x01,
					},
				},
			},
			args: args{
				other: &Contact{
					Type:           0x02,
					CallsignLength: 0x02,
					Callsign:       []byte{0x01, 0x02},
					NumberIPs:      0x01,
					IPs: []ContactIP{
						{
							Data:   []byte{0x01, 0x02},
							Length: 0x02,
							Type:   0x01,
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Callsign mismatch",
			fields: fields{
				Type:           0x01,
				CallsignLength: 0x02,
				Callsign:       []byte{0x01, 0x03},
				NumberIPs:      0x01,
				IPs: []ContactIP{
					{
						Data:   []byte{0x01, 0x02},
						Length: 0x02,
						Type:   0x01,
					},
				},
			},
			args: args{
				other: &Contact{
					Type:           0x01,
					CallsignLength: 0x02,
					Callsign:       []byte{0x01, 0x02},
					NumberIPs:      0x01,
					IPs: []ContactIP{
						{
							Data:   []byte{0x01, 0x02},
							Length: 0x02,
							Type:   0x01,
						},
					},
				},
			},
			want: false,
		},
		{
			name: "IPs mismatch",
			fields: fields{
				Type:           0x01,
				CallsignLength: 0x02,
				Callsign:       []byte{0x01, 0x02},
				NumberIPs:      0x01,
				IPs: []ContactIP{
					{
						Data:   []byte{0x01, 0x02},
						Length: 0x02,
						Type:   0x01,
					},
				},
			},
			args: args{
				other: &Contact{
					Type:           0x01,
					CallsignLength: 0x02,
					Callsign:       []byte{0x01, 0x02},
					NumberIPs:      0x01,
					IPs: []ContactIP{
						{
							Data:   []byte{0x01, 0x03},
							Length: 0x02,
							Type:   0x01,
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Contact{
				Type:           tt.fields.Type,
				CallsignLength: tt.fields.CallsignLength,
				Callsign:       tt.fields.Callsign,
				NumberIPs:      tt.fields.NumberIPs,
				IPs:            tt.fields.IPs,
			}
			if got := c.Compare(tt.args.other); got != tt.want {
				t.Errorf("Contact.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
