package client

var Skeleton string = `<!DOCTYPE html>
	<html>
	<head>
		<script src="/resources/htmx.js"></script>
		<script src="/resources/json-enc.js"></script>
		<script src="/resources/script.js"></script>
		<link rel="stylesheet" type="text/css" href="/resources/style.css">
		<title>Beango Messenger</title>
	</head>
	<body>
		<div id="header">{{block "header" .}}{{end}}</div>
		<div id="content">{{template "content" .}}</div>
		<div id="footer">{{block "header" .}}{{end}}</div>
	</body>
	</html>`

var LoginPage string = `{{define "content"}}<span class="logo"><span>> Beango Messenger </span></span>
	<div id="login-form">
		<form hx-ext="json-enc">
			<div class="form-row">
				<label for="username">Username:</label>
				<input type="text" name="username">
			</div>
			<div class="form-row">
				<label for="password">Password:</label>
				<input type="password" name="password">
			</div>
			<div class="form-row button-row">
				<button hx-post="/login/login" type="submit" hx-swap="none">Log In</button>
				<button hx-post="/login/signup" type="submit" hx-swap="none">Sign Up</button>
			</div>
			<div id="errors" class="error"></div>
		</form>
	</div>{{end}}`

var HomePage string = `{{define "content"}}<div id="chat-container">
	<div id="sidebar">
		<span class="heading-1">Chats</span>
		<div id=chat-list class="homepage-column chat-list">
			{{ range .Chats }}
			{{ block "chat-list" .}}
				<div 
					class="chat-selector list-item"
					hx-get="/home/chat/{{ .ID }}?name={{ .Name }}" 
					hx-target="#chat"
				>
					[{{ .Type}}] <b>{{ .Name }}</b>
				</div>
			{{ end }}
			{{ end }}
		</div>
	</div>
	<div id="chat"></div>
	</div>{{end}}`

var MessagePane string = `<span class="heading-1">{{ .Name }}</span>
	<table
		id="message-table"
		class="homepage-column message-list"
	>` + messageRows + `</table>
	<div class="message-bar">
		<span class="message-prompt">> </span>
		<textarea
			class="message-input"
			placeholder="Type your message"
			hx-on:keypress="sendMessageOnEnter(event)"
			name="content"
			hx-post="/home/chat/{{ .ID }}/sendMessage"
			hx-trigger="send-message"
			hx-swap="none"
			hx-on::after-request="if(event.detail.successful) this.value = '';"
			hx-ext="json-enc"
		></textarea>
	</div>` + newMessageFetcher

var MessagePaneRefresh string = newMessageFetcher + `
	<table id="message-table" hx-swap-oob="afterbegin">` + messageRows + `</table>`

var MessagePaneScroll string = `<table id="message-table" hx-swap-oob="beforeend">` +
	messageRows + `</table>`

var messageRows string = `
	{{ range $i, $m := .Messages }}
		{{ if and (eq $i 0) (not (or $.IsRefresh false)) }}
			<tr
				hx-get="/home/chat/{{ $.ID }}/scrollUp?to={{ $.ToMessageID }}"
				hx-swap="none"
				hx-trigger="intersect once"
				class="list-item"
			>
				<td class="cue">{{ $m.UserDisplayName }}</td>
				<td class="message">{{ $m.Content }}</td>
			</tr>
		{{ else }}
			<tr class="list-item">
				<td class="cue">{{ $m.UserDisplayName }}</td>
				<td class="message">{{ $m.Content }}</td>
			</tr>
		{{ end }}
	{{ end }}`

var newMessageFetcher string = `<div 
	hx-get="/home/chat/{{ .ID }}/refresh?from={{ .FromMessageID }}"
	hx-swap="outerHTML"
	class="chat-selector list-item"
	hx-trigger="chat-refresh from:document, every 5s"
	/>`
