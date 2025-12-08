module github.com/eclipse-softworks/luna-sdk/cli

go 1.21

require (
	github.com/spf13/cobra v1.10.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	al.essio.dev/pkg/shellescape v1.5.1 // indirect
	github.com/danieljoos/wincred v1.2.2 // indirect
	github.com/eclipse-softworks/luna-sdk-go v0.0.0-00010101000000-000000000000
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/zalando/go-keyring v0.2.6
	golang.org/x/sys v0.26.0 // indirect
)

replace github.com/eclipse-softworks/luna-sdk-go => ../packages/go
