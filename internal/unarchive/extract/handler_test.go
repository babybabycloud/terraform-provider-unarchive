package extract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetZipHandler(t *testing.T) {
	handlerType := ".zip"
	handler := getHandler(handlerType)

	_, ok := handler.(*zipHandler)
	assert.True(t, ok)
}

func TestGetTarHandler(t *testing.T) {
	handlerType := ".tar"
	handler := getHandler(handlerType)

	_, ok := handler.(*tarHandler)
	assert.True(t, ok)
}

func TestGetTarGzHandler(t *testing.T) {
	handlerType := ".tar.gz"
	handler := getHandler(handlerType)

	_, ok := handler.(*targzHandler)
	assert.True(t, ok)
}

func TestGetUnknownHandler(t *testing.T) {
	handlerType := "any"
	handler := getHandler(handlerType)

	_, ok := handler.(*unknownHandler)
	assert.True(t, ok)
}
