package generated

import (
	bookmarkProto "github.com/forum-gamers/nine-tails-fox/generated/bookmark"
	commentProto "github.com/forum-gamers/nine-tails-fox/generated/comment"
	postProto "github.com/forum-gamers/nine-tails-fox/generated/post"
	"github.com/forum-gamers/nine-tails-fox/pkg/comment"
	"github.com/forum-gamers/nine-tails-fox/pkg/post"
)

func ParsePostRespToProto(datas []post.PostResponse) (result []*postProto.PostResponse) {
	for _, data := range datas {
		medias := make([]*postProto.Media, 0)

		if data.Media != nil || len(data.Media) > 0 {
			for _, media := range data.Media {
				medias = append(medias, &postProto.Media{
					Url:  media.Url,
					Type: media.Type,
					Id:   media.Id,
				})
			}
		}
		result = append(result, &postProto.PostResponse{
			XId:          data.Id.Hex(),
			UserId:       data.UserId,
			Text:         data.Text,
			Media:        medias,
			AllowComment: data.AllowComment,
			CreatedAt:    data.CreatedAt.String(),
			UpdatedAt:    data.UpdatedAt.String(),
			CountLike:    int64(data.CountLike),
			CountShare:   int64(data.CountShare),
			IsLiked:      data.IsLiked,
			IsShared:     data.IsShared,
			Tags:         data.Tags,
			Privacy:      data.Privacy,
			TotalData:    int64(data.TotalData),
		})
	}
	return
}

func ParseBookmarkPostRespToProto(datas []post.PostResponse) (result []*bookmarkProto.PostResponse) {
	for _, data := range datas {
		medias := make([]*bookmarkProto.Media, 0)

		if data.Media != nil || len(data.Media) > 0 {
			for _, media := range data.Media {
				medias = append(medias, &bookmarkProto.Media{
					Url:  media.Url,
					Type: media.Type,
					Id:   media.Id,
				})
			}
		}
		result = append(result, &bookmarkProto.PostResponse{
			XId:          data.Id.Hex(),
			UserId:       data.UserId,
			Text:         data.Text,
			Media:        medias,
			AllowComment: data.AllowComment,
			CreatedAt:    data.CreatedAt.String(),
			UpdatedAt:    data.UpdatedAt.String(),
			CountLike:    int64(data.CountLike),
			CountShare:   int64(data.CountShare),
			IsLiked:      data.IsLiked,
			IsShared:     data.IsShared,
			Tags:         data.Tags,
			Privacy:      data.Privacy,
			TotalData:    int64(data.TotalData),
		})
	}
	return
}

func ParseTagsRespToProto(datas []post.TopTags) (result []*postProto.TopTag) {
	for _, data := range datas {
		postIds := make([]string, 0)
		if len(data.Posts) > 0 {
			for _, postId := range data.Posts {
				postIds = append(postIds, postId.Hex())
			}
		}
		result = append(result, &postProto.TopTag{
			XId:   data.Id,
			Posts: postIds,
			Count: int64(data.Count),
		})
	}
	return
}

func ParseCommentRespToProto(datas []comment.CommentResponse) (result []*commentProto.CommentResp) {
	for _, data := range datas {
		replies := make([]*commentProto.Reply, 0)
		if len(data.Reply) > 0 {
			for _, reply := range data.Reply {
				replies = append(replies, &commentProto.Reply{
					XId:       reply.Id.Hex(),
					UserId:    reply.UserId,
					Text:      reply.Text,
					CreatedAt: reply.CreatedAt.String(),
					UpdatedAt: reply.UpdatedAt.String(),
				})
			}
		}
		result = append(result, &commentProto.CommentResp{
			XId:       data.Id.Hex(),
			UserId:    data.UserId,
			Text:      data.Text,
			PostId:    data.PostId.Hex(),
			CreatedAt: data.CreatedAt.String(),
			UpdatedAt: data.UpdatedAt.String(),
			Reply:     replies,
			TotalData: int64(data.TotalData),
		})
	}
	return
}
