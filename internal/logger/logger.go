package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Enviroment string
}

func New(cfg Config) (*zap.Logger, error) {
	var (
		l   *zap.Logger
		err error
	)

	if cfg.Enviroment == "" {
		cfg.Enviroment = os.Getenv("APP_ENV")
	}
	if cfg.Enviroment == "dev" {
		l, err = zap.NewDevelopment()
	} else {
		enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), zap.InfoLevel)
		l = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	}
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(l)
	return l, nil
}
