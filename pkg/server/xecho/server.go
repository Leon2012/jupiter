// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xecho

import (
	"context"
	"os"

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
)

// Server ...
type Server struct {
	*echo.Echo
	config *Config
}

func newServer(config *Config) *Server {
	return &Server{
		Echo:   echo.New(),
		config: config,
	}
}

// Server implements server.Server interface.
func (s *Server) Serve() error {
	s.Echo.Logger.SetOutput(os.Stdout)
	s.Echo.Debug = s.config.Debug
	s.Echo.HideBanner = true
	s.Echo.StdLogger = xlog.JupiterLogger.StdLog()
	//s.Use(s.config.RecoverMiddleware())
	//s.Use(s.config.AccessLogger())
	//s.Use(s.config.PrometheusServerInterceptor())
	for _, route := range s.Echo.Routes() {
		s.config.logger.Info("echo add route", xlog.FieldMethod(route.Method), xlog.String("path", route.Path))
	}
	return s.Echo.Start(s.config.Address())
}

// Stop implements server.Server interface
// it will terminate echo server immediately
func (s *Server) Stop() error {
	return s.Echo.Close()
}

// GracefulStop implements server.Server interface
// it will stop echo server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}

// Info returns server info, used by governor and consumer balancer
// TODO(gorexlv): implements government protocol with juno
func (s *Server) Info() *server.ServiceInfo {
	return &server.ServiceInfo{
		Name:      pkg.Name(),
		Scheme:    "http",
		IP:        s.config.Host,
		Port:      s.config.Port,
		Weight:    0.0,
		Enable:    false,
		Healthy:   false,
		Metadata:  map[string]string{},
		Region:    "",
		Zone:      "",
		GroupName: "",
	}
}