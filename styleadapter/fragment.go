package styleadapter

const (
	TypeKey = "__type__"

	TypeMotionDiv     = "div"
	TypeSpan          = "span"
	TypeRect          = "rect"
	TypePath          = "path"
	TypeForeignObject = "foreignObject"
	TypeDiv           = "div" // inner XHTML inside foreignObject
	TypeFilter        = "filter"
	TypeMask          = "mask"
	TypeG             = "g"
)

// StyleFragment is a flat map with mandatory __type__ and target properties.
type StyleFragment map[string]any

func NewFragment(typ string) StyleFragment {
	return StyleFragment{TypeKey: typ}
}

func MergeFragments(fragments []StyleFragment) []StyleFragment {
	byType := make(map[string]StyleFragment)
	var order []string
	for _, f := range fragments {
		if f == nil {
			continue
		}
		typ, _ := f[TypeKey].(string)
		if typ == "" {
			continue
		}
		if existing, ok := byType[typ]; ok {
			for k, v := range f {
				if k != TypeKey {
					existing[k] = v
				}
			}
		} else {
			cp := make(StyleFragment, len(f))
			for k, v := range f {
				cp[k] = v
			}
			byType[typ] = cp
			order = append(order, typ)
		}
	}
	out := make([]StyleFragment, 0, len(order))
	for _, typ := range order {
		out = append(out, byType[typ])
	}
	return out
}

func (f StyleFragment) Set(prop string, value any) {
	if f == nil {
		return
	}
	f[prop] = value
}

func (f StyleFragment) GetString(prop string) string {
	if f == nil {
		return ""
	}
	v, ok := f[prop]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
