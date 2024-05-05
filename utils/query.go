package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func NewQueryUtils() QueryUtils {
	return &QueryUtilsImpl{}
}

func (q *QueryUtilsImpl) NewLookup(from, localField, foreignField, as string) bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: from},
				{Key: "localField", Value: localField},
				{Key: "foreignField", Value: foreignField},
				{Key: "as", Value: as},
			},
		},
	}
}

func (q *QueryUtilsImpl) NewSkip(val int) bson.D {
	return bson.D{{Key: "$skip", Value: val}}
}

func (q *QueryUtilsImpl) NewRawUnwind(val string) bson.D {
	return bson.D{{Key: "$unwind", Value: val}}
}

func (q *QueryUtilsImpl) NewLimit(val int) bson.D {
	return bson.D{{Key: "$limit", Value: val}}
}

func (q *QueryUtilsImpl) IsDo(key, input, userId string) bson.E {
	return bson.E{Key: key,
		Value: bson.D{
			{Key: "$reduce",
				Value: bson.D{
					{Key: "input", Value: input},
					{Key: "initialValue", Value: false},
					{Key: "in",
						Value: bson.D{
							{Key: "$cond",
								Value: bson.A{
									bson.D{
										{Key: "$eq",
											Value: bson.A{
												"$$this.userId",
												userId,
											},
										},
									},
									true,
									"$$value",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (q *QueryUtilsImpl) NewCountComment(key, field string) bson.E {
	return bson.E{Key: key,
		Value: bson.D{
			{Key: "$sum",
				Value: bson.A{
					bson.D{{Key: "$size", Value: field}},
					bson.D{{Key: "$size", Value: fmt.Sprintf("%s.reply", field)}},
				},
			},
		},
	}
}
