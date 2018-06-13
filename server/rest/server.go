package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liyue201/go-logger"
	"net/http"
	"strconv"
	"wallet-scan/btc-node-service/btc"
	"wallet-scan/btc-node-service/logic"
)

type RestSever struct {
	httpServer *http.Server
	txmgr      *logic.TxMgr
	btcCli     *btc.Client
}

func NewHttpServer(port int, txmgr *logic.TxMgr, btcClient *btc.Client) *RestSever {
	gin.SetMode(gin.ReleaseMode)
	engin := gin.Default()

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	httpServer := &http.Server{Addr: addr, Handler: engin}
	server := &RestSever{
		httpServer: httpServer,
		txmgr:      txmgr,
		btcCli:     btcClient,
	}
	server.initRoute(engin)

	return server
}

func (s *RestSever) Run() {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		logger.Errorf("RestSever.Run %s", err)
	}
}

func (s *RestSever) Stop() {
	s.httpServer.Shutdown(context.Background())
}

func (s *RestSever) initRoute(r gin.IRouter) {
	r.GET("/api/v1/txs", s.GetAddressTransations)
	r.GET("/api/v1/addr/:addr/utxo", s.GetAddressUtxo)
	r.POST("/api/v1/tx/send", s.SendTx)
}

func (s *RestSever) GetAddressTransations(c *gin.Context) {
	addr := c.Query("address")

	strLimit := c.Query("limit")
	if strLimit == "" {
		strLimit = "100"
	}
	limit, _ := strconv.Atoi(strLimit)
	strOffsert := c.Query("offset")
	if strOffsert == "" {
		strOffsert = "0"
	}
	offset, _ := strconv.Atoi(strOffsert)
	nextNodeId := c.Query("next_node_id")

	logger.Debugf("addr:%s, limit:%d, offset:%d", addr, limit, offset)

	if addr == "" {
		RespJson(c, BadRequest, nil)
		return
	}

	txs, nextNodeId, _ := s.txmgr.GetAddressTxs(addr, uint(offset), uint(limit), -1, nextNodeId)

	ret := struct {
		Txs        interface{} `json:"txs"`
		NextNodeId string      `json:"nextNodeId"`
	}{
		Txs:        txs,
		NextNodeId: nextNodeId,
	}

	RespJson(c, OK, ret)
}

func (s *RestSever) GetAddressUtxo(c *gin.Context) {
	addr := c.Param("addr")

	logger.Debugf("addr:%s", addr)

	utxos, _ := s.txmgr.GetAddressUtxos(addr)

	RespJson(c, OK, utxos)
}

func (s *RestSever) SendTx(c *gin.Context) {
	req := struct {
		Rawtx string `json:"rawtx" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("SendTx %s", err)
		RespJson(c, BadRequest, nil)
		return
	}

	txid, err := s.btcCli.SendTx(req.Rawtx)
	if err != nil {
		RespJson(c, SendTxFail, nil)
		return
	}
	ret := struct {
		Txid string `json:"txid"`
	}{
		Txid: txid,
	}
	RespJson(c, OK, ret)
}
