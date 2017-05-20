package templates

const starlingStr = `<TextureAtlas imagePath="{{.ImageFilename}}">
{{- range .Sprites}}
    <SubTexture name="{{.Name}}" x="{{.Left}}" y="{{.Top}}" width="{{.Width}}" height="{{.Height}}"/>
{{- end}}
</TextureAtlas>`
