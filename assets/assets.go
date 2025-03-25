package assets

import (
	"embed"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	//go:embed manifests/*
	manifests  embed.FS
	appsScheme = runtime.NewScheme()
	appsCodecs = serializer.NewCodecFactory(appsScheme)
)

func init() {
	if err := appsv1.AddToScheme(appsScheme); err != nil {
		panic(err)
	}
}

func GetDeploymentFromFile(name string) (*appsv1.Deployment, error) {
	deploymentBytes, err := manifests.ReadFile(name)
	if err != nil {
		return nil, err
	}
	deploymentObject, err := runtime.Decode(
		appsCodecs.UniversalDecoder(appsv1.SchemeGroupVersion),
		deploymentBytes,
	)
	if err != nil {
		return nil, err
	}
	return deploymentObject.(*appsv1.Deployment), nil
}
