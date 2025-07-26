package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"shorty/internal/common"
	"shorty/internal/config"
	"shorty/internal/models"
	"shorty/internal/service"
	"shorty/pkg/jwt"
	"shorty/pkg/logger"
	"shorty/pkg/middleware"
	"shorty/pkg/req"
	"shorty/pkg/res"
)

// AdminHandlerDeps holds the dependencies required to initialize an AdminHandler.
type AdminHandlerDeps struct {
	Config      *config.Config
	UserService service.UserServ
	LinkService service.LinkServ
	StatService service.StatServ
	JWTService  *jwt.JWT
}

// AdminHandler handles admin-related routes and operations.
type AdminHandler struct {
	Config      *config.Config
	UserService service.UserServ
	LinkService service.LinkServ
	StatService service.StatServ
	JWTService  *jwt.JWT
}

// NewAdminHandler registers admin-related routes and attaches them to AdminHandler methods.
func NewAdminHandler(router *http.ServeMux, deps AdminHandlerDeps) {
	handler := &AdminHandler{
		Config:      deps.Config,
		UserService: deps.UserService,
		LinkService: deps.LinkService,
		StatService: deps.StatService,
		JWTService:  deps.JWTService,
	}

	adminMiddleware := middleware.AdminMiddleware(deps.JWTService, deps.UserService)

	// User management
	router.Handle("GET /admin/users", adminMiddleware(handler.GetUsers()))
	router.Handle("GET /admin/users/{id}", adminMiddleware(handler.GetUser()))
	router.Handle("PATCH /admin/users/{id}", adminMiddleware(handler.UpdateUser()))
	router.Handle("DELETE /admin/users/{id}", adminMiddleware(handler.DeleteUser()))
	router.Handle("PATCH /admin/users/{id}/block", adminMiddleware(handler.BlockUser()))
	router.Handle("PATCH /admin/users/{id}/unblock", adminMiddleware(handler.UnblockUser()))
	router.Handle("GET /admin/users/blocked/count", adminMiddleware(handler.GetBlockedUsersCount()))

	// Link management
	router.Handle("PATCH /admin/links/{id}/block", adminMiddleware(handler.BlockLink()))
	router.Handle("PATCH /admin/links/{id}/unblock", adminMiddleware(handler.UnblockLink()))
	router.Handle("DELETE /admin/links/{id}", adminMiddleware(handler.DeleteLink()))
	router.Handle("GET /admin/links/blocked/count", adminMiddleware(handler.GetBlockedLinksCount()))
	router.Handle("GET /admin/links/deleted/count", adminMiddleware(handler.GetDeletedLinksCount()))
	router.Handle("GET /admin/links/created/count", adminMiddleware(handler.GetTotalLinks()))

	// Statistics
	router.Handle("GET /admin/stats", adminMiddleware(handler.GetClickedLinkStats()))
	router.Handle("GET /admin/stats/links", adminMiddleware(handler.GetAllLinksStats()))
}

// GetUsers method to retrieve the list of users.
func (h *AdminHandler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		limit := h.Config.DefaultLimit
		if limit <= 0 {
			limit = 5
		}
		if lStr := r.URL.Query().Get("limit"); lStr != "" {
			if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
				limit = l
			}
		}

		page := 1
		if pStr := r.URL.Query().Get("page"); pStr != "" {
			if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
				page = p
			}
		}
		offset := (page - 1) * limit

		total, err := h.UserService.Count(r.Context())
		if err != nil {
			res.ERROR(w, common.ErrorGetUsers, http.StatusInternalServerError)
			return
		}
		users, err := h.UserService.GetAll(r.Context(), limit, offset)
		if err != nil {
			res.ERROR(w, common.ErrorGetUsers, http.StatusInternalServerError)
			return
		}

		totalPages := 1
		if limit > 0 {
			totalPages = int((total + int64(limit) - 1) / int64(limit))
		}

		var next, prev interface{}
		if page < totalPages {
			next = makePageURL(r, page+1, totalPages)
		}
		if page > 1 {
			prev = makePageURL(r, page-1, totalPages)
		}

		resp := map[string]interface{}{
			"info": map[string]interface{}{
				"count": total,
				"pages": totalPages,
				"next":  next,
				"prev":  prev,
			},
			"results": users,
		}
		res.JSON(w, resp, http.StatusOK)
	}
}

// GetUser method to retrieve a user by their ID.
func (h *AdminHandler) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("Invalid user ID", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		user, err := h.UserService.GetByID(ctx, uint(userID))
		if err != nil {
			logger.Error("Error searching for user", zap.Uint("userID", userID), zap.Error(err))
			res.ERROR(w, common.ErrUserNotFound, http.StatusNotFound)
			return
		}
		logger.Info("User found", zap.Uint("id", userID))
		res.JSON(w, user, http.StatusOK)
	}
}

