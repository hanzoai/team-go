package metrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	metric "github.com/luxfi/metric"
)

// The vars in this package are built by the package-level constructors, which
// register into the default registry, and /v1/metrics serves metric.Handler().
// Those two ends have to name the same registry or the endpoint answers 200 with
// nothing in it — a well-formed reply that no scraper can tell from a service
// that is simply idle. So this records through the real vars and reads the bytes
// back off the real handler.
func TestTheEndpointCarriesTheMetricsThisPackageRecords(t *testing.T) {
	buildInfo.WithLabelValues("v-test").Set(1)
	reqTotal.WithLabelValues("/v1/iam", http.MethodGet, "200").Inc()
	reqDuration.WithLabelValues("/v1/iam", http.MethodGet).Observe(0.25)

	rec := httptest.NewRecorder()
	metric.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/metrics", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if body == "" {
		t.Fatal("/v1/metrics is empty: the handler and the constructors hold different registries")
	}
	for _, name := range []string{
		"team_build_info",
		"team_http_requests_total",
		"team_http_request_duration_seconds",
	} {
		if !strings.Contains(body, name) {
			t.Errorf("the exposition does not carry %q", name)
		}
	}
	if !strings.Contains(body, `team_build_info{version="v-test"} 1`) {
		t.Errorf("team_build_info is named but carries no value.\n%s", body)
	}
}
