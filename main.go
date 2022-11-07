package main

import (
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/kustomize"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		region := "LON1"

		network, err := civo.NewNetwork(ctx, "network", &civo.NetworkArgs{
			Label:  pulumi.String("k8s_social"),
			Region: pulumi.String(region),
		})
		if err != nil {
			return err
		}

		firewall, err := civo.NewFirewall(ctx, "firewall", &civo.FirewallArgs{
			NetworkId:          network.ID(),
			Region:             pulumi.String(region),
			CreateDefaultRules: pulumi.Bool(false),
			IngressRules: civo.FirewallIngressRuleArray{
				&civo.FirewallIngressRuleArgs{
					Label:     pulumi.String("k8s"),
					Protocol:  pulumi.String("tcp"),
					PortRange: pulumi.String("6443"),
					Cidrs:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Action:    pulumi.String("allow"),
				},
				&civo.FirewallIngressRuleArgs{
					Label:     pulumi.String("https-tcp"),
					Protocol:  pulumi.String("tcp"),
					PortRange: pulumi.String("443"),
					Cidrs:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Action:    pulumi.String("allow"),
				},
				&civo.FirewallIngressRuleArgs{
					Label:     pulumi.String("https-udp"),
					Protocol:  pulumi.String("udp"),
					PortRange: pulumi.String("443"),
					Cidrs:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Action:    pulumi.String("allow"),
				},
				&civo.FirewallIngressRuleArgs{
					Label:     pulumi.String("http-tcp"),
					Protocol:  pulumi.String("tcp"),
					PortRange: pulumi.String("80"),
					Cidrs:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Action:    pulumi.String("allow"),
				},
				&civo.FirewallIngressRuleArgs{
					Label:     pulumi.String("http-udp"),
					Protocol:  pulumi.String("udp"),
					PortRange: pulumi.String("80"),
					Cidrs:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Action:    pulumi.String("allow"),
				},
			},
			EgressRules: civo.FirewallEgressRuleArray{
				&civo.FirewallEgressRuleArgs{
					Label:     pulumi.String("all-tcp"),
					Protocol:  pulumi.String("tcp"),
					PortRange: pulumi.String("1-65535"),
					Cidrs:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Action:    pulumi.String("allow"),
				},
				&civo.FirewallEgressRuleArgs{
					Label:     pulumi.String("all-udp"),
					Protocol:  pulumi.String("udp"),
					PortRange: pulumi.String("1-65535"),
					Cidrs:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Action:    pulumi.String("allow"),
				},
			},
		})
		if err != nil {
			return err
		}

		k8SCluster, err := civo.NewKubernetesCluster(ctx, "k8sCluster", &civo.KubernetesClusterArgs{
			Name:       pulumi.String("k8s.social"),
			FirewallId: firewall.ID(),
			NetworkId:  network.ID(),
			Pools: &civo.KubernetesClusterPoolsArgs{
				NodeCount: pulumi.Int(3),
				Size:      pulumi.String("g4s.kube.small"),
			},
			Region:       pulumi.String(region),
			Cni:          pulumi.String("cilium"),
			Applications: pulumi.String("Nginx"),
		})
		if err != nil {
			return err
		}

		objectStoreCredsResource, err := civo.NewObjectStoreCredential(ctx, "objectStoreCreds", &civo.ObjectStoreCredentialArgs{
			Region: pulumi.String(region),
			Name:   pulumi.String("k8s.social-media"),
		})
		if err != nil {
			return err
		}

		_, err = civo.NewObjectStore(ctx, "media", &civo.ObjectStoreArgs{
			MaxSizeGb:   pulumi.Int(500),
			Region:      pulumi.String(region),
			Name:        pulumi.String("k8s.social-media"),
			AccessKeyId: objectStoreCredsResource.AccessKeyId,
		})
		if err != nil {
			return err
		}

		_, err = civo.NewDnsDomainName(ctx, "domainName", &civo.DnsDomainNameArgs{
			Name: pulumi.String("k8s.social"),
		})
		if err != nil {
			return err
		}

		kubernetesProvider, err := kubernetes.NewProvider(ctx, "kubernetesProvider", &kubernetes.ProviderArgs{
			Kubeconfig: k8SCluster.Kubeconfig,
		})
		if err != nil {
			return err
		}

		_, err = kustomize.NewDirectory(ctx, "argoCD",
			kustomize.DirectoryArgs{
				Directory: pulumi.String("./manifests/argocd"),
			},
			pulumi.Provider(kubernetesProvider),
		)
		if err != nil {
			return err
		}

		ctx.Export("objectStoreCreds.accessKeyId", objectStoreCredsResource.AccessKeyId)
		ctx.Export("objectStoreCreds.secretAccessKey", objectStoreCredsResource.SecretAccessKey)
		ctx.Export("kubernetesEndpoint", k8SCluster.ApiEndpoint)
		ctx.Export("kubeconfig", k8SCluster.Kubeconfig)

		return nil
	})
}
