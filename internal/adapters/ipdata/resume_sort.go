// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.


package ipdata

import "github.com/eterline/micro-utils/internal/models"

type ResumeInfo struct {
	ResolveDurationMs int64                  `json:"resolve_duration_ms" yaml:"resolve_duration_ms"`
	Resumes           []models.ResumeAboutIP `json:"resumes,omitempty" yaml:"resumes,omitempty"`
	NameServers       []string               `json:"ns,omitempty" yaml:"ns,omitempty"`
	ErrorIPs          string                 `json:"ip_error,omitempty" yaml:"ip_error,omitempty"`
	ErrorNS           string                 `json:"ns_error,omitempty" yaml:"ns_error,omitempty"`
}

func SortResolvedAndResume(res map[string]models.AboutResolve, rsvl []models.ResumeAboutIP) map[string]ResumeInfo {
	result := make(map[string]ResumeInfo)

	for key, resolve := range res {
		info := ResumeInfo{
			ResolveDurationMs: resolve.ResolveDurationMs,
			NameServers:       resolve.NameServers,
			ErrorIPs:          resolve.ErrorIPs,
			ErrorNS:           resolve.ErrorNS,
		}

		var matched []models.ResumeAboutIP
		for _, resume := range rsvl {
			for _, ip := range resolve.IPs {
				if resume.RequestIP.Equal(ip) {
					matched = append(matched, resume)
					break
				}
			}
		}
		info.Resumes = matched

		result[key] = info
	}

	return result
}
