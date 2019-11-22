package services

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var TopicLikeService = newTopicLikeService()

func newTopicLikeService() *topicLikeService {
	return &topicLikeService{}
}

type topicLikeService struct {
}

func (this *topicLikeService) Get(id int64) *model.TopicLike {
	return repositories.TopicLikeRepository.Get(simple.DB(), id)
}

func (this *topicLikeService) Take(where ...interface{}) *model.TopicLike {
	return repositories.TopicLikeRepository.Take(simple.DB(), where...)
}

func (this *topicLikeService) Find(cnd *simple.SqlCnd) []model.TopicLike {
	return repositories.TopicLikeRepository.Find(simple.DB(), cnd)
}

func (this *topicLikeService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.TopicLike) {
	cnd.FindOne(db, &ret)
	return
}

func (this *topicLikeService) FindPageByParams(params *simple.QueryParams) (list []model.TopicLike, paging *simple.Paging) {
	return repositories.TopicLikeRepository.FindPageByParams(simple.DB(), params)
}

func (this *topicLikeService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.TopicLike, paging *simple.Paging) {
	return repositories.TopicLikeRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *topicLikeService) Create(t *model.TopicLike) error {
	return repositories.TopicLikeRepository.Create(simple.DB(), t)
}

func (this *topicLikeService) Update(t *model.TopicLike) error {
	return repositories.TopicLikeRepository.Update(simple.DB(), t)
}

func (this *topicLikeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicLikeRepository.Updates(simple.DB(), id, columns)
}

func (this *topicLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicLikeRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *topicLikeService) Delete(id int64) {
	repositories.TopicLikeRepository.Delete(simple.DB(), id)
}

// 统计数量
func (this *topicLikeService) Count(topicId int64) int64 {
	var count int64 = 0
	simple.DB().Model(&model.TopicLike{}).Where("topic_id = ?", topicId).Count(&count)
	return count
}

func (this *topicLikeService) Like(userId int64, topicId int64) error {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return errors.New("话题不存在")
	}

	// 判断是否已经点赞了
	topicLike := repositories.TopicLikeRepository.Take(simple.DB(), "user_id = ? and topic_id = ?", userId, topicId)
	if topicLike != nil {
		return errors.New("已点赞")
	}

	return simple.Tx(simple.DB(), func(tx *gorm.DB) error {
		// 点赞
		err := repositories.TopicLikeRepository.Create(tx, &model.TopicLike{
			UserId:     userId,
			TopicId:    topicId,
			CreateTime: simple.NowTimestamp(),
		})
		if err != nil {
			return err
		}

		// 更新帖子点赞数
		return simple.DB().Exec("update t_topic set like_count = like_count + 1 where id = ?", topicId).Error
	})
}
