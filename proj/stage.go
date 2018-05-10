package proj

type StageType int8

const (
	DevStage  StageType = 1 // 开发
	TestStage StageType = 2 // 测试
	PreStage  StageType = 3
	ProdStage StageType = 4 // 正式
)

func (m StageType) String() string {
	switch m {
	case DevStage:
		return "dev"
	case TestStage:
		return "test"
	case PreStage:
		return "pre"
	case ProdStage:
		return "prod"
	}
	return ""
}
