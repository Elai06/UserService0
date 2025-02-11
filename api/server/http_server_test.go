package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"userService/internal/repository"
	"userService/internal/repository/mocks"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIUserService(ctrl)
	mockRepo := NewHttpHandler(mockService)

	fakeData := repository.Data{
		UserId: gofakeit.Int64(),
		Name:   gofakeit.FirstName()}

	tests := []struct {
		name           string
		expectedResult createResult
		inputData      repository.Data
		setupMock      func()
		expectedError  bool
	}{
		{
			name: "Success create user",
			expectedResult: createResult{
				Message: "User created successfully",
				Result:  &mongo.InsertOneResult{InsertedID: ""},
			},
			inputData: fakeData,
			setupMock: func() {
				mockService.EXPECT().CreateUser(fakeData).Return(&mongo.InsertOneResult{InsertedID: ""}, nil)
			},
			expectedError: false,
		},
		{
			name: "Failed create user",
			expectedResult: createResult{
				Message: "",
				Result:  &mongo.InsertOneResult{InsertedID: ""},
			},
			inputData: fakeData,
			setupMock: func() {
				mockService.EXPECT().CreateUser(fakeData).Return(&mongo.InsertOneResult{InsertedID: ""}, fmt.Errorf("some error"))
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
			resultData := createResult{}
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

	mockService := mocks.NewMockIUserService(ctrl)
	mockRepo := NewHttpHandler(mockService)

	fakeData := repository.Data{
		UserId: gofakeit.Int64(),
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
				UserId: fakeData.UserId,
				Name:   fakeData.Name,
			},
			inputData: fakeData.UserId,
			setupMock: func() {
				mockService.EXPECT().GetUserByID(fakeData.UserId).Return(&fakeData, nil)
			},
			expectedError: false,
		},
		{
			name:           "Failed get user",
			expectedResult: fakeData,
			inputData:      fakeData.UserId,
			setupMock: func() {
				mockService.EXPECT().GetUserByID(fakeData.UserId).Return(&fakeData, fmt.Errorf("some error"))
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
			mockRepo.getUserById(rec, req)
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

	mockService := mocks.NewMockIUserService(ctrl)
	mockRepo := NewHttpHandler(mockService)

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
			mockRepo := mocks.NewMockIUserService(ctrl)
			th := NewHttpHandler(mockRepo)
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
		UserId: gofakeit.Int64(),
		Name:   gofakeit.Name()}
}
