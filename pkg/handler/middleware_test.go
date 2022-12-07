package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/azonnix/todo-app/pkg/service"
	mock_service "github.com/azonnix/todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, tocken string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		tocken               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer tocken",
			tocken:      "tocken",
			mockBehavior: func(s *mock_service.MockAuthorization, tocken string) {
				s.EXPECT().ParseToken(tocken).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "No Header",
			headerName:           "",
			mockBehavior:         func(s *mock_service.MockAuthorization, tocken string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		{
			name:                 "Invalid Bearer",
			headerName:           "Authorization",
			headerValue:          "Bearr tocken",
			tocken:               "tocken",
			mockBehavior:         func(s *mock_service.MockAuthorization, tocken string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                 "No Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mock_service.MockAuthorization, tocken string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Service Failure",
			headerName:  "Authorization",
			headerValue: "Bearer tocken",
			tocken:      "tocken",
			mockBehavior: func(s *mock_service.MockAuthorization, tocken string) {
				s.EXPECT().ParseToken(tocken).Return(1, errors.New("faild to parse tocken"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"faild to parse tocken"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			auth := mock_service.NewMockAuthorization(controller)
			testCase.mockBehavior(auth, testCase.tocken)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			server := gin.New()
			server.GET("/protected", handler.UserIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/protected", nil)
			request.Header.Set(testCase.headerName, testCase.headerValue)

			server.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), testCase.expectedResponseBody)
		})
	}
}
