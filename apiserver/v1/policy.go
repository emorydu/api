// Copyright 2023 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"encoding/json"
	"fmt"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"github.com/ory/ladon"
	"gorm.io/gorm"
)

// AuthzPolicy defines iam policy type.
type AuthzPolicy struct {
	ladon.DefaultPolicy
}

// Policy represents a policy restful resource, include a ladon policy.
// It is also used as gorm model.
type Policy struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// The user of the policy
	Username string `json:"username" gorm:"column:username" validate:"omitempty"`

	// AuthzPolicy policy, will not be stored in db.
	Policy AuthzPolicy `json:"policy,omitempty" gorm:"-" validate:"omitempty"`

	// The ladon policy content, just a string format of ladon.DefaultPolicy. DO NOT modify directly.
	PolicyShadow string `json:"-" gorm:"column:policyShadow" validate:"omitempty"`
}

// PolicyList is the whole list of all policies which have been stored in storage.
type PolicyList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	metav1.ListMeta `json:",inline"`

	// List of policies.
	Items []*Policy `json:"items"`
}

// TableName maps to mysql table name.
func (p *Policy) TableName() string {
	return "policy"
}

func (ap AuthzPolicy) String() string {
	data, _ := json.Marshal(ap)

	return string(data)
}

// BeforeCreate run before create database record.
func (p *Policy) BeforeCreate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeCreate(tx); err != nil {
		return fmt.Errorf("failed to run `BeforeCreate` hook: %w", err)
	}

	p.Policy.ID = p.Name
	p.PolicyShadow = p.Policy.String()

	return nil
}

// AfterCreate run after create database record.
func (p *Policy) AfterCreate(tx *gorm.DB) error {
	p.InstanceID = idutil.GetInstanceID(p.ID, "policy-")

	return tx.Save(p).Error
}

// BeforeUpdate run before update database record.
func (p *Policy) BeforeUpdate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeUpdate(tx); err != nil {
		return fmt.Errorf("failed to run `BeforeUpdate` hook: %w", err)
	}

	p.Policy.ID = p.Name
	p.PolicyShadow = p.Policy.String()

	return nil
}

// AfterFind run after find to unmarshal a policy string into ladon.DefaultPolicy struct.
func (p *Policy) AfterFind(tx *gorm.DB) error {
	if err := p.ObjectMeta.AfterFind(tx); err != nil {
		return fmt.Errorf("failed to run `AfterFind` hook: %w", err)
	}

	if err := json.Unmarshal([]byte(p.PolicyShadow), &p.Policy); err != nil {
		return fmt.Errorf("failed to unmarshal policyShadow: %w", err)
	}

	return nil
}
