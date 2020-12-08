package qrcode

import (
	"database/sql"
	"encoding/base64"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type QrcodeHandler struct {
	store QrcodeStore
	tmpl  *template.Template
}

func NewQrcodeHandler(store QrcodeStore, tmpl *template.Template) *QrcodeHandler {
	if tmpl == nil {
		panic("provided template is nil")
	}
	return &QrcodeHandler{
		store: store,
		tmpl:  tmpl,
	}
}

func (h *QrcodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "index", map[string]string{
		"Csrf": csrf.Token(r),
	})
}

func (h *QrcodeHandler) Manage(w http.ResponseWriter, r *http.Request) {
	pw := mux.Vars(r)["password"]
	qr, err := h.store.GetQrcodeByPassword(r.Context(), pw)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	qrPng, err := qrcode.Encode(r.Host+"/link/"+strconv.Itoa(qr.Id), qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	tmplData := map[string]string{
		"QrCode": base64.StdEncoding.EncodeToString(qrPng),
		"Target": qr.Target,
		"Link":   "/link/" + strconv.Itoa(qr.Id),
		"Csrf":   csrf.Token(r),
	}

	h.tmpl.ExecuteTemplate(w, "manage", tmplData)
}

func (h *QrcodeHandler) Store(w http.ResponseWriter, r *http.Request) {
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
		log.Print(err)
		return
	}

	// redirect to management page
	http.Redirect(w, r, "/manage/"+password, http.StatusSeeOther)
}

func (h *QrcodeHandler) Update(w http.ResponseWriter, r *http.Request) {

	pw := mux.Vars(r)["password"]
	qr, err := h.store.GetQrcodeByPassword(r.Context(), pw)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Print(err)
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
		log.Print(err)
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
		log.Print(err)
		return
	}

	http.Redirect(w, r, qr.Target, http.StatusFound)
}