// UpdateUser method to update a user by their ID.
func (h *AdminHandler) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("Invalid user ID", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		body, err := req.HandleBody[models.User](&w, r)
		if err != nil {
			logger.Error("Error processing request body", zap.Error(err))
			res.ERROR(w, common.ErrRequestBodyParse, http.StatusBadRequest)
			return
		}
		body.ID = uint(userID)
		updatedUser, err := h.UserService.Update(ctx, body)
		if err != nil {
			logger.Error("Error updating user", zap.Uint("userID", userID), zap.Error(err))
			res.ERROR(w, common.ErrUserUpdateFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("User successfully updated", zap.Uint("userID", userID))
		res.JSON(w, updatedUser, http.StatusOK)
	}
}

// DeleteUser method to delete a user by their ID.
func (h *AdminHandler) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("Invalid user ID", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		err = h.UserService.Delete(ctx, uint(userID))
		if err != nil {
			logger.Error("Error deleting user", zap.Uint("userID", userID), zap.Error(err))
			res.ERROR(w, common.ErrUserDeleteFailed, http.StatusInternalServerError)
			return
		}
		logger.Info("User successfully deleted", zap.Uint("userID", userID))
		res.JSON(w, map[string]string{"message": "User deleted"}, http.StatusOK)
	}
}

// BlockUser method to block a user by their ID.
func (h *AdminHandler) BlockUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("User ID parsing error", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}
		user, err := h.UserService.GetByID(ctx, uint(id))
		if err != nil {
			logger.Error("Error when searching for a user", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		updateUser, err := h.UserService.Block(ctx, user.ID)
		if err != nil {
			logger.Error("Error when blocking the user", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("The user has been successfully blocked", zap.Uint("id", updateUser.ID))
		res.JSON(w, updateUser, http.StatusOK)
	}
}

// UnblockUser is a handler method that unblocks a user by their ID.
func (h *AdminHandler) UnblockUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("User ID parsing error", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		user, err := h.UserService.GetByID(ctx, uint(id))
		if err != nil {
			logger.Error("Error when searching for a user", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}
		updatedUser, err := h.UserService.UnBlock(ctx, user.ID)
		if err != nil {
			logger.Error("Error when unblocking the user", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrUnBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("The user has been successfully unblocked", zap.Uint("id", updatedUser.ID))
		res.JSON(w, updatedUser, http.StatusOK)
	}
}

// BlockLink is a handler method that blocks a link by its ID.
func (h *AdminHandler) BlockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("Link ID parsing error", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		// Получаем ссылку из базы
		link, err := h.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Error when searching for a link", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		// Блокируем ссылку
		updatedLink, err := h.LinkService.Block(ctx, link.ID)
		if err != nil {
			logger.Error("Error when blocking the link", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("The link has been successfully blocked", zap.Uint("id", updatedLink.ID))
		res.JSON(w, updatedLink, http.StatusOK)
	}
}

// UnblockLink is a handler method that unblocks a link by its ID.
func (h *AdminHandler) UnblockLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("Link ID parsing error", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		// Получаем ссылку из базы
		link, err := h.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("Error when searching for a link", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrNotFound, http.StatusNotFound)
			return
		}

		// Блокируем ссылку
		updatedLink, err := h.LinkService.UnBlock(ctx, link.ID)
		if err != nil {
			logger.Error("Error when unblocking the link", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrUnBlockFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("The link has been successfully unblocked", zap.Uint("id", updatedLink.ID))
		res.JSON(w, updatedLink, http.StatusOK)
	}
}

// DeleteLink handles deletion of a link by its ID.
func (h *AdminHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := h.parseIDFromPath(r)
		if err != nil {
			logger.Error("Invalid ID for deleting a link", zap.Error(err))
			res.ERROR(w, common.ErrInvalidID, http.StatusBadRequest)
			return
		}

		_, err = h.LinkService.FindByID(ctx, uint(id))
		if err != nil {
			logger.Error("The link could not be found for deletion", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkNotFound, http.StatusNotFound)
			return
		}

		err = h.LinkService.Delete(ctx, uint(id))
		if err != nil {
			logger.Error("Error when deleting a link", zap.Uint("id", uint(id)), zap.Error(err))
			res.ERROR(w, common.ErrLinkDeleteFailed, http.StatusInternalServerError)
			return
		}

		logger.Info("The link was successfully deleted", zap.Uint("id", uint(id)))
		res.JSON(w, map[string]string{"message": "link deleted"}, http.StatusOK)
	}
}

// GetBlockedUsersCount returns the total number of blocked users.
func (h *AdminHandler) GetBlockedUsersCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		count, err := h.UserService.GetBlockedUsersCount(ctx)
		if err != nil {
			logger.Error("Error when getting the number of blocked users", zap.Error(err))
			res.ERROR(w, common.ErrInternal, http.StatusInternalServerError)
			return
		}
		res.JSON(w, map[string]int64{"blocked_users": count}, http.StatusOK)
	}
}

// GetBlockedLinksCount returns the total number of blocked links.
func (h *AdminHandler) GetBlockedLinksCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		count, err := h.LinkService.GetBlockedLinksCount(ctx)
		if err != nil {
			logger.Error("Error when getting the number of blocked links", zap.Error(err))
			res.ERROR(w, common.ErrInternal, http.StatusInternalServerError)
			return
		}
		res.JSON(w, map[string]int64{"blocked_links": count}, http.StatusOK)
	}
}

