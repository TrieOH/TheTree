package backend

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"testing"

	"lib/authz"
	libjobs "lib/jobs"
	"lib/jobs/httpjobs"

	"github.com/go-chi/chi/v5"
)

func recorder(log *[]string, name string, openErr error) Module {
	return Module{
		Name: name,
		Open: func(context.Context, *Deps) (Close, error) {
			if openErr != nil {
				return nil, openErr
			}
			*log = append(*log, "open "+name)
			return func(context.Context) error {
				*log = append(*log, "close "+name)
				return nil
			}, nil
		},
	}
}

func TestModulesOpenInOrderAndCloseInReverse(t *testing.T) {
	var log []string
	closeAll, err := openModules(context.Background(), &Deps{}, []Module{
		recorder(&log, "a", nil), recorder(&log, "b", nil), recorder(&log, "c", nil),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = closeAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"open a", "open b", "open c", "close c", "close b", "close a"}
	if !slices.Equal(log, want) {
		t.Fatalf("log = %v, want %v", log, want)
	}
}

func TestFailedOpenClosesWhatOpened(t *testing.T) {
	var log []string
	_, err := openModules(context.Background(), &Deps{}, []Module{
		recorder(&log, "a", nil), recorder(&log, "b", nil), recorder(&log, "c", errors.New("boom")),
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	want := []string{"open a", "open b", "close b", "close a"}
	if !slices.Equal(log, want) {
		t.Fatalf("log = %v, want %v", log, want)
	}
}

type fakeJobs struct {
	libjobs.DelayedEnqueuer

	served  bool
	handler http.Handler
}

func (f *fakeJobs) Serve(*libjobs.Registry, ...libjobs.Periodic) { f.served = true }
func (f *fakeJobs) Served() bool                                 { return f.served }
func (f *fakeJobs) Handler() http.Handler                        { return f.handler }

func TestMountJobsRefusesAnUnservedCarrier(t *testing.T) {
	err := mountJobs(&Deps{Jobs: &fakeJobs{}}, chi.NewRouter())
	if err == nil {
		t.Fatal("expected an error when Serve was never called")
	}
}

func TestMountJobsMountsTheHandler(t *testing.T) {
	called := false
	r := chi.NewRouter()
	err := mountJobs(&Deps{Jobs: &fakeJobs{served: true, handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	})}}, r)
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, httpjobs.Path, nil)
	r.ServeHTTP(noopWriter{}, req)
	if !called {
		t.Fatal("the jobs handler was not mounted at httpjobs.Path")
	}
}

type noopWriter struct{}

func (noopWriter) Header() http.Header         { return http.Header{} }
func (noopWriter) Write(b []byte) (int, error) { return len(b), nil }
func (noopWriter) WriteHeader(int)             {}

const spec = `
openapi: 3.1.0
paths:
  /public:
    get:
      operationId: getPublic
      security: []
  /private:
    get:
      operationId: getPrivate
      security:
        - bearerAuth: []
`

func TestChainsFollowTheSpec(t *testing.T) {
	pass := func(next http.Handler) http.Handler { return next }
	d := &Deps{spec: []byte(spec)}
	chains, err := d.Chains(authz.Primitives{JWT: pass, APIKey: pass, Any: pass})
	if err != nil {
		t.Fatal(err)
	}
	if len(chains["GetPublic"]) != 0 {
		t.Fatalf("public operation has %d middlewares, want 0", len(chains["GetPublic"]))
	}
	if len(chains["GetPrivate"]) == 0 {
		t.Fatal("bearer operation has no auth middleware")
	}
}

func TestChainsFailWithoutAPrimitive(t *testing.T) {
	d := &Deps{spec: []byte(spec)}
	_, err := d.Chains(authz.Primitives{})
	if err == nil {
		t.Fatal("expected an error: the spec needs a JWT primitive")
	}
}
