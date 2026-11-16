package server

import (
	"errors"
)

type Options struct {
	Address             string           `hcl:"address,optional" flagName:"address" flagSName:"a" flagDescribe:"IP address to listen" default:"localhost"`
	Port                string           `hcl:"port,optional" flagName:"port" flagSName:"p" flagDescribe:"Port number to liten" default:"8080"`
	PermitWrite         bool             `hcl:"permit_write,optional" flagName:"permit-write" flagSName:"w" flagDescribe:"Permit clients to write to the TTY (BE CAREFUL)" default:"false"`
	EnableBasicAuth     bool             `hcl:"enable_basic_auth,optional" default:"false"`
	Credential          string           `hcl:"credential,optional" flagName:"credential" flagSName:"c" flagDescribe:"Credential for Basic Authentication (ex: user:pass, default disabled)" default:""`
	EnableRandomUrl     bool             `hcl:"enable_random_url,optional" flagName:"random-url" flagSName:"r" flagDescribe:"Add a random string to the URL" default:"false"`
	RandomUrlLength     int              `hcl:"random_url_length,optional" flagName:"random-url-length" flagDescribe:"Random URL length" default:"8"`
	EnableTLS           bool             `hcl:"enable_tls,optional" flagName:"tls" flagSName:"t" flagDescribe:"Enable TLS/SSL" default:"false"`
	TLSCrtFile          string           `hcl:"tls_crt_file,optional" flagName:"tls-crt" flagDescribe:"TLS/SSL certificate file path" default:"~/.gotty.crt"`
	TLSKeyFile          string           `hcl:"tls_key_file,optional" flagName:"tls-key" flagDescribe:"TLS/SSL key file path" default:"~/.gotty.key"`
	EnableTLSClientAuth bool             `hcl:"enable_tls_client_auth,optional" default:"false"`
	TLSCACrtFile        string           `hcl:"tls_ca_crt_file,optional" flagName:"tls-ca-crt" flagDescribe:"TLS/SSL CA certificate file for client certifications" default:"~/.gotty.ca.crt"`
	IndexFile           string           `hcl:"index_file,optional" flagName:"index" flagDescribe:"Custom index.html file" default:""`
	TitleFormat         string           `hcl:"title_format,optional" flagName:"title-format" flagSName:"" flagDescribe:"Title format of browser window" default:"{{ .command }}@{{ .hostname }}"`
	EnableReconnect     bool             `hcl:"enable_reconnect,optional" flagName:"reconnect" flagDescribe:"Enable reconnection" default:"false"`
	ReconnectTime       int              `hcl:"reconnect_time,optional" flagName:"reconnect-time" flagDescribe:"Time to reconnect" default:"10"`
	MaxConnection       int              `hcl:"max_connection,optional" flagName:"max-connection" flagDescribe:"Maximum connection to gotty" default:"0"`
	Once                bool             `hcl:"once,optional" flagName:"once" flagDescribe:"Accept only one client and exit on disconnection" default:"false"`
	Timeout             int              `hcl:"timeout,optional" flagName:"timeout" flagDescribe:"Timeout seconds for waiting a client(0 to disable)" default:"0"`
	PermitArguments     bool             `hcl:"permit_arguments,optional" flagName:"permit-arguments" flagDescribe:"Permit clients to send command line arguments in URL (e.g. http://example.com:8080/?arg=AAA&arg=BBB)" default:"true"`
	Preferences         *HtermPrefernces `hcl:"preferences,block"`
	Width               int              `hcl:"width,optional" flagName:"width" flagDescribe:"Static width of the screen, 0(default) means dynamically resize" default:"0"`
	Height              int              `hcl:"height,optional" flagName:"height" flagDescribe:"Static height of the screen, 0(default) means dynamically resize" default:"0"`
	WSOrigin            string           `hcl:"ws_origin,optional" flagName:"ws-origin" flagDescribe:"A regular expression that matches origin URLs to be accepted by WebSocket. No cross origin requests are acceptable by default" default:""`
	Term                string           `hcl:"term,optional" flagName:"term" flagDescribe:"Terminal name to use on the browser, one of xterm or hterm." default:"xterm"`

	TitleVariables map[string]any
}

func (options *Options) Validate() error {
	if options.EnableTLSClientAuth && !options.EnableTLS {
		return errors.New("TLS client authentication is enabled, but TLS is not enabled")
	}
	return nil
}

