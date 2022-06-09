#include <stdint.h>
#include <x264.h>

extern void go_nalu_process(x264_t *h, x264_nal_t *nal, void *opaque);

void nalu_process(x264_t *h, x264_nal_t *nal, void *opaque) { go_nalu_process(h, nal, opaque); }
