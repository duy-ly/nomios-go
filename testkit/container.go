package testkit

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
)

func CreateMysqlContainer(resourcePath string) testcontainers.Container {
	ctx := context.Background()

	image := "arm64v8/mysql:oracle"
	if strings.Contains(runtime.GOARCH, "amd") {
		image = "mysql:8.0"
	}

	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithImage(image),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Env: map[string]string{
					"TZ": "Asia/Saigon",
				},
				Files: []testcontainers.ContainerFile{
					{
						HostFilePath:      resourcePath,
						ContainerFilePath: "/docker-entrypoint-initdb.d/schema.sql",
						FileMode:          0755,
					},
				},
			},
		}),
		mysqlcontainer.WithDatabase("nomios_db"),
		mysqlcontainer.WithUsername("root"),
		mysqlcontainer.WithPassword("12345678"),
		mysqlcontainer.WithConfigFile(filepath.Join(GetProjectRoot(), "testkit/resources/conf.d/mysql.cnf")),
	}

	container, err := mysqlcontainer.RunContainer(ctx, opts...)
	if err != nil {
		panic(fmt.Sprintf("can not run mysql container %v", err))
	}

	return container
}
