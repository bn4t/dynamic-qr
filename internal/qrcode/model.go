package qrcode

import "context"

type Qrcode struct {
	Id       int
	Target   string
	Password string
}

type QrcodeStore interface {
	NewQrcode(ctx context.Context, password string, target string) error
	GetQrcode(ctx context.Context, id int) (Qrcode, error)
	GetQrcodeByPassword(ctx context.Context, password string) (Qrcode, error)
	UpdateTargetUrl(ctx context.Context, id int, newTargetUrl string) error
	Close() error
}
