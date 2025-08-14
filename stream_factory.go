package log

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type steamFactory flam.Factory[Stream]

type steamFactoryArgs struct {
	dig.In

	Creators      []StreamCreator `group:"flam.log.streams.creator"`
	FactoryConfig flam.FactoryConfig
}

func newStreamFactory(
	args steamFactoryArgs,
) (steamFactory, error) {
	var creators []flam.ResourceCreator[Stream]
	for _, creator := range args.Creators {
		creators = append(creators, creator)
	}

	return flam.NewFactory(
		creators,
		PathStreams,
		args.FactoryConfig,
		nil,
	)
}
