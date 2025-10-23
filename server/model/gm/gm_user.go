package gm

import (
	"time"

	"gorm.io/gorm"
)

// GMUser 游戏用户模型
type GMUser struct {
	ID          uint           `json:"id" gorm:"primarykey"`                   // 主键ID
	GameUserId  string         `json:"gameUserId" gorm:"uniqueIndex;not null"` // 游戏用户ID
	UserName    string         `json:"userName" gorm:"uniqueIndex;not null"`   // 用户名
	NickName    string         `json:"nickName" gorm:"not null"`               // 昵称
	Password    string         `json:"-" gorm:"not null"`                      // 密码
	Phone       string         `json:"phone"`                                  // 手机号
	Email       string         `json:"email"`                                  // 邮箱
	HeaderImg   string         `json:"headerImg"`                              // 头像
	Enable      int            `json:"enable" gorm:"default:1"`                // 启用状态 1:启用 2:禁用
	AuthorityId string         `json:"authorityId"`                            // 权限ID
	CreatedAt   time.Time      `json:"createdAt"`                              // 创建时间
	UpdatedAt   time.Time      `json:"updatedAt"`                              // 更新时间
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`                         // 删除时间
}

// TableName 表名
func (GMUser) TableName() string {
	return "player_info"
}

type BanInfo struct {
	Status  bool  `json:"status" bson:"status"`
	BanTime int32 `json:"ban_time" bson:"ban_time"`
}

type MainTaskInfo struct {
	TaskId   int32 `json:"task_id" bson:"task_id"`
	Progress int32 `json:"progress" bson:"progress"`
}

type TaskProgressInfo struct {
	ExpireTime int32 `json:"expire_time" bson:"expire_time"`
	Value      int32 `json:"value" bson:"value"`
}

type EquipInfo struct {
	EquipId int32            `json:"equip_id" bson:"equip_id,omitempty"`
	Lv      int32            `json:"lv" bson:"lv,omitempty"`
	Attrs   map[string]int64 `json:"attrs" bson:"attrs,omitempty"`
	InsId   string           `json:"ins_id" bson:"ins_id,omitempty"`
}

type EquipsInfo struct {
	SlotInfo map[int32]EquipInfo  `json:"slot_info" bson:"slot_info"`
	EquipMap map[string]EquipInfo `json:"equip_map" bson:"equip_map"`
	Attrs    map[string]int64     `json:"attrs" bson:"attrs"`
}

type AuraInfo struct {
	PutOnSkinId  int32            `json:"put_on_skin_id" bson:"put_on_skin_id"`
	Lv           int32            `json:"lv" bson:"lv"`
	Exp          int32            `json:"exp" bson:"exp"`
	Rank         int32            `json:"rank" bson:"rank"`
	AuraSkinInfo map[string]int32 `json:"aura_skin_info" bson:"aura_skin_info"`
	Attrs        map[string]int64 `json:"attrs" bson:"attrs"`
}

type BlessInfo struct {
	Lv             int32 `json:"lv" bson:"lv"`
	UpLvEndTime    int32 `json:"up_lv_end_time" bson:"up_lv_end_time"`
	LastUseSubTime int32 `json:"last_use_sub_time" bson:"last_use_sub_time"`
}

// AttendantInfo 伙伴信息
type AttendantInfo struct {
	DeployList   []string             `json:"deploy_list" bson:"deploy_list"`
	Attrs        map[string]int64     `json:"attrs" bson:"attrs"`
	RecruitInfo  AttendantRecruitInfo `json:"recruit_info" bson:"recruit_info"`
	Book         AttendantBook        `json:"book" bson:"book"`
	PaidSkinList []int32              `json:"paid_skin_list" bson:"paid_skin_list"`
}

type AttendantRecruitInfo struct {
	RecruitList   []AttendantRecruitItem `json:"recruit_list" bson:"recruit_list"`
	PrivilegeLeft int32                  `json:"privilege_left" bson:"privilege_left"`
	GuaranteeId   int32                  `json:"guarantee_id" bson:"guarantee_id"`
	GuaranteeLeft int32                  `json:"guarantee_left" bson:"guarantee_left"`
}

type AttendantRecruitItem struct {
	Id          int32 `json:"id"  bson:"id"`
	IsRecruited bool  `json:"is_recruited" bson:"is_recruited"`
}

type AttendantBook struct {
	AttendantIds []int32 `json:"attendant_ids" bson:"attendant_ids,omitempty"`
	RewardIds    []int32 `json:"reward_ids" bson:"reward_ids,omitempty"`
}

type Attendant struct {
	InsId            string                  `json:"ins_id" bson:"ins_id"`
	Id               int32                   `json:"id" bson:"id"`
	Lv               int32                   `json:"lv" bson:"lv"`
	Star             int32                   `json:"star" bson:"star"`
	AwakeLv          int32                   `json:"awake_lv" bson:"awake_lv"`
	TalentList       []AttendantTalent       `json:"talent_list" bson:"talent_list"`
	RefineTalentList []AttendantTalent       `json:"refine_talent_list" bson:"refine_talent_list,omitempty"`
	AwakeTalentList  []AttendantTalent       `json:"awake_talent_list" bson:"awake_talent_list,omitempty"`
	SkinMap          map[int32]AttendantSkin `json:"skin_map" bson:"skin_map"`
}

type AttendantTalent struct {
	Id       int32 `json:"id" bson:"id"`
	Lv       int32 `json:"lv" bson:"lv"`
	Disabled bool  `json:"disabled" bson:"disabled"`
}

type AttendantSkin struct {
	Id     int32 `json:"id" bson:"id"`
	IsWear bool  `json:"is_wear" bson:"is_wear"`
}

type ExpLevelInfo struct {
	UsedPoint int32            `json:"used_point" bson:"used_point"`
	PointInfo map[string]int32 `json:"point_info" bson:"point_info"`
	Attrs     map[string]int64 `json:"attrs" bson:"attrs"`
}

type PlayerGuildInfo struct {
	GuildId              int32 `json:"guild_id" bson:"guild_id"`
	ExitTime             int32 `json:"exit_time" bson:"exit_time"`
	Position             int32 `json:"position" bson:"position"`
	LastRewardBossId     int32 `json:"last_reward_boss_id" bson:"last_reward_boss_id"`
	RewardExpireTime     int32 `json:"reward_expire_time" bson:"reward_expire_time"`
	RankRewardStatus     bool  `json:"rank_reward_status" bson:"rank_reward_status"`
	ClaimLastRewardCount int32 `json:"claim_last_reward_count" bson:"claim_last_reward_count"`
}

type Avatar struct {
	Id         int32 `json:"id" bson:"id"`
	ExpireTime int32 `json:"expire_time" bson:"expire_time"`
	IsRead     bool  `json:"is_read" bson:"is_read"`
}

type AvatarInfo struct {
	WearId    int32            `json:"wear_id" bson:"wear_id"`
	UnlockMap map[int32]Avatar `json:"unlock_map" bson:"unlock_map"`
}

// GemsInlayInfo 宝石镶嵌信息
type GemsInlayInfo struct {
	SlotInfo map[string]int32 `json:"slot_info" bson:"slot_info"`
	Attrs    map[string]int64 `json:"attrs" bson:"attrs"`
}

// FacadeMaskInfo 幻化信息
type FacadeMaskInfo struct {
	SlotInfo         map[string]int32 `json:"slot_info" bson:"slot_info"`
	UnlockInfo       map[string]int32 `json:"unlock_info" bson:"unlock_info"`
	UnlockRewardList []int32          `json:"unlock_reward_list" bson:"unlock_reward_list"`
}

// MonsterBossInfo 巨兽boss信息
type MonsterBossInfo struct {
	BossId     int32            `json:"boss_id" bson:"boss_id"`
	BossDamage map[string]int64 `json:"boss_damage" bson:"boss_damage"`
	UnlockTime int32            `json:"unlock_time" bson:"unlock_time"`
	UpdateTime int32            `json:"update_time" bson:"update_time"`
}

// DragonLevelInfo 巨龙信息
type DragonLevelInfo struct {
	LevelId      int32     `json:"level_id" bson:"level_id"`
	GroupConfigs [][]int32 `json:"group_configs" bson:"group_configs"`
}

// Collection 藏品
type Collection struct {
	Id       int32 `json:"id" bson:"id"`
	Lv       int32 `json:"lv" bson:"lv"`
	Star     int32 `json:"star" bson:"star"`
	SkillNum int32 `json:"skill_num" bson:"skill_num"`
}

// CollectionInfo 藏品信息
type CollectionInfo struct {
	GuaranteeUsed int32            `json:"guarantee_used" bson:"guarantee_used"`
	Attrs         map[string]int64 `json:"attrs" bson:"attrs"`
	Suits         map[string]int32 `json:"suits" bson:"suits"`
}

// MonsterWaveInfo 怪物波数数据
type MonsterWaveInfo struct {
	Number        int32 `json:"number" bson:"number"`
	MonsterId     int32 `json:"monster_id" bson:"monster_id"`
	IntervalPixel int32 `json:"interval_pixel" bson:"interval_pixel"`
}

// MonsterGroupData 怪物组波数数据
type MonsterGroupData struct {
	WavesInfo map[string]MonsterWaveInfo `json:"waves_info" bson:"waves_info"`
	Type      int32                      `json:"type" bson:"type"`
	Range     int32                      `json:"range" bson:"range"`
}

// LevelInfo 主线关卡信息
type LevelInfo struct {
	LevelProgressId int32              `json:"level_progress_id" bson:"level_progress_id"`
	NodeProgressId  int32              `json:"node_progress_id" bson:"node_progress_id"`
	MonsterData     []MonsterGroupData `json:"monster_data" bson:"monster_data"`
	MapEvent        []int32            `json:"map_event" bson:"map_event"`
	DayItemResource map[string]int32   `json:"day_item_resource" bson:"day_item_resource"`
	ExpireTime      int32              `json:"expire_time" bson:"expire_time"`
	EquipVitality   int32              `json:"equip_vitality" bson:"equip_vitality"`
	ExpVitality     int32              `json:"exp_vitality" bson:"exp_vitality"`
}

type RewardInfo struct {
	Type  int32  `json:"type" bson:"type"`
	Id    int32  `json:"id"  bson:"id"`
	Num   int64  `json:"num" bson:"num"`
	InsId string `json:"ins_id" bson:"ins_id,omitempty"`
}

// EmailInfo 邮件信息
type EmailInfo struct {
	EmailId     int64             `json:"email_id" bson:"email_id"`
	Type        int32             `json:"type" bson:"type"`
	SenderI18n  map[string]string `json:"senderI18n" bson:"senderI18n"`
	TitleI18n   map[string]string `json:"titleI18n" bson:"titleI18n"`
	ContentI18n map[string]string `json:"contentI18n" bson:"contentI18n"`
	TplId       int32             `json:"tpl_id" bson:"tpl_id"`
	TplParams   []string          `json:"tpl_params" bson:"tpl_params"`
	Attachments []RewardInfo      `json:"attachments" bson:"attachments"`
	Remark      string            `json:"remark" bson:"remark"`
	CreateTime  int32             `json:"create_time" bson:"create_time"`
	StartTime   int32             `json:"start_time" bson:"start_time"`
	EndTime     int32             `json:"end_time" bson:"end_time"`
	IsRead      bool              `json:"is_read" bson:"is_read"`
	IsRewarded  bool              `json:"is_rewarded" bson:"is_rewarded"`
}

type MineEventInfo struct {
	MineId    int32   `json:"mine_id" bson:"mine_id"`
	EventId   int32   `json:"event_id" bson:"event_id"`
	WorkerIds []int32 `json:"worker_ids" bson:"worker_ids"`
	EndTime   int32   `json:"end_time" bson:"end_time"`
}

// MineInfo 矿洞信息
type MineInfo struct {
	MineIdList        []int32                  `json:"mine_id_list" bson:"mine_id_list"`
	MineWorkerList    []int32                  `json:"mine_worker_list" bson:"mine_worker_list"`
	MineEventInfo     map[string]MineEventInfo `json:"mine_event_info" bson:"mine_event_info"`
	MineCdInfo        map[string]int32         `json:"mine_cd_info" bson:"mine_cd_info"`
	FinishEventIds    []int32                  `json:"finish_event_ids" bson:"finish_event_ids"`
	ExploreEnergy     int32                    `json:"explore_energy" bson:"explore_energy"`
	LastRecoverTime   int32                    `json:"last_recover_time" bson:"last_recover_time"`
	EventRefreshCount map[string]int32         `json:"event_refresh_count" bson:"event_refresh_count"`
	AutoEndTime       int32                    `json:"auto_end_time" bson:"auto_end_time"`
	DoMineNum         int32                    `json:"do_mine_num" bson:"do_mine_num"`
}

// SkillRuneInfo 符文信息
type SkillRuneInfo struct {
	Id   int32 `json:"id" bson:"id"`
	Lv   int32 `json:"lv" bson:"lv"`
	Rank int32 `json:"rank" bson:"rank"`
	Step int32 `json:"step" bson:"step"`
}

// MysterySkillInfo 奥义信息
type MysterySkillInfo struct {
	UsedMysterySkill      map[string]int32         `json:"used_mystery_skill" bson:"used_mystery_skill"`
	PutOnRuneInfo         map[string]int32         `json:"put_on_rune_info" bson:"put_on_rune_info"`
	RuneInfo              map[string]SkillRuneInfo `json:"rune_info" bson:"rune_info"`
	PassiveSkillInfo      map[string]int32         `json:"passive_skill_info" bson:"passive_skill_info"`
	ResonanceRuneGroupIds []int32                  `json:"resonance_rune_groupIds" bson:"resonance_rune_groupIds"`
	ResonanceLv           map[string]int32         `json:"resonance_lv" bson:"resonance_lv"`
	Attrs                 map[string]int64         `json:"attrs" bson:"attrs"`
}

// MysteryDrawRuneInfo 奥义符文抽奖信息
type MysteryDrawRuneInfo struct {
	ExtractDrawTimes int32 `json:"extract_draw_times" bson:"extract_draw_times"`
	ExtractType      int32 `json:"extract_type" bson:"extract_type"`
	DrawTimes        int32 `json:"draw_times" bson:"draw_times"`
}

// CrossPkInfo 跨服排位赛信息
type CrossPkInfo struct {
	Season    string   `json:"season" bson:"season"`
	Tier      int32    `json:"tier" bson:"tier"`
	CrossId   int32    `json:"cross_id" bson:"cross_id"`
	ZoneId    int32    `json:"zone_id" bson:"zone_id"`
	PrevTier  int32    `json:"prev_tier" bson:"prev_tier"`
	RivalList []string `json:"rival_list" bson:"rival_list"`
}

// Task 任务状态
type Task struct {
	RewardNum  int32 `json:"reward_num" bson:"reward_num"`
	ExpireTime int32 `json:"expire_time" bson:"expire_time"`
}

// TaskInfo 任务信息
type TaskInfo struct {
	TaskMap map[string]Task `json:"task_map" bson:"task_map"`
}

// RecoverInfo 恢复类道具信息
type RecoverInfo struct {
	Id             int32 `json:"id" bson:"id"`
	Value          int32 `json:"value" bson:"value"`
	LastUpdateTime int32 `json:"last_update_time" bson:"last_update_time"`
}

// MazeCrystalInfo 水晶信息，仅主格子存储
type MazeCrystalInfo struct {
	RewardStatus  bool   `json:"reward_status" bson:"reward_status"`
	CrystalStatus []bool `json:"crystal_status" bson:"crystal_status"`
}

// MazeSite 迷宫格子信息
type MazeSite struct {
	Site        int32           `json:"site" bson:"site"`
	SiteId      int32           `json:"site_id" bson:"site_id"`
	Status      bool            `json:"status" bson:"status"`
	CoverOptNum int32           `json:"cover_opt_num" bson:"cover_opt_num"`
	MasterSite  int32           `json:"master_site" bson:"master_site"`
	CrystalInfo MazeCrystalInfo `json:"crystal_info" bson:"crystal_info"`
}

// MazeFriend 迷宫挚友信息
type MazeFriend struct {
	Id       int32 `json:"id" bson:"id"`
	Lv       int32 `json:"lv" bson:"lv"`
	Exp      int32 `json:"exp" bson:"exp"`
	RewardLv int32 `json:"reward_lv" bson:"reward_lv"`
}

// MazeInfo 迷宫其它信息
type MazeInfo struct {
	Attrs map[string]int64 `json:"attrs" bson:"attrs"`
}

// MagicSlot 魔能槽位
type MagicSlot struct {
	Type  int32  `json:"type" bson:"type"`
	Lv    int32  `json:"lv" bson:"lv"`
	InsId string `json:"ins_id" bson:"ins_id"`
	Id    int32  `json:"id" bson:"id"`
}

// MagicEquipAttr 魔能装备属性词条
type MagicEquipAttr struct {
	Id          int32 `json:"id" bson:"id"`
	AttrId      int32 `json:"attr_id" bson:"attr_id"`
	AttrVal     int64 `json:"attr_val" bson:"attr_val"`
	RefineCount int32 `json:"refine_count" bson:"refine_count"`
}

// MagicEquip 魔能装备
type MagicEquip struct {
	InsId              string                    `json:"ins_id" bson:"ins_id"`
	Id                 int32                     `json:"id" bson:"id"`
	Lv                 int32                     `json:"lv" bson:"lv"`
	IsReaded           bool                      `json:"is_readed" bson:"is_readed"`
	IsLocked           bool                      `json:"is_locked" bson:"is_locked"`
	AttrEntries        []MagicEquipAttr          `json:"attr_entries" bson:"attr_entries"`
	RefinedAttrEntries map[string]MagicEquipAttr `json:"refined_attr_entries" bson:"refined_attr_entries,omitempty"`
}

// MagicEquipInfo 魔能装备信息
type MagicEquipInfo struct {
	Attrs          map[string]int64 `json:"attrs" bson:"attrs"`
	BookMap        map[string]bool  `json:"book_map" bson:"book_map"`
	Capacity       int32            `json:"capacity" bson:"capacity"`
	CapacityExpand int32            `json:"capacity_expand" bson:"capacity_expand"`
}

// MagicRealmInfo 魔能秘境关卡信息
type MagicRealmInfo struct {
	LevelId        int32                       `json:"level_id" bson:"level_id"`
	Gifts          [][]int32                   `json:"gifts" bson:"gifts"`
	Multiple       int32                       `json:"multiple" bson:"multiple"`
	MapMonsterData []map[string]int32          `json:"map_monster_data" bson:"map_monster_data"`
	Maps           [][]int32                   `json:"maps" bson:"maps"`
	Rewards        map[string]map[string]int32 `json:"rewards" bson:"rewards"`
	AtkLevelId     int32                       `json:"atk_level_id" bson:"atk_level_id"`
	SoulReward     int32                       `json:"soul_reward" bson:"soul_reward"`
	BossData       map[string]int32            `json:"boss_data" bson:"boss_data"`
	SelectedMapIds []int32                     `json:"selected_map_ids" bson:"selected_map_ids"`
}

// GoodsInfo 商品信息
type GoodsInfo struct {
	BuyCount   int32 `json:"buy_count" bson:"buy_count"`
	ExpireTime int32 `json:"expire_time" bson:"expire_time"`
}

// FundInfo 基金信息
type FundInfo struct {
	Type                  int32 `json:"type" bson:"type,omitempty"`
	MaxFreeRewardProgress int32 `json:"max_free_reward_progress" bson:"max_free_reward_progress"`
	MaxPayRewardProgress  int32 `json:"max_pay_reward_progress" bson:"max_pay_reward_progress"`
}

// PlayerBattleInfo 玩家战斗数据
type PlayerBattleInfo struct {
	PlayerId         string               `json:"player_id" bson:"player_id"`
	Nickname         string               `json:"nickname" bson:"nickname"`
	Lv               int32                `json:"lv" bson:"lv"`
	BreakLv          int32                `json:"break_lv" bson:"break_lv"`
	Power            int64                `json:"power" bson:"power"`
	AuraPutOnId      int32                `json:"aura_put_on_id" bson:"aura_put_on_id"`
	FacadeSlotInfo   map[string]int32     `json:"facade_slot_info" bson:"facade_slot_info"`
	Attrs            map[string]int64     `json:"attrs" bson:"attrs"`
	AttendantInfo    AttendantInfo        `json:"attendant_info" bson:"attendant_info"`
	Attendants       map[string]Attendant `json:"attendants" bson:"attendants"`
	MagicSlotMap     map[string]MagicSlot `json:"magic_slot_map" bson:"magic_slot_map"`
	MysterySkillInfo MysterySkillInfo     `json:"mystery_skill_info" bson:"mystery_skill_info"`
	UniqueId         string               `json:"unique_id" bson:"unique_id"`
	AvatarInfoMap    map[int32]AvatarInfo `json:"avatar_info_map" bson:"avatar_info_map"`
	AreaId           int32                `json:"area_id" bson:"area_id"`
	LoginTime        int32                `json:"login_time" bson:"login_time"`
}

// GuildBattleFightInfo 攻防战上一场战斗信息战斗包
type GuildBattleFightInfo struct {
	Seed    int32            `json:"seed" bson:"seed"`
	AtkInfo PlayerBattleInfo `json:"atk_info" bson:"atk_info"`
	DefInfo PlayerBattleInfo `json:"def_info" bson:"def_info"`
	RobotId int32            `json:"robot_id" bson:"robot_id"`
}

// UserGuildBattleInfo 公会攻防战
type UserGuildBattleInfo struct {
	FightTimes     []int32              `json:"fight_times" bson:"fight_times"`
	LastJoinDate   string               `json:"last_join_date" bson:"last_join_date"`
	LastFightInfo  GuildBattleFightInfo `json:"last_fight_info" bson:"last_fight_info"`
	LastShowResult int32                `json:"last_show_result" bson:"last_show_result"`
}

// SignInInfo 签到信息
type SignInInfo struct {
	Cycle      int32   `json:"cycle" bson:"cycle"`
	DayInCycle int32   `json:"day_in_cycle" bson:"day_in_cycle"`
	DayStatus  []int32 `json:"day_status" bson:"day_status"`
}

// ActiveInfo 活跃度信息
type ActiveInfo struct {
	Count                 int32 `json:"count" bson:"count"`
	ExpireTime            int32 `json:"expire_time" bson:"expire_time"`
	MaxFreeRewardProgress int32 `json:"max_free_reward_progress" bson:"max_free_reward_progress"`
	MaxPayRewardProgress  int32 `json:"max_pay_reward_progress" bson:"max_pay_reward_progress"`
	IsPay                 bool  `json:"is_pay" bson:"is_pay"`
}

// SecretBookInfo 历练秘典信息
type SecretBookInfo struct {
	Progress       int32            `json:"progress" bson:"progress"`
	DailyRewardMap map[string]int32 `json:"daily_reward_map" bson:"daily_reward_map"`
	ExpireTime     int32            `json:"expire_time" bson:"expire_time"`
}

// PrivilegeCard 特权卡信息
type PrivilegeCard struct {
	EndTime       int32 `json:"end_time" bson:"end_time"`
	LastClaimTime int32 `json:"last_claim_time" bson:"last_claim_time"`
}

// LuckyPotionsInfo 幸运魔炉信息
type LuckyPotionsInfo struct {
	LuckyPotionsNum     int32 `json:"lucky_potions_num" bson:"lucky_potions_num"`
	PaidLuckyPotionsNum int32 `json:"paid_lucky_potions_num" bson:"paid_lucky_potions_num"`
	GuaranteeUsed       int32 `json:"guarantee_used" bson:"guarantee_used"`
}

// FirstRechargeInfo 首充信息
type FirstRechargeInfo struct {
	ClaimCount int32 `json:"claim_count" bson:"claim_count"`
	BuyTime    int32 `json:"buy_time" bson:"buy_time"`
}

// CarnivalActivityInfo 七日嘉年华信息
type CarnivalActivityInfo struct {
	ScoreRewardRecord []int32 `json:"score_reward_record" bson:"score_reward_record"`
	IsSendEmail       int32   `json:"is_send_email" bson:"is_send_email"`
}

// SprintActivityInfo 冲刺活动信息
type SprintActivityInfo struct {
	ActId     int32 `json:"act_id" bson:"act_id"`
	StartTime int32 `json:"start_time" bson:"start_time"`
	EndTime   int32 `json:"end_time" bson:"end_time"`
}

// RechargeGiftInfo 连天贺礼
type RechargeGiftInfo struct {
	LastUpdateTime int32   `json:"last_update_time" bson:"last_update_time"`
	Days           int32   `json:"days" bson:"days"`
	GetRecord      []int32 `json:"get_record" bson:"get_record"`
}

// ActSweepstakeInfo 抽奖活动信息
type ActSweepstakeInfo struct {
	ActId             int32            `json:"act_id" bson:"act_id"`
	StartTime         int32            `json:"start_time" bson:"start_time"`
	EndTime           int32            `json:"end_time" bson:"end_time"`
	GuaranteeTimes    map[string]int32 `json:"guarantee_times" bson:"guarantee_times"`
	SelectedInfo      map[string]int32 `json:"selected_info" bson:"selected_info"`
	DrawCount         int32            `json:"draw_count" bson:"draw_count"`
	ProgressGetRecord []int32          `json:"progress_get_record" bson:"progress_get_record"`
}

// ActFundInfo 活动战令信息
type ActFundInfo struct {
	ActId                int32   `json:"act_id" bson:"act_id"`
	StartTime            int32   `json:"start_time" bson:"start_time"`
	EndTime              int32   `json:"end_time" bson:"end_time"`
	FreeRewardGetRecord  []int32 `json:"free_reward_get_record" bson:"free_reward_get_record"`
	MidRewardGetRecord   []int32 `json:"mid_reward_get_record" bson:"mid_reward_get_record"`
	HighRewardGetRecord  []int32 `json:"high_reward_get_record" bson:"high_reward_get_record"`
	ExtraScoreGetCount   int32   `json:"extra_score_get_count" bson:"extra_score_get_count"`
	AddUpRewardGetRecord []int32 `json:"add_up_reward_get_record" bson:"add_up_reward_get_record"`
	Progress             int32   `json:"progress" bson:"progress"`
	IsBuyMid             bool    `json:"is_buy_mid" bson:"is_buy_mid"`
	IsBuyHigh            bool    `json:"is_buy_high" bson:"is_buy_high"`
}

// ActLoginInfo 活动登录信息
type ActLoginInfo struct {
	ActId          int32   `json:"act_id" bson:"act_id"`
	StartTime      int32   `json:"start_time" bson:"start_time"`
	EndTime        int32   `json:"end_time" bson:"end_time"`
	LoginGetRecord []int32 `json:"login_get_record" bson:"login_get_record"`
}

// ActShootGameInfo 活动射箭小游戏
type ActShootGameInfo struct {
	ActId           int32            `json:"act_id" bson:"act_id"`
	StartTime       int32            `json:"start_time" bson:"start_time"`
	EndTime         int32            `json:"end_time" bson:"end_time"`
	ScoreMap        map[string]int32 `json:"score_map" bson:"score_map"`
	RewardGetRecord []int32          `json:"reward_get_record" bson:"reward_get_record"`
}

// ActFruitGameInfo 活动水果杀手小游戏
type ActFruitGameInfo struct {
	ActId           int32            `json:"act_id" bson:"act_id"`
	StartTime       int32            `json:"start_time" bson:"start_time"`
	EndTime         int32            `json:"end_time" bson:"end_time"`
	ScoreMap        map[string]int32 `json:"score_map" bson:"score_map"`
	RewardGetRecord []int32          `json:"reward_get_record" bson:"reward_get_record"`
}

// PlayerFashionInfo 玩家时装信息
type PlayerFashionInfo struct {
	UnlockInfo map[string]int32 `json:"unlock_info" bson:"unlock_info"`
	StarInfo   map[string]int32 `json:"star_info" bson:"star_info"`
	Attrs      map[string]int64 `json:"attrs" bson:"attrs"`
}

// GTRivalInfo 对手信息
type GTRivalInfo struct {
	RivalPlayerId string `json:"rival_player_id" bson:"rival_player_id"`
	Score         int32  `json:"score" bson:"score"`
}

// BattleInfo 玩家战斗信息
type BattleInfo struct {
	Seed       int32            `json:"seed" bson:"seed"`
	AtkInfo    PlayerBattleInfo `json:"atk_info" bson:"atk_info"`
	DefInfo    PlayerBattleInfo `json:"def_info" bson:"def_info"`
	RivalIndex int32            `json:"rival_index" bson:"rival_index"`
}

// PlayerGTInfo 玩家夺旗信息
type PlayerGTInfo struct {
	RivalGuildId     int32            `json:"rival_guild_id" bson:"rival_guild_id"`
	RivalGuildName   string           `json:"rival_guild_name" bson:"rival_guild_name"`
	BattleInfo       BattleInfo       `json:"battle_info" bson:"battle_info"`
	Energy           int32            `json:"energy" bson:"energy"`
	RivalList        []GTRivalInfo    `json:"rival_list" bson:"rival_list"`
	CombatStatus     int32            `json:"combat_status" bson:"combat_status"`
	BattleType       int32            `json:"battle_type" bson:"battle_type"`
	BuffGoodsId      int32            `json:"buff_goods_id" bson:"buff_goods_id"`
	DefeatPlayerList []string         `json:"defeat_player_list" bson:"defeat_player_list"`
	GainScore        int32            `json:"gain_score" bson:"gain_score"`
	UseItems         map[string]int32 `json:"use_items" bson:"use_items"`
	EndTime          int32            `json:"end_time" bson:"end_time"`
	FundInfo         FundInfo         `json:"fund_info" bson:"fund_info"`
}

// PlayerBaseInfo 玩家基础信息
type PlayerBaseInfo struct {
	PlayerId        string               `json:"player_id" bson:"player_id"`
	Nickname        string               `json:"nickname" bson:"nickname"`
	Lv              int32                `json:"lv" bson:"lv"`
	BreakLv         int32                `json:"break_lv" bson:"break_lv"`
	Power           int64                `json:"power" bson:"power"`
	AreaId          int32                `json:"area_id" bson:"area_id"`
	UniqueId        string               `json:"unique_id" bson:"unique_id"`
	AuraPutOnId     int32                `json:"aura_put_on_id" bson:"aura_put_on_id"`
	FacadeSlotInfo  map[int32]int32      `json:"facade_slot_info" bson:"facade_slot_info"`
	AvatarInfoMap   map[int32]AvatarInfo `json:"avatar_info_map" bson:"avatar_info_map"`
	LoginTime       int32                `json:"login_time" bson:"login_time"`
	PlayerGuildInfo PlayerGuildInfo      `json:"player_guild_info" bson:"player_guild_info"`
}

type PlayerInfo struct {
	DataStructVersion     int32                           `json:"data_struct_version" bson:"data_struct_version"`
	UserId                string                          `json:"user_id" bson:"user_id"`
	PlayerId              string                          `json:"player_id" bson:"player_id"`
	Nickname              string                          `json:"nickname" bson:"nickname"`
	Lv                    int32                           `json:"lv" bson:"lv"`
	Exp                   int32                           `json:"exp" bson:"exp"`
	BreakLv               int32                           `json:"break_lv" bson:"break_lv"`
	BreakTaskGetList      []int32                         `json:"break_task_get_list" bson:"break_task_get_list"`
	BreakTaskProgress     map[string]int32                `json:"break_task_progress" bson:"break_task_progress"`
	Attrs                 map[string]int64                `json:"attrs" bson:"attrs"`
	Power                 int64                           `json:"power" bson:"power"`
	Sex                   int32                           `json:"sex" bson:"sex"`
	AreaId                int32                           `json:"area_id" bson:"area_id"`
	BanInfo               BanInfo                         `json:"ban_info" bson:"ban_info"`
	UniqueId              string                          `json:"unique_id" bson:"unique_id"`
	RedeemCode            []string                        `json:"redeem_code" bson:"redeem_code"`
	LoginTime             int32                           `json:"login_time" bson:"login_time"`
	RegisterTime          int32                           `json:"register_time" bson:"register_time"`
	LastCheckTime         int32                           `json:"last_check_time" bson:"last_check_time"`
	GuideIdList           []int32                         `json:"guide_id_list" bson:"guide_id_list"`
	MainTask              MainTaskInfo                    `json:"main_task" bson:"main_task"`
	Items                 map[string]int32                `json:"items" bson:"items"`
	TasksCategoriesInfo   map[string]TaskProgressInfo     `json:"tasks_categories_info" bson:"tasks_categories_info"`
	EquipsInfo            EquipsInfo                      `json:"equips_info" bson:"equips_info"`
	AuraInfo              AuraInfo                        `json:"aura_info" bson:"aura_info"`
	BlessInfo             BlessInfo                       `json:"bless_info" bson:"bless_info"`
	AttendantInfo         AttendantInfo                   `json:"attendant_info" bson:"attendant_info"`
	Attendants            map[string]Attendant            `json:"attendants" bson:"attendants"`
	ExpLevelInfo          ExpLevelInfo                    `json:"explevel_info" bson:"explevel_info"`
	Gems                  map[string]int32                `json:"gems" bson:"gems"`
	GemsInlayInfo         GemsInlayInfo                   `json:"gems_inlay_info" bson:"gems_inlay_info"`
	FacadeMaskInfo        FacadeMaskInfo                  `json:"facade_mask_info" bson:"facade_mask_info"`
	MonsterBossInfo       MonsterBossInfo                 `json:"monster_boss_info" bson:"monster_boss_info"`
	DragonLevelInfo       DragonLevelInfo                 `json:"dragon_level_info" bson:"dragon_level_info"`
	CollectionPieces      map[string]int32                `json:"collection_pieces" bson:"collection_pieces"`
	Collections           map[string]Collection           `json:"collections" bson:"collections"`
	CollectionInfo        CollectionInfo                  `json:"collection_info" bson:"collection_info"`
	LevelInfo             LevelInfo                       `json:"level_info" bson:"level_info"`
	Emails                map[string]EmailInfo            `json:"emails" bson:"emails"`
	PendingEmails         map[string]EmailInfo            `json:"pending_emails" bson:"pending_emails"`
	LastPersonalEmailId   int64                           `json:"last_personal_email_id" bson:"last_personal_email_id"`
	LastSystemEmailId     int64                           `json:"last_system_email_id" bson:"last_system_email_id"`
	MineInfo              MineInfo                        `json:"mine_info" bson:"mine_info"`
	DailyTimes            map[string]int32                `json:"daily_times" bson:"daily_times"`
	DailyTimesExpireAt    int32                           `json:"daily_times_expire_at" bson:"daily_times_expire_at"`
	MysteryRunes          map[string]int32                `json:"mystery_runes" bson:"mystery_runes"`
	MysterySkillInfo      MysterySkillInfo                `json:"mystery_skill_info" bson:"mystery_skill_info"`
	MysteryDrawRuneInfo   MysteryDrawRuneInfo             `json:"mystery_draw_rune_info" bson:"mystery_draw_rune_info"`
	CrossPkInfo           CrossPkInfo                     `json:"cross_pk_info" bson:"cross_pk_info"`
	TaskInfoMap           map[string]TaskInfo             `json:"task_info_map" bson:"task_info_map"`
	PlayerGuildInfo       PlayerGuildInfo                 `json:"player_guild_info" bson:"player_guild_info"`
	RecoverMap            map[string]RecoverInfo          `json:"recover_map" bson:"recover_map"`
	MazeSiteMap           map[string]MazeSite             `json:"maze_site_map" bson:"maze_site_map"`
	MazeFriendMap         map[string]MazeFriend           `json:"maze_friend_map" bson:"maze_friend_map"`
	MazeInfo              MazeInfo                        `json:"maze_info" bson:"maze_info"`
	SystemUnlockMap       map[string]int32                `json:"system_unlock_map" bson:"system_unlock_map"`
	MagicSlotMap          map[string]MagicSlot            `json:"magic_slot_map" bson:"magic_slot_map"`
	MagicEquipMap         map[string]MagicEquip           `json:"magic_equip_map" bson:"magic_equip_map"`
	MagicEquipInfo        MagicEquipInfo                  `json:"magic_equip_info" bson:"magic_equip_info"`
	MagicRealmInfo        MagicRealmInfo                  `json:"magic_realm_info" bson:"magic_realm_info"`
	ShopMap               map[string]map[string]GoodsInfo `json:"shop_map" bson:"shop_map"`
	SystemUnlockAwardList []int32                         `json:"system_unlock_award_list" bson:"system_unlock_award_list"`
	FundMap               map[string]FundInfo             `json:"fund_map" bson:"fund_map"`
	UserGuildBattleInfo   UserGuildBattleInfo             `json:"user_guild_battle_info" bson:"user_guild_battle_info"`
	SignInInfo            SignInInfo                      `json:"sign_in_info" bson:"sign_in_info"`
	ActiveMap             map[string]ActiveInfo           `json:"active_map" bson:"active_map"`
	AvatarInfoMap         map[int32]AvatarInfo            `json:"avatar_info_map" bson:"avatar_info_map"`
	SecretBookInfo        SecretBookInfo                  `json:"secret_book_info" bson:"secret_book_info"`
	RechargeScore         int64                           `json:"recharge_score" bson:"recharge_score"`
	PrivilegeCardMap      map[string]PrivilegeCard        `json:"privilege_card_map" bson:"privilege_card_map"`
	LuckyPotionsInfo      LuckyPotionsInfo                `json:"lucky_potions_info" bson:"lucky_potions_info"`
	FirstRechargeMap      map[string]FirstRechargeInfo    `json:"first_recharge_map" bson:"first_recharge_map"`
	CarnivalActivityInfo  CarnivalActivityInfo            `json:"carnival_activity_info" bson:"carnival_activity_info"`
	SprintActivityInfo    SprintActivityInfo              `json:"sprint_activity_info" bson:"sprint_activity_info"`
	RechargeGiftInfo      RechargeGiftInfo                `json:"recharge_gift_info" bson:"recharge_gift_info"`
	ActSweepstakeInfoMap  map[string]ActSweepstakeInfo    `json:"act_sweepstake_info_map" bson:"act_sweepstake_info_map"`
	ActFundInfoMap        map[string]ActFundInfo          `json:"act_fund_info_map" bson:"act_fund_info_map"`
	ActLoginInfoMap       map[string]ActLoginInfo         `json:"act_login_info_map" bson:"act_login_info_map"`
	ActShootGameInfoMap   map[string]ActShootGameInfo     `json:"act_shootgame_info_map" bson:"act_shootgame_info_map"`
	ActFruitGameInfoMap   map[string]ActFruitGameInfo     `json:"act_fruitgame_info_map" bson:"act_fruitgame_info_map"`
	PlayerFashionInfo     PlayerFashionInfo               `json:"player_fashion_info" bson:"player_fashion_info"`
	PlayerGTInfo          PlayerGTInfo                    `json:"player_gt_info" bson:"player_gt_info"`
	DailyResourceExpireAt int32                           `json:"daily_resource_expire_at" bson:"daily_resource_expire_at"`
}

// PlayerInfoWithStatus 带状态的玩家信息（用于API返回）
type PlayerInfoWithStatus struct {
	*PlayerInfo
	Online   bool `json:"online"`    // 在线状态
	BanLogin bool `json:"ban_login"` // 封号状态
	BanChat  bool `json:"ban_chat"`  // 禁言状态
}

// ToPlayerInfoWithStatus 将PlayerInfo转换为带状态的结构体
func (p *PlayerInfo) ToPlayerInfoWithStatus() *PlayerInfoWithStatus {
	return &PlayerInfoWithStatus{
		PlayerInfo: p,
		Online:     false, // 默认值，需要从外部设置
		BanLogin:   false,
		BanChat:    false,
	}
}
