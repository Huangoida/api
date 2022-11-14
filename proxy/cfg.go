package proxy

import "time"

type Option struct {
	LimitCountDispatchWorker   uint64
	LimitCountCopyWorker       uint64
	LimitCountHeathCheckWorker int
	LimitCountConn             int
	LimitIntervalHeathCheck    time.Duration
	LimitDurationConnKeepalive time.Duration
	LimitDurationConnIdle      time.Duration
	LimitTimeoutWrite          time.Duration
	LimitTimeoutRead           time.Duration
	LimitBufferRead            int
	LimitBufferWrite           int
	LimitBytesBody             int
	LimitBytesCaching          uint64

	JWTCfgFile   string
	CrossCfgFile string

	EnableWebSocket              bool
	EnableJSPlugin               bool
	DisableHeaderNameNormalizing bool
}

type Cfg struct {
	Addr string

	Option *Option
}
