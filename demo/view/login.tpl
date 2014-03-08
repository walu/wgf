{{wgfInclude "header.tpl" .}}
<form action="{{.urlLogin}}" method="post">
	<input type="text" placeholder="user name" name = "uname" />
	<input type="submit" />
</form>
{{wgfInclude "footer.tpl" .}}
