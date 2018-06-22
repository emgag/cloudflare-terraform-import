package dns

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/emgag/cloudflare-terraform-import/internal/lib/util"
	"io"
	"strings"
	"text/template"
)

var recordTemplate = `
resource "cloudflare_record" "{{ .ResourceName }}" {
  // record id: {{ .Zone }}/{{ .RR.ID }}
  domain = "{{ .Zone }}"
  name   = "{{ .RR.Name | stripdomain }}"
  type   = "{{ .RR.Type }}"
  {{ if eq .RR.Type "CAA" }}data = {
    flags = "{{ index .RR.Data "flags" }}"
    tag = "{{ index .RR.Data "tag" }}"
    value = "{{ index .RR.Data "value" }}"
  }
  {{- else -}}
  value  = "{{ .RR.Content | escape }}"{{ end }}
  ttl    = {{ .RR.TTL }}
  {{- if eq .RR.Proxied true }}
  proxied = {{ .RR.Proxied }}
  {{- else }}

  {{- end }}
  {{- if gt .RR.Priority 0 }}
  priority = {{ .RR.Priority }}
  {{- else }}

  {{- end }}
}
`

type Exporter struct {
	API *cloudflare.API
}

func (e *Exporter) DumpZone(zone string, tfWriter io.Writer, shWriter io.Writer) error {
	id, err := e.API.ZoneIDByName(zone)

	if err != nil {
		return err
	}

	recs, err := e.API.DNSRecords(id, cloudflare.DNSRecord{})

	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"escape": func(in string) string { return strings.Replace(in, "\"", "\\\"", -1) },
		"stripdomain": func(in string) string {
			return strings.Replace(in, "."+zone, "", 1)
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).Parse(recordTemplate)

	rrCounter := make(map[string]int)

	for _, r := range recs {
		t := fmt.Sprintf(
			"%s-%s",
			r.Name,
			r.Type,
		)

		rrCounter[t]++

		if rrCounter[t] > 1 {
			t = t + "-" + fmt.Sprintf("%d", rrCounter[t])
		}

		t = util.ToResourceName(t)

		err := tmpl.Execute(tfWriter, struct {
			Zone         string
			RR           cloudflare.DNSRecord
			ResourceName string
		}{zone, r, t})

		if err != nil {
			return err
		}

		fmt.Fprintf(shWriter, "terraform import cloudflare_record.%s %s/%s\n", t, zone, r.ID)
	}

	return nil
}

func NewExporter(api *cloudflare.API) *Exporter {
	return &Exporter{API: api}
}
