package server_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/testutil"
)

type testData struct {
	StartFn func(testData, context.Context) error
	StopFn  func(testData, context.Context) error
	server.ServerUtils
}

func (t testData) Start(ctx context.Context) error {
	return t.StartFn(t, ctx)
}

func (t testData) Stop(ctx context.Context) error {
	return t.StopFn(t, ctx)
}

func (testData) GetName() string {
	return "testData"
}

func (testData) GetServerConfig() *server.ServerConfig {
	return &server.ServerConfig{
		Type: server.ProcessServerType,
	}
}

var _ server.Serveable = testData{}

var TIMEOUT = 10 * time.Second

func WithStop(fn func(t testData, ctx context.Context) error) testData {
	return testData{
		StartFn: func(td testData, ctx context.Context) error {
			return nil
		},
		StopFn:      fn,
		ServerUtils: server.NewTestServerUtils(),
	}
}

func WithStart(fn func(t testData, ctx context.Context) error) testData {
	return testData{
		StartFn: fn,
		StopFn: func(td testData, ctx context.Context) error {
			return nil
		},
		ServerUtils: server.NewTestServerUtils(),
	}
}

var tests = []struct {
	name              string
	testData          testData
	expectsStartError bool
	expectsStopError  bool
}{
	{
		name: "AsyncAtStart",
		testData: WithStart(func(t testData, ctx context.Context) error {
			errch1 := t.RunAsync(func(ctx context.Context) error {
				time.Sleep(1 * time.Second)
				return nil
			}, ctx, TIMEOUT)
			return <-errch1
		}),
	},

	{
		name:              "AsyncErrorAtStart",
		expectsStartError: true,
		testData: WithStart(func(t testData, ctx context.Context) error {
			errch1 := t.RunAsync(func(ctx context.Context) error {
				return errors.New("example error")
			}, ctx, TIMEOUT)
			return <-errch1
		}),
	},

	{
		name:              "AsyncOnOfMultipleReturnsErrorAtStart",
		expectsStartError: true,
		testData: WithStart(func(t testData, ctx context.Context) error {
			errch1 := t.RunAsync(func(ctx context.Context) error {
				return nil
			}, ctx, TIMEOUT)
			errch2 := t.RunAsync(func(ctx context.Context) error {
				return errors.New("example error")
			}, ctx, TIMEOUT)
			errch3 := t.RunAsync(func(ctx context.Context) error {
				return nil
			}, ctx, TIMEOUT)
			if err := <-errch1; err != nil {
				return err
			}
			if err := <-errch2; err != nil {
				return err
			}
			if err := <-errch3; err != nil {
				return err
			}
			return nil
		}),
	},

	{
		name:             "AsyncErrorAtStop",
		expectsStopError: true,
		testData: WithStop(func(t testData, ctx context.Context) error {
			errch1 := t.RunAsync(func(ctx context.Context) error {
				return errors.New("example error")
			}, ctx, TIMEOUT)
			return <-errch1
		}),
	},

	{
		name: "AsyncAtEnd",
		testData: WithStop(func(t testData, ctx context.Context) error {
			errch1 := t.RunAsync(func(ctx context.Context) error {
				time.Sleep(1 * time.Second)
				return nil
			}, ctx, TIMEOUT)
			return <-errch1
		}),
	},

	{
		name: "MultipleAsyncAtStart",
		testData: WithStart(func(t testData, ctx context.Context) error {
			errch1 := t.RunAsync(func(ctx context.Context) error {
				time.Sleep(1 * time.Second)
				return nil
			}, ctx, TIMEOUT)
			errch2 := t.RunAsync(func(ctx context.Context) error {
				time.Sleep(1 * time.Second)
				return nil
			}, ctx, TIMEOUT)
			<-errch1
			<-errch2
			return nil
		}),
	},
}

func TestRunAsyncAtStart(t *testing.T) {
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			appConfig := config.NewLocalTestAppConfig()
			srv := server.NewServer(func(su server.ServerUtils) testData {
				return v.testData
			}, appConfig)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if v.expectsStartError {
				log.Print("")
			}
			sErr := srv.Start(ctx)
			if v.expectsStartError {
				testutil.AssertNonNil(sErr, t)
			} else {
				testutil.AssertNilError(sErr, t)
			}

			eErr := srv.Stop(ctx)
			if v.expectsStopError {
				testutil.AssertNonNil(eErr, t)
			} else {
				testutil.AssertNilError(eErr, t)
			}
		})
	}
}
