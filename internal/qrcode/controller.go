package qrcode

import (
	"encoding/base64"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

type QrcodeHandler struct {
	store QrcodeStore
	tmpl  *template.Template
}

func (h *QrcodeHandler) Manage(w http.ResponseWriter, r *http.Request) {
	pw := mux.Vars(r)["password"]
	qr, err := h.store.GetQrcodeByPassword(r.Context(), pw)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	qrPng, err := qrcode.Encode(qr.Target, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	managementData.QrCode = base64.StdEncoding.EncodeToString(qrPng)

	h.tmpl.ExecuteTemplate(w, "manage", nil)
}

func (h *QrcodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form provided", http.StatusBadRequest)
		return
	}

	// check if the url supplied is valid
	targetUri, err := url.Parse(r.PostFormValue("target"))
	if err != nil {
		http.Error(w, "invalid target url supplied", http.StatusBadRequest)
		return
	}

	// concat three uuids to create a secure enough string
	password := uuid.NewV4().String() + uuid.NewV4().String() + uuid.NewV4().String()

	if err := h.store.NewQrcode(r.Context(), password, targetUri.String()); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	// redirect to management page
	http.Redirect(w, r, "/manage/"+password, http.StatusSeeOther)
}

func (h *QrcodeHandler) Update(w http.ResponseWriter, r *http.Request) {

	pw := mux.Vars(r)["password"]
	qr, err := h.store.GetQrcodeByPassword(r.Context(), pw)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form provided", http.StatusBadRequest)
		return
	}

	// check if the new url supplied is valid
	targetUri, err := url.Parse(r.PostFormValue("target"))
	if err != nil {
		http.Error(w, "invalid target url supplied", http.StatusBadRequest)
		return
	}

	// update the target url
	if err := h.store.UpdateTargetUrl(r.Context(), qr.Id, targetUri.String()); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// redirect back to management page
	http.Redirect(w, r, "/manage/"+pw, http.StatusSeeOther)
}

func (h *QrcodeHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	qr, err := h.store.GetQrcode(r.Context(), id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, qr.Target, http.StatusFound)
}