// GetDeletedLinksCount returns the number of deleted links.
func (h *AdminHandler) GetDeletedLinksCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		count, err := h.LinkService.GetDeletedLinksCount(ctx)
		if err != nil {
			logger.Error("Error when receiving the number of deleted links", zap.Error(err))
			res.ERROR(w, common.ErrInternal, http.StatusInternalServerError)
			return
		}
		res.JSON(w, map[string]int64{"deleted_links": count}, http.StatusOK)
	}
}

// GetTotalLinks returns the total number of created links.
func (h *AdminHandler) GetTotalLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		count, err := h.LinkService.GetTotalLinks(ctx)
		if err != nil {
			logger.Error("Error when getting the number of links created", zap.Error(err))
			res.ERROR(w, common.ErrInternal, http.StatusInternalServerError)
			return
		}
		res.JSON(w, map[string]int64{"created_links": count}, http.StatusOK)
	}
}

// GetClickedLinkStats returns statistics for link clicks, grouped by the specified interval (e.g., day, month).
func (h *AdminHandler) GetClickedLinkStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		from, to, by, err := h.parseStatParams(r)
		if err != nil {
			logger.Error("Error parsing stat parameters", zap.Error(err))
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}

		logger.Info("Getting statistics", zap.String("by", by), zap.Time("from", from), zap.Time("to", to))
		stats := h.StatService.GetClickedLinkStats(ctx, by, from, to)
		logger.Info("Statistics received successfully", zap.Int("record_count", len(stats)))
		res.JSON(w, stats, http.StatusOK)
	}
}

// GetAllLinksStats returns full statistics for links between two dates.
func (h *AdminHandler) GetAllLinksStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fromStr := r.URL.Query().Get("from")
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}

		toStr := r.URL.Query().Get("to")
		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			res.ERROR(w, common.ErrInvalidParam, http.StatusBadRequest)
			return
		}
		stats := h.StatService.GetAllLinksStats(ctx, from, to)
		res.JSON(w, stats, http.StatusOK)
	}
}

// parseIDFromPath parses the "id" path parameter from the request and returns it as uint.
func (h *AdminHandler) parseIDFromPath(r *http.Request) (uint, error) {
	id := r.PathValue("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	return uint(userID), nil
}

// parseStatParams parses query parameters "from", "to", and "by" used for statistics requests.
// Returns time range and grouping type (e.g., "day", "month").
func (h *AdminHandler) parseStatParams(r *http.Request) (from, to time.Time, by string, err error) {
	fromStr := r.URL.Query().Get("from")
	if fromStr == "" {
		err = fmt.Errorf("'from' parameter is required")
		return
	}

	toStr := r.URL.Query().Get("to")
	if toStr == "" {
		err = fmt.Errorf("'to' parameter is required")
		return
	}

	by = r.URL.Query().Get("by")
	if by == "" {
		by = common.GroupByDay
	}

	from, err = time.Parse("2006-01-02", fromStr)
	if err != nil {
		err = fmt.Errorf("invalid 'from' date format: %w", err)
		return
	}

	to, err = time.Parse("2006-01-02", toStr)
	if err != nil {
		err = fmt.Errorf("invalid 'to' date format: %w", err)
		return
	}

	if by != common.GroupByDay && by != common.GroupByMonth {
		err = fmt.Errorf("invalid 'by' parameter, must be 'day' or 'month'")
		return
	}

	if from.After(to) {
		err = fmt.Errorf("'from' date must be before or equal to 'to' date")
		return
	}

	return from, to, by, nil
}

// getScheme attempts to determine the original request scheme (http or https),
// taking into account reverse proxies.
func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	return "http"
}

// makePageURL constructs a paginated URL for the given page and totalPages,
// preserving other query parameters and request path.
func makePageURL(r *http.Request, page, totalPages int) interface{} {
	if page < 1 || page > totalPages {
		return nil
	}
	u := *r.URL
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	return fmt.Sprintf("%s://%s%s", scheme, r.Host, u.RequestURI())
}
