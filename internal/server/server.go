package server

import (
	"encoding/json"
	"github.com/Rosi-Eliz/distributed-kv-store/internal/raft"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	router   *gin.Engine
	raftNode *raft.Node
}

func NewServer(raftNode *raft.Node) *Server {
	s := &Server{
		router:   gin.Default(),
		raftNode: raftNode,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.POST("/set", s.handleSet())
	s.router.GET("/get/:key", s.handleGet())
	s.router.DELETE("/delete/:key", s.handleDelete())
	s.router.POST("/join", s.handleJoin())
	s.router.POST("/transaction", s.handleTransaction()) // New transaction endpoint
}

func (s *Server) handleSet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Key   string        `json:"key"`
			Value string        `json:"value"`
			TTL   time.Duration `json:"ttl"` // TTL in seconds
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cmd := raft.Command{
			Op:    "set",
			Key:   req.Key,
			Value: req.Value,
			TTL:   req.TTL,
		}
		data, err := json.Marshal(cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		future := s.raftNode.Raft.Apply(data, 10*time.Second)
		if future.Error() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": future.Error().Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (s *Server) handleGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		value, ok := s.raftNode.FSM.Get(key)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})
	}
}

func (s *Server) handleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")

		cmd := raft.Command{
			Op:  "delete",
			Key: key,
		}
		data, err := json.Marshal(cmd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		future := s.raftNode.Raft.Apply(data, 10*time.Second)
		if future.Error() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": future.Error().Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (s *Server) handleJoin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID      string `json:"id"`
			Address string `json:"address"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.raftNode.AddVoter(req.ID, req.Address); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "node added"})
	}
}

func (s *Server) handleTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tx raft.Transaction
		if err := c.BindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data, err := json.Marshal(tx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		future := s.raftNode.Raft.Apply(data, 10*time.Second)
		if future.Error() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": future.Error().Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "transaction applied"})
	}
}

func (s *Server) Run(addr string) {
	s.router.Run(addr)
}
