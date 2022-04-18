package alpha

type Env[K comparable, V any] struct {
	parent   *Env[K, V]
	bindings map[K]V
}

func NewEnv[K comparable, V any](parent *Env[K, V]) *Env[K, V] {
	return &Env[K, V]{
		parent:   parent,
		bindings: map[K]V{},
	}
}

func zero[V any]() V {
	var result V
	return result
}

func (env Env[K, V]) LookUp(key K) (V, bool) {
	if v, ok := env.bindings[key]; ok {
		return v, ok
	} else if env.parent != nil {
		return env.parent.LookUp(key)
	} else {
		return zero[V](), false
	}
}

func (env Env[K, V]) Add(key K, value V) {
	env.bindings[key] = value
}

// Contain returns whether the key occurs at the current level.
func (env Env[K, V]) Contain(key K) bool {
	_, existed := env.bindings[key]
	return existed
}
