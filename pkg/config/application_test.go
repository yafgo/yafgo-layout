package config

import (
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ApplicationTestSuite struct {
	suite.Suite
	testCount uint32

	appName   string
	envPrefix string
}

func TestApplicationTestSuite(t *testing.T) {
	s := &ApplicationTestSuite{
		appName:   "go-toy",
		envPrefix: "MYAPP",
	}

	cfgMode := "dev"
	cfgDir := "./.."
	cfgFile := filepath.Join(cfgDir, cfgMode+".yaml")
	color.Greenln("[cfgFile]", cfgFile)

	// 创建配置文件
	assert.Nil(t, s.createDemoConfig(cfgFile))

	// 初始化配置
	SetEnvPrefix(s.envPrefix)
	SetConfigDir(cfgDir)
	SetupConfig(cfgMode)

	// 执行测试
	suite.Run(t, s)

	// 删除配置文件
	assert.Equal(t, true, s.removeFile(cfgFile))

}

func (s *ApplicationTestSuite) SetupTest() {
	color.Grayln("SetupTest...")

	// 预先设置 env
	os.Setenv(s.envPrefix+"_APP_NAME", s.appName)
	os.Setenv(s.envPrefix+"_DB_HOST", "127.0.0.1")

	color.Grayln("SetupTest Done")
}

func (s *ApplicationTestSuite) TearDownTest() {
	s.testCount++
	color.Grayf("TearDownTest test count:%d\n", s.testCount)
}

func (s *ApplicationTestSuite) TestEnv() {
	// 测试获取 env 中的值
	conf := Config()
	s.Equal(s.appName, conf.GetString("APP_NAME"))
	s.Equal("127.0.0.1", conf.GetString("DB_HOST"))
}

func (s *ApplicationTestSuite) TestSet() {
	conf := Config()
	conf.Set("app", map[string]any{
		"env": "local",
	})

	s.Equal("local", conf.GetString("app.env"))
}

func (s *ApplicationTestSuite) TestGet() {
	conf := Config()
	s.Equal(s.appName, conf.Get("app.name"))
}

func (s *ApplicationTestSuite) TestGetString() {
	conf := Config()
	conf.Set("database", map[string]any{
		"default": "mysql",
		"connections": map[string]any{
			"mysql": map[string]any{
				"host": "127.0.0.1",
			},
		},
	})

	s.Equal(s.appName, conf.GetString("APP_NAME"))
	s.Equal("127.0.0.1", conf.GetString("database.connections.mysql.host"))
	s.Equal("mysql", conf.GetString("database.default"))
	color.Infoln("[database]", conf.Get("database"))
}

func (s *ApplicationTestSuite) TestGetInt() {
	conf := Config()
	s.Equal(3306, conf.GetInt("db.port"))
	s.Equal(3306, conf.GetInt("db.PORT")) // 大小写不敏感
}

func (s *ApplicationTestSuite) TestGetFloat64() {
	conf := Config()
	s.Equal(3.1415926, conf.GetFloat64("MY_PI"))
	s.Equal(3.1415926, conf.GetFloat64("my_pi")) // 大小写不敏感
}

func (s *ApplicationTestSuite) TestGetBool() {
	conf := Config()
	s.Equal(true, conf.GetBool("APP.DEBUG"))
	s.Equal(true, conf.GetBool("app.debug")) // 大小写不敏感
}

func (s *ApplicationTestSuite) TestGetDuration() {
	conf := Config()
	s.Equal(time.Second*30, conf.GetDuration("MY_TIMEOUT"))
	s.Equal((time.Hour*2)+(time.Minute*30), conf.GetDuration("CACHE_DURATION"))
}

// createDemoConfig 创建示例配置文件
func (s *ApplicationTestSuite) createDemoConfig(configFile string) (err error) {
	content := []byte(`app:
  name: go-toy
  env: local
  key: "app-key"
  debug: true
  url: http://localhost
  host: 0.0.0.0:3000

db:
  connection: mysql
  host: "127.0.0.1"
  port: 3306
  database: "dbname"
  username: "dbuser"
  password: "dbpassword"

redis:
  host: "127.0.0.1"
  port: 6379
  password: "redisPassword"

log_channel: single
my_pi: 3.1415926
my_timeout: 30s
cache_duration: 2.5h
`)
	err = s.createFile(configFile, content)
	return
}

// createFile 创建文件, 可选是否在创建时写入内容
func (s *ApplicationTestSuite) createFile(file string, content ...[]byte) (err error) {
	err = os.MkdirAll(path.Dir(file), os.ModePerm)
	if err != nil {
		return
	}

	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()

	// Write file content
	if len(content) > 0 {
		_, err = f.Write(content[0])
	}

	return
}

func (s *ApplicationTestSuite) removeFile(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		dir, err := os.ReadDir(file)

		if err != nil {
			return false
		}

		for _, d := range dir {
			err := os.RemoveAll(path.Join([]string{file, d.Name()}...))
			if err != nil {
				return false
			}
		}
	}

	err = os.Remove(file)
	return err == nil
}
