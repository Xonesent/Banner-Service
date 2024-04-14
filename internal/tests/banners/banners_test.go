package banners

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"net/http"
	"testing"
)

type TestStruct struct {
	name   string
	method string

	paramsInput string
	endpoint    string

	headers map[string]string
	reqBody map[string]interface{}

	prepare bool

	statusCode       int
	responseBody     map[string]interface{}
	manyResponseBody []map[string]interface{}
}

func Test_MiddleWare(t *testing.T) {
	testsMiddleWare := []TestStruct{
		{
			name:     "Unauthorized",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
			},
			reqBody: map[string]interface{}{},

			statusCode: 401,
			responseBody: map[string]interface{}{
				"error_place": "MDWManager.CheckAuthToken.NilToken",
				"error_value": "%!!(MISSING)s(\u003cnil\u003e)",
			},
		},
		{
			name:     "Forbidden",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "user_token",
			},
			reqBody: map[string]interface{}{},

			statusCode: 403,
			responseBody: map[string]interface{}{
				"error_place": "MDWManager.CheckAuthToken.Forbidden",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
		{
			name:     "NotAllowed",
			method:   http.MethodPatch,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{},

			statusCode: 405,
			responseBody: map[string]interface{}{
				"error_place": "Method Not Allowed",
				"error_value": "nil",
			},
		},
	}
	for _, test := range testsMiddleWare {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func Test_AddBanner(t *testing.T) {
	testsSignUp := []TestStruct{
		{
			name:     "BadRequestNotEnoughArguments",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"feature_id": 1,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersHandlers.AddBanner.ReadRequest",
				"error_value": "Key: 'AddBannerRequest.TagIds' Error:Field validation for 'TagIds' failed on the 'required' tag",
			},
		},
		{
			name:     "OK",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 1,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			statusCode: 201,
			responseBody: map[string]interface{}{
				"banner_id": 1,
			},
		},
		{
			name:     "BadRequestNilTagIds",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{},
				"feature_id": 1,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersHandlers.AddBanner.NilTagIds",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
		{
			name:     "BadRequestAlreadyExist",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 1,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersUC.AddBanner.AlreadyExists",
				"error_value": "these banners already exists [{1 1} {2 1} {3 1}]",
			},
		},
	}

	for _, test := range testsSignUp {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func Test_PatchBanner(t *testing.T) {
	testsPatchBanner := []TestStruct{
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 2,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:        "OK",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "2",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{4, 5, 6},
				"feature_id": 2,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			statusCode: 200,
			responseBody: map[string]interface{}{
				"message": "Success",
			},
		},
		{
			name:        "NothingToUpdate",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "2",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{4, 5, 6},
				"feature_id": 2,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersUC.PatchBanner.NothingToUpdate",
				"error_value": "nothing to update",
			},
		},
		{
			name:        "NotFound",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "100",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{},

			statusCode: 404,
			responseBody: map[string]interface{}{
				"error_place": "BannersUC.PatchBanner.ErrNotFound",
				"error_value": "impossible to update, banner with id 100 doesnt exist",
			},
		},
		{
			name:        "BadRequestParams",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "a",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersHandlers.PatchBanner.WrongBannerParams",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 3,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:        "AlreadyExist",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "2",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 3,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersUC.PatchBanner.AlreadyExists",
				"error_value": "these banners already exists {tag_id feature_id} [{1 3} {2 3} {3 3}]",
			},
		},
	}

	for _, test := range testsPatchBanner {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func Test_DeleteBanner(t *testing.T) {
	testsDeleteBanner := []TestStruct{
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 4,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:        "OK",
			method:      http.MethodDelete,
			endpoint:    "/banner",
			paramsInput: "4",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{},

			statusCode:   204,
			responseBody: nil,
		},
		{
			name:        "WrongParams",
			method:      http.MethodDelete,
			endpoint:    "/banner",
			paramsInput: "a",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersHandlers.DeleteBanner.WrongBannerParams",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
		{
			name:        "WrongParams",
			method:      http.MethodDelete,
			endpoint:    "/banner",
			paramsInput: "100",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{},

			statusCode: 404,
			responseBody: map[string]interface{}{
				"error_place": "BannersUC.DeleteBanner.DoNotExist",
				"error_value": "impossible to delete, banner with id 100 doesnt exist",
			},
		},
	}

	for _, test := range testsDeleteBanner {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func Test_GetBanner(t *testing.T) {
	testsGetBanner := []TestStruct{
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 5,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": false,
			},

			prepare: true,
		},
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 6,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:     "OK",
			method:   http.MethodGet,
			endpoint: "/user_banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_id":           1,
				"feature_id":       5,
				"use_last_version": false,
			},

			statusCode: 200,
			responseBody: map[string]interface{}{
				"title": "some_title",
				"text":  "some_text",
				"url":   "some_url",
			},
		},
		{
			name:     "BadRequest",
			method:   http.MethodGet,
			endpoint: "/user_banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_id":           1,
				"use_last_version": false,
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersHandlers.GetBanner.ReadRequest",
				"error_value": "Key: 'GetBannerRequest.FeatureId' Error:Field validation for 'FeatureId' failed on the 'required' tag",
			},
		},
		{
			name:     "NotFoundNotActivePostgres",
			method:   http.MethodGet,
			endpoint: "/user_banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "user_token",
			},
			reqBody: map[string]interface{}{
				"tag_id":           1,
				"feature_id":       5,
				"use_last_version": true,
			},

			statusCode: 404,
			responseBody: map[string]interface{}{
				"error_place": "BannersUC.GetBanner.NotAdmin",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
		{
			name:     "NotFoundNotActiveRedis",
			method:   http.MethodGet,
			endpoint: "/user_banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "user_token",
			},
			reqBody: map[string]interface{}{
				"tag_id":           1,
				"feature_id":       5,
				"use_last_version": true,
			},

			statusCode: 404,
			responseBody: map[string]interface{}{
				"error_place": "BannersUC.GetBanner.NotAdmin",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
	}

	for _, test := range testsGetBanner {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func Test_GetManyBanner(t *testing.T) {
	testsGetBanner := []TestStruct{
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{100, 101, 102},
				"feature_id": 7,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{100, 101, 102},
				"feature_id": 8,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{100, 101, 102},
				"feature_id": 9,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:     "OK",
			method:   http.MethodGet,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_id": 100,
				"limit":  0,
				"offset": 0,
			},

			statusCode: 200,
			manyResponseBody: []map[string]interface{}{
				{
					"banner_id":  8,
					"tag_ids":    []int64{100, 101, 102},
					"feature_id": 7,
					"content":    map[string]interface{}{"title": "some_title", "text": "some_text", "url": "some_url"},
					"is_active":  true,
					"created_at": "2024-04-14T22:24:20.550189Z",
					"updated_at": "2024-04-14T22:24:20.550189Z",
					"version":    1,
				},
				{
					"banner_id":  9,
					"tag_ids":    []int64{100, 101, 102},
					"feature_id": 8,
					"content":    map[string]interface{}{"title": "some_title", "text": "some_text", "url": "some_url"},
					"is_active":  true,
					"created_at": "2024-04-14T22:24:23.139481Z",
					"updated_at": "2024-04-14T22:24:23.139481Z",
					"version":    1,
				},
				{
					"banner_id":  10,
					"tag_ids":    []int64{100, 101, 102},
					"feature_id": 9,
					"content":    map[string]interface{}{"title": "some_title", "text": "some_text", "url": "some_url"},
					"is_active":  true,
					"created_at": "2024-04-14T22:24:25.536019Z",
					"updated_at": "2024-04-14T22:24:25.536019Z",
					"version":    1,
				},
			},
		},
	}

	for _, test := range testsGetBanner {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func Test_ViewVersions(t *testing.T) {
	testsViewVersions := []TestStruct{
		{
			name:     "AddBanner",
			method:   http.MethodPost,
			endpoint: "/banner",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids":    []int64{1, 2, 3},
				"feature_id": 10,
				"content": map[string]string{
					"title": "some_title",
					"text":  "some_text",
					"url":   "some_url",
				},
				"is_active": true,
			},

			prepare: true,
		},
		{
			name:        "PatchBanner",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "10",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids": []int64{2, 3, 4},
			},

			prepare: true,
		},
		{
			name:        "PatchBanner",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "10",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids": []int64{1, 2},
			},

			prepare: true,
		},
		{
			name:        "PatchBanner",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "10",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids": []int64{1},
			},

			prepare: true,
		},
		{
			name:        "PatchBanner",
			method:      http.MethodPatch,
			endpoint:    "/banner",
			paramsInput: "10",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},
			reqBody: map[string]interface{}{
				"tag_ids": []int64{1, 2, 3},
			},

			prepare: true,
		},
		{
			name:        "GetBannerVersions",
			method:      http.MethodGet,
			endpoint:    "/banner_versions",
			paramsInput: "10",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},

			statusCode: 200,
			manyResponseBody: []map[string]interface{}{
				{
					"BannerId":  10,
					"TagIds":    []int64{1, 2, 3},
					"FeatureId": 7,
					"Content":   map[string]interface{}{"Title": "some_title", "Text": "some_text", "Url": "some_url"},
					"IsActive":  true,
					"CreatedAt": "2024-04-14T21:25:53.796887Z",
					"UpdatedAt": "2024-04-14T22:06:27.103372Z",
					"Version":   4,
				},
				{
					"BannerId":  10,
					"TagIds":    []int64{1},
					"FeatureId": 7,
					"Content":   map[string]interface{}{"Title": "some_title", "Text": "some_text", "Url": "some_url"},
					"IsActive":  true,
					"CreatedAt": "2024-04-14T21:25:53.796887Z",
					"UpdatedAt": "2024-04-14T22:06:06.935378Z",
					"Version":   3,
				},
				{
					"BannerId":  10,
					"TagIds":    []int64{1, 2},
					"FeatureId": 7,
					"Content":   map[string]interface{}{"Title": "some_title", "Text": "some_text", "Url": "some_url"},
					"IsActive":  true,
					"CreatedAt": "2024-04-14T21:25:53.796887Z",
					"UpdatedAt": "2024-04-14T22:06:04.373926Z",
					"Version":   2,
				},
				{
					"BannerId":  10,
					"TagIds":    []int64{2, 3, 4},
					"FeatureId": 7,
					"Content":   map[string]interface{}{"Title": "some_title", "Text": "some_text", "Url": "some_url"},
					"IsActive":  true,
					"CreatedAt": "2024-04-14T21:25:53.796887Z",
					"UpdatedAt": "2024-04-14T22:05:28.280286Z",
					"Version":   1,
				},
			},
		},
	}

	for _, test := range testsViewVersions {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func Test_BannerRollBack(t *testing.T) {
	testsBannerRollBack := []TestStruct{
		{
			name:        "BannerRollBack",
			method:      http.MethodPut,
			endpoint:    "/banner_rollback",
			paramsInput: "10/2",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},

			statusCode: 200,
			responseBody: map[string]interface{}{
				"message": "Success",
			},
		},
		{
			name:        "BannerRollBack",
			method:      http.MethodPut,
			endpoint:    "/banner_rollback",
			paramsInput: "a/2",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersHandlers.ViewVersions.WrongBannerParams",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
		{
			name:        "BannerRollBack",
			method:      http.MethodPut,
			endpoint:    "/banner_rollback",
			paramsInput: "11/a",

			headers: map[string]string{
				"Content-Type": "application/json",
				"token":        "admin_token",
			},

			statusCode: 400,
			responseBody: map[string]interface{}{
				"error_place": "BannersHandlers.ViewVersions.WrongVersionParams",
				"error_value": "%!!(MISSING)s(<nil>)",
			},
		},
	}

	for _, test := range testsBannerRollBack {
		t.Run(test.name, func(t *testing.T) {
			runTest(test, t)
		})
	}
}

func runTest(test TestStruct, t *testing.T) {
	client := &http.Client{}

	requestBody, err := json.Marshal(test.reqBody)
	if err != nil {
		log.Fatalln(err)
	}
	if test.paramsInput != "" {
		test.endpoint = fmt.Sprintf("%s/%s", test.endpoint, test.paramsInput)
	}
	req, err := http.NewRequest(test.method, fmt.Sprintf("http://localhost:8892%s", test.endpoint), bytes.NewBuffer(requestBody))
	for key, value := range test.headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if test.prepare {
		return
	}

	utils.AssertEqual(t, test.statusCode, resp.StatusCode, "StatusCode")
	if resp.StatusCode == 204 {
		return
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var body interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			log.Fatalln(err)
		}

		if test.responseBody == nil {
			utils.AssertEqual(t, test.manyResponseBody, body, "ResponseBody")
		}

		utils.AssertEqual(t, test.responseBody, body, "ResponseBody")
	} else {
		var body map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			log.Fatalln(err)
		}
		utils.AssertEqual(t, test.responseBody["error_place"], body["error_place"], "ResponseBody")
		utils.AssertEqual(t, test.responseBody["error_value"], body["error_value"], "ResponseBody")
	}
}
