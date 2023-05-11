package bili

import (
	"fmt"
	"net/url"
	"strings"
)

func GetBVIDFromURL(rawURL string) (string, error) {
	if strings.HasPrefix(rawURL, "BV") && len(rawURL) == 12 {
		return rawURL, nil
	}
	// https://www.bilibili.com/video/BV1NP411o7MB/?spm_id_from=333.1007.top_right_bar_window_custom_collection.content.click

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(u.Host, ".bilibili.com") {
		return "", fmt.Errorf("invalid bilibili url")
	}

	path := strings.Trim(u.Path, "/")
	if !strings.HasPrefix(path, "video/") {
		return "", fmt.Errorf("invalid bilibili video url")
	}

	pathItems := strings.Split(path, "/")
	if len(pathItems) < 2 {
		return "", fmt.Errorf("invalid bilibili video url : %s", rawURL)
	}

	return pathItems[1], nil
}

func CleanName(raw string) string {
	replacer := strings.NewReplacer(
		" ", "　",
		"(", "（",
		")", "）",
		"/", "／",
		";", "",
		"~", "",
		"$", "＄",
		".", "．",
		",", "，",
		"#", "＃",
		":", "：",
		"\"", "＂",
		"'", "＇",
		"?", "？",
		"-", "－",
		"@", "＠",
		"%", "％",
		"^", "＾",
		"&", "＆",
		"*", "＊",
		"!", "！",
		"+", "＋",
		"=", "＝",
		"|", "｜",
		"\\", "＼",
		"<", "＜",
		">", "＞",
		"[", "［",
		"]", "］",
		"{", "｛",
		"}", "｝",
		"_", "＿",
	)
	return replacer.Replace(raw)
}
