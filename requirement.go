package version

var (
	requirementOperation = map[string]operatorFunc{}
)

const ()

func init() {

}

type Requirements struct {
	requirements [][]requirement
}

type requirement struct {
	version  Version
	operator operatorFunc
	original string
}

func NewRequirements(v string) (Requirements, error) {
	return Requirements{}, nil
}

func (rs Requirements) Check(v Version) bool {
	for _, r := range rs.requirements {
		if andRequirementCheck(v, r) {
			return true
		}
	}
	return false
}

func andRequirementCheck(v Version, requirements []requirement) bool {
	for _, c := range requirements {
		if !c.check(v) {
			return false
		}
	}
	return true
}

func (r requirement) check(v Version) bool {
	return r.operator(v, r.version)
}
