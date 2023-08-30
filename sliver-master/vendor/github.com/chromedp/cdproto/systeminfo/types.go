package systeminfo

// Code generated by cdproto-gen. DO NOT EDIT.

import (
	"fmt"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// GPUDevice describes a single graphics processor (GPU).
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-GPUDevice
type GPUDevice struct {
	VendorID      float64 `json:"vendorId"`           // PCI ID of the GPU vendor, if available; 0 otherwise.
	DeviceID      float64 `json:"deviceId"`           // PCI ID of the GPU device, if available; 0 otherwise.
	SubSysID      float64 `json:"subSysId,omitempty"` // Sub sys ID of the GPU, only available on Windows.
	Revision      float64 `json:"revision,omitempty"` // Revision of the GPU, only available on Windows.
	VendorString  string  `json:"vendorString"`       // String description of the GPU vendor, if the PCI ID is not available.
	DeviceString  string  `json:"deviceString"`       // String description of the GPU device, if the PCI ID is not available.
	DriverVendor  string  `json:"driverVendor"`       // String description of the GPU driver vendor.
	DriverVersion string  `json:"driverVersion"`      // String description of the GPU driver version.
}

// Size describes the width and height dimensions of an entity.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-Size
type Size struct {
	Width  int64 `json:"width"`  // Width in pixels.
	Height int64 `json:"height"` // Height in pixels.
}

// VideoDecodeAcceleratorCapability describes a supported video decoding
// profile with its associated minimum and maximum resolutions.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-VideoDecodeAcceleratorCapability
type VideoDecodeAcceleratorCapability struct {
	Profile       string `json:"profile"`       // Video codec profile that is supported, e.g. VP9 Profile 2.
	MaxResolution *Size  `json:"maxResolution"` // Maximum video dimensions in pixels supported for this |profile|.
	MinResolution *Size  `json:"minResolution"` // Minimum video dimensions in pixels supported for this |profile|.
}

// VideoEncodeAcceleratorCapability describes a supported video encoding
// profile with its associated maximum resolution and maximum framerate.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-VideoEncodeAcceleratorCapability
type VideoEncodeAcceleratorCapability struct {
	Profile                 string `json:"profile"`               // Video codec profile that is supported, e.g H264 Main.
	MaxResolution           *Size  `json:"maxResolution"`         // Maximum video dimensions in pixels supported for this |profile|.
	MaxFramerateNumerator   int64  `json:"maxFramerateNumerator"` // Maximum encoding framerate in frames per second supported for this |profile|, as fraction's numerator and denominator, e.g. 24/1 fps, 24000/1001 fps, etc.
	MaxFramerateDenominator int64  `json:"maxFramerateDenominator"`
}

// SubsamplingFormat yUV subsampling type of the pixels of a given image.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-SubsamplingFormat
type SubsamplingFormat string

// String returns the SubsamplingFormat as string value.
func (t SubsamplingFormat) String() string {
	return string(t)
}

// SubsamplingFormat values.
const (
	SubsamplingFormatYuv420 SubsamplingFormat = "yuv420"
	SubsamplingFormatYuv422 SubsamplingFormat = "yuv422"
	SubsamplingFormatYuv444 SubsamplingFormat = "yuv444"
)

// MarshalEasyJSON satisfies easyjson.Marshaler.
func (t SubsamplingFormat) MarshalEasyJSON(out *jwriter.Writer) {
	out.String(string(t))
}

// MarshalJSON satisfies json.Marshaler.
func (t SubsamplingFormat) MarshalJSON() ([]byte, error) {
	return easyjson.Marshal(t)
}

// UnmarshalEasyJSON satisfies easyjson.Unmarshaler.
func (t *SubsamplingFormat) UnmarshalEasyJSON(in *jlexer.Lexer) {
	v := in.String()
	switch SubsamplingFormat(v) {
	case SubsamplingFormatYuv420:
		*t = SubsamplingFormatYuv420
	case SubsamplingFormatYuv422:
		*t = SubsamplingFormatYuv422
	case SubsamplingFormatYuv444:
		*t = SubsamplingFormatYuv444

	default:
		in.AddError(fmt.Errorf("unknown SubsamplingFormat value: %v", v))
	}
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (t *SubsamplingFormat) UnmarshalJSON(buf []byte) error {
	return easyjson.Unmarshal(buf, t)
}

// ImageType image format of a given image.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-ImageType
type ImageType string

// String returns the ImageType as string value.
func (t ImageType) String() string {
	return string(t)
}

// ImageType values.
const (
	ImageTypeJpeg    ImageType = "jpeg"
	ImageTypeWebp    ImageType = "webp"
	ImageTypeUnknown ImageType = "unknown"
)

// MarshalEasyJSON satisfies easyjson.Marshaler.
func (t ImageType) MarshalEasyJSON(out *jwriter.Writer) {
	out.String(string(t))
}

// MarshalJSON satisfies json.Marshaler.
func (t ImageType) MarshalJSON() ([]byte, error) {
	return easyjson.Marshal(t)
}

// UnmarshalEasyJSON satisfies easyjson.Unmarshaler.
func (t *ImageType) UnmarshalEasyJSON(in *jlexer.Lexer) {
	v := in.String()
	switch ImageType(v) {
	case ImageTypeJpeg:
		*t = ImageTypeJpeg
	case ImageTypeWebp:
		*t = ImageTypeWebp
	case ImageTypeUnknown:
		*t = ImageTypeUnknown

	default:
		in.AddError(fmt.Errorf("unknown ImageType value: %v", v))
	}
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (t *ImageType) UnmarshalJSON(buf []byte) error {
	return easyjson.Unmarshal(buf, t)
}

// ImageDecodeAcceleratorCapability describes a supported image decoding
// profile with its associated minimum and maximum resolutions and subsampling.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-ImageDecodeAcceleratorCapability
type ImageDecodeAcceleratorCapability struct {
	ImageType     ImageType           `json:"imageType"`     // Image coded, e.g. Jpeg.
	MaxDimensions *Size               `json:"maxDimensions"` // Maximum supported dimensions of the image in pixels.
	MinDimensions *Size               `json:"minDimensions"` // Minimum supported dimensions of the image in pixels.
	Subsamplings  []SubsamplingFormat `json:"subsamplings"`  // Optional array of supported subsampling formats, e.g. 4:2:0, if known.
}

// GPUInfo provides information about the GPU(s) on the system.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-GPUInfo
type GPUInfo struct {
	Devices              []*GPUDevice                        `json:"devices"` // The graphics devices on the system. Element 0 is the primary GPU.
	AuxAttributes        easyjson.RawMessage                 `json:"auxAttributes,omitempty"`
	FeatureStatus        easyjson.RawMessage                 `json:"featureStatus,omitempty"`
	DriverBugWorkarounds []string                            `json:"driverBugWorkarounds"` // An optional array of GPU driver bug workarounds.
	VideoDecoding        []*VideoDecodeAcceleratorCapability `json:"videoDecoding"`        // Supported accelerated video decoding capabilities.
	VideoEncoding        []*VideoEncodeAcceleratorCapability `json:"videoEncoding"`        // Supported accelerated video encoding capabilities.
	ImageDecoding        []*ImageDecodeAcceleratorCapability `json:"imageDecoding"`        // Supported accelerated image decoding capabilities.
}

// ProcessInfo represents process info.
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/SystemInfo#type-ProcessInfo
type ProcessInfo struct {
	Type    string  `json:"type"`    // Specifies process type.
	ID      int64   `json:"id"`      // Specifies process id.
	CPUTime float64 `json:"cpuTime"` // Specifies cumulative CPU usage in seconds across all threads of the process since the process start.
}
