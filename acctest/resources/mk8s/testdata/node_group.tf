resource "mws_mk8s_node_group" "node_group" {
  cluster_name = "{{.ClusterName}}"
  node_group_name = "{{.Name}}"

  metadata = {
    display_name = "{{.Name}}"
    name         = "{{.Name}}"
    description  = "Test Node Group"
  }

  subnet = {
    ref = "{{.SubnetID}}"
  }

  vm_type = {
    ref = "compute/vmTypes/gen-2-8"
  }

  scale = {
    autoscaling = {
      min = {{.Autoscalling.Min}}
      max = {{.Autoscalling.Max}}
    }
  }

  service_account = {
    ref = "iam/projects/mws-terraform-testing/serviceAccounts/gitlab-acctest"
  }

  rollout_strategy = {
    max_surge       = {{.Rollout.Surge}}
    max_unavailable = {{.Rollout.Unavailable}}
  }

  version_control = {
    release_channel = "{{.Channel}}"
    auto_update = true
    version = "{{.Version}}"
    maintenance_window = {
      weekly = {
        days    = {{.GetDays}}
        hour    = 3
      }
    }
  }

  image_storage_size = "{{.StorageSize}}"

  taints = []

  zone = "{{.Zone}}"
}
