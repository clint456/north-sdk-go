package interfaces

// NorthServiceSDK
// 上游发布数据+接收平台下行命令
type NorthServiceSDK interface {
	// ========= North 通道实例生命周期（对应device） =========
	// AddNorth 添加上游通道实例
	AddNorth(north models.North) (string, error)
	// Norths 返回所有管理的上游通道实例
	Norths() []models.North
	// GetNorthByName 根据名称获取上游通道实例
	GetNorthByName(name string) (models.North, error)
	// UpdateNorth 更新上游通道实例
	UpdateNorth(north models.North) error
	// RemoveNorthByName 根据名称移除上游通道实例
	RemoveNorthByName(name string) error
	// 判断指定的上游通道实例是否存在
	ExistsNorthByName(name string) bool
	// PatchNorth 按需打补丁更新北向通道
	PatchNorth(north models.North) error

	// ========= NorthProfile 北向点位元数据生命周期（对应deivceprofileofile=====
	AddNorthProfile(profile models.NorthProfile) (string, error)
	NorthProfiles() []models.NorthProfile
	GetNorthProfileByName(name string) (models.NorthProfile, error)
	UpdateNorthProfile(profile models.NorthProfile) error
	RemoveNorthProfileByName(name string) error
	ExistsNorthProfileByName(name string) bool
	PatchNorthProfile(profile models.NorthProfile) error

	// ========= 点位/命令查询（对应deviceResource/DeviceCommand） =========
	// NorthResource 根据北向通道名称与资源名。返回缓存中的北向点位
	NorthResource(northName string, northResourceName string) (models.NorthResource, bool)
	// NorthCommand 根据北向通道名称与命令，返回缓存中的北向命令
	NorthCommand(northName string, northCommandName string) (models.NorthCommand, bool)

	// ================= 上游发布（数据出：向平台推送，核心反转点） =================
	// PublishEvent 向外部北向平台发布一个事件（值+品质+时标）
	PublishEvent(event models.Event) error
	// PublishReading 向外部平台发布单点读数
	PublishReading(reading models.Reading) error
	// PublishReport 发布协议特有的报告（IEC61850 Report 等），由协议适配器消费
	PublishReport(northName string, report sdkModels.Report) error
	// BatchPublishEvents 批量上送，适配遥测整包上送
	BatchPublishEvents(events []models.Event) error
	// AsyncPublishingEnabled 返回是否启用异步上送
	AsyncPublishingEnabled() bool
	// AsyncEventsChannel 返回异步上送通道，SDK 内部将数据推给北向协议驱动
	AsyncEventsChannel() chan *sdkModels.UpstreamEvent

	// ================= 下游命令（命令进：接收平台指令，第二个反转点） =================
	// RegisterCommandHandler 注册平台下行命令回调（遥控/读点/取代/定值读写）
	// northName 为北向通道名；handler 收到命令后完成路由并返回结果
	RegisterCommandHandler(northName string, handler func(cmd models.NorthCommand) (CommandResult, error)) error
	// UnregisterCommandHandler 注销平台下行命令回调
	UnregisterCommandHandler(northName string) error
	// CommandResponse 异步回执平台命令（结果/失败原因）
	CommandResponse(reqId string, result CommandResult) error

	// ================= 通道运行状态（对应 Device OperatingState） =================
	// UpdateNorthOperatingState 更新北向通道运行状态（UP/DOWN/PAUSED）
	UpdateNorthOperatingState(name string, state models.OperatingState) error

	// ================= 自动发布（对应 AutoEvent，定时周期上送） =================
	// AddAutoPublish 给北向通道新增一个定时上送任务
	AddAutoPublish(northName string, autoPublish models.AutoPublish) error
	// RemoveAutoPublish 移除北向通道的定时上送任务
	RemoveAutoPublish(northName string, autoPublish models.AutoPublish) error

	// ================= 通道发现（可选，对应 ProvisionWatcher） =================
	// AddProvisionWatcher 添加一个通道发现器
	AddProvisionWatcher(watcher models.ProvisionWatcher) (string, error)
	ProvisionWatchers() []models.ProvisionWatcher
}
