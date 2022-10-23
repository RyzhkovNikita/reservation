package storage

type Extension int8
type ImageExtension Extension

const (
	UNKNOWN Extension      = 1
	JPEG    ImageExtension = 2
	PNG     ImageExtension = 3
	PDF     Extension      = 4
)

func (e Extension) ToString() string {
	switch e {
	case Extension(JPEG):
		return "jpeg"
	case Extension(PNG):
		return "png"
	case PDF:
		return "pdf"
	default:
		return ""
	}
}

func (e ImageExtension) ToString() string {
	return Extension(e).ToString()
}
