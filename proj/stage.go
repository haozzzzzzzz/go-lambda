package proj

type StageType int8

const (
	TestStage StageType = 1 // 测试
	ProdStage StageType = 2 // 正式
)

func (m StageType) String() string {
	switch m {
	case TestStage:
		return "test"
	case ProdStage:
		return "prod"
	}
	return ""
}
