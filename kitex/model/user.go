/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id             int32  `json:"id" column:"id"`
	Name           string `json:"name" column:"name"`
	CollegeName    string `json:"college_name" column:"college_name"`
	CollegeAddress string `json:"college_address" column:"college_address"`
	Emails         string `json:"emails" column:"emails"`
}

func (u *User) TableName() string {
	return "students"
}
