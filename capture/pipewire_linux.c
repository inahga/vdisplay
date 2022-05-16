#include <spa/debug/types.h>
#include <spa/param/video/format-utils.h>
#include <spa/param/video/type-info.h>

#include <pipewire/pipewire.h>

// weird ccls bug
#ifndef NULL
#define NULL 0
#endif

struct pipewire_data {
	struct pw_context *context;
	struct pw_core *core;
	struct pw_stream *stream;
	struct pw_main_loop *loop;
	struct spa_hook core_listener;
	struct spa_hook stream_listener;
	struct spa_video_info format;

	uint32_t fd;
	uint32_t node_id;
};

extern void pipewire_receive_buffer(uint32_t, struct spa_video_info *, struct pw_buffer *);

static void pipewire_on_process(void *userdata)
{
	struct pipewire_data *data = userdata;
	struct pw_buffer *b;
	struct spa_buffer *buf;

	if ((b = pw_stream_dequeue_buffer(data->stream)) == NULL) {
		fprintf(stderr, "[pipewire] cgo: out of buffers: %m\n");
		return;
	}

	buf = b->buffer;
	if (buf->datas[0].data == NULL) {
		// why would there be an empty buffer?
		fprintf(stderr, "[pipewire] cgo: skipping empty buffer\n");
		return;
	}

	fprintf(stderr, "[pipewire] cgo: got a frame of size %d\n", buf->datas[0].chunk->size);
	pipewire_receive_buffer(data->node_id, &data->format, b);

	pw_stream_queue_buffer(data->stream, b);
}

static void pipewire_on_param_changed(void *userdata, uint32_t id, const struct spa_pod *param)
{
	struct pipewire_data *data = userdata;

	// what does this even do?
	if (param == NULL || id != SPA_PARAM_Format)
		return;
	if (spa_format_parse(param, &data->format.media_type, &data->format.media_subtype) < 0)
		return;
	if (data->format.media_type != SPA_MEDIA_TYPE_video ||
	    data->format.media_subtype != SPA_MEDIA_SUBTYPE_raw)
		return;
	if (spa_format_video_raw_parse(param, &data->format.info.raw) < 0)
		return;

	fprintf(stderr, "[pipewire] cgo: got video format:\n");
	fprintf(stderr, "  format: %d (%s)\n", data->format.info.raw.format,
		spa_debug_type_find_name(spa_type_video_format, data->format.info.raw.format));
	fprintf(stderr, "  size: %dx%d\n", data->format.info.raw.size.width,
		data->format.info.raw.size.height);
	fprintf(stderr, "  framerate: %d/%d\n", data->format.info.raw.framerate.num,
		data->format.info.raw.framerate.denom);
}

static const struct pw_stream_events pipewire_stream_events = {
    PW_VERSION_STREAM_EVENTS,
    .param_changed = pipewire_on_param_changed,
    .process = pipewire_on_process,
};

static const struct pw_core_events pipewire_core_events = {
    PW_VERSION_CORE_EVENTS,
};

int pipewire_run_loop(uint32_t fd, uint32_t node_id, uint32_t framerate)
{
	struct pipewire_data *data = calloc(1, sizeof(struct pipewire_data));
	const struct spa_pod *params[1];
	uint8_t params_buffer[1024];
	struct spa_pod_builder pod_builder;

	data->fd = fd;
	data->node_id = node_id;

	pw_init(NULL, NULL);

	data->loop = pw_main_loop_new(NULL);
	data->context = pw_context_new(pw_main_loop_get_loop(data->loop), NULL, 0);

	data->core = pw_context_connect_fd(data->context, data->fd, NULL, 0);
	if (!data->core) {
		return -2;
	}
	fprintf(stderr, "[pipewire] cgo: connected to fd\n");

	pw_core_add_listener(data->core, &data->core_listener, &pipewire_core_events, data);

	data->stream =
	    pw_stream_new(data->core, "vdisplay pipewire stream",
			  pw_properties_new(PW_KEY_MEDIA_TYPE, "Video", PW_KEY_MEDIA_CATEGORY,
					    "Capture", PW_KEY_MEDIA_ROLE, "Screen", NULL));
	pw_stream_add_listener(data->stream, &data->stream_listener, &pipewire_stream_events, data);
	fprintf(stderr, "[pipewire] cgo: created stream %p\n", data->stream);

	pod_builder = SPA_POD_BUILDER_INIT(params_buffer, sizeof(params_buffer));
	params[0] = spa_pod_builder_add_object(
	    &pod_builder,

	    SPA_TYPE_OBJECT_Format, SPA_PARAM_EnumFormat,

	    SPA_FORMAT_mediaType, SPA_POD_Id(SPA_MEDIA_TYPE_video),

	    SPA_FORMAT_mediaSubtype, SPA_POD_Id(SPA_MEDIA_SUBTYPE_raw),

	    SPA_FORMAT_VIDEO_format,
	    SPA_POD_CHOICE_ENUM_Id(7, SPA_VIDEO_FORMAT_RGB, SPA_VIDEO_FORMAT_RGB,
				   SPA_VIDEO_FORMAT_RGBA, SPA_VIDEO_FORMAT_RGBx,
				   SPA_VIDEO_FORMAT_BGRx, SPA_VIDEO_FORMAT_YUY2,
				   SPA_VIDEO_FORMAT_I420),

	    SPA_FORMAT_VIDEO_size,
	    SPA_POD_CHOICE_RANGE_Rectangle(&SPA_RECTANGLE(320, 240), &SPA_RECTANGLE(1, 1),
					   &SPA_RECTANGLE(4096, 4096)),

	    SPA_FORMAT_VIDEO_framerate,
	    SPA_POD_CHOICE_RANGE_Fraction(&SPA_FRACTION(framerate, 1), &SPA_FRACTION(0, 1),
					  &SPA_FRACTION(framerate, 1)));

	if (pw_stream_connect(data->stream, PW_DIRECTION_INPUT, data->node_id,
			      PW_STREAM_FLAG_AUTOCONNECT | PW_STREAM_FLAG_MAP_BUFFERS, params,
			      1) < 0) {
		return -3;
	};
	fprintf(stderr, "[pipewire] cgo: connected to stream\n");

	fprintf(stderr, "[pipewire] cgo: starting pipewire thread loop\n");
	pw_main_loop_run(data->loop);

	// todo: cleanup
	return 0;
}
