package webui

import (
	"embed"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/radityprtama/proxygate/v6/internal/config"
	"github.com/radityprtama/proxygate/v6/internal/util"
	coreauth "github.com/radityprtama/proxygate/v6/sdk/cliproxy/auth"
)

//go:embed all
var templates embed.FS

const templateDir = "templates"

// Handler provides Web UI handlers.
type Handler struct {
	cfg         *config.Config
	authManager *coreauth.Manager
}

// NewHandler creates a new Web UI handler.
func NewHandler(cfg *config.Config, authManager *coreauth.Manager) *Handler {
	return &Handler{
		cfg:         cfg,
		authManager: authManager,
	}
}

// RegisterRoutes registers Web UI routes on the Gin engine.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	if h.cfg == nil || !h.cfg.WebUI.Enabled {
		return
	}

	basePath := "/ui"
	if h.cfg.WebUI.Path != "" {
		basePath = h.cfg.WebUI.Path
	}

	ui := r.Group(basePath)
	{
		ui.GET("/", h.Dashboard)
		ui.GET("/auth", h.AuthList)
		ui.GET("/config", h.Config)
		ui.GET("/logs", h.Logs)
		ui.DELETE("/auth/:id", h.AuthDelete)
	}
}

// StaticDir returns the path to static assets.
func StaticDir() string {
	return filepath.Join(util.WritablePathOrDefault(), "webui")
}

// Dashboard renders the main dashboard.
func (h *Handler) Dashboard(c *gin.Context) {
	h.renderTemplate(c, "dashboard.html", gin.H{
		"Title": "ProxyGate Dashboard",
	})
}

// AuthList renders the auth management page.
func (h *Handler) AuthList(c *gin.Context) {
	// Get auth files count
	auths := h.getAuthCount()
	
	h.renderTemplate(c, "auth.html", gin.H{
		"Title": "Auth Management",
		"AuthCount": auths,
	})
}

// Config renders the configuration page.
func (h *Handler) Config(c *gin.Context) {
	h.renderTemplate(c, "config.html", gin.H{
		"Title": "Configuration",
	})
}

// Logs renders the logs page.
func (h *Handler) Logs(c *gin.Context) {
	h.renderTemplate(c, "logs.html", gin.H{
		"Title": "Logs",
	})
}

// AuthDelete deletes an auth file.
func (h *Handler) AuthDelete(c *gin.Context) {
	// TODO: Implement deletion logic
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// getAuthCount returns the number of auth files.
func (h *Handler) getAuthCount() int {
	return util.CountAuthFiles(h.cfg.AuthDir)
}

// renderTemplate renders an HTML template.
func (h *Handler) renderTemplate(c *gin.Context, name string, data gin.H) {
	data["WebUIPath"] = "/ui"
	
	// Try embedded templates first, then fall back to file system
	if content, err := templates.ReadFile("templates/" + name); err == nil {
		c.Data(http.StatusOK, "text/html", content)
		return
	}
	
	// Fall back to file system
	templatePath := filepath.Join(StaticDir(), templateDir, name)
	c.File(templatePath)
}

