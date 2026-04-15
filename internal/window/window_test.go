package window

import (
	"testing"
	"time"
)

func TestNew_InvalidSize(t *testing.T) {
	_, err := New(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for size=0")
	}
}

func TestNew_InvalidDuration(t *testing.T) {
	_, err := New(5, 0)
	if err == nil {
		t.Fatal("expected error for duration=0")
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := New(5, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil slider")
	}
}

func TestAdd_IncreasesTotal(t *testing.T) {
	s, _ := New(10, time.Minute)
	s.Add(3)
	s.Add(7)
	if got := s.Total(); got != 10 {
		t.Fatalf("expected total 10, got %d", got)
	}
}

func TestTotal_EmptyWindow(t *testing.T) {
	s, _ := New(5, time.Minute)
	if got := s.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestBuckets_ReturnsCopy(t *testing.T) {
	s, _ := New(5, time.Minute)
	s.Add(1)
	b1 := s.Buckets()
	b1[0].Count = 999
	b2 := s.Buckets()
	if b2[0].Count == 999 {
		t.Fatal("Buckets should return an independent copy")
	}
}

func TestAdd_EvictsExpiredBuckets(t *testing.T) {
	s, _ := New(5, 50*time.Millisecond)
	s.Add(10)
	time.Sleep(60 * time.Millisecond)
	s.Add(2)
	if got := s.Total(); got != 2 {
		t.Fatalf("expected 2 after eviction, got %d", got)
	}
}

func TestAdd_RespectsSizeLimit(t *testing.T) {
	s, _ := New(3, time.Second)
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Millisecond)
		s.Add(1)
	}
	if len(s.Buckets()) > 3 {
		t.Fatalf("bucket count exceeded size limit")
	}
}
