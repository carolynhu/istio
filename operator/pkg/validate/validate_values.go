// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validate

import (
	"fmt"

	"github.com/ghodss/yaml"

	"istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	"istio.io/istio/operator/pkg/util"
)

var (
	// DefaultValuesValidations maps a data path to a validation function.
	DefaultValuesValidations = map[string]ValidatorFunc{
		"global.proxy.includeIPRanges":     validateIPRangesOrStar,
		"global.proxy.excludeIPRanges":     validateIPRangesOrStar,
		"global.proxy.includeInboundPorts": validateStringList(validatePortNumberString),
		"global.proxy.excludeInboundPorts": validateStringList(validatePortNumberString),
		"meshConfig":                       validateMeshConfig,
	}
)

// CheckValues validates the values in the given tree, which follows the Istio values.yaml schema.
func CheckValues(root interface{}) util.Errors {
	vs, err := yaml.Marshal(root)
	if err != nil {
		return util.Errors{err}
	}
	val := &v1alpha1.Values{}
	if err := util.UnmarshalWithJSONPB(string(vs), val, false); err != nil {
		return util.Errors{err}
	}
	if val.Global.DefaultPodDisruptionBudget.Enabled.Value && val.Pilot.AutoscaleEnabled.Value {
		// 2 derived from: the default PDB minAvaliable 1 plus 1
		if val.Pilot.AutoscaleMin < 2 || val.Pilot.ReplicaCount < 2 {
			return util.NewErrs(fmt.Errorf("Istiod HorizontalPodAutoscaler MinReplica is violating PDB minAvailable 1."))

		}
	}
	if val.Global.DefaultPodDisruptionBudget.Enabled.Value && val.Gateways.IstioIngressgateway.AutoscaleEnabled.Value {
		if val.Gateways.IstioIngressgateway.AutoscaleMin < 2 {
			return util.NewErrs(fmt.Errorf("IstioIngressgateway HorizontalPodAutoscaler MinReplica is violating PDB minAvailable 1."))

		}
	}
	// similar for IstioEgressgateway

	return ValuesValidate(DefaultValuesValidations, root, nil)
}

// ValuesValidate validates the values of the tree using the supplied Func
func ValuesValidate(validations map[string]ValidatorFunc, node interface{}, path util.Path) (errs util.Errors) {
	pstr := path.String()
	scope.Debugf("ValuesValidate %s", pstr)
	vf := validations[pstr]
	if vf != nil {
		errs = util.AppendErrs(errs, vf(path, node))
	}

	nn, ok := node.(map[string]interface{})
	if !ok {
		// Leaf, nothing more to recurse.
		return errs
	}
	for k, v := range nn {
		errs = util.AppendErrs(errs, ValuesValidate(validations, v, append(path, k)))
	}

	return errs
}
