package handler

import "net/http"

type AdminHandl interface {
	GetUsers() http.HandlerFunc
	GetUser() http.HandlerFunc
	UpdateUser() http.HandlerFunc
	DeleteUser() http.HandlerFunc
	BlockUser() http.HandlerFunc
	UnblockUser() http.HandlerFunc
	BlockLink() http.HandlerFunc
	UnblockLink() http.HandlerFunc
	DeleteLink() http.HandlerFunc
	GetBlockedUsersCount() http.HandlerFunc
	GetBlockedLinksCount() http.HandlerFunc
	GetDeletedLinksCount() http.HandlerFunc
	GetTotalLinks() http.HandlerFunc
	GetClickedLinkStats() http.HandlerFunc
	GetAllLinksStats() http.HandlerFunc
}

type AuthHandl interface {
	SignUp() http.HandlerFunc
	SignIn() http.HandlerFunc
}

type UserHandl interface {
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	CreateLink() http.HandlerFunc
	GetLinks() http.HandlerFunc
	UpdateLink() http.HandlerFunc
	DeleteLink() http.HandlerFunc
	Redirect() http.HandlerFunc
}
