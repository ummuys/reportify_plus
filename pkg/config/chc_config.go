package config

// type RepCacheConfig struct {
// 	Addr     string
// 	Password string
// 	DB       int
// 	Exp      int
// }

// func ParseReportCacheEnv() (RepCacheConfig, error) {
// 	var errs []string

// 	add := func(err error) {
// 		if err != nil {
// 			errs = append(errs, err.Error())
// 		}
// 	}

// 	addr, err := parseStr("CACHE_ADDR")
// 	add("cache_addr")

// 	pass, err := parseStr("CACHE_PASSWORD")

// 	add("cache_password")

// 	db, err := parseInt("CACHE_DB", true)

// 	add("cache_db")

// 	exp, err := parseInt("CACHE_EXPIRE_TIME", true)

// 	add("cache_db")

// 	if len(errs) > 0 {
// 		return {}, errors.New(strings.Join(errs, ", "))
// 	}

// 	return RepCacheConfig{
// 		Addr:     addr,
// 		Password: pass,
// 		DB:       db,
// 		Exp:      exp,
// 	}, nil
// }
