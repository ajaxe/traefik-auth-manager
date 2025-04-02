package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"github.com/ajaxe/traefik-auth-manager/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestHostedAppHandler_Validation_OK(t *testing.T) {
	sut := &hostedAppHandler{}
	data := correctHostedAppModel()

	if err := sut.validate(data); err != nil {
		t.Error(err)
	}
}

func TestHostedAppHandler_ValidateName_ReturnBadRequest(t *testing.T) {
	sut := &hostedAppHandler{}
	data := correctHostedAppModel()
	data.Name = ""

	err := sut.validate(data)

	assertRequriedErrMessage(t, "app name", err)

	assertHttpError(t, http.StatusBadRequest, err)
}
func TestHostedAppHandler_ValidateServiceToken_ReturnBadRequest(t *testing.T) {
	sut := &hostedAppHandler{}
	data := correctHostedAppModel()
	data.ServiceToken = ""

	err := sut.validate(data)

	assertRequriedErrMessage(t, "service token", err)

	assertHttpError(t, http.StatusBadRequest, err)
}
func TestHostedAppHandler_ValidateServiceURL_ReturnBadRequest(t *testing.T) {
	sut := &hostedAppHandler{}
	data := correctHostedAppModel()
	data.ServiceURL = ""

	err := sut.validate(data)

	assertRequriedErrMessage(t, "url", err)

	assertHttpError(t, http.StatusBadRequest, err)
}

func correctHostedAppModel() models.HostedApplication {
	return models.HostedApplication{
		ID:           bson.NewObjectID(),
		Name:         "test",
		ServiceToken: "1234567",
		ServiceURL:   "https://foo.com/",
		Active:       true,
	}
}

func assertRequriedErrMessage(t *testing.T, errToken string, err error) {
	m := err.Error()
	if !(strings.Contains(m, "is required.") && strings.Contains(m, errToken)) {
		t.Errorf("Expect error to contain '%s', got: '%s'", errToken, m)
	}
}

func assertHttpError(t *testing.T, httpStatus int, err error) {
	if err == nil {
		t.Error("Expect error")
	}
	ap, ok := err.(*helpers.AppError)
	if !ok {
		t.Error("Expect error of type 'helpers.AppError'")
	}
	if ap.HTTPStatus() != httpStatus {
		t.Errorf("Expect http status: %v got: %v", httpStatus, ap.HTTPStatus())
	}
}
