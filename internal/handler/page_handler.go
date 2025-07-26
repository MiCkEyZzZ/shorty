package handler

import (
	"html/template"
	"net/http"
	"strings"

	"shorty/internal/service"
	"shorty/pkg/jwt"
)

type TemplateData struct {
	Title           string
	IsAuthenticated bool
	Role            string
	Page            string // "index" или "stats"
}

type PageHandler struct {
	jwtService  *jwt.JWT
	linkService service.LinkServ
}

func NewPageHandler(jwtSvc *jwt.JWT, linkSvc service.LinkServ) *PageHandler {
	return &PageHandler{jwtService: jwtSvc, linkService: linkSvc}
}

func (h *PageHandler) renderLayout(w http.ResponseWriter, data TemplateData) {
	// парсим layout + header + оба контент-шаблона
	paths := []string{
		"web/templates/layout.html",
		"web/templates/header.html",
		"web/templates/index.html",
		"web/templates/stats.html",
		"web/templates/settings.html",
		"web/templates/login.html",
		"web/templates/register.html",
	}
	tmpl := template.Must(template.ParseFiles(paths...))

	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PageHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	// Если путь не ровно "/", пробуем редирект по хешу
	if r.URL.Path != "/" {
		hash := strings.TrimPrefix(r.URL.Path, "/")
		if link, err := h.linkService.GetByHash(r.Context(), hash); err == nil {
			http.Redirect(w, r, link.Url, http.StatusFound)
			return
		}
		http.NotFound(w, r)
		return
	}

	// Иначе — рендерим главную страницу
	data := h.getAuthData(r)
	data.Title = "Shorty"
	data.Page = "index"
	h.renderLayout(w, data)
}

func (h *PageHandler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	data := h.getAuthData(r)
	data.Title = "Настройки"
	data.Page = "settings"
	h.renderLayout(w, data)
}

func (h *PageHandler) StatsPage(w http.ResponseWriter, r *http.Request) {
	data := h.getAuthData(r)
	if !data.IsAuthenticated || data.Role != "admin" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	data.Title = "Статистика"
	data.Page = "stats"
	h.renderLayout(w, data)
}

// LoginPage и RegisterPage остаются без layout
func (h *PageHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := h.getAuthData(r)
	data.Title = "Вход"
	data.Page = "login"
	h.renderLayout(w, data)
}
func (h *PageHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := h.getAuthData(r)
	data.Title = "Регистрация"
	data.Page = "register"
	h.renderLayout(w, data)
}

func (h *PageHandler) getAuthData(r *http.Request) TemplateData {
	td := TemplateData{}
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return td
	}
	claims, err := h.jwtService.ParseToken(cookie.Value)
	if err != nil {
		return td
	}
	td.IsAuthenticated = true
	td.Role = string(claims.Role)
	return td
}
