package main

type BaseStation struct {
	id                         int // from 0 to 19
	num_free_channels          int
	num_free_reserved_channels int // for the FCA schemes to be tested
	using_reserved_fca_scheme  bool
}

func (bs *BaseStation) assign_non_reserved_channel() {
	bs.num_free_channels -= 1
}

func (bs *BaseStation) has_available_non_reserved_channel() bool {
	return bs.num_free_channels > 0
}

func (bs *BaseStation) free_up_non_reserved_channel() {
	bs.num_free_channels += 1
}

// For reserved channels when testing the FCA schemes
func (bs *BaseStation) assign_reserved_channel() {
	bs.num_free_reserved_channels -= 1
}

func (bs *BaseStation) has_available_reserved_channel() bool {
	return bs.num_free_reserved_channels > 0
}

func (bs *BaseStation) free_up_reserved_channel() {
	bs.num_free_reserved_channels += 1
}