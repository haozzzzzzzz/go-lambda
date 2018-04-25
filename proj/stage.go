package proj

type StageType int8

const (
	DevStage  StageType = 1 // 开发
	TestStage StageType = 2 // 测试
	ProdStage StageType = 3 // 正式
)

func (m StageType) String() string {
	switch m {
	case DevStage:
		return "dev"
	case TestStage:
		return "test"
	case ProdStage:
		return "prod"
	}
	return ""
}
