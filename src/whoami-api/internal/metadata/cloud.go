package metadata

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/robmoraes/dslab/assets/video-05/app/internal/config"
)

const imdsBaseURL = "http://169.254.169.254"

type CloudCollector struct {
	cfg    config.AWSMetadataConfig
	client *http.Client

	mu         sync.Mutex
	cached     CloudInfo
	cacheUntil time.Time
	token      string
	tokenUntil time.Time
}

func NewCloudCollector(cfg config.AWSMetadataConfig) *CloudCollector {
	return &CloudCollector{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *CloudCollector) Collect(parent context.Context) CloudInfo {
	now := time.Now().UTC()
	if c.cfg.Mode == "disabled" {
		return CloudInfo{
			Enabled:    false,
			Available:  false,
			CheckedAt:  now,
			CacheUntil: now,
			Error:      "disabled",
		}
	}

	c.mu.Lock()
	if now.Before(c.cacheUntil) {
		cached := c.cached
		c.mu.Unlock()
		return cached
	}
	c.mu.Unlock()

	ctx, cancel := context.WithTimeout(parent, c.cfg.Timeout)
	defer cancel()

	info := c.fetch(ctx, now)
	info.CacheUntil = now.Add(c.cfg.CacheTTL)

	c.mu.Lock()
	c.cached = info
	c.cacheUntil = info.CacheUntil
	c.mu.Unlock()

	return info
}

func (c *CloudCollector) fetch(ctx context.Context, checkedAt time.Time) CloudInfo {
	info := CloudInfo{
		Enabled:   true,
		Provider:  "aws-ec2-imds",
		CheckedAt: checkedAt,
	}

	token, _ := c.tokenFor(ctx, checkedAt)

	instanceID, err := c.getMetadata(ctx, token, "/latest/meta-data/instance-id")
	if err != nil {
		info.Error = err.Error()
		return info
	}

	info.Available = true
	info.InstanceID = instanceID
	info.InstanceType, _ = c.getMetadata(ctx, token, "/latest/meta-data/instance-type")
	info.AvailabilityZone, _ = c.getMetadata(ctx, token, "/latest/meta-data/placement/availability-zone")
	info.PrivateIPv4, _ = c.getMetadata(ctx, token, "/latest/meta-data/local-ipv4")
	info.PublicIPv4, _ = c.getMetadata(ctx, token, "/latest/meta-data/public-ipv4")
	info.Region = regionFromAvailabilityZone(info.AvailabilityZone)

	return info
}

func (c *CloudCollector) tokenFor(ctx context.Context, now time.Time) (string, error) {
	c.mu.Lock()
	if c.token != "" && now.Before(c.tokenUntil) {
		token := c.token
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, imdsBaseURL+"/latest/api/token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "60")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("imds token status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", err
	}

	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", fmt.Errorf("imds token empty")
	}

	c.mu.Lock()
	c.token = token
	c.tokenUntil = now.Add(55 * time.Second)
	c.mu.Unlock()

	return token, nil
}

func (c *CloudCollector) getMetadata(ctx context.Context, token string, path string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imdsBaseURL+path, nil)
	if err != nil {
		return "", err
	}
	if token != "" {
		req.Header.Set("X-aws-ec2-metadata-token", token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return "", nil
	default:
		return "", fmt.Errorf("imds metadata %s status %d", path, resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func regionFromAvailabilityZone(az string) string {
	if len(az) < 2 {
		return ""
	}
	last := az[len(az)-1]
	if last < 'a' || last > 'z' {
		return ""
	}
	return az[:len(az)-1]
}
