package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type reqParams struct {
	Method    string `json:"method,required"`
	BizParams string `json:"biz_params,required"`
}

var ServiceMap sync.Map

func HttpGateway(ctx context.Context, c *app.RequestContext) {
	serviceName := c.Param("svc")
	cli, ok := ServiceMap.Load(serviceName)
	if !ok {
		c.JSON(http.StatusNotFound, fmt.Errorf("Service %v not found", serviceName))
		return
	}
	var req reqParams
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.Warnf("BindAndValidate failed: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request parameters"})
		return
	}

	httpReq, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte(req.BizParams)))
	if err != nil {
		hlog.Warnf("Failed to create new HTTP request: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create new HTTP request"})
	}
	httpReq.URL.Path = fmt.Sprintf("/%s/%s", serviceName, req.Method)
	customReq, err := generic.FromHTTPRequest(httpReq)
	if err != nil {
		hlog.Errorf("Failed to convert HTTP request to custom request: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process request"})
		return
	}
	resp, err := cli.(genericclient.Client).GenericCall(ctx, "", customReq)
	if err != nil {
		hlog.Errorf("GenericCall error: %v", err)
		bizErr, ok := kerrors.FromBizStatusError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server handle error"})
			return
		}
		respMap := map[string]interface{}{
			"code":    bizErr.BizStatusCode(),
			"message": bizErr.BizMessage(),
		}
		c.JSON(http.StatusOK, respMap)
		return
	}
	realResp, ok := resp.(*generic.HTTPResponse)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid server response"})
		return
	}

	respMap := map[string]interface{}{
		"code":    0,
		"message": "ok",
		"data":    realResp.Body,
	}
	c.JSON(http.StatusOK, respMap)
}
