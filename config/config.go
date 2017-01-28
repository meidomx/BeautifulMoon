package config

type GlobalConfig struct {
	DisplayConfig DisplayConfig
}

type DisplayConfig struct {
	DisplayResolution DisplayResolutionType
	WindowResizable   bool
	OpenGLConfig      OpenGLConfig
	FullScreen        bool
}

type OpenGLConfig struct {
	MAJOR_VERSION int
	MINOR_VERSION int
}

type DisplayResolutionType struct {
	Width     int
	Height    int
	Frequency int
}
