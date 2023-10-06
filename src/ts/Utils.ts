const SEC = {
	ONE_MINUTE : 60, // seconds
	ONE_HOUR   : 3600,
	ONE_DAY    : 86400,
	ONE_WEEK   : 604800,
	ONE_MONTH  : 2419200,
	ONE_YEAR   : 31536000
}

export function RenderDateShort(UnixDateInt : number) : string {
	if(!UnixDateInt) return "-1y";
	const UnixDateNow = new Date().getTime();
	const DiffSeconds = Math.floor(UnixDateNow/1000 - UnixDateInt)

	if (DiffSeconds < 60) return `${DiffSeconds}s`
	if (DiffSeconds < SEC.ONE_HOUR) return `${Math.floor(DiffSeconds/SEC.ONE_MINUTE)}m`
	if (DiffSeconds < SEC.ONE_DAY) return `${Math.floor(DiffSeconds/SEC.ONE_HOUR)}h`
	if (DiffSeconds < SEC.ONE_WEEK) return `${Math.floor(DiffSeconds/SEC.ONE_DAY)}d`
	if (DiffSeconds < SEC.ONE_MONTH) return `${Math.floor(DiffSeconds/SEC.ONE_WEEK)}w`
	if (DiffSeconds < SEC.ONE_YEAR) return `${Math.floor(DiffSeconds/SEC.ONE_MONTH)}m`
	return `${Math.floor(DiffSeconds/SEC.ONE_YEAR)}y`
}


	


// "renderDate": func(unixDateInt int) string {
// 	var visualAmount int
// 	var pluralS string

// 	if unixDateInt == 0 {
// 		return "never"
// 	}

// 	now := time.Now()

// 	unixDate := time.Unix(int64(unixDateInt), 0)

// 	diffSeconds := int(math.Floor(now.Sub(unixDate).Seconds()))

// 	if diffSeconds < 60 {
// 		visualAmount = diffSeconds
// 		if visualAmount != 1 {
// 			pluralS = "s"
// 		}
// 		return fmt.Sprintf("%d second%s ago", visualAmount, pluralS)
// 	} else if diffSeconds < 3600 {
// 		visualAmount = int(math.Floor(float64(diffSeconds) / float64(60)))
// 		if visualAmount != 1 {
// 			pluralS = "s"
// 		}
// 		return fmt.Sprintf("%d minute%s ago", visualAmount, pluralS)
// 	} else if diffSeconds < 86400 {
// 		visualAmount = int(math.Floor(float64(diffSeconds) / float64(3600)))
// 		if visualAmount != 1 {
// 			pluralS = "s"
// 		}
// 		return fmt.Sprintf("%d hour%s ago", visualAmount, pluralS)
// 	} else if diffSeconds < 604800 {
// 		visualAmount = int(math.Floor(float64(diffSeconds) / float64(86400)))
// 		if visualAmount != 1 {
// 			pluralS = "s"
// 		}
// 		return fmt.Sprintf("%d day%s ago", visualAmount, pluralS)
// 	} else if diffSeconds < 2419200 {
// 		visualAmount = int(math.Floor(float64(diffSeconds) / float64(604800)))
// 		if visualAmount != 1 {
// 			pluralS = "s"
// 		}
// 		return fmt.Sprintf("%d week%s ago, at %d:%d", visualAmount, pluralS, unixDate.Hour(), unixDate.Minute())
// 	} else {
// 		return fmt.Sprintf("%s %d, %d", unixDate.Month(), unixDate.Day(), unixDate.Year())
// 	}
// },