/*
Copyright 2018 The Kubernetes Authors All rights reserved.

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

package util

import (
	"testing"
	"time"
)

func TestGetStartTime(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name              string
		uptime            time.Duration
		lookback          string
		delay             string
		expectErr         bool
		expectedStartTime time.Time
	}{
		{
			name:              "bad lookback value",
			uptime:            0,
			lookback:          "abc",
			delay:             "",
			expectErr:         true,
			expectedStartTime: time.Time{},
		},
		{
			name:              "bad delay value",
			uptime:            0,
			lookback:          "",
			delay:             "abc",
			expectErr:         true,
			expectedStartTime: time.Time{},
		},
		{
			name:              "node is just up, no lookback and delay",
			uptime:            0,
			lookback:          "",
			delay:             "",
			expectErr:         false,
			expectedStartTime: now,
		},
		{
			name:              "no delay, lookback > uptime",
			uptime:            5 * time.Second,
			lookback:          "7s",
			delay:             "",
			expectErr:         false,
			expectedStartTime: now.Add(-5 * time.Second),
		},
		{
			name:              "no delay, lookback < uptime",
			uptime:            5 * time.Second,
			lookback:          "3s",
			delay:             "",
			expectErr:         false,
			expectedStartTime: now.Add(-3 * time.Second),
		},
		{
			name:              "no lookback, delay > uptime",
			uptime:            5 * time.Second,
			lookback:          "",
			delay:             "7s",
			expectErr:         false,
			expectedStartTime: now.Add(2 * time.Second),
		},
		{
			name:              "no lookback, delay < uptime",
			uptime:            5 * time.Second,
			lookback:          "",
			delay:             "3s",
			expectErr:         false,
			expectedStartTime: now,
		},
		{
			name:              "uptime < delay",
			uptime:            10 * time.Second,
			lookback:          "6s",
			delay:             "12s",
			expectErr:         false,
			expectedStartTime: now.Add(2 * time.Second),
		},
		{
			name:              "uptime > delay, uptime < lookback",
			uptime:            10 * time.Second,
			lookback:          "12s",
			delay:             "7s",
			expectErr:         false,
			expectedStartTime: now.Add(-3 * time.Second),
		},
		{
			name:              "uptime > delay, uptime > lookback, lookback > uptime - delay",
			uptime:            10 * time.Second,
			lookback:          "6s",
			delay:             "7s",
			expectErr:         false,
			expectedStartTime: now.Add(-3 * time.Second),
		},
		{
			name:              "uptime > delay, uptime > lookback, lookback < uptime - delay",
			uptime:            10 * time.Second,
			lookback:          "2s",
			delay:             "7s",
			expectErr:         false,
			expectedStartTime: now.Add(-2 * time.Second),
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			startTime, err := GetStartTime(now, test.uptime, test.lookback, test.delay)
			if test.expectErr && err == nil {
				t.Fatalf("Expect to get error, but got no returned error.")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("Expect to get no error, but got returned error: %v", err)
			}
			if test.expectedStartTime != startTime {
				t.Fatalf("Expect to get start time %v, but got %v", test.expectedStartTime, startTime)
			}
		})
	}
}

func TestGetOSVersion(t *testing.T) {
	testCases := []struct {
		name              string
		fakeOSReleasePath string
		expectedOSVersion string
		expectErr         bool
	}{
		{
			name:              "COS",
			fakeOSReleasePath: "testdata/os-release-cos",
			expectedOSVersion: "cos 77-12293.0.0",
			expectErr:         false,
		},
		{
			name:              "Debian",
			fakeOSReleasePath: "testdata/os-release-debian",
			expectedOSVersion: "debian 9 (stretch)",
			expectErr:         false,
		},
		{
			name:              "Ubuntu",
			fakeOSReleasePath: "testdata/os-release-ubuntu",
			expectedOSVersion: "ubuntu 16.04.6 LTS (Xenial Xerus)",
			expectErr:         false,
		},
		{
			name:              "centos",
			fakeOSReleasePath: "testdata/os-release-centos",
			expectedOSVersion: "centos 7 (Core)",
			expectErr:         false,
		},
		{
			name:              "rhel",
			fakeOSReleasePath: "testdata/os-release-rhel",
			expectedOSVersion: "rhel 7.7 (Maipo)",
			expectErr:         false,
		},
		{
			name:              "Unknown",
			fakeOSReleasePath: "testdata/os-release-unknown",
			expectedOSVersion: "",
			expectErr:         true,
		},
		{
			name:              "Empty",
			fakeOSReleasePath: "testdata/os-release-empty",
			expectedOSVersion: "",
			expectErr:         true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			originalOSReleasePath := osReleasePath
			defer func() {
				osReleasePath = originalOSReleasePath
			}()

			osReleasePath = test.fakeOSReleasePath
			osVersion, err := GetOSVersion()

			if test.expectErr && err == nil {
				t.Errorf("Expect to get error, but got no returned error.")
			}
			if !test.expectErr && err != nil {
				t.Errorf("Expect to get no error, but got returned error: %v", err)
			}
			if !test.expectErr && osVersion != test.expectedOSVersion {
				t.Errorf("Wanted: %+v. \nGot: %+v", test.expectedOSVersion, osVersion)
			}
		})
	}
}
