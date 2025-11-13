package common

// Coalesce 返回第一個非空字串，若全空則返回最後一個（fallback）
func Coalesce(strs ...string) string {
    for _, s := range strs {
        if s != "" {
            return s
        }
    }
    return strs[len(strs)-1]
}