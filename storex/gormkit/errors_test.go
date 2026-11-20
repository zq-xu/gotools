package gormkit

import (
	"net/http"
	"testing"

	"gorm.io/gorm"
)

func TestNewErrorInfoForGetErrorNotFound(t *testing.T) {
	ei := NewErrorInfoForGetError(gorm.ErrRecordNotFound)
	if ei == nil {
		t.Fatalf("expect a non-nil ErrorInfo for ErrRecordNotFound")
	}

	if ei.StatusCode() != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", ei.StatusCode(), http.StatusNotFound)
	}
}

func TestNewErrorInfoForGetErrorNil(t *testing.T) {
	if ei := NewErrorInfoForGetError(nil); ei != nil {
		t.Fatalf("expect nil ErrorInfo for nil error, got %v", ei)
	}
}
