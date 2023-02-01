package main

type BaseStation struct {
	num_free_channels int
}

func (bs *BaseStation) assign_channel() {
	bs.num_free_channels -= 1
}

func (bs *BaseStation) has_available_channel() bool {
	return bs.num_free_channels > 0
}

func (bs *BaseStation) free_up_channel() {
	bs.num_free_channels += 1
}