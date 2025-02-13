// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/externalversions/core/v1beta1"
	gardenmanagedseedinformers "github.com/gardener/gardener/pkg/client/seedmanagement/informers/externalversions/seedmanagement/v1alpha1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getGardenMetricsDefinitions() map[string]*prometheus.Desc {
	return map[string]*prometheus.Desc{
		metricGardenManagedSeedInfo: prometheus.NewDesc(
			metricGardenManagedSeedInfo,
			"Information about a managed seed.",
			[]string{
				"name",
				"shoot",
			},
			nil,
		),

		metricGardenOperationsTotal: prometheus.NewDesc(
			metricGardenOperationsTotal,
			"Count of ongoing operations.",
			[]string{
				"operation",
				"state",
				"iaas",
				"seed",
				"version",
				"region",
			},
			nil,
		),

		metricGardenPlantCondition: prometheus.NewDesc(
			metricGardenPlantCondition,
			"Condition state of a Plant. Possible values: -1=Unknown|0=Unhealthy|1=Healthy|2=Progressing",
			[]string{
				"name",
				"project",
				"condition",
			},
			nil,
		),

		metricGardenPlantInfo: prometheus.NewDesc(
			metricGardenPlantInfo,
			"Information about a plant.",
			[]string{
				"name",
				"project",
				"provider",
				"region",
				"version",
			},
			nil,
		),

		metricGardenProjectsStatus: prometheus.NewDesc(
			metricGardenProjectsStatus,
			"Status of projects.",
			[]string{
				"name",
				"cluster",
				"phase",
			},
			nil,
		),

		metricGardenSeedCondition: prometheus.NewDesc(
			metricGardenSeedCondition,
			"Condition state of a Seed.",
			[]string{
				"name",
				"condition",
			},
			nil,
		),

		metricGardenSeedInfo: prometheus.NewDesc(
			metricGardenSeedInfo,
			"Information about a Seed.",
			[]string{
				"name",
				"namespace",
				"iaas",
				"region",
				"visible",
				"protected",
			},
			nil,
		),

		metricGardenSeedCapacity: prometheus.NewDesc(
			metricGardenSeedCapacity,
			"Seed capacity.",
			[]string{
				"name",
				"namespace",
				"iaas",
				"region",
				"visible",
				"protected",
				"resource",
			},
			nil,
		),

		metricGardenSeedUsage: prometheus.NewDesc(
			metricGardenSeedUsage,
			"Seed usage.",
			[]string{
				"name",
				"namespace",
				"iaas",
				"region",
				"visible",
				"protected",
				"resource",
			},
			nil,
		),

		metricGardenShootCondition: prometheus.NewDesc(
			metricGardenShootCondition,
			"Condition state of a Shoot. Possible values: -1=Unknown|0=Unhealthy|1=Healthy|2=Progressing",
			[]string{
				"name",
				"project",
				"condition",
				"operation",
				"purpose",
				"is_seed",
				"iaas",
				"seed",
				"seed_iaas",
				"seed_region",
				"uid",
				"has_user_errors",
				"is_compliant",
			},
			nil,
		),

		metricGardenShootCreation: prometheus.NewDesc(
			metricGardenShootCreation,
			"Timestamp of the shoot creation.",
			[]string{
				"name",
				"project",
				"uid",
			},
			nil,
		),

		metricGardenShootHibernated: prometheus.NewDesc(
			metricGardenShootHibernated,
			"Hibernation status of a shoot.",
			[]string{
				"name",
				"project",
				"uid",
			},
			nil,
		),

		metricGardenShootInfo: prometheus.NewDesc(
			metricGardenShootInfo,
			"Information about a Shoot.",
			[]string{
				"name",
				"project",
				"iaas",
				"version",
				"region",
				"seed",
				"is_seed",
				"seed_iaas",
				"seed_region",
			},
			nil,
		),

		metricGardenShootNodeMaxTotal: prometheus.NewDesc(
			metricGardenShootNodeMaxTotal,
			"Max node count of a Shoot.",
			[]string{"name",
				"project",
			},
			nil,
		),

		metricGardenShootNodeMinTotal: prometheus.NewDesc(
			metricGardenShootNodeMinTotal,
			"Min node count of a Shoot.",
			[]string{"name",
				"project"},
			nil,
		),

		metricGardenShootNodeInfo: prometheus.NewDesc(
			metricGardenShootNodeInfo,
			"Information about the nodes in a Shoot.",
			[]string{
				"name",
				"project",
				"worker_group",
				"image",
				"version",
				"cri",
				"container_runtimes",
			},
			nil,
		),

		metricGardenShootOperationProgressPercent: prometheus.NewDesc(
			metricGardenShootOperationProgressPercent,
			"Operation progress percent of a Shoot.",
			[]string{
				"name",
				"project",
				"operation",
			},
			nil,
		),

		metricGardenShootOperationState: prometheus.NewDesc(
			metricGardenShootOperationState,
			"Operation state of a Shoot.",
			[]string{"name",
				"project",
				"operation",
			},
			nil,
		),

		metricGardenUsersSum: prometheus.NewDesc(
			metricGardenUsersSum,
			"Count of users.",
			[]string{
				"kind",
			},
			nil,
		),
	}
}

type gardenMetricsCollector struct {
	managedSeedInformer gardenmanagedseedinformers.ManagedSeedInformer
	shootInformer       gardencoreinformers.ShootInformer
	seedInformer        gardencoreinformers.SeedInformer
	projectInformer     gardencoreinformers.ProjectInformer
	plantInformer       gardencoreinformers.PlantInformer
	descs               map[string]*prometheus.Desc
	logger              *logrus.Logger
}

// Describe implements the prometheus.Describe interface, which intends the gardenMetricsCollector to be a Prometheus collector.
func (c *gardenMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.descs {
		ch <- desc
	}
	registerShootCustomizationMetrics(ch)
}

// Collect implements the prometheus.Collect interface, which intends the gardenMetricsCollector to be a Prometheus collector.
// TODO Can we run the collectors in parallel?
func (c *gardenMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	c.collectManagedSeedMetrics(ch)
	c.collectProjectMetrics(ch)
	c.collectShootMetrics(ch)
	c.collectSeedMetrics(ch)
	c.collectPlantMetrics(ch)
}

// SetupMetricsCollector takes informers to configure the metrics collectors.
func SetupMetricsCollector(shootInformer gardencoreinformers.ShootInformer, seedInformer gardencoreinformers.SeedInformer, projectInformer gardencoreinformers.ProjectInformer, plantInformer gardencoreinformers.PlantInformer, managedSeedInformer gardenmanagedseedinformers.ManagedSeedInformer, logger *logrus.Logger) {
	metricsCollector := gardenMetricsCollector{
		managedSeedInformer: managedSeedInformer,
		shootInformer:       shootInformer,
		seedInformer:        seedInformer,
		projectInformer:     projectInformer,
		plantInformer:       plantInformer,
		descs:               getGardenMetricsDefinitions(),
		logger:              logger,
	}
	prometheus.MustRegister(&metricsCollector)
	prometheus.MustRegister(ScrapeFailures)
}
