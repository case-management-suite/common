package service_test

import (
	"testing"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/service"
	"github.com/case-management-suite/testutil"
)

type SomethingElse struct {
	OtherName string
}

type MyData struct {
	Name string
	SE   SomethingElse
	service.ServiceUtils
}

func TestXxx(t *testing.T) {
	s := MyData{}
	ms := service.NewService(config.NewLocalTestAppConfig(), "myservice", &s)
	ms.Value.Logger.Info().Msg("Worked!")
	testutil.AssertTrue(ms.Value.IsSet, t)
}
