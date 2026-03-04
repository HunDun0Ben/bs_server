package conf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadAllConfig_WithMockFS(t *testing.T) {
	// 1. 准备内存文件系统
	mockFS := afero.NewMemMapFs()
	SetFs(mockFS)
	defer SetFs(afero.NewOsFs())

	// 2. 创建模拟配置目录
	mockDir := "/etc/bs_server/conf"
	_ = mockFS.MkdirAll(mockDir, 0755)

	appContent := `
jwt:
  secret: "mock-secret"
otel:
  enable: true
  strict: true
  endpoint: "localhost:4317"
`
	_ = afero.WriteFile(mockFS, filepath.Join(mockDir, "application.yaml"), []byte(appContent), 0644)

	// 3. 设置环境变量并执行加载
	os.Setenv("APP_CONF", mockDir)
	defer os.Unsetenv("APP_CONF")

	err := InitConfig()

	// 4. 断言
	assert.NoError(t, err)
	assert.Equal(t, "mock-secret", AppConfig.JWT.Secret)
	assert.True(t, AppConfig.OTEL.Enable)
	assert.True(t, AppConfig.OTEL.Strict)
	assert.Equal(t, "localhost:4317", AppConfig.OTEL.Endpoint)
}

func TestLoadAllConfig_DefaultPathMissing(t *testing.T) {
	// 清除环境变量，让它尝试读取默认的 ./conf
	os.Unsetenv("APP_CONF")

	// 在一个保证没有 ./conf 的隔离环境下运行
	err := loadAllConfig()

	// 应该报错
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config directory not found")
}

func TestLoadConfigFiles_CorrectNamespace(t *testing.T) {
	mockFS := afero.NewMemMapFs()
	SetFs(mockFS)
	defer SetFs(afero.NewOsFs())

	mockDir := "/tmp/testconf"
	_ = mockFS.MkdirAll(mockDir, 0755)

	// 模拟多个配置文件
	_ = afero.WriteFile(mockFS, filepath.Join(mockDir, "mongodb.yaml"), []byte("mongodb:\n  uri: \"mongo-test\""), 0644)

	err := loadConfigFiles(mockDir)
	assert.NoError(t, err)

	// 验证 Namespace 映射
	v, ok := GetConfig("mongodb")
	assert.True(t, ok)
	assert.Equal(t, "mongo-test", v.GetString("mongodb.uri"))
}

func TestInitConfig_ErrorCase(t *testing.T) {
	// 模拟找不到目录的情况
	os.Setenv("APP_CONF", "/non-existent-directory")
	defer os.Unsetenv("APP_CONF")

	err := InitConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config directory not found")
}
