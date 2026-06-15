resource "mws_mk8s_cluster" "cluster" {
  availability = {
  {{ if .IsHA }}
    zonal_ha = {
  {{ else }}
    standalone = {
  {{ end }}
      zone = "{{.Zone}}"
    }
  }

  cluster_name = "{{.Name}}"

  metadata = {
    description  = "Kubernetes cluster"
    display_name = "Example Cluster"
  }

  network = {
    pods_cidr       = "{{.PodCIDR}}"
    services_cidr   = "{{.CIDR}}"
    primary_endpoint = {
      ref = "{{.Endpoint}}"
    }
  }

  version_control = {
    release_channel = "{{.Channel}}"
    maintenance_window = {
      weekly = {
        days    = {{.GetDays}}
        hour    = 3
      }
    }
    version = "v1.34.1-mws.1"
  }
}
