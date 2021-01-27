/*
Copyright 2021 The Tekton Authors

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

package main

import (
	"context"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/configmap"
)

// TODO: move this to its own package, with tests.

type cfgKey struct{}

// Config holds the collection of configurations that we attach to contexts.
type Config struct {
	Concurrency *Concurrency
}

// FromContext extracts a Config from the provided context.
func FromContext(ctx context.Context) *Config {
	x, ok := ctx.Value(cfgKey{}).(*Config)
	if ok {
		return x
	}
	return nil
}

// FromContextOrDefaults is like FromContext, but when no Config is attached it
// returns a Config populated with the defaults for each of the Config fields.
func FromContextOrDefaults(ctx context.Context) *Config {
	if cfg := FromContext(ctx); cfg != nil {
		return cfg
	}
	conc, _ := NewConcurrencyFromMap(map[string]string{})
	return &Config{
		Concurrency: conc,
	}
}

// ToContext attaches the provided Config to the provided context, returning the
// new context with the Config attached.
func ToContext(ctx context.Context, c *Config) context.Context {
	return context.WithValue(ctx, cfgKey{}, c)
}

// Store is a typed wrapper around configmap.Untyped store to handle our configmaps.
type Store struct {
	*configmap.UntypedStore
}

// NewStore creates a new store of Configs and optionally calls functions when ConfigMaps are updated.
func NewStore(logger configmap.Logger, onAfterStore ...func(name string, value interface{})) *Store {
	return &Store{
		UntypedStore: configmap.NewUntypedStore(
			"concurrency",
			logger,
			configmap.Constructors{
				"config-concurrency": NewConcurrencyFromConfigMap,
			},
			onAfterStore...,
		),
	}
}

// ToContext attaches the current Config state to the provided context.
func (s *Store) ToContext(ctx context.Context) context.Context {
	return ToContext(ctx, s.Load())
}

// Load creates a Config from the current config state of the Store.
func (s *Store) Load() *Config {
	conc := s.UntypedLoad("config-concurrency")
	if conc == nil {
		conc, _ = NewConcurrencyFromMap(map[string]string{})
	}

	return &Config{
		Concurrency: conc.(*Concurrency).DeepCopy(),
	}
}

// NewConcurrencyFromConfigMap returns a Config for the given configmap
func NewConcurrencyFromConfigMap(config *corev1.ConfigMap) (*Concurrency, error) {
	return NewConcurrencyFromMap(config.Data)
}

// NewConcurrencyFromMap returns a Config given a map corresponding to a ConfigMap
func NewConcurrencyFromMap(cfgMap map[string]string) (*Concurrency, error) {
	c := Concurrency{
		Key:   "",
		Limit: -1,
	}

	if v, ok := cfgMap["concurrency-key"]; ok {
		c.Key = v
	}
	if v, ok := cfgMap["concurrency-limit"]; ok {
		i, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to parse limit: %v", err)
		}
		c.Limit = int(i)
	}
	return &c, nil
}

type Concurrency struct {
	Key   string
	Limit int
}

func (c *Concurrency) DeepCopy() *Concurrency {
	if c == nil {
		c, _ = NewConcurrencyFromMap(map[string]string{})
	}
	return &Concurrency{
		Key:   c.Key,
		Limit: c.Limit,
	}
}
