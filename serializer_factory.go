package log

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type serializerFactory flam.Factory[Serializer]

type serializerFactoryArgs struct {
	dig.In

	Creators      []SerializerCreator `group:"flam.log.serializers.creator"`
	FactoryConfig flam.FactoryConfig
}

func newSerializerFactory(
	args serializerFactoryArgs,
) (serializerFactory, error) {
	var creators []flam.ResourceCreator[Serializer]
	for _, creator := range args.Creators {
		creators = append(creators, creator)
	}

	return flam.NewFactory(
		creators,
		PathSerializers,
		args.FactoryConfig,
		nil,
	)
}
