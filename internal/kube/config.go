package kube

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Kubeconfig loads the rest configuration needed by k8s clients to interact with clusters.
// Loading rules are based on standard defined kubernetes config loading.
func Kubeconfig(url, file string) (*rest.Config, error) {
	// Default PathOptions gets kubeconfig in this order: the explicit path given, KUBECONFIG current context, recommentded file path
	po := clientcmd.NewDefaultPathOptions()
	po.LoadingRules.ExplicitPath = file

	return clientcmd.BuildConfigFromKubeconfigGetter(url, po.GetStartingConfig)
}

// Append adds the provided kubeconfig in the []byte to the Kubeconfig in the target path without altering other existing conifgs.
// If the target path is empty, standard kubeconfig loading rules apply.
func AppendConfig(cfg []byte, target string) error {
	s, err := clientcmd.Load(cfg)
	if err != nil {
		return err
	}

	// Default PathOptions gets kubeconfig in this order: the explicit path given, KUBECONFIG current context, recommentded file path
	po := clientcmd.NewDefaultPathOptions()
	po.LoadingRules.ExplicitPath = target

	t, err := po.GetStartingConfig()
	if err != nil {
		return err
	}

	// append contexts
	for k, v := range s.Contexts {
		t.Contexts[k] = v
	}

	// append clusters
	for k, v := range s.Clusters {
		t.Clusters[k] = v
	}

	// append authinfos
	for k, v := range s.AuthInfos {
		t.AuthInfos[k] = v
	}

	t.CurrentContext = s.CurrentContext

	// write config back
	return clientcmd.ModifyConfig(po, *t, false)
}

// RemoveConfig remoes the provided kubeconfig in the []byte from the Kubeconfig in the target path without altering other existing conifgs.
// If the target path is empty, standard kubeconfig loading rules apply.
func RemoveConfig(cfg []byte, target string) error {
	s, err := clientcmd.Load(cfg)
	if err != nil {
		return err
	}

	// Default PathOptions gets kubeconfig in this order: the explicit path given, KUBECONFIG current context, recommentded file path
	po := clientcmd.NewDefaultPathOptions()
	po.LoadingRules.ExplicitPath = target

	t, err := po.GetStartingConfig()
	if err != nil {
		return err
	}

	// remove contexts
	for k, _ := range s.Contexts {
		delete(t.Contexts, k)
	}

	// remove clusters
	for k, _ := range s.Clusters {
		delete(t.Clusters, k)
	}

	// remove authinfos
	for k, _ := range s.AuthInfos {
		delete(t.AuthInfos, k)
	}

	t.CurrentContext = ""

	// write config back
	return clientcmd.ModifyConfig(po, *t, false)
}