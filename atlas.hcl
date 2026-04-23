variable "envfile" {
  type    = string
  default = ".env"
}

locals {
  envfile = {
    for line in split("\n", file(var.envfile)) : split("=", line)[0] => trim(regex("=(.*)", line)[0], "\"")
    if !startswith(line, "#") && length(split("=", line)) > 1
  }
}

env "dev" {
  src = "ent://backend/ent/schema"
  url = local.envfile["DATABASE_URL"]
  dev = local.envfile["ATLAS_DATABASE_URL"]

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \" \" }}"
    }
  }
}

env "prod" {
  src = "ent://backend/ent/schema"
  url = local.envfile["DATABASE_URL"]
  dev = local.envfile["ATLAS_DATABASE_URL"]

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \" \" }}"
    }
  }
}
