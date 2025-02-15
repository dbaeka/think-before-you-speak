package common

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testMsg struct {
	Msg string `json:"msg"`
}

func TestResponse(t *testing.T) {
	testCases := []struct {
		name     string
		code     int
		err      error
		success  bool
		redirect bool
		resp     interface{}
	}{
		{
			name: "with code 400",
			code: 400,
			err:  errors.New("some err"),
			resp: testMsg{
				Msg: "some err",
			},
		},
		{
			name: "with code 401",
			code: 401,
			err:  errors.New("some err"),
			resp: testMsg{
				Msg: "some err",
			},
		},
		{
			name: "with error",
			err:  errors.New("some err"),
			code: 500,
			resp: testMsg{
				Msg: "some err",
			},
		},
		{
			name:    "response success",
			success: true,
			code:    200,
			resp: testMsg{
				Msg: "success",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			e := echo.New()
			c := e.NewContext(req, w)

			if tc.success {
				ResponseSuccess(c, testMsg{
					Msg: "success",
				})
			} else {
				ResponseFailed(c, tc.code, tc.err)
			}

			resp := testMsg{}
			err := json.NewDecoder(w.Body).Decode(&resp)
			assert.Nil(t, err)
			assert.Equal(t, tc.resp, resp)
		})
	}
}

func TestResponseSuccessStream(t *testing.T) {
	server := sse.New()
	server.AutoReplay = false
	stream := server.CreateStream("test")
	if stream == nil {
		t.Fatal("failed to create stream")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := echo.New()
		c := e.NewContext(r, w)
		ResponseSuccessStream(c, server)
	}))
	defer ts.Close()

	client := sse.NewClient(ts.URL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgBuffer := make(chan []byte, 1)
	errChan := make(chan error, 1)
	go func(sseClient *sse.Client) {
		err := sseClient.SubscribeWithContext(ctx, "test", func(msg *sse.Event) {
			msgBuffer <- msg.Data
			cancel()
		})
		if err != nil {
			errChan <- err
		}
	}(client)

	select {
	case err := <-errChan:
		t.Fatalf("SSE subscription failed: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	go func(s *sse.Server) {
		s.Publish("test", &sse.Event{
			Data: []byte("test event"),
		})
	}(server)

	select {
	case msg := <-msgBuffer:
		assert.Equal(t, "test event", string(msg))
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for ResponseSuccessStream to complete")
	}
}
