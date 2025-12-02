{{/*
Basename for mcmtp resources
{{ include "mcmtp.basename" . }}
*/}}
{{- define "mcmtp.basename" -}}
mcmtp-{{ include "basename" . }}
{{- end -}}

{{/*
Kubernetes standard labels for mcmtp
labels: {{ include "mcmtp.labels" . | nindent 4 }}
*/}}
{{- define "mcmtp.labels" -}}
{{ include "labels" . }}
app.kubernetes.io/component: mcmtp
# app.metaways.net/software: 
app.kubernetes.io/version: {{ .Values.mcmtp.deployment.image.tag | quote }}
{{- end -}}

{{/*
Labels used to match mcmtp labels
selector: {{- include "mcmtp.matchLabels" . | nindent 4 }}
*/}}
{{- define "mcmtp.matchLabels" -}}
{{include "matchLabels" . }}
app.kubernetes.io/component: mcmtp
{{- end -}}

{{/*
Secret name
secretName: {{ template "mcmtp.secret" . }}
*/}}
{{- define "mcmtp.secret" -}}
{{- if .Values.mcmtp.secrets.existingSecret -}}
{{ .Values.mcmtp.secrets.existingSecret }}
{{- else -}}
{{ include "mcmtp.basename" . }}
{{- end -}}
{{- end -}}

{{/*
Pvc name
claimName: {{ template "mcmtp.pvc" . }}
*/}}
{{- define "mcmtp.pvc" -}}
{{- if .Values.mcmtp.persistence.existingPvc -}}
{{ .Values.mcmtp.persistence.existingPvc }}
{{- else -}}
{{ include "mcmtp.basename" . }}
{{- end -}}
{{- end -}}