package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"gobase/cmd"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	cmd.Execute()
}
