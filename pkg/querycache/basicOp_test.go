package querycache

import (
	"reflect"
	"testing"
)

func getCacheStore() *CacheStore {
	return NewCacheStore("test", "localhost", 6379)
}

func TestCacheStore_SetGetDelBytes(t *testing.T) {
	type args struct {
		key       string
		valSet    []byte
		valGetset []byte
	}
	tests := []struct {
		name          string
		args          args
		wantGetset    []byte
		wantGet       []byte
		wantDel       int
		wantSetErr    bool
		wantGetsetErr bool
		wantGetErr    bool
		wantDelErr    bool
	}{
		{
			name:          "regular",
			args:          args{"regular", []byte("regular-value1"), []byte("regular-value2")},
			wantGetset:    []byte("regular-value1"),
			wantGet:       []byte("regular-value2"),
			wantDel:       1,
			wantSetErr:    false,
			wantGetsetErr: false,
			wantGetErr:    false,
			wantDelErr:    false,
		},
		{
			name:          "with-space",
			args:          args{"with-space", []byte("spaced value 1"), []byte("spaced value 2")},
			wantGetset:    []byte("spaced value 1"),
			wantGet:       []byte("spaced value 2"),
			wantDel:       1,
			wantSetErr:    false,
			wantGetsetErr: false,
			wantGetErr:    false,
			wantDelErr:    false,
		},
	}
	cs := getCacheStore()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cs.SetBytes(tt.args.key, tt.args.valSet); (err != nil) != tt.wantSetErr {
				t.Errorf("CacheStore.SetBytes() error = %v, wantSetErr %v", err, tt.wantSetErr)
			}
			gotGetset, err := cs.GetsetBytes(tt.args.key, tt.args.valGetset)
			if (err != nil) != tt.wantGetsetErr {
				t.Errorf("CacheStore.GetsetBytes() error = %v, wantGetsetErr %v", err, tt.wantGetsetErr)
				return
			}
			if !reflect.DeepEqual(gotGetset, tt.wantGetset) {
				t.Errorf("CacheStore.GetsetBytes() gotGetset = %v, wantGetset %v", gotGetset, tt.wantGetset)
			}
			gotGet, err := cs.GetBytes(tt.args.key)
			if (err != nil) != tt.wantGetErr {
				t.Errorf("CacheStore.GetBytes() error = %v, wantGetErr %v", err, tt.wantGetErr)
				return
			}
			if !reflect.DeepEqual(gotGet, tt.wantGet) {
				t.Errorf("CacheStore.GetBytes() gotGet = %v, wantGet %v", gotGet, tt.wantGet)
			}
			gotDel, err := cs.Delete(tt.args.key)
			if (err != nil) != tt.wantDelErr {
				t.Errorf("CacheStore.Delete() error = %v, wantDelErr %v", err, tt.wantDelErr)
				return
			}
			if gotDel != tt.wantDel {
				t.Errorf("CacheStore.Delete() gotDel = %v, wantDel %v", gotDel, tt.wantDel)
			}
		})
	}
}