// TODO(ghthor): add defaults to all these options that match the example .gotty file
type HtermPrefernces struct {
	AltGrMode                     *string                      `hcl:"alt_gr_mode,optional" json:"alt-gr-mode,omitempty"`
	AltBackspaceIsMetaBackspace   bool                         `hcl:"alt_backspace_is_meta_backspace,optional" json:"alt-backspace-is-meta-backspace,omitempty"`
	AltIsMeta                     bool                         `hcl:"alt_is_meta,optional" json:"alt-is-meta,omitempty"`
	AltSendsWhat                  string                       `hcl:"alt_sends_what,optional" json:"alt-sends-what,omitempty"`
	AudibleBellSound              string                       `hcl:"audible_bell_sound,optional" json:"audible-bell-sound,omitempty"`
	DesktopNotificationBell       bool                         `hcl:"desktop_notification_bell,optional" json:"desktop-notification-bell,omitempty"`
	BackgroundColor               string                       `hcl:"background_color,optional" json:"background-color,omitempty"`
	BackgroundImage               string                       `hcl:"background_image,optional" json:"background-image,omitempty"`
	BackgroundSize                string                       `hcl:"background_size,optional" json:"background-size,omitempty"`
	BackgroundPosition            string                       `hcl:"background_position,optional" json:"background-position,omitempty"`
	BackspaceSendsBackspace       bool                         `hcl:"backspace_sends_backspace,optional" json:"backspace-sends-backspace,omitempty"`
	CharacterMapOverrides         map[string]map[string]string `hcl:"character_map_overrides,optional" json:"character-map-overrides,omitempty"`
	CloseOnExit                   bool                         `hcl:"close_on_exit,optional" json:"close-on-exit,omitempty"`
	CursorBlink                   bool                         `hcl:"cursor_blink,optional" json:"cursor-blink,omitempty"`
	CursorBlinkCycle              [2]int                       `hcl:"cursor_blink_cycle,optional" json:"cursor-blink-cycle,omitempty"`
	CursorColor                   string                       `hcl:"cursor_color,optional" json:"cursor-color,omitempty"`
	ColorPaletteOverrides         []*string                    `hcl:"color_palette_overrides,optional" json:"color-palette-overrides,omitempty"`
	CopyOnSelect                  bool                         `hcl:"copy_on_select,optional" json:"copy-on-select,omitempty"`
	UseDefaultWindowCopy          bool                         `hcl:"use_default_window_copy,optional" json:"use-default-window-copy,omitempty"`
	ClearSelectionAfterCopy       bool                         `hcl:"clear_selection_after_copy,optional" json:"clear-selection-after-copy,omitempty"`
	CtrlPlusMinusZeroZoom         bool                         `hcl:"ctrl_plus_minus_zero_zoom,optional" json:"ctrl-plus-minus-zero-zoom,omitempty"`
	CtrlCCopy                     bool                         `hcl:"ctrl_c_copy,optional" json:"ctrl-c-copy,omitempty"`
	CtrlVPaste                    bool                         `hcl:"ctrl_v_paste,optional" json:"ctrl-v-paste,omitempty"`
	EastAsianAmbiguousAsTwoColumn bool                         `hcl:"east_asian_ambiguous_as_two_column,optional" json:"east-asian-ambiguous-as-two-column,omitempty"`
	Enable8BitControl             *bool                        `hcl:"enable_8_bit_control,optional" json:"enable-8-bit-control,omitempty"`
	EnableBold                    *bool                        `hcl:"enable_bold,optional" json:"enable-bold,omitempty"`
	EnableBoldAsBright            bool                         `hcl:"enable_bold_as_bright,optional" json:"enable-bold-as-bright,omitempty"`
	EnableClipboardNotice         bool                         `hcl:"enable_clipboard_notice,optional" json:"enable-clipboard-notice,omitempty"`
	EnableClipboardWrite          bool                         `hcl:"enable_clipboard_write,optional" json:"enable-clipboard-write,omitempty"`
	EnableDec12                   bool                         `hcl:"enable_dec12,optional" json:"enable-dec12,omitempty"`
	EnableWebGL                   bool                         `hcl:"enable_webgl,optional" json:"enable-webgl,omitempty"`
	Environment                   map[string]string            `hcl:"environment,optional" json:"environment,omitempty"`
	FontFamily                    string                       `hcl:"font_family,optional" json:"font-family,omitempty"`
	FontSize                      int                          `hcl:"font_size,optional" json:"font-size,omitempty"`
	FontSmoothing                 string                       `hcl:"font_smoothing,optional" json:"font-smoothing,omitempty"`
	ForegroundColor               string                       `hcl:"foreground_color,optional" json:"foreground-color,omitempty"`
	HomeKeysScroll                bool                         `hcl:"home_keys_scroll,optional" json:"home-keys-scroll,omitempty"`
	Keybindings                   map[string]string            `hcl:"keybindings,optional" json:"keybindings,omitempty"`
	MaxStringSequence             int                          `hcl:"max_string_sequence,optional" json:"max-string-sequence,omitempty"`
	MediaKeysAreFkeys             bool                         `hcl:"media_keys_are_fkeys,optional" json:"media-keys-are-fkeys,omitempty"`
	MetaSendsEscape               bool                         `hcl:"meta_sends_escape,optional" json:"meta-sends-escape,omitempty"`
	MousePasteButton              *int                         `hcl:"mouse_paste_button,optional" json:"mouse-paste-button,omitempty"`
	PageKeysScroll                bool                         `hcl:"page_keys_scroll,optional" json:"page-keys-scroll,omitempty"`
	PassAltNumber                 *bool                        `hcl:"pass_alt_number,optional" json:"pass-alt-number,omitempty"`
	PassCtrlNumber                *bool                        `hcl:"pass_ctrl_number,optional" json:"pass-ctrl-number,omitempty"`
	PassMetaNumber                *bool                        `hcl:"pass_meta_number,optional" json:"pass-meta-number,omitempty"`
	PassMetaV                     bool                         `hcl:"pass_meta_v,optional" json:"pass-meta-v,omitempty"`
	ReceiveEncoding               string                       `hcl:"receive_encoding,optional" json:"receive-encoding,omitempty"`
	ScrollOnKeystroke             bool                         `hcl:"scroll_on_keystroke,optional" json:"scroll-on-keystroke,omitempty"`
	ScrollOnOutput                bool                         `hcl:"scroll_on_output,optional" json:"scroll-on-output,omitempty"`
	ScrollbarVisible              bool                         `hcl:"scrollbar_visible,optional" json:"scrollbar-visible,omitempty"`
	ScrollWheelMoveMultiplier     int                          `hcl:"scroll_wheel_move_multiplier,optional" json:"scroll-wheel-move-multiplier,omitempty"`
	SendEncoding                  string                       `hcl:"send_encoding,optional" json:"send-encoding,omitempty"`
	ShiftInsertPaste              bool                         `hcl:"shift_insert_paste,optional" json:"shift-insert-paste,omitempty"`
	UserCss                       string                       `hcl:"user_css,optional" json:"user-css,omitempty"`
}
