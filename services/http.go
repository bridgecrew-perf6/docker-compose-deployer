package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const BEARER_SCHEMA string = "Bearer "

type deployResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (a *deployResponse) write(w http.ResponseWriter) {
	data, err := json.Marshal(a)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if a.Code >= 100 && a.Code < 600 {
		w.WriteHeader(a.Code)
	}
	fmt.Fprintln(w, string(data))
}

type deployHandler struct {
	*deployer
	secret string
}

func (a *deployHandler) getSecret(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(a.secret), nil
}

func (a *deployHandler) getClaims(r *http.Request) (*jwt.RegisteredClaims, error) {
	s := r.Header.Get("Authorization")
	if s == "" {
		return nil, fmt.Errorf("Unauthorized")
	}
	claims := new(jwt.RegisteredClaims)
	_, err := jwt.ParseWithClaims(s[len(BEARER_SCHEMA):], claims, a.getSecret)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (h *deployHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claims, err := h.getClaims(r)
	if err != nil {
		httpError(w, err, http.StatusUnauthorized)
		return
	}
	if claims.ID == "" || claims.Subject == "" {
		httpError(w, nil, http.StatusBadRequest)
		return
	}

	err = h.deployer.Run(claims.Subject, claims.ID)
	if err != nil {
		logrus.Errorf("deploy failed: %v", err)
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	resp := &deployResponse{
		Code:    0,
		Message: "OK",
		Data: map[string]interface{}{
			"svc": claims.Subject,
			"tag": claims.ID,
		},
	}
	resp.write(w)
}

func httpError(w http.ResponseWriter, err error, code int) {
	resp := new(deployResponse)
	resp.Code = code
	if err != nil {
		resp.Message = err.Error()
	} else {
		resp.Message = http.StatusText(code)
	}
	resp.write(w)
}

func NewDeployHandler() (*deployHandler, error) {
	deployer, err := NewDeployer()
	if err != nil {
		return nil, err
	}
	secret := viper.GetString("secret")
	return &deployHandler{deployer, secret}, nil
}
