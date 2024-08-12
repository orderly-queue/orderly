env "postgres" {
  src = "file://schema.hcl"
  dev = "docker://postgres/latest/orderly"
  migration {
    dir    = "file://migrations/files"
    format = golang-migrate
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
