package main

import (
	"fmt"
	"github.com/CaliOpen/gofido/config"
	"github.com/CaliOpen/gofido/store"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/tstranex/u2f"
	"net/http"
)

type FidoServer struct {
	Store  store.StoreInterface
	Config config.ServerConfig
}

type RegisterResponseData struct {
	Challenge        string `json:"challenge" binding:"required"`
	RegistrationData string `json:"registrationData" binding:"required"`
	ClientData       string `json:"clientData" binding:"required"`
	Version          string `json:"version" binding:"required"`
}

// gin handler for u2f considerations
func U2fMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "sameorigin")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	}
}

func (server *FidoServer) Initialize(config *config.Config, store store.StoreInterface) error {
	server.Store = store
	server.Config = config.Server
	log.Info("Server initialized to run AppId ", server.Config.AppId)
	return nil
}

func (server *FidoServer) Run() error {
	// s := u2f.StdServer(&FidoStore{}, server.AppId)

	// HTTP router
	router := gin.Default()
	router.Use(U2fMiddleware())

	// Local static test page enable ?
	if server.Config.Static.Enable {
		router.Static("/static", server.Config.Static.Directory)
	}

	// API
	api := router.Group("api")
	{
		api.GET("/:user/register", server.RegisterRequest)
		api.POST("/:user/register", server.RegisterResponse)

		api.GET("/:user/sign", server.SignRequest)
		api.POST("/:user/sign", server.SignResponse)
	}

	hostname := fmt.Sprintf("%s:%d", server.Config.ListenInterface, server.Config.ListenPort)

	if server.Config.Tls.Enable == true {
		log.Info("Server listening on https://", hostname)
		err := router.RunTLS(hostname, server.Config.Tls.CertFile, server.Config.Tls.KeyFile)
		if err != nil {
			log.Fatal("Start TLS failed: ", err)
		}
	} else {
		log.Info("Server listening on http://", hostname)
		router.Run(hostname)
	}
	return nil
}

func (server *FidoServer) RegisterRequest(c *gin.Context) {
	user := c.Param("user")
	challenge, err := server.Store.NewChallenge(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	registrations, err := server.Store.GetRegistrations(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req := u2f.NewWebRegisterRequest(&challenge, registrations)
	c.JSON(http.StatusOK, req)
}

func (server *FidoServer) RegisterResponse(c *gin.Context) {
	user := c.Param("user")

	var data RegisterResponseData
	var resp u2f.RegisterResponse

	if err := c.ShouldBind(&data); err != nil {
		log.Error("Invalid data ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp.RegistrationData = data.RegistrationData
	resp.ClientData = data.ClientData
	resp.Version = data.Version

	challenge, err := server.Store.GetChallenge(user, data.Challenge)
	if err != nil {
		log.Error("GetChallenge failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = server.Store.NewRegistration(user, challenge, resp)
	if err != nil {
		log.Error("Registration failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
	return
}

func (server *FidoServer) SignRequest(c *gin.Context) {
	user := c.Param("user")
	registrations, err := server.Store.GetRegistrations(user)
	if err != nil {
		log.Error("GetRegistrations failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	challenge, err := server.Store.NewChallenge(user)
	if err != nil {
		log.Error("NewChallenge failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req := challenge.SignRequest(registrations)
	for _, reg := range registrations {
		err = server.Store.InsertKeyChallenge(user, reg.KeyHandle, challenge)
		if err != nil {
			log.Error("InsertKeyChallenge failed: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, req)
}

func (server *FidoServer) SignResponse(c *gin.Context) {
	user := c.Param("user")

	signResp := &u2f.SignResponse{}
	if err := c.ShouldBind(&signResp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registrations, err := server.Store.GetRegistrations(user)
	if err != nil {
		log.Error("GetRegistrations failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, reg := range registrations {
		for _, key_challenge := range server.Store.GetKeyChallenges(user, reg.KeyHandle) {
			challenge_str := store.EncodeBase64(key_challenge.Challenge)
			challenge, err := server.Store.GetChallenge(user, challenge_str)
			if err != nil {
				log.Error("Challenge ", challenge_str, " not found for user ", user)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			counter, err := server.Store.GetKeyCounter(user, reg.KeyHandle)
			if err != nil {
				log.Error("No counter for key ", reg.KeyHandle)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Info("Authenticate key with challenge ", challenge_str, " and counter ", counter.Counter)
			newCounter, authErr := reg.Authenticate(*signResp, challenge, counter.Counter)
			if authErr == nil {
				err = server.Store.UpdateCounter(user, reg.KeyHandle, newCounter)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, reg)
				return
			}
		}
	}
	log.Error("No known registration authenticate sign response")
	c.JSON(http.StatusBadRequest, gin.H{"error": "authentication failed"})
}
