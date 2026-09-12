package main

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestShutdownResources_ClosesHTTPBeforeMongo(t *testing.T) {
	var calls []string
	err := shutdownResources(
		context.Background(),
		func(context.Context) error {
			calls = append(calls, "http")
			return nil
		},
		func(context.Context) error {
			calls = append(calls, "mongo")
			return nil
		},
	)
	if err != nil {
		t.Fatalf("shutdownResources() error = %v", err)
	}
	if want := []string{"http", "mongo"}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("shutdown order = %v, want %v", calls, want)
	}
}

func TestShutdownResources_AttemptsBothAndJoinsErrors(t *testing.T) {
	httpErr := errors.New("http shutdown failed")
	mongoErr := errors.New("mongo close failed")
	mongoCalled := false

	err := shutdownResources(
		context.Background(),
		func(context.Context) error { return httpErr },
		func(context.Context) error {
			mongoCalled = true
			return mongoErr
		},
	)
	if !mongoCalled {
		t.Fatal("Mongo close was skipped after HTTP shutdown failure")
	}
	if !errors.Is(err, httpErr) || !errors.Is(err, mongoErr) {
		t.Fatalf("shutdownResources() error = %v, want both errors", err)
	}
}
