package douyu

import (
	"fmt"
	"github.com/Sora233/DDBOT/lsp/concern"
	"github.com/Sora233/DDBOT/lsp/concern_type"
	"github.com/Sora233/DDBOT/lsp/mmsg"
	localutils "github.com/Sora233/DDBOT/utils"
	"github.com/Sora233/MiraiGo-Template/utils"
	"reflect"
)

var logger = utils.GetModuleLogger("douyu-concern")

const (
	Live concern_type.Type = "live"
)

type Concern struct {
	*StateManager
}

func (c *Concern) Site() string {
	return Site
}

func (c *Concern) Types() []concern_type.Type {
	return []concern_type.Type{Live}
}

func (c *Concern) ParseId(s string) (interface{}, error) {
	return ParseUid(s)
}

func (c *Concern) GetStateManager() concern.IStateManager {
	return c.StateManager
}

func (c *Concern) Stop() {
	logger.Trace("正在停止douyu concern")
	logger.Trace("正在停止douyu StateManager")
	c.StateManager.Stop()
	logger.Trace("douyu StateManager已停止")
	logger.Trace("douyu concern已停止")
}

func (c *Concern) Start() error {
	c.StateManager.UseNotifyGeneratorFunc(c.notifyGenerator())
	c.StateManager.UseFreshFunc(c.fresh())
	return c.StateManager.Start()
}

func (c *Concern) Add(ctx mmsg.IMsgCtx, groupCode int64, _id interface{}, ctype concern_type.Type) (concern.IdentityInfo, error) {
	id := _id.(int64)
	var err error
	log := logger.WithFields(localutils.GroupLogFields(groupCode)).WithField("id", id)

	err = c.StateManager.CheckGroupConcern(groupCode, id, ctype)
	if err != nil {
		return nil, err
	}
	liveInfo, err := c.FindOrLoadRoom(id)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("查询房间信息失败 %v - %v", id, err)
	}
	_, err = c.StateManager.AddGroupConcern(groupCode, id, ctype)
	if err != nil {
		return nil, err
	}
	return concern.NewIdentity(id, liveInfo.Nickname), nil
}

func (c *Concern) Remove(ctx mmsg.IMsgCtx, groupCode int64, _id interface{}, ctype concern_type.Type) (concern.IdentityInfo, error) {
	id := _id.(int64)
	identity, _ := c.Get(id)
	_, err := c.StateManager.RemoveGroupConcern(groupCode, id, ctype)
	return identity, err
}

func (c *Concern) Get(id interface{}) (concern.IdentityInfo, error) {
	liveInfo, err := c.FindOrLoadRoom(id.(int64))
	if err != nil {
		return nil, err
	}
	return concern.NewIdentity(liveInfo.GetRoomId(), liveInfo.GetNickname()), nil
}

func (c *Concern) notifyGenerator() concern.NotifyGeneratorFunc {
	return func(groupCode int64, event concern.Event) []concern.Notify {
		switch info := event.(type) {
		case *LiveInfo:
			if info.Living() {
				info.Logger().WithFields(localutils.GroupLogFields(groupCode)).Debug("living notify")
			} else {
				info.Logger().WithFields(localutils.GroupLogFields(groupCode)).Debug("noliving notify")
			}
			return []concern.Notify{NewConcernLiveNotify(groupCode, info)}
		default:
			logger.Errorf("unknown EventType %+v", event)
			return nil
		}
	}
}

func (c *Concern) fresh() concern.FreshFunc {
	return c.EmitQueueFresher(func(ctype concern_type.Type, id interface{}) ([]concern.Event, error) {
		var result []concern.Event
		roomid, ok := id.(int64)
		if !ok {
			return nil, fmt.Errorf("cast fresh id type<%v> to int64 failed", reflect.ValueOf(id).Type().String())
		}
		if ctype.ContainAll(Live) {
			oldInfo, _ := c.FindRoom(roomid, false)
			liveInfo, err := c.FindRoom(roomid, true)
			if err != nil {
				return nil, fmt.Errorf("load liveinfo failed %v", err)
			}
			if oldInfo == nil {
				liveInfo.liveStatusChanged = true
			}
			if oldInfo != nil && oldInfo.Living() != liveInfo.Living() {
				liveInfo.liveStatusChanged = true
			}
			if oldInfo != nil && oldInfo.RoomName != liveInfo.RoomName {
				liveInfo.liveTitleChanged = true
			}
			if oldInfo == nil || oldInfo.Living() != liveInfo.Living() || oldInfo.RoomName != liveInfo.RoomName {
				result = append(result, liveInfo)
			}
		}
		return result, nil
	})
}

func (c *Concern) FindRoom(id int64, load bool) (*LiveInfo, error) {
	var liveInfo *LiveInfo
	if load {
		betardResp, err := Betard(id)
		if err != nil {
			return nil, err
		}
		liveInfo = &LiveInfo{
			Nickname:   betardResp.GetRoom().GetNickname(),
			RoomId:     betardResp.GetRoom().GetRoomId(),
			RoomName:   betardResp.GetRoom().GetRoomName(),
			RoomUrl:    betardResp.GetRoom().GetRoomUrl(),
			ShowStatus: betardResp.GetRoom().GetShowStatus(),
			VideoLoop:  betardResp.GetRoom().GetVideoLoop(),
			Avatar:     betardResp.GetRoom().GetAvatar(),
		}
		_ = c.StateManager.AddLiveInfo(liveInfo)
	}
	if liveInfo != nil {
		return liveInfo, nil
	}
	return c.StateManager.GetLiveInfo(id)
}

func (c *Concern) FindOrLoadRoom(roomId int64) (*LiveInfo, error) {
	info, _ := c.FindRoom(roomId, false)
	if info == nil {
		return c.FindRoom(roomId, true)
	}
	return info, nil
}

func NewConcern(notify chan<- concern.Notify) *Concern {
	c := &Concern{
		StateManager: NewStateManager(notify),
	}
	return c
}
