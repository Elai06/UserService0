package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"userService/internal/repository"
	"userService/internal/repository/mocks"
	"userService/model"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserService(ctrl)
	mockRepo := NewHTTPHandler(mockService)
	fakeData := repository.Data{
		UserID: gofakeit.Int64(),
		Name:   gofakeit.FirstName()}

	tests := []struct {
		name           string
		expectedResult model.CreateResult
		inputData      repository.Data
		setupMock      func()
		expectedError  bool
	}{
		{
			name: "Success create user",
			expectedResult: model.CreateResult{
				Message: "User created successfully",
				Result:  &mongo.InsertOneResult{InsertedID: ""},
			},
			inputData: fakeData,
			setupMock: func() {
				mockService.EXPECT().CreateUser(gomock.Any(), fakeData).Return(&mongo.InsertOneResult{InsertedID: ""}, nil)
			},
			expectedError: false,
		},
		{
			name: "Failed create user",
			expectedResult: model.CreateResult{
				Message: "",
				Result:  &mongo.InsertOneResult{InsertedID: ""},
			},
			inputData: fakeData,
			setupMock: func() {
				mockService.EXPECT().CreateUser(gomock.Any(), fakeData).Return(&mongo.InsertOneResult{InsertedID: ""}, fmt.Errorf("some error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.inputData)
			req := httptest.NewRequest(http.MethodPost, "/createTask", bytes.NewReader(reqBody))
			rec := httptest.NewRecorder()

			if tt.setupMock != nil {
				tt.setupMock()
			}

			mockRepo.createUser(rec, req)

			resultData := model.CreateResult{}
			err := json.Unmarshal(rec.Body.Bytes(), &resultData)

			if tt.expectedError {
				assert.NotEqual(t, tt.expectedResult, resultData)
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedResult, resultData)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	mockRepo := NewHTTPHandler(mockService)

	fakeData := repository.Data{
		UserID: gofakeit.Int64(),
		Name:   gofakeit.Name(),
	}

	tests := []struct {
		name           string
		expectedResult repository.Data
		inputData      int64
		setupMock      func()
		expectedError  bool
	}{
		{
			name: "Success get user",
			expectedResult: repository.Data{
				UserID: fakeData.UserID,
				Name:   fakeData.Name,
			},
			inputData: fakeData.UserID,
			setupMock: func() {
				mockService.EXPECT().GetUserByID(gomock.Any(), fakeData.UserID).Return(&fakeData, nil)
			},
			expectedError: false,
		},
		{
			name:           "Failed get user",
			expectedResult: fakeData,
			inputData:      fakeData.UserID,
			setupMock: func() {
				mockService.EXPECT().GetUserByID(gomock.Any(), fakeData.UserID).Return(&fakeData, fmt.Errorf("some error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/getUser?id="+strconv.FormatInt(tt.inputData, 10), nil)
			rec := httptest.NewRecorder()

			if tt.setupMock != nil {
				tt.setupMock()
			}

			mockRepo.getUserByID(rec, req)

			resultData := repository.Data{}
			err := json.Unmarshal(rec.Body.Bytes(), &resultData)

			if tt.expectedError {
				assert.NotEqual(t, tt.expectedResult, resultData)
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedResult, resultData)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserService(ctrl)
	mockRepo := NewHTTPHandler(mockService)

	fakeDatas := []repository.Data{
		getFakeData(),
		getFakeData(),
		getFakeData(),
		getFakeData(),
		getFakeData(),
	}

	tests := []struct {
		name           string
		expectedResult []repository.Data
		inputData      int64
		setupMock      func()
		expectedError  bool
	}{
		{
			name:           "Success get user",
			expectedResult: fakeDatas,
			setupMock: func() {
				mockService.EXPECT().GetUsers().Return(&fakeDatas, nil)
			},
			expectedError: false,
		},
		{
			name:           "Success get user",
			expectedResult: fakeDatas,
			setupMock: func() {
				mockService.EXPECT().GetUsers().Return(&fakeDatas, fmt.Errorf("some error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/getUsers", nil)
			rec := httptest.NewRecorder()

			if tt.setupMock != nil {
				tt.setupMock()
			}

			mockRepo.getUsers(rec, req)

			var resultData []repository.Data
			err := json.Unmarshal(rec.Body.Bytes(), &resultData)

			if tt.expectedError {
				assert.NotEqual(t, tt.expectedResult, resultData)
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedResult, resultData)
				assert.NoError(t, err)
			}
		})
	}
}

func TestStartServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		port        string
		expectError bool
	}{
		{"Valid Port", ":8081", false},
		{"Invalid Port", ":invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserService(ctrl)
			th := NewHTTPHandler(mockRepo)

			_, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				err := th.StartServer()
				if tt.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				cancel()
			}()

			time.Sleep(200 * time.Millisecond)

			cancel()
		})
	}
}

func getFakeData() repository.Data {
	return repository.Data{
		UserID: gofakeit.Int64(),
		Name:   gofakeit.Name()}
}
