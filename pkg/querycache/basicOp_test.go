package querycache

import (
	"reflect"
	"testing"
)

func getCacheStore() *CacheStore {
	cs := &CacheStore{
		env: "test",
	}
	cs.InitConn("localhost", 6379)
	return cs
}

func TestCacheStore_SetGetDelBytes(t *testing.T) {
	type args struct {
		key string
		val []byte
	}
	tests := []struct {
		name       string
		args       args
		wantGet    []byte
		wantDel    int
		wantSetErr bool
		wantGetErr bool
		wantDelErr bool
	}{
		{
			name:       "regular",
			args:       args{"regular", []byte("regular-value")},
			wantGet:    []byte("regular-value"),
			wantDel:    1,
			wantSetErr: false,
			wantGetErr: false,
			wantDelErr: false,
		},
		{
			name:       "with-space",
			args:       args{"with-space", []byte("spaced value")},
			wantGet:    []byte("spaced value"),
			wantDel:    1,
			wantSetErr: false,
			wantGetErr: false,
			wantDelErr: false,
		},
	}
	cs := getCacheStore()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cs.SetBytes(tt.args.key, tt.args.val); (err != nil) != tt.wantSetErr {
				t.Errorf("CacheStore.SetBytes() error = %v, wantSetErr %v", err, tt.wantSetErr)
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
