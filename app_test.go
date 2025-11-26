package onion

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/HnH/di"
	"github.com/stretchr/testify/suite"

	"github.com/HnH/onion/internal/mocks"
)

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

type ApplicationSuite struct {
	suite.Suite
}

func (suite *ApplicationSuite) TestNew() {
	var app = New(context.Background())
	suite.Require().IsType(&application{}, app)
}

func (suite *ApplicationSuite) TestNewOptions() {
	var (
		cnt = di.NewContainer()
		cfg = &mocks.Provider{ID: "something"}
		lay = &mocks.Layer{ID: "infra"}
		app = New(
			context.Background(),
			WithContainer(cnt),
			WithProviders(cfg),
			WithLayers(lay),
		).(*application)
	)

	suite.Require().Equal(app.container, cnt)
	suite.Require().Equal(len(app.providers), 1)
	suite.Require().Equal(app.providers[0], cfg)
	suite.Require().Equal(len(app.layers), 1)
	suite.Require().Equal(app.layers[0], lay)
}

func (suite *ApplicationSuite) TestListen() {
	var app = New(
		context.Background(),
		WithProviders(&mocks.Provider{ID: "something"}),
		WithLayers(&mocks.Layer{ID: "infra"}, &mocks.Layer{ID: "domain"}),
	).(*application)

	var chanErrors = make(chan error)
	go func() {
		if err := app.Listen(); err != nil {
			chanErrors <- err
		} else {
			close(chanErrors)
		}
	}()

	time.Sleep(500 * time.Millisecond) // ensure that application subscribed signals
	suite.Require().NoError(syscall.Kill(syscall.Getpid(), syscall.SIGINT))
	suite.Require().NoError(syscall.Kill(syscall.Getpid(), syscall.SIGINT))
	suite.Require().NoError(syscall.Kill(syscall.Getpid(), syscall.SIGINT))

	var out, ok = <-chanErrors
	suite.Require().NoError(out)
	suite.Require().False(ok)
}
