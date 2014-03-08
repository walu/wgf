{{wgfInclude "header.tpl" .}}

<h2>uname: {{.sessionUname}}</h2>

<h3>Links</h3>
<ul>
{{range .links}}
	<li><a href="{{.href}}" target="_blank">{{.name}}</a></li>
{{end}}
</ul>
{{wgfInclude "footer.tpl" .}}
