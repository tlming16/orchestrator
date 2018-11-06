/*
   Copyright 2014 Outbrain Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package logic

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var recoveryDisableHints *cache.Cache
var recentRecoveryAuthorizationRequests *cache.Cache

func SetupLeaderStates(checkAndRecoverWaitPeriod time.Duration) {
	recoveryDisableHints = cache.New(checkAndRecoverWaitPeriod, time.Second)
	recentRecoveryAuthorizationRequests = cache.New(time.Minute, time.Second)
}

func OnLeaderStateChange(isTurnedLeader bool) {
	if isTurnedLeader {
		recoveryDisableHints.Add("turned-leader", true, cache.DefaultExpiration)
		recoveryDisableHints.Delete("turned-non-leader")
	} else {
		for key, _ := range recentRecoveryAuthorizationRequests.Items() {
			recentRecoveryAuthorizationRequests.Delete(key)
		}
		recoveryDisableHints.Add("turned-non-leader", true, cache.NoExpiration)
	}
}

func recoveryDisabledHint() *string {
	items := recoveryDisableHints.Items()
	for k, _ := range items {
		return &k
	}
	return nil
}
