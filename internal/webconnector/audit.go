package webconnector

import (
	"log"
	"strings"
)

type Audit struct{}

func NewAudit() *Audit {
	return &Audit{}
}

func (a *Audit) Startup(addr, baseDir, allowedOrigin string) {
	log.Printf("[audit] startup version=%s addr=%s base_dir=%s allowed_origin=%s",
		Version, addr, baseDir, allowedOrigin)
}

func (a *Audit) PairSuccess(origin, remote string) {
	log.Printf("[audit] pair_success origin=%s remote=%s", origin, remote)
}

func (a *Audit) PairFailure(origin, remote string) {
	log.Printf("[audit] pair_failure origin=%s remote=%s", origin, remote)
}

func (a *Audit) AuthFailure(path, remote string) {
	log.Printf("[audit] auth_failure path=%s remote=%s", path, remote)
}

func (a *Audit) SandboxViolation(path, remote string) {
	log.Printf("[audit] sandbox_violation path=%s remote=%s", truncatePath(path), remote)
}

func (a *Audit) MimeRejected(path, mime, remote string) {
	log.Printf("[audit] mime_rejected path=%s mime=%s remote=%s", truncatePath(path), mime, remote)
}

func (a *Audit) ConfigUpdated(fields []string) {
	log.Printf("[audit] config_updated fields=%s", strings.Join(fields, ","))
}

func (a *Audit) PairCodeRefreshed() {
	log.Printf("[audit] pair_code_refreshed")
}

func (a *Audit) ConfigReloaded(configPath string, portChanged bool) {
	log.Printf("[audit] config_reloaded path=%s port_changed=%v", configPath, portChanged)
}

func (a *Audit) OriginRejected(origin, path string) {
	log.Printf("[audit] origin_rejected origin=%s path=%s", origin, path)
}

func truncatePath(p string) string {
	if len(p) <= 120 {
		return p
	}
	return p[:60] + "..." + p[len(p)-40:]
}
