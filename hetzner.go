// Package hetzner is a Hetzner Cloud deploy driver for togo. It provisions a
// cloud server that installs Docker and runs the app image via cloud-init.
// Select with deploy.provider=hetzner; needs HCLOUD_TOKEN.
package hetzner

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/togo-framework/deploy"
	"github.com/togo-framework/togo"
)

func init() { deploy.RegisterDriver("hetzner", New) }

// New builds the Hetzner deployer (HCLOUD_TOKEN required).
func New(_ *togo.Kernel) (deploy.Deployer, error) {
	tok := os.Getenv("HCLOUD_TOKEN")
	if tok == "" {
		return nil, errors.New("deploy-hetzner: HCLOUD_TOKEN not set")
	}
	return &driver{c: hcloud.NewClient(hcloud.WithToken(tok))}, nil
}

type driver struct{ c *hcloud.Client }

func cloudInit(image string) string {
	return fmt.Sprintf("#cloud-config\nruncmd:\n  - curl -fsSL https://get.docker.com | sh\n  - docker run -d --name app --restart always -p 80:8080 %s\n", image)
}

func (d *driver) Provision(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	loc := spec.Region
	if loc == "" {
		loc = "nbg1"
	}
	st := "cx22"
	if v, ok := spec.Options["size"].(string); ok && v != "" {
		st = v
	}
	res, _, err := d.c.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:       spec.App,
		ServerType: &hcloud.ServerType{Name: st},
		Image:      &hcloud.Image{Name: "docker-ce"},
		Location:   &hcloud.Location{Name: loc},
		UserData:   cloudInit(spec.Image),
	})
	if err != nil {
		return nil, fmt.Errorf("hetzner provision: %w", err)
	}
	ip := ""
	if res.Server != nil && res.Server.PublicNet.IPv4.IP != nil {
		ip = res.Server.PublicNet.IPv4.IP.String()
	}
	return &deploy.Result{URL: "http://" + ip, Message: "server created; app boots via cloud-init", Raw: map[string]any{"id": res.Server.ID, "ip": ip}}, nil
}

func (d *driver) Deploy(ctx context.Context, spec deploy.Spec) (*deploy.Result, error) {
	s, _, err := d.c.Server.GetByName(ctx, spec.App)
	if err != nil || s == nil {
		return d.Provision(ctx, spec)
	}
	return &deploy.Result{URL: "http://" + s.PublicNet.IPv4.IP.String(), Message: "redeploy: restart the app container via your CI/SSH; server up", Raw: map[string]any{"id": s.ID}}, nil
}

func (d *driver) Destroy(ctx context.Context, spec deploy.Spec) error {
	s, _, err := d.c.Server.GetByName(ctx, spec.App)
	if err != nil || s == nil {
		return err
	}
	_, _, err = d.c.Server.DeleteWithResult(ctx, s)
	return err
}

func (d *driver) Status(ctx context.Context, spec deploy.Spec) (*deploy.Status, error) {
	s, _, err := d.c.Server.GetByName(ctx, spec.App)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return &deploy.Status{Healthy: false, Detail: "no server"}, nil
	}
	return &deploy.Status{Healthy: s.Status == hcloud.ServerStatusRunning, Detail: string(s.Status), Raw: map[string]any{"ip": s.PublicNet.IPv4.IP.String()}}, nil
}
