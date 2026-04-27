package httpserver

import "testing"

func TestServerConfig_ApplyDefaults(t *testing.T) {
	t.Run("Нулевой ServerConfig после ApplyDefaults: ListenHost 127.0.0.1, ExternalHost localhost, Port 80, Timeout 60", func(t *testing.T) {
		srv := ServerConfig{}
		srv.ApplyDefaults()
		if srv.ListenHost != "127.0.0.1" {
			t.Fatalf("ListenHost=%q", srv.ListenHost)
		}
		if srv.ExternalHost != "localhost" {
			t.Fatalf("ExternalHost=%q", srv.ExternalHost)
		}
		if srv.Port != 80 {
			t.Fatalf("Port=%d, ожидали 80", srv.Port)
		}
		if srv.Timeout != 60 {
			t.Fatalf("Timeout=%d, ожидали 60", srv.Timeout)
		}
	})

	t.Run("Все поля заданы явно: ApplyDefaults ничего не меняет", func(t *testing.T) {
		srv := ServerConfig{
			ListenHost:   "0.0.0.0",
			ExternalHost: "api.example",
			Port:         8080,
			Timeout:      30,
		}
		srv.ApplyDefaults()
		if srv.ListenHost != "0.0.0.0" || srv.ExternalHost != "api.example" || srv.Port != 8080 || srv.Timeout != 30 {
			t.Fatalf("поля изменились: %+v", srv)
		}
	})
}

func TestServerConfig_Validate(t *testing.T) {
	t.Run("TLS enabled: без cert_file — ошибка", func(t *testing.T) {
		srv := ServerConfig{
			Port: 8080,
			TLS: &TLSListenConfig{
				Enabled:      true,
				KeyFile:      "/x/k",
				ClientCAFile: "/x/ca",
			},
		}
		errs := srv.Validate()
		if len(errs) != 1 {
			t.Fatalf("ожидали 1 ошибку, получили %d: %v", len(errs), errs)
		}
	})

	t.Run("Port 1, 80 и 65535: Validate без ошибок", func(t *testing.T) {
		for _, port := range []int{1, 80, 65535} {
			errs := (&ServerConfig{Port: port}).Validate()
			if len(errs) != 0 {
				t.Fatalf("port %d: %v", port, errs)
			}
		}
	})

	t.Run("Port 0, -1 и 65536: по одной ошибке Validate на каждый", func(t *testing.T) {
		for _, port := range []int{0, -1, 65536} {
			errs := (&ServerConfig{Port: port}).Validate()
			if len(errs) != 1 {
				t.Fatalf("port %d: ожидали 1 ошибку, получили %v", port, errs)
			}
		}
	})
}

func TestServerConfig_ApplyDefaults_then_Validate(t *testing.T) {
	t.Run("Пустой ServerConfig после ApplyDefaults: Validate без ошибок (Port 80 в диапазоне)", func(t *testing.T) {
		srv := ServerConfig{}
		srv.ApplyDefaults()
		errs := srv.Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})
}
