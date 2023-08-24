// See https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/
package storage

import (
	"os"
	"screenshoter/utils"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	godotenv.Load(utils.RootDir() + "/.env")
	code := m.Run()
	os.Exit(code)
}

func TestUpload(t *testing.T) {
	err := Upload("wikipedia.org.png", utils.RootDir()+"/fixtures/wikipedia.org.png")
	assert.NoError(t, err)
}

func TestGetPresignURL(t *testing.T) {
	got, err := GetPresignURL("wikipedia.org.png")
	assert.NoError(t, err)
	log.Info(got)
	assert.NotEmpty(t, got)
}
