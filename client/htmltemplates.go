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
	<body hx-on::before-request="clearErrorNodes();">
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
				<input type="text" name="username" maxlength="25">
			</div>
			<div class="form-row">
				<label for="password">Password:</label>
				<input type="password" name="password" maxlength="25">
			</div>
			<div class="form-row button-row">
				<button hx-post="/login/login" type="submit" hx-swap="none" class="form-button">Log In</button>
				<button hx-post="/login/signup" type="submit" hx-swap="none" class="form-button">Sign Up</button>
			</div>
			<div id="errors" class="error"></div>
		</form>
	</div>{{end}}`

var HomePage string = `{{define "content"}}
	<div id="errors" class="error"></div>
	<div class="chat-container">
		<div class="sidebar">
			<div class="heading-1">
				<span>Chats</span>
				<button type="submit" class="new-button" hx-get="/home/newChat" hx-target="#main-pane">
					Create
				</button>
			</div>
			<div id=chat-list class="homepage-column chat-list">
				{{ range .Chats }}
				{{ block "chat-list" .}}
					<div 
						class="chat-selector list-item"
						hx-get="/home/chat/{{ .ID }}?name={{ .Name }}" 
						hx-target="#main-pane"
					>
						[{{ .Type}}] <b>{{ .Name }}</b>
					</div>
				{{ end }}
				{{ end }}
			</div>
		</div>
		<div id="main-pane" class="main-pane"></div>
	</div>{{end}}`

var MessagePane string = `<span class="heading-1">{{ .Name }}</span>
	<table
		id="message-table"
		class="homepage-column message-list"
	>` + messageRows + `</table>
	<div class="input-bar">
		<span class="input-prompt">> </span>
		<textarea
			class="input-value message-input"
			placeholder="Type your message"
			hx-on:keypress="sendMessageOnEnter(event)"
			name="content"
			maxlength="5000"
			hx-post="/home/chat/{{ .ID }}/sendMessage"
			hx-trigger="send-message consume"
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

var NewChatPane string = `<span class="heading-1">Create a new chat</span>
	<div class="input-bar">
		<span class="input-prompt">> </span>
		<textarea
			class="input-value"
			placeholder="Search for a user to chat with"
			name="query"
			maxlength="25"
			hx-post="/home/newChat/search"
			hx-trigger="keyup changed delay:500ms"
			hx-target="#search-results"
			hx-ext="json-enc"
		></textarea>
		<div id="search-results"/>
	</div>`

var UserSearchResults string = `
	{{ if not .Users }}
		<span class="info">No results.</span>
	{{ else }}
		{{ range .Users }}
			<div 
				class="chat-selector list-item"
				hx-post="/home/newChat/create"
				hx-target="#main-pane"
				hx-vals='{"userID": {{ .ID }}}'
				hx-ext="json-enc" 
			>
				<b>{{ .Username }}</b> {{ .DisplayName }}
			</div>
		{{ end }}
	{{ end }}
`
