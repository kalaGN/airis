/*
 * @Author: afei
 * @Date: 2022-07-14 09:47:44
 * @LastEditors: afei
 * @LastEditTime: 2022-07-14 13:54:32
 * @Description:
 *
 * Copyright (c) 2022 , All Rights Reserved.
 */
package config

// ConfigFunc 动态加载配置信息
type ConfigFunc func() map[string]interface{}

// ConfigFuncs 先加载到此数组，loadConfig 再动态生成配置信息
var ConfigFuncs map[string]ConfigFunc

func GetProduction() {
	// 已废弃，使用 app_config.go 中的新配置管理
}
