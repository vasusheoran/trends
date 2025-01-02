package metrics

const (
	TickerUpdateLatency    = "ticker_update_latency"
	TickerInitLatency      = "ticker_init_latency"
	TickerGetLatency       = "ticker_get_latency"
	ParseFileLatency       = "parse_file_latency"
	DeleteTickerLatency    = "db_delete_ticker_latency"
	GetUniqueTickerLatency = "db_unique_ticker_latency"
	GetTickerByNameLatency = "db_ticker_by_name_latency"
	PaginateTickerLatency  = "db_paginate_ticker_latency"
	SaveTickerLatenct      = "db_save_ticker_latency"
)
