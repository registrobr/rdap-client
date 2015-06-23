package rdap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/registrobr/rdap-client/Godeps/_workspace/src/github.com/aryann/difflib"
)

func diff(a, b interface{}) []difflib.DiffRecord {
	return difflib.Diff(strings.Split(spew.Sdump(a), "\n"),
		strings.Split(spew.Sdump(b), "\n"))
}

func createTestServers(object interface{}, entry string, bootstrapStatus, rdapStatus int) (*httptest.Server, *httptest.Server) {
	rdapTS := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rdapStatus > 0 {
				w.WriteHeader(rdapStatus)
				return
			}

			json.NewEncoder(w).Encode(object)
		}),
	)

	return rdapTS, httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if bootstrapStatus > 0 {
				w.WriteHeader(bootstrapStatus)
				return
			}

			registry := serviceRegistry{
				Version: "1.0",
				Services: []service{
					{
						{entry},
						{rdapTS.URL},
					},
				},
			}

			json.NewEncoder(w).Encode(registry)
		}),
	)
}

func objType(object interface{}) string {
	typ := reflect.TypeOf(object)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ.Name()
}
