package util

import (
	"bytes"
	"fmt"
	log "github.com/inconshreveable/log15"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05.999"
	moduleKey  = "module"
	sessionKey = "session"
)

var ip net.IP

func formatLogfmtValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	if t, ok := value.(time.Time); ok {
		// Performance optimization: No need for escaping since the provided
		// timeFormat doesn't have any escape characters, and escaping is
		// expensive.
		return t.Format(timeFormat)
	}
	value = formatShared(value)
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 3, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', 3, 64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case string:
		return (v)
	default:
		return (fmt.Sprintf("%+v", value))
	}
}

func formatShared(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

func GetKey(key string, ctx []interface{}) string {
	for i := 0; i < len(ctx); i += 2 {
		k, ok := ctx[i].(string)
		if !ok {
			continue
		}
		if k == key {
			return formatLogfmtValue(ctx[i+1])
		}
	}
	return ""
}

func logfmt(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		if i != 0 {
			buf.WriteByte(' ')
		}

		k, ok := ctx[i].(string)
		v := formatLogfmtValue(ctx[i+1])
		if !ok {
			k, v = "log", formatLogfmtValue(k)
		}

		if k == moduleKey || k == sessionKey {
			continue
		} else {
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
}

func CustomFormat() log.Format {
	return log.FormatFunc(func(r *log.Record) []byte {
		b := &bytes.Buffer{}
		lvl := strings.ToUpper(r.Lvl.String())
		fmt.Fprintf(b, "[%s][%s][%s]", r.Time.Format(timeFormat), lvl, GetIP())
		module := GetKey(moduleKey, r.Ctx)

		if module != "" {
			fmt.Fprintf(b, "[%s]", module)
		}

		session := GetKey(sessionKey, r.Ctx)

		if session != "" {
			fmt.Fprintf(b, "[%s]", session)
		}

		logfmt(b, r.Ctx)
		fmt.Fprintf(b, " %s\n", r.Msg)
		return b.Bytes()
	})
}

func GetIP() net.IP {
	if ip != nil {
		return ip
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IPv4zero
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
	return net.IPv4zero
}
