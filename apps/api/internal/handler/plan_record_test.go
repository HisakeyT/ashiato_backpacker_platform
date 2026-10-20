package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/gin-gonic/gin"
)

func recordBodyContext(body string) *gin.Context {
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	context.Request.Header.Set("Content-Type", "application/json")
	return context
}

func TestRecordBodyRejectsServerOwnedFields(t *testing.T) {
	for _, body := range []string{
		`{"name":"寺","cost":100}`, `{"name":"寺","sort_order":2}`,
		`{"name":"寺","travel_plan_id":4}`, `{"name":"寺","meal_type":"lunch"}`,
	} {
		_, _, ok := readRecordBody(recordBodyContext(body), model.RecordSightseeing, true)
		if ok {
			t.Fatalf("accepted %s", body)
		}
	}
}

func TestRecordBodyDistinguishesMissingAndNull(t *testing.T) {
	_, patch, ok := readRecordBody(recordBodyContext(`{"location_text":null,"map_url":null,"meal_type":null}`), model.RecordMeal, false)
	if !ok || patch.Name != nil || patch.LocationText == nil || *patch.LocationText != nil || patch.MapURL == nil || *patch.MapURL != nil || patch.MealType == nil || *patch.MealType != nil {
		t.Fatalf("patch = %+v, ok = %v", patch, ok)
	}
}
