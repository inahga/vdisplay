package encoder

/*
#cgo pkg-config: x264
#include <stdint.h>
#include <x264.h>

typedef void (*nalu_process_cb_t)(x264_t *, x264_nal_t *, void *);
void nalu_process(x264_t *, x264_nal_t *, void *);
*/
import "C"
import (
	"fmt"
	"image"
)

type H264 struct {
	h        *C.x264_t
	params   C.x264_param_t
	userdata uintptr
}

func NewH264() (*H264, error) {
	return nil, nil
}

func (h *H264) init(img image.Image) error {
	if err := C.x264_param_default_preset(&h.params, C.CString("superfast"), C.CString("zerolatency")); err < 0 {
		return fmt.Errorf("x264_param_default_preset(): return code %d", err)
	}

	h.params.i_bitdepth = 8          // bad assumption
	h.params.i_csp = C.X264_CSP_BGRA // bad assumption
	h.params.i_width = C.int(img.Bounds().Dx())
	h.params.i_height = C.int(img.Bounds().Dy())

	// Magic values... probably depends on the receiving decoder.
	h.params.b_vfr_input = 0
	h.params.b_repeat_headers = 1
	h.params.b_annexb = 1

	// might break something...
	h.params.b_sliced_threads = 1

	h.params.nalu_process = C.nalu_process_cb_t(C.nalu_process)

	if err := C.x264_param_apply_profile(&h.params, C.CString("high")); err < 0 {
		return fmt.Errorf("x264_param_apply_profile(): return code %d", err)
	}
	if h.h = C.x264_encoder_open(&h.params); h.h == nil {
		return fmt.Errorf("x264_encoder_open(): failed")
	}
	return nil
}

func (h *H264) Encode(img image.Image) error {
	if h.h == nil {
		if err := h.init(img); err != nil {
			return err
		}
	}

	return nil
}

func (h *H264) Close() {
	if h.h != nil {
		C.x264_encoder_close(h.h)
	}
}

//export go_nalu_process
func go_nalu_process(h *C.x264_t, nal *C.x264_nal_t, opaque interface{}) {
}
