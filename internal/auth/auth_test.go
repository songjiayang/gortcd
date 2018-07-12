package auth

import (
	"testing"

	"github.com/gortc/stun"
)

func TestStatic_Auth(t *testing.T) {
	var (
		s = NewStatic([]StaticCredential{
			{Username: "username", Realm: "realm", Password: "password"},
		})
		i = stun.NewLongTermIntegrity("username", "realm", "password")
		u = stun.NewUsername("username")
	)
	t.Run("ZeroAlloc", func(t *testing.T) {
		m := stun.MustBuild(stun.BindingRequest, u, i)
		if testing.AllocsPerRun(10, func() {
			if _, err := s.Auth(m); err != nil {
				t.Fatal(err)
			}
		}) > 0 {
			t.Fatal("unexpected allocations")
		}
	})
	for _, tc := range []struct {
		name string
		m    *stun.Message
		ok   bool
	}{
		{
			name: "positive",
			m:    stun.MustBuild(stun.BindingRequest, u, i),
			ok:   true,
		},
		{
			name: "negative",
			m: stun.MustBuild(stun.BindingRequest, u,
				stun.NewLongTermIntegrity("username", "realm", "password2"),
			),
			ok: false,
		},
		{
			name: "bad username",
			m:    stun.MustBuild(stun.BindingRequest, stun.NewUsername("user"), i),
			ok:   false,
		},
		{
			name: "no username",
			m:    stun.MustBuild(stun.BindingRequest, i),
			ok:   false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotI, err := s.Auth(tc.m)
			if !tc.ok {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Error(err)
			}
			r := stun.MustBuild(tc.m, u, gotI)
			if _, err = s.Auth(r); err != nil {
				t.Error(err)
			}
		})
	}
}

func BenchmarkStatic_Auth(b *testing.B) {
	var (
		s = NewStatic([]StaticCredential{
			{Username: "username", Realm: "realm", Password: "password"},
		})
		i = stun.NewLongTermIntegrity("username", "realm", "password")
		u = stun.NewUsername("username")
		m = stun.MustBuild(stun.BindingRequest, u, i)
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := s.Auth(m)
		if err != nil {
			b.Fatal(err)
		}
	}
}