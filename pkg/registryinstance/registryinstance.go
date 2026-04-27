package registryinstance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Config struct {
	ManagerBaseURL string
	Secret         string
	ServiceName    string
	InstanceID     string
	GrpcPublicAddr string
	HTTPBase       string
	Interval       time.Duration
	HTTPClient     *http.Client
}

func Run(ctx context.Context, log *zap.Logger, cfg Config) error {
	cfg.ManagerBaseURL = strings.TrimSuffix(strings.TrimSpace(cfg.ManagerBaseURL), "/")
	if cfg.ManagerBaseURL == "" {
		return nil
	}
	if strings.TrimSpace(cfg.GrpcPublicAddr) == "" && strings.TrimSpace(cfg.HTTPBase) == "" {
		return fmt.Errorf("registryinstance: grpc_public_addr or http_base required when manager_base_url is set")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 10 * time.Second
	}
	id := strings.TrimSpace(cfg.InstanceID)
	if id == "" {
		id = uuid.NewString()
	}
	cli := cfg.HTTPClient
	if cli == nil {
		cli = &http.Client{Timeout: 15 * time.Second}
	}
	for {
		if err := register(ctx, cli, cfg, id); err != nil {
			log.Warn("registry register retry", zap.Error(err))
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(2 * time.Second):
			}
			continue
		}
		break
	}
	t := time.NewTicker(cfg.Interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			if err := heartbeat(ctx, cli, cfg, id); err != nil {
				log.Warn("registry heartbeat", zap.Error(err))
			}
		}
	}
}

type registerBody struct {
	InstanceID  string `json:"instance_id"`
	ServiceName string `json:"service_name"`
	GrpcAddr    string `json:"grpc_addr,omitempty"`
	HTTPBase    string `json:"http_base,omitempty"`
}

func register(ctx context.Context, cli *http.Client, cfg Config, id string) error {
	body, err := json.Marshal(registerBody{
		InstanceID:  id,
		ServiceName: cfg.ServiceName,
		GrpcAddr:    cfg.GrpcPublicAddr,
		HTTPBase:    cfg.HTTPBase,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.ManagerBaseURL+"/api/v1/instances", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s := strings.TrimSpace(cfg.Secret); s != "" {
		req.Header.Set("X-Registry-Token", s)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("register status %d", resp.StatusCode)
	}
	return nil
}

func heartbeat(ctx context.Context, cli *http.Client, cfg Config, id string) error {
	u := fmt.Sprintf("%s/api/v1/instances/%s/heartbeat", cfg.ManagerBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return err
	}
	if s := strings.TrimSpace(cfg.Secret); s != "" {
		req.Header.Set("X-Registry-Token", s)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat status %d", resp.StatusCode)
	}
	return nil
}
