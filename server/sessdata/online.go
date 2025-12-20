package sessdata

import (
	"bytes"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/bjdgyc/anylink/dbdata"
	"github.com/bjdgyc/anylink/pkg/utils"
)

type Online struct {
	Token             string    `json:"token"`
	Username          string    `json:"username"`
	Nickname          string    `json:"nickname"`
	Email             string    `json:"email"`
	Group             string    `json:"group"`
	MacAddr           string    `json:"mac_addr"`
	UniqueMac         bool      `json:"unique_mac"`
	Ip                net.IP    `json:"ip"`
	RemoteAddr        string    `json:"remote_addr"`
	TransportProtocol string    `json:"transport_protocol"`
	TunName           string    `json:"tun_name"`
	Mtu               int       `json:"mtu"`
	Client            string    `json:"client"`
	BandwidthUp       string    `json:"bandwidth_up"`
	BandwidthDown     string    `json:"bandwidth_down"`
	BandwidthUpAll    string    `json:"bandwidth_up_all"`
	BandwidthDownAll  string    `json:"bandwidth_down_all"`
	LastLogin         time.Time `json:"last_login"`
}

type Onlines []Online

// UserExtraInfo 用户扩展信息（姓名和邮箱）
type UserExtraInfo struct {
	Nickname string
	Email    string
}

func (o Onlines) Len() int {
	return len(o)
}

func (o Onlines) Less(i, j int) bool {
	return bytes.Compare(o[i].Ip, o[j].Ip) < 0
}

func (o Onlines) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func OnlineSess() []Online {
	return GetOnlineSess("", "", false)
}

/**
 * @Description: GetOnlineSess
 * @param search_cate 分类：用户名、登录组、MAC地址、IP地址、远端地址
 * @param search_text 关键字，模糊搜索
 * @param show_sleeper 是否显示休眠用户
 * @return []Online
 */
func GetOnlineSess(search_cate string, search_text string, show_sleeper bool) []Online {
	var datas Onlines
	if strings.TrimSpace(search_text) == "" {
		search_cate = ""
	}

	// 批量获取本地用户扩展信息
	userInfoMap := getLocalUserInfo()

	// 构造在线用户列表
	sessMux.Lock()
	defer sessMux.Unlock()

	for _, v := range sessions {
		v.mux.Lock()
		cSess := v.CSess
		if cSess == nil {
			cSess = &ConnSession{}
		}

		// 获取用户扩展信息
		userInfo := userInfoMap[v.Username]

		// 选择需要比较的字符串（用于搜索）
		var compareText string
		switch search_cate {
		case "username":
			compareText = v.Username
		case "nickname":
			compareText = userInfo.Nickname
		case "email":
			compareText = userInfo.Email
		case "group":
			compareText = v.Group
		case "mac_addr":
			compareText = v.MacAddr
		case "ip":
			if cSess != nil {
				compareText = cSess.IpAddr.String()
			}
		case "remote_addr":
			if cSess != nil {
				compareText = cSess.RemoteAddr
			}
		}

		if search_cate != "" && !strings.Contains(compareText, search_text) {
			v.mux.Unlock()
			continue
		}

		if show_sleeper || v.IsActive {
			transportProtocol := "TCP"
			dSess := cSess.GetDtlsSession()
			if dSess != nil {
				transportProtocol = "UDP"
			}
			val := Online{
				Token:             v.Token,
				Ip:                cSess.IpAddr,
				Username:          v.Username,
				Nickname:          userInfo.Nickname,
				Email:             userInfo.Email,
				Group:             v.Group,
				MacAddr:           v.MacAddr,
				UniqueMac:         v.UniqueMac,
				RemoteAddr:        cSess.RemoteAddr,
				TransportProtocol: transportProtocol,
				TunName:           cSess.IfName,
				Mtu:               cSess.Mtu,
				Client:            cSess.Client,
				BandwidthUp:       utils.HumanByte(cSess.BandwidthUpPeriod.Load()) + "/s",
				BandwidthDown:     utils.HumanByte(cSess.BandwidthDownPeriod.Load()) + "/s",
				BandwidthUpAll:    utils.HumanByte(cSess.BandwidthUpAll.Load()),
				BandwidthDownAll:  utils.HumanByte(cSess.BandwidthDownAll.Load()),
				LastLogin:         v.LastLogin,
			}
			datas = append(datas, val)
		}
		v.mux.Unlock()
	}
	sort.Sort(&datas)
	return datas
}

// getLocalUserInfo 批量获取本地认证用户的扩展信息（nickname和email）
func getLocalUserInfo() map[string]UserExtraInfo {
	userInfoMap := make(map[string]UserExtraInfo)

	// 收集所有本地认证用户的username
	sessMux.Lock()
	usernameSet := make(map[string]bool)
	for _, v := range sessions {
		v.mux.Lock()
		if v.IsActive && v.AuthType == "local" {
			usernameSet[v.Username] = true
		}
		v.mux.Unlock()
	}
	sessMux.Unlock()

	// 批量查询数据库
	if len(usernameSet) > 0 {
		usernames := make([]string, 0, len(usernameSet))
		for username := range usernameSet {
			usernames = append(usernames, username)
		}

		var users []dbdata.User
		err := dbdata.GetXdb().In("username", usernames).Cols("username", "nickname", "email").Find(&users)
		if err == nil {
			for _, u := range users {
				userInfoMap[u.Username] = UserExtraInfo{
					Nickname: u.Nickname,
					Email:    u.Email,
				}
			}
		}
	}

	return userInfoMap
}
