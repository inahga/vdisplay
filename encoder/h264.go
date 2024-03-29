package encoder

/*
#cgo pkg-config: x264
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <x264.h>

typedef void (*nalu_process_cb_t)(x264_t *, x264_nal_t *, void *);
void nalu_process(x264_t *, x264_nal_t *, void *);
*/
import "C"
import (
	"fmt"
	"image"
	"log"
	"unsafe"
)

type H264 struct {
	h        *C.x264_t
	params   C.x264_param_t
	picture  C.x264_picture_t
	pts      C.long
	userdata uintptr
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

	if err := C.x264_param_apply_profile(&h.params, C.CString("high444")); err < 0 {
		return fmt.Errorf("x264_param_apply_profile(): return code %d", err)
	}
	if err := C.x264_picture_alloc(&h.picture, h.params.i_csp, h.params.i_width, h.params.i_height); err < 0 {
		return fmt.Errorf("x264_picture_alloc(): return code %d", err)
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

	rgba, ok := img.(*image.RGBA)
	if !ok {
		return fmt.Errorf("underlying type of img is not *image.RGBA")
	}
	h.picture.i_pts = h.pts
	C.memcpy(unsafe.Pointer(h.picture.img.plane[0]), unsafe.Pointer(&rgba.Pix[0]), C.ulong(len(rgba.Pix)))

	var (
		nal    *C.x264_nal_t
		nnal   C.int
		outpic C.x264_picture_t
	)
	size := C.x264_encoder_encode(h.h, &nal, &nnal, &h.picture, &outpic)
	if size < 0 {
		return fmt.Errorf("x264_encoder_encode(): return code %d", size)
	}
	log.Printf("[h264] encode: size = %d, nnal = %d, out_pts = %d", size, nnal, outpic.i_pts)
	h.pts++
	return nil
}

func (h *H264) Close() {
	C.x264_picture_clean(&h.picture)
	if h.h != nil {
		C.x264_encoder_close(h.h)
	}
}

//export go_nalu_process
func go_nalu_process(h *C.x264_t, nal *C.x264_nal_t, opaque interface{}) {
	log.Printf("[h264] go_nalu_process: type = %d, size = %d", nal.i_type, nal.i_payload)
}
