package main

type (
	// Config for the plugin.
	Config struct {
		Region string
	}

	// Plugin values.
	Plugin struct {
		Config Config
	}
)

// Exec executes the plugin.
func (p Plugin) Exec() error {
	return nil
}
