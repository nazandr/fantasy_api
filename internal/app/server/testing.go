package server

import "testing"

func TestServer(t *testing.T) *APIServer {
	t.Helper()
	conf := NewConfig()
	conf.Store.Database_url = "mongodb://localhost:27017/"
	conf.Store.DbName = "test_db"
	s := New(conf)
	if err := s.storeConfig(); err != nil {
		t.Fatal(err)
	}

	return s
}
