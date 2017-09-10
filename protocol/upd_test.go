package protocol

import (
	"reflect"
	"testing"
)

func TestUpdPayload_Bytes(t *testing.T) {
	type fields struct {
		Operation  uint8
		DataLength uint16
		Data       []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Basic encoding",
			fields: fields{
				Data:       []byte{1, 2, 3, 4, 5},
				DataLength: 5,
				Operation:  4,
			},
			want: []byte{4, 5, 0, 1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UpdPayload{
				Operation:  tt.fields.Operation,
				DataLength: tt.fields.DataLength,
				Data:       tt.fields.Data,
			}
			if got := u.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdPayload.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUpdPayload(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *UpdPayload
		wantErr bool
	}{
		{
			name: "Basic parse",
			args: args{
				buf: []byte{4, 5, 0, 1, 2, 3, 4, 5},
			},
			want: &UpdPayload{
				Data:       []byte{1, 2, 3, 4, 5},
				DataLength: 5,
				Operation:  4,
			},
			wantErr: false,
		},
		{
			name: "Invalid data",
			args: args{
				buf: []byte{4, 5, 0},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUpdPayload(tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUpdPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUpdPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdPayloadCacheRequest_Bytes(t *testing.T) {
	type fields struct {
		NumEntries uint32
		Entries    []UpdRequestCacheEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Basic cache entry",
			fields: fields{
				NumEntries: 1,
				Entries: []UpdRequestCacheEntry{
					{
						SeqCounter: 1,
						Source: Contact{
							Type:           1,
							CallsignLength: 0,
							NumberIPs:      0,
						},
					},
				},
			},
			want: []byte{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UpdPayloadCacheRequest{
				NumEntries: tt.fields.NumEntries,
				Entries:    tt.fields.Entries,
			}
			if got := r.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdPayloadCacheRequest.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCacheEntry(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name  string
		args  args
		want  UpdRequestCacheEntry
		want1 []byte
	}{
		{
			name: "Basic cache entry parse",
			args: args{
				buf: []byte{1, 0, 0, 0, 1, 0, 0, 1, 2, 3, 4, 5},
			},
			want: UpdRequestCacheEntry{
				SeqCounter: 1,
				Source: Contact{
					Type:           1,
					CallsignLength: 0,
					NumberIPs:      0,
				},
			},
			want1: []byte{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseCacheEntry(tt.args.buf)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCacheEntry() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseCacheEntry() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParsePayloadCacheRequest(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want UpdPayloadCacheRequest
	}{
		{
			name: "Basic payload cache request",
			args: args{
				buf: []byte{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0},
			},
			want: UpdPayloadCacheRequest{
				NumEntries: 1,
				Entries: []UpdRequestCacheEntry{
					{
						SeqCounter: 1,
						Source: Contact{
							Type:           1,
							CallsignLength: 0,
							NumberIPs:      0,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePayloadCacheRequest(tt.args.buf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePayloadCacheRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdRequestCacheEntry_Bytes(t *testing.T) {
	type fields struct {
		SeqCounter uint32
		Source     Contact
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &UpdRequestCacheEntry{
				SeqCounter: tt.fields.SeqCounter,
				Source:     tt.fields.Source,
			}
			if got := e.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdRequestCacheEntry.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdPayloadCacheResponse_Bytes(t *testing.T) {

	mbuf := []byte{1, 0, 0, 0}

	msg := Message{
		Flags:         0,
		Path:          "",
		PathLength:    0,
		PayloadLenght: 0,
		PayloadType:   2,
		Payload:       []byte{},
		SeqCounter:    23,
		Source: Contact{
			CallsignLength: 0,
			Callsign:       []byte{},
			IPs:            []ContactIP{},
			Type:           12,
			NumberIPs:      0,
		},
		TTL:     23,
		Version: 12,
	}

	mbuf = append(mbuf, msg.Bytes()...)

	type fields struct {
		NumEntries uint32
		Entries    []Message
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Basic cache response",
			fields: fields{
				NumEntries: 1,
				Entries: []Message{
					{
						Flags:         0,
						Path:          "",
						PathLength:    0,
						PayloadLenght: 0,
						PayloadType:   2,
						Payload:       []byte{},
						SeqCounter:    23,
						Source: Contact{
							CallsignLength: 0,
							Callsign:       []byte{},
							IPs:            []ContactIP{},
							Type:           12,
							NumberIPs:      0,
						},
						TTL:     23,
						Version: 12,
					},
				},
			},
			want: mbuf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UpdPayloadCacheResponse{
				NumEntries: tt.fields.NumEntries,
				Entries:    tt.fields.Entries,
			}
			if got := r.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdPayloadCacheResponse.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePayloadCacheResponse(t *testing.T) {
	mbuf := []byte{1, 0, 0, 0}

	msg := Message{
		Flags:         0,
		Path:          "",
		PathLength:    0,
		PayloadLenght: 0,
		PayloadType:   2,
		Payload:       []byte{},
		SeqCounter:    23,
		Source: Contact{
			CallsignLength: 0,
			Callsign:       []byte{},
			IPs:            []ContactIP{},
			Type:           12,
			NumberIPs:      0,
		},
		TTL:     23,
		Version: 12,
	}

	mbuf = append(mbuf, msg.Bytes()...)

	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want UpdPayloadCacheResponse
	}{
		{
			name: "Basic parse payload cache response",
			args: args{
				buf: mbuf,
			},
			want: UpdPayloadCacheResponse{
				NumEntries: 1,
				Entries: []Message{
					{
						Flags:         0,
						Path:          "",
						PathLength:    0,
						PayloadLenght: 0,
						PayloadType:   2,
						Payload:       []byte{},
						SeqCounter:    23,
						Source: Contact{
							CallsignLength: 0,
							Callsign:       []byte{},
							IPs:            []ContactIP{},
							Type:           12,
							NumberIPs:      0,
						},
						TTL:     23,
						Version: 12,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePayloadCacheResponse(tt.args.buf); !reflect.DeepEqual(got.Bytes(), tt.want.Bytes()) {
				t.Errorf("ParsePayloadCacheResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
